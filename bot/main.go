package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/bbolt"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
	"github.com/spaincoin/spaincoin/exchange/database"
)

var (
	hotWalletKey  *crypto.PrivateKey
	hotWalletAddr crypto.Address
	nodeRPCURL    string
)

const telegramAPI = "https://api.telegram.org/bot"

var (
	botToken   string
	orderDB    *database.OrderDB
	adminIDs   map[int64]string // chatID -> role ("super" or "admin")
	bankInfo   string           // IBAN or payment info
	priceTiers []priceTier      // automatic pricing
)

type priceTier struct {
	SoldUpTo float64
	Price    float64
}

// Default price tiers — override with /settiers
var defaultTiers = []priceTier{
	{500, 0.05},
	{1000, 0.08},
	{2500, 0.12},
	{5000, 0.18},
	{10000, 0.25},
	{25000, 0.40},
	{50000, 0.70},
	{100000, 1.00},
	{500000, 2.00},
	{1000000, 5.00},
}

// ==========================================
// Telegram types
// ==========================================

type TelegramUpdate struct {
	UpdateID      int64             `json:"update_id"`
	Message       *TelegramMessage  `json:"message"`
	CallbackQuery *TelegramCallback `json:"callback_query"`
}

type TelegramCallback struct {
	ID      string           `json:"id"`
	From    *TelegramUser    `json:"from"`
	Message *TelegramMessage `json:"message"`
	Data    string           `json:"data"`
}

type TelegramMessage struct {
	MessageID int64         `json:"message_id"`
	From      *TelegramUser `json:"from"`
	Chat      *TelegramChat `json:"chat"`
	Text      string        `json:"text"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type TelegramChat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

func main() {
	botToken = os.Getenv("SPC_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("SPC_BOT_TOKEN required")
	}

	// Parse admin IDs: "123:super,456:admin,789:admin"
	adminIDs = make(map[int64]string)
	adminStr := os.Getenv("SPC_ADMIN_IDS")
	if adminStr != "" {
		for _, entry := range strings.Split(adminStr, ",") {
			parts := strings.SplitN(strings.TrimSpace(entry), ":", 2)
			if len(parts) >= 1 {
				id, err := strconv.ParseInt(parts[0], 10, 64)
				if err != nil {
					continue
				}
				role := "admin"
				if len(parts) == 2 {
					role = parts[1]
				}
				adminIDs[id] = role
			}
		}
	}
	// Backward compat: single admin
	if len(adminIDs) == 0 {
		singleAdmin := os.Getenv("SPC_ADMIN_CHAT_ID")
		if singleAdmin != "" {
			id, _ := strconv.ParseInt(singleAdmin, 10, 64)
			if id > 0 {
				adminIDs[id] = "super"
			}
		}
	}

	bankInfo = os.Getenv("SPC_BANK_INFO")
	if bankInfo == "" {
		bankInfo = "(no configurado — usa /setbank IBAN)"
	}

	dataDir := os.Getenv("SPC_DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	os.MkdirAll(dataDir, 0700)

	db, err := bbolt.Open(filepath.Join(dataDir, "bot.db"), 0600, nil)
	if err != nil {
		log.Fatalf("open bot.db: %v", err)
	}
	defer db.Close()

	orderDB, err = database.NewOrderDB(db)
	if err != nil {
		log.Fatalf("init order db: %v", err)
	}

	priceTiers = defaultTiers

	// Group chat ID for daily reports (set via SPC_GROUP_CHAT_ID)
	groupChatStr := os.Getenv("SPC_GROUP_CHAT_ID")
	var groupChatID int64
	if groupChatStr != "" {
		groupChatID, _ = strconv.ParseInt(groupChatStr, 10, 64)
	}

	// Initialize hot wallet for automatic SPC sending
	nodeRPCURL = os.Getenv("SPC_NODE_URL")
	if nodeRPCURL == "" {
		nodeRPCURL = "http://204.168.176.40:8545"
	}
	hotKeyHex := os.Getenv("SPC_HOT_WALLET_KEY")
	if hotKeyHex != "" {
		priv, pub, err := crypto.PrivateKeyFromHex(hotKeyHex)
		if err != nil {
			log.Printf("[WARN] Invalid hot wallet key: %v", err)
		} else {
			hotWalletKey = priv
			hotWalletAddr = pub.ToAddress()
			log.Printf("Hot wallet: %s", hotWalletAddr.String())
		}
	} else {
		log.Println("[WARN] SPC_HOT_WALLET_KEY not set — auto-send disabled")
	}

	// Write price on startup so web always has it
	startPrice := getCurrentPrice()
	log.Printf("SpainCoin Bot starting... price=%.4f€", startPrice)
	log.Printf("Admins: %v", adminIDs)
	log.Printf("Group chat: %d", groupChatID)

	offset := int64(0)
	client := &http.Client{Timeout: 35 * time.Second}

	// Daily report at 9:00 AM (Europe/Madrid)
	go func() {
		for {
			now := time.Now()
			// Calculate next 9:00 AM
			next := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			time.Sleep(time.Until(next))

			if groupChatID != 0 {
				sendDailyReport(client, groupChatID)
			}
			// Also send to admins
			for id := range adminIDs {
				sendDailyReport(client, id)
			}

			time.Sleep(1 * time.Minute) // avoid double-send
		}
	}()

	for {
		updates, err := getUpdates(client, offset)
		if err != nil {
			log.Printf("getUpdates error: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		for _, u := range updates {
			if u.CallbackQuery != nil {
				handleCallback(client, u.CallbackQuery)
			} else if u.Message != nil {
				handleMessage(client, u.Message)
			}
			offset = u.UpdateID + 1
		}
	}
}

func getUpdates(client *http.Client, offset int64) ([]TelegramUpdate, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d&timeout=30", telegramAPI, botToken, offset)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool             `json:"ok"`
		Result []TelegramUpdate `json:"result"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Result, nil
}

func sendMessage(client *http.Client, chatID int64, text string) {
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPI, botToken)
	body := fmt.Sprintf(`{"chat_id":%d,"text":%s,"parse_mode":"HTML"}`, chatID, jsonStr(text))
	http.Post(url, "application/json", strings.NewReader(body))
}

// sendMessageAndTrack sends a message and returns the message ID for later deletion.
func sendMessageAndTrack(client *http.Client, chatID int64, text string, buttons [][]InlineButton) int64 {
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPI, botToken)
	kb, _ := json.Marshal(map[string]interface{}{"inline_keyboard": buttons})
	body := fmt.Sprintf(`{"chat_id":%d,"text":%s,"parse_mode":"HTML","reply_markup":%s}`, chatID, jsonStr(text), string(kb))
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageID int64 `json:"message_id"`
		} `json:"result"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Result.MessageID
}

func sendMessageWithButtons(client *http.Client, chatID int64, text string, buttons [][]InlineButton) {
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPI, botToken)
	kb, _ := json.Marshal(map[string]interface{}{"inline_keyboard": buttons})
	body := fmt.Sprintf(`{"chat_id":%d,"text":%s,"parse_mode":"HTML","reply_markup":%s}`, chatID, jsonStr(text), string(kb))
	http.Post(url, "application/json", strings.NewReader(body))
}

func answerCallback(client *http.Client, callbackID string) {
	url := fmt.Sprintf("%s%s/answerCallbackQuery", telegramAPI, botToken)
	body := fmt.Sprintf(`{"callback_query_id":"%s"}`, callbackID)
	http.Post(url, "application/json", strings.NewReader(body))
}

// InlineButton for Telegram inline keyboard.
type InlineButton struct {
	Text string `json:"text"`
	Data string `json:"callback_data,omitempty"`
	URL  string `json:"url,omitempty"`
}

func mainMenuButtons() [][]InlineButton {
	return [][]InlineButton{
		{{Text: "🔐 Crear Wallet", Data: "wallet"}, {Text: "👛 Mi Saldo", Data: "saldo"}},
		{{Text: "💶 Comprar SPC", Data: "comprar"}, {Text: "💰 Vender SPC", Data: "vender"}},
		{{Text: "📊 Precio", Data: "precio"}, {Text: "📖 Cómo comprar", Data: "como"}},
		{{Text: "🌐 Web", URL: "https://spaincoin.es"}, {Text: "👥 Comunidad", URL: "https://t.me/+m5O2f0sZaZBkYmJk"}},
	}
}

func jsonStr(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func deleteMessage(client *http.Client, chatID int64, messageID int64) {
	url := fmt.Sprintf("%s%s/deleteMessage", telegramAPI, botToken)
	body := fmt.Sprintf(`{"chat_id":%d,"message_id":%d}`, chatID, messageID)
	http.Post(url, "application/json", strings.NewReader(body))
}

func isAdmin(chatID int64) bool {
	_, ok := adminIDs[chatID]
	return ok
}

func isSuper(chatID int64) bool {
	role, ok := adminIDs[chatID]
	return ok && role == "super"
}

func handleCallback(client *http.Client, cb *TelegramCallback) {
	chatID := cb.From.ID
	userName := cb.From.FirstName
	if cb.From.Username != "" {
		userName = "@" + cb.From.Username
	}

	answerCallback(client, cb.ID)

	switch cb.Data {
	case "wallet":
		handleMiWallet(client, chatID)
	case "saldo":
		handleSaldo(client, chatID)
	case "comprar":
		sendMessageWithButtons(client, chatID,
			fmt.Sprintf("💶 <b>¿Cuánto quieres comprar?</b>\n\nPrecio actual: <b>%.4f€</b>/SPC", getCurrentPrice()),
			[][]InlineButton{
				{{Text: "10€", Data: "buy_10"}, {Text: "25€", Data: "buy_25"}, {Text: "50€", Data: "buy_50"}},
				{{Text: "100€", Data: "buy_100"}, {Text: "250€", Data: "buy_250"}, {Text: "500€", Data: "buy_500"}},
				{{Text: "Otra cantidad → escribe /comprar 75", Data: "noop"}},
				{{Text: "← Menú", Data: "menu"}},
			},
		)
	case "vender":
		sendMessageWithButtons(client, chatID,
			fmt.Sprintf("💰 <b>¿Cuántos SPC quieres vender?</b>\n\nPrecio actual: <b>%.4f€</b>/SPC", getCurrentPrice()),
			[][]InlineButton{
				{{Text: "50 SPC", Data: "sell_50"}, {Text: "100 SPC", Data: "sell_100"}, {Text: "500 SPC", Data: "sell_500"}},
				{{Text: "Otra cantidad → escribe /vender 200", Data: "noop"}},
				{{Text: "← Menú", Data: "menu"}},
			},
		)
	case "precio":
		handlePrecio(client, chatID)
	case "como":
		handleComoComprar(client, chatID)
	case "menu":
		handleStart(client, chatID, userName)
	case "noop":
		// Do nothing

	default:
		// Handle buy/sell amount buttons
		if strings.HasPrefix(cb.Data, "buy_") {
			amountStr := strings.TrimPrefix(cb.Data, "buy_")
			handleComprar(client, chatID, userName, chatID, "/comprar "+amountStr)
		} else if strings.HasPrefix(cb.Data, "sell_") {
			amountStr := strings.TrimPrefix(cb.Data, "sell_")
			handleVender(client, chatID, userName, chatID, "/vender "+amountStr)
		}
	}
}

func handleSaldo(client *http.Client, chatID int64) {
	walletAddr, _ := orderDB.GetAdminValue(fmt.Sprintf("wallet_%d", chatID))
	if walletAddr == "" {
		sendMessageWithButtons(client, chatID,
			"No tienes wallet registrada. Crea una primero:",
			[][]InlineButton{
				{{Text: "🔐 Crear Wallet", Data: "wallet"}},
				{{Text: "← Menú", Data: "menu"}},
			},
		)
		return
	}

	msg := fmt.Sprintf(`👛 <b>Tu wallet:</b>
<code>%s</code>

Para ver tu saldo en tiempo real, abre:
spaincoin.es/#/wallet`, walletAddr)
	sendMessageWithButtons(client, chatID, msg, [][]InlineButton{
		{{Text: "🌐 Ver en la web", URL: "https://spaincoin.es/#/wallet"}},
		{{Text: "💶 Comprar SPC", Data: "comprar"}, {Text: "💰 Vender SPC", Data: "vender"}},
		{{Text: "← Menú", Data: "menu"}},
	})
}

func notifyAdmins(client *http.Client, text string, excludeChat int64) {
	for id := range adminIDs {
		if id != excludeChat {
			sendMessage(client, id, text)
		}
	}
}

// getCurrentPrice returns the price based on auto-tiers or manual override.
func getCurrentPrice() float64 {
	// Check for manual override
	manualPrice, _ := orderDB.GetPrice()
	if manualPrice > 0 {
		autoMode, _ := orderDB.GetAdminValue("price_mode")
		if autoMode == "manual" {
			writePriceFile(manualPrice)
			return manualPrice
		}
	}

	// Auto-pricing based on total sold
	totalSPC, _, _, _ := orderDB.GetStats()
	price := 0.05 // default
	for i := len(priceTiers) - 1; i >= 0; i-- {
		if totalSPC >= priceTiers[i].SoldUpTo && i < len(priceTiers)-1 {
			price = priceTiers[i+1].Price
			break
		}
	}
	if price == 0.05 && len(priceTiers) > 0 {
		price = priceTiers[0].Price
	}

	writePriceFile(price)
	return price
}

// writePriceFile writes the current price to a shared file so the web API can read it.
func writePriceFile(price float64) {
	dataDir := os.Getenv("SPC_DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	data, _ := json.Marshal(price)
	os.WriteFile(dataDir+"/spc_price.json", data, 0644)
}

// ==========================================
// Message handler
// ==========================================

func handleMessage(client *http.Client, msg *TelegramMessage) {
	text := strings.TrimSpace(msg.Text)
	chatID := msg.Chat.ID
	userName := msg.From.FirstName
	if msg.From.Username != "" {
		userName = "@" + msg.From.Username
	}

	isGroup := msg.Chat.Type == "group" || msg.Chat.Type == "supergroup"
	userChatID := msg.From.ID // private chat with user

	log.Printf("[MSG] from=%s chat=%d group=%v text=%s", userName, chatID, isGroup, text)

	// Strip @botname from commands in groups
	if isGroup && strings.Contains(text, "@") {
		text = strings.Split(text, "@")[0]
	}

	// In groups: bot ignores EVERYTHING. Only sends daily report automatically.
	if isGroup {
		return
	}

	switch {
	case text == "/start" || text == "/menu":
		if isGroup {
			sendMessageWithButtons(client, chatID,
				"🇪🇸 <b>SpainCoin</b> — escríbeme por privado para operar:",
				[][]InlineButton{{{Text: "Abrir bot →", URL: "https://t.me/spaincoin_bot?start=go"}}},
			)
		} else {
			handleStart(client, chatID, userName)
		}
	case strings.HasPrefix(text, "/start buy_SPC"):
		// Auto-register wallet from onboarding link
		addr := strings.TrimPrefix(text, "/start buy_")
		if len(addr) == 43 && strings.HasPrefix(addr, "SPC") {
			key := fmt.Sprintf("wallet_%d", chatID)
			orderDB.SetAdminValue(key, addr)
			registerWalletAPI(client, addr)
			log.Printf("[AUTO-REG] %s → %s", userName, addr)
			for id := range adminIDs {
				sendMessage(client, id, fmt.Sprintf("📝 Auto-registro desde web: %s → <code>%s</code>", userName, addr))
			}
		}
		handleStart(client, chatID, userName)
	case strings.HasPrefix(text, "/start"):
		handleStart(client, chatID, userName)
	case text == "/precio" || text == "/price":
		handlePrecio(client, chatID) // precio es público, se puede ver en grupo
	case strings.HasPrefix(text, "/comprar"):
		if isGroup {
			sendMessage(client, chatID, fmt.Sprintf("👤 %s — te he enviado un mensaje privado para completar la compra.", userName))
			sendMessage(client, userChatID, "Has pedido comprar SPC desde el grupo. Escríbeme aquí:")
		}
		handleComprar(client, userChatID, userName, userChatID, text)
	case strings.HasPrefix(text, "/vender"):
		if isGroup {
			sendMessage(client, chatID, fmt.Sprintf("👤 %s — te he enviado un mensaje privado para completar la venta.", userName))
			sendMessage(client, userChatID, "Has pedido vender SPC desde el grupo. Escríbeme aquí:")
		}
		handleVender(client, userChatID, userName, userChatID, text)
	case strings.HasPrefix(text, "/registro"):
		if isGroup {
			sendMessage(client, chatID, "🔐 El registro se hace por privado → @spaincoin_bot")
			return
		}
		handleRegistro(client, chatID, userName, text)
	case text == "/miwallet":
		if isGroup {
			sendMessage(client, chatID, "📱 Escríbeme por privado → @spaincoin_bot")
			return
		}
		handleMiWallet(client, chatID)
	case text == "/comocomprar":
		handleComoComprar(client, chatID)
	case text == "/ayuda" || text == "/help":
		handleAyuda(client, chatID)

	// Admin commands
	case isAdmin(chatID) && text == "/ventas":
		handleVentas(client, chatID)
	case isAdmin(chatID) && strings.HasPrefix(text, "/confirmar"):
		handleConfirmar(client, chatID, text, userName)
	case isAdmin(chatID) && strings.HasPrefix(text, "/cancelar"):
		handleCancelar(client, chatID, text)
	case isAdmin(chatID) && text == "/stats":
		handleStats(client, chatID)

	// Super admin only (no price control via Telegram — only auto-tiers or SSH)
	case isSuper(chatID) && strings.HasPrefix(text, "/addadmin"):
		handleAddAdmin(client, chatID, text)
	case isSuper(chatID) && text == "/admins":
		handleListAdmins(client, chatID)
	case isSuper(chatID) && strings.HasPrefix(text, "/setvendidos"):
		handleSetVendidos(client, chatID, text)
	case isSuper(chatID) && text == "/myid":
		sendMessage(client, chatID, fmt.Sprintf("Tu chat ID: <code>%d</code>", chatID))
	case isSuper(chatID) && text == "/reporte":
		groupStr, _ := orderDB.GetAdminValue("group_chat_id")
		if groupStr == "" {
			groupStr = os.Getenv("SPC_GROUP_CHAT_ID")
		}
		gid, _ := strconv.ParseInt(groupStr, 10, 64)
		if gid != 0 {
			sendDailyReport(client, gid)
			sendMessage(client, chatID, "✅ Reporte enviado al grupo.")
		} else {
			sendDailyReport(client, chatID)
			sendMessage(client, chatID, "⚠️ No hay grupo configurado. Reporte enviado aquí.")
		}
	case isSuper(chatID) && text == "/tiers":
		handleShowTiers(client, chatID)

	default:
		if strings.HasPrefix(text, "/") {
			sendMessage(client, chatID, "Comando no reconocido. Escribe /ayuda")
		}
	}
}

// ==========================================
// Public commands
// ==========================================

func handleStart(client *http.Client, chatID int64, name string) {
	price := getCurrentPrice()
	totalSPC, _, _, _ := orderDB.GetStats()
	nextTier := ""
	for _, t := range priceTiers {
		if totalSPC < t.SoldUpTo {
			remaining := t.SoldUpTo - totalSPC
			nextTier = fmt.Sprintf("\n⏰ Faltan %.0f SPC para subir a %.2f€", remaining, getNextPrice())
			break
		}
	}

	msg := fmt.Sprintf(`¡Hola %s! 🇪🇸

Bienvenido a <b>SpainCoin ($SPC)</b>
La primera blockchain española

💰 Precio actual: <b>%.4f€</b>%s`, name, price, nextTier)

	sendMessageWithButtons(client, chatID, msg, mainMenuButtons())
}

func handlePrecio(client *http.Client, chatID int64) {
	price := getCurrentPrice()
	totalSPC, _, count, _ := orderDB.GetStats()
	walletCount := getWalletCount(client)

	nextPrice := getNextPrice()
	nextTierInfo := ""
	for _, t := range priceTiers {
		if totalSPC < t.SoldUpTo {
			remaining := t.SoldUpTo - totalSPC
			nextTierInfo = fmt.Sprintf("\n\n⏰ <b>Siguiente subida:</b>\nCuando se vendan %.0f SPC más → precio sube a <b>%.2f€</b>", remaining, nextPrice)
			break
		}
	}

	walletLine := ""
	if walletCount > 0 {
		walletLine = fmt.Sprintf("\n🇪🇸 Comunidad: %d wallets", walletCount)
	}

	msg := fmt.Sprintf(`💰 <b>SpainCoin ($SPC)</b>

Precio: <b>%.4f€</b>
Vendidos: %.2f SPC
Operaciones: %d%s%s

🌐 spaincoin.es`, price, totalSPC, count, walletLine, nextTierInfo)
	sendMessageWithButtons(client, chatID, msg, [][]InlineButton{
		{{Text: "💶 Comprar", Data: "comprar"}, {Text: "💰 Vender", Data: "vender"}},
		{{Text: "← Menú", Data: "menu"}},
	})
}

func getNextPrice() float64 {
	totalSPC, _, _, _ := orderDB.GetStats()
	for _, t := range priceTiers {
		if totalSPC < t.SoldUpTo {
			// Find the next tier
			for _, t2 := range priceTiers {
				if t2.SoldUpTo > t.SoldUpTo {
					return t2.Price
				}
			}
			return t.Price
		}
	}
	return priceTiers[len(priceTiers)-1].Price
}

func handleComprar(client *http.Client, chatID int64, userName string, userID int64, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, `💶 <b>¿Cuánto quieres comprar?</b>

Escribe /comprar seguido de la cantidad en euros.

Ejemplos:
<code>/comprar 10</code> — Comprar 10€ de SPC
<code>/comprar 50</code> — Comprar 50€ de SPC
<code>/comprar 100</code> — Comprar 100€ de SPC`)
		return
	}

	amountEUR, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || amountEUR < 1 || amountEUR > 50000 {
		sendMessage(client, chatID, "Cantidad inválida. Mínimo 1€, máximo 50.000€.\nEjemplo: /comprar 50")
		return
	}

	price := getCurrentPrice()
	amountSPC := math.Round((amountEUR/price)*10000) / 10000

	// Check for wallet address — from command or registered
	walletAddr := ""
	if len(parts) >= 3 && strings.HasPrefix(parts[2], "SPC") {
		walletAddr = parts[2]
	}
	if walletAddr == "" {
		// Try registered wallet
		saved, _ := orderDB.GetAdminValue(fmt.Sprintf("wallet_%d", chatID))
		if saved != "" {
			walletAddr = saved
		}
	}

	if walletAddr == "" {
		sendMessage(client, chatID, fmt.Sprintf(`💰 <b>Compra de %.2f€ de SPC</b>

Al precio actual recibirías: <b>%.4f SPC</b>

Primero registra tu wallet:
<code>/registro SPCtu_direccion</code>

Después simplemente:
<code>/comprar %.0f</code>

¿No tienes wallet? → /miwallet`, amountEUR, amountSPC, amountEUR))
		return
	}

	if len(walletAddr) != 43 || !strings.HasPrefix(walletAddr, "SPC") {
		sendMessage(client, chatID, "❌ Dirección SPC inválida. Debe empezar por SPC y tener 43 caracteres.\nEjemplo: SPCa1b2c3d4e5f6...")
		return
	}

	order := &database.Order{
		Type:       "buy",
		UserName:   userName,
		UserChatID: chatID,
		WalletAddr: walletAddr,
		AmountEUR:  amountEUR,
		AmountSPC:  amountSPC,
		PriceEUR:   price,
		Status:     database.OrderPending,
	}
	orderDB.CreateOrder(order)

	msg := fmt.Sprintf(`✅ <b>Orden #%d creada</b>

💶 Pagas: <b>%.2f€</b>
💰 Recibes: <b>%.4f SPC</b>
📍 Precio: %.4f€/SPC

<b>👉 Siguiente paso:</b>
Haz una transferencia de <b>%.2f€</b> al IBAN de abajo con el concepto indicado.

⏱️ Normalmente en menos de 10 minutos recibirás tus SPC.`, order.ID, amountEUR, amountSPC, price, amountEUR)
	sendMessage(client, chatID, msg)

	// Send IBAN in separate message so it's easy to copy
	ibanMsg := fmt.Sprintf(`🏦 <b>Datos de pago:</b>

IBAN (toca para copiar):
<code>%s</code>

Concepto (toca para copiar):
<code>SPC-%d</code>`, bankInfo, order.ID)
	sendMessage(client, chatID, ibanMsg)

	// Notify all admins
	adminMsg := fmt.Sprintf(`🔔 <b>NUEVA COMPRA #%d</b>

👤 %s
💶 %.2f€ → %.4f SPC
📍 Precio: %.4f€
📬 <code>%s</code>
💳 Concepto: SPC-%d

/confirmar %d
/cancelar %d`, order.ID, userName, amountEUR, amountSPC, price, walletAddr, order.ID, order.ID, order.ID)
	for id := range adminIDs {
		sendMessage(client, id, adminMsg)
	}
}

func handleVender(client *http.Client, chatID int64, userName string, userID int64, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, `💰 <b>¿Cuántos SPC quieres vender?</b>

Escribe /vender seguido de la cantidad de SPC.

Ejemplos:
<code>/vender 100</code> — Vender 100 SPC
<code>/vender 500</code> — Vender 500 SPC`)
		return
	}

	amountSPC, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || amountSPC < 0.01 {
		sendMessage(client, chatID, "Cantidad inválida. Ejemplo: /vender 100")
		return
	}

	price := getCurrentPrice()
	amountEUR := math.Round(amountSPC*price*100) / 100

	adminAddr, _ := orderDB.GetAdminValue("hot_wallet_address")
	if adminAddr == "" {
		adminAddr = "(pendiente — un admin lo configurará)"
	}

	order := &database.Order{
		Type:       "sell",
		UserName:   userName,
		UserChatID: chatID,
		AmountEUR:  amountEUR,
		AmountSPC:  amountSPC,
		PriceEUR:   price,
		Status:     database.OrderPending,
	}
	orderDB.CreateOrder(order)

	msg := fmt.Sprintf(`✅ <b>Orden de venta #%d creada</b>

💰 Vendes: <b>%.4f SPC</b>
💶 Recibes: <b>%.2f€</b>
📍 Precio: %.4f€/SPC

<b>👉 Siguiente paso:</b>
Envía <b>%.4f SPC</b> a esta dirección:
<code>%s</code>

Cuando verifiquemos la recepción, te haremos transferencia de %.2f€.`, order.ID, amountSPC, amountEUR, price, amountSPC, adminAddr, amountEUR)
	sendMessage(client, chatID, msg)

	adminMsg := fmt.Sprintf(`🔔 <b>NUEVA VENTA #%d</b>

👤 %s
💰 %.4f SPC → %.2f€

/confirmar %d
/cancelar %d`, order.ID, userName, amountSPC, amountEUR, order.ID, order.ID)
	for id := range adminIDs {
		sendMessage(client, id, adminMsg)
	}
}

func handleRegistro(client *http.Client, chatID int64, userName, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 || !strings.HasPrefix(parts[1], "SPC") || len(parts[1]) != 43 {
		sendMessageWithButtons(client, chatID, `📝 <b>Registra tu wallet</b>

Escribe /registro seguido de tu dirección SPC:

<code>/registro SPCtu_direccion_aqui</code>

Así no tendrás que pegarla cada vez que compres.`, [][]InlineButton{
			{{Text: "🔐 No tengo wallet → Crear una", URL: "https://spaincoin.es/#/wallet"}},
			{{Text: "← Menú", Data: "menu"}},
		})
		return
	}

	addr := parts[1]
	key := fmt.Sprintf("wallet_%d", chatID)
	orderDB.SetAdminValue(key, addr)

	// Register in exchange API for wallet counter
	registerWalletAPI(client, addr)

	count := getWalletCount(client)
	countMsg := ""
	if count > 0 {
		countMsg = fmt.Sprintf("\n\n🇪🇸 Ya somos <b>%d wallets</b> en SpainCoin", count)
	}

	sendMessageWithButtons(client, chatID, fmt.Sprintf(`✅ <b>Wallet registrada</b>

📬 <code>%s</code>%s`, addr, countMsg), [][]InlineButton{
		{{Text: "💶 Comprar SPC", Data: "comprar"}},
		{{Text: "← Menú", Data: "menu"}},
	})

	// Notify admins
	for id := range adminIDs {
		sendMessage(client, id, fmt.Sprintf("📝 Nuevo registro: %s → <code>%s</code>", userName, addr))
	}
}

// registerWalletAPI calls the exchange API to register a wallet for the counter.
func registerWalletAPI(client *http.Client, addr string) {
	apiURL := os.Getenv("SPC_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}
	body := fmt.Sprintf(`{"address":"%s"}`, addr)
	http.Post(apiURL+"/api/wallets/register", "application/json", strings.NewReader(body))
}

// getWalletCount fetches the current wallet count from the exchange API.
func getWalletCount(client *http.Client) int {
	apiURL := os.Getenv("SPC_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}
	resp, err := client.Get(apiURL + "/api/wallets/count")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	var result struct {
		Total int `json:"total"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Total
}

func handleComoComprar(client *http.Client, chatID int64) {
	msg := fmt.Sprintf(`📖 <b>Cómo comprar SpainCoin en 3 pasos</b>

<b>Paso 1 — Crea tu wallet</b>
Entra en spaincoin.es desde tu móvil.
Pulsa "Wallet" → "Crear Wallet".
Te dará una dirección SPCxxx... y una clave privada.
⚠️ Guarda la clave privada en papel.

<b>Paso 2 — Haz el pedido</b>
Escribe aquí:
<code>/comprar 50 SPCtu_direccion</code>
(cambia 50 por la cantidad en € y pega tu dirección)

<b>Paso 3 — Paga</b>
Haz una transferencia a:
%s
Con el concepto que te indique el bot.

¡Listo! Recibirás tus SPC en minutos. 🎉

💰 Precio actual: <b>%.4f€</b> por SPC`, bankInfo, getCurrentPrice())
	sendMessage(client, chatID, msg)
}

func handleMiWallet(client *http.Client, chatID int64) {
	msg := `🔐 <b>Crear tu wallet SpainCoin</b>

<b>Desde el móvil (recomendado):</b>
1. Abre spaincoin.es en tu navegador
2. Toca "Wallet"
3. Toca "Crear Wallet"
4. ¡Listo! Tu dirección SPCxxx... aparece en pantalla

⚠️ <b>MUY IMPORTANTE:</b>
• Guarda tu clave privada EN PAPEL
• Si la pierdes, pierdes tus fondos PARA SIEMPRE
• NUNCA se la des a nadie

🌐 spaincoin.es/#/wallet`
	sendMessage(client, chatID, msg)
}

func handleAyuda(client *http.Client, chatID int64) {
	msg := `📖 <b>Comandos SpainCoin</b>

💰 <b>Operar:</b>
/registro SPCxxx — Registrar tu wallet (una sola vez)
/comprar 50 — Comprar 50€ de SPC
/vender 100 — Vender 100 SPC
/precio — Ver precio actual

📱 <b>Wallet:</b>
/miwallet — Cómo crear tu wallet
/comocomprar — Guía paso a paso

ℹ️ <b>Info:</b>
/ayuda — Este mensaje

🌐 spaincoin.es`

	if isAdmin(chatID) {
		msg += `

🔑 <b>Admin:</b>
/ventas — Órdenes pendientes
/confirmar 23 — Confirmar orden
/cancelar 23 — Cancelar orden
/stats — Estadísticas`
	}

	if isSuper(chatID) {
		msg += `

👑 <b>Super Admin:</b>
/addadmin 123456 — Añadir admin
/admins — Ver admins
/tiers — Ver escalones de precio
/myid — Ver tu chat ID
(Precio solo modificable via SSH)`
	}

	sendMessage(client, chatID, msg)
}

// ==========================================
// Admin commands
// ==========================================

func handleVentas(client *http.Client, chatID int64) {
	orders, _ := orderDB.GetPendingOrders()
	if len(orders) == 0 {
		sendMessage(client, chatID, "✅ No hay órdenes pendientes.")
		return
	}

	msg := fmt.Sprintf("📋 <b>%d órdenes pendientes:</b>\n", len(orders))
	for _, o := range orders {
		emoji := "🟢 COMPRA"
		if o.Type == "sell" {
			emoji = "🔴 VENTA"
		}
		msg += fmt.Sprintf("\n%s #%d\n👤 %s\n💶 %.2f€ / %.4f SPC\n→ /confirmar %d\n", emoji, o.ID, o.UserName, o.AmountEUR, o.AmountSPC, o.ID)
	}
	sendMessage(client, chatID, msg)
}

func handleConfirmar(client *http.Client, chatID int64, text string, adminName string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, "Uso: /confirmar 23")
		return
	}
	id, _ := strconv.ParseInt(parts[1], 10, 64)
	order, err := orderDB.GetOrder(id)
	if err != nil {
		sendMessage(client, chatID, fmt.Sprintf("Orden #%d no encontrada.", id))
		return
	}
	if order.Status != database.OrderPending {
		sendMessage(client, chatID, fmt.Sprintf("Orden #%d ya está %s.", id, order.Status))
		return
	}

	// For buy orders: send SPC automatically
	txID := ""
	if order.Type == "buy" && order.WalletAddr != "" && hotWalletKey != nil {
		var err error
		txID, err = sendSPCToWallet(client, order.WalletAddr, order.AmountSPC)
		if err != nil {
			sendMessage(client, chatID, fmt.Sprintf("⚠️ Orden confirmada pero error enviando SPC: %v\nEnvía manualmente %.4f SPC a <code>%s</code>", err, order.AmountSPC, order.WalletAddr))
			orderDB.ConfirmOrder(id)
			return
		}
	}

	orderDB.ConfirmOrder(id)

	// Update price immediately after sale (writes spc_price.json)
	getCurrentPrice()

	// Notify confirming admin
	txInfo := ""
	if txID != "" {
		txInfo = fmt.Sprintf("\n📦 TX: <code>%s</code>", txID[:16]+"...")
	}
	sendMessage(client, chatID, fmt.Sprintf("✅ Orden #%d confirmada. SPC enviados automáticamente.%s", id, txInfo))

	// Notify other admins
	notifyAdmins(client, fmt.Sprintf("✅ Orden #%d confirmada por %s. SPC enviados.%s", id, adminName, txInfo), chatID)

	// Notify user
	if order.UserChatID > 0 {
		var userMsg string
		if order.Type == "buy" {
			txLine := ""
			if txID != "" {
				txLine = fmt.Sprintf("\n📦 TX: <code>%s</code>", txID)
			}
			userMsg = fmt.Sprintf(`🎉 <b>¡Compra completada!</b>

Has recibido <b>%.4f SPC</b> en tu wallet:
<code>%s</code>
%s
Precio: %.4f€/SPC
Total pagado: %.2f€

¡Bienvenido a SpainCoin! 🇪🇸
Comparte con tus amigos → cuantos más seamos, más vale $SPC.`, order.AmountSPC, order.WalletAddr, txLine, order.PriceEUR, order.AmountEUR)
		} else {
			userMsg = fmt.Sprintf(`🎉 <b>¡Venta completada!</b>

Te hemos enviado <b>%.2f€</b> por transferencia.

Gracias por confiar en SpainCoin 🇪🇸`, order.AmountEUR)
		}
		sendMessage(client, order.UserChatID, userMsg)
	}

	// Check if price should auto-update
	newPrice := getCurrentPrice()
	oldPrice := order.PriceEUR
	if newPrice > oldPrice {
		for id := range adminIDs {
			sendMessage(client, id, fmt.Sprintf("📈 <b>Precio automático actualizado:</b> %.4f€ → %.4f€", oldPrice, newPrice))
		}
	}
}

func handleCancelar(client *http.Client, chatID int64, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, "Uso: /cancelar 23")
		return
	}
	id, _ := strconv.ParseInt(parts[1], 10, 64)
	order, err := orderDB.GetOrder(id)
	if err != nil {
		sendMessage(client, chatID, fmt.Sprintf("Orden #%d no encontrada.", id))
		return
	}
	orderDB.CancelOrder(id)
	sendMessage(client, chatID, fmt.Sprintf("❌ Orden #%d cancelada.", id))

	if order.UserChatID > 0 {
		sendMessage(client, order.UserChatID, fmt.Sprintf("❌ Tu orden #%d ha sido cancelada. Si crees que es un error, escribe /ayuda.", id))
	}
}

func handleStats(client *http.Client, chatID int64) {
	totalSPC, totalEUR, count, _ := orderDB.GetStats()
	price := getCurrentPrice()

	// Find current tier
	currentTier := ""
	for i, t := range priceTiers {
		if totalSPC < t.SoldUpTo {
			currentTier = fmt.Sprintf("Tier %d: hasta %.0f SPC a %.2f€", i+1, t.SoldUpTo, t.Price)
			break
		}
	}

	msg := fmt.Sprintf(`📊 <b>Estadísticas SpainCoin</b>

💰 Precio actual: %.4f€
📈 Total SPC vendidos: %.2f
💶 Total EUR recaudado: %.2f€
🔄 Operaciones completadas: %d
📦 Valor SPC vendido: %.2f€

🏷️ %s
🏦 Banco: %s`, price, totalSPC, totalEUR, count, totalSPC*price, currentTier, bankInfo)
	sendMessage(client, chatID, msg)
}

// ==========================================
// Super admin commands
// ==========================================

func handleSetPrecio(client *http.Client, chatID int64, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, "Uso: /setprecio 0.15")
		return
	}
	price, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || price <= 0 {
		sendMessage(client, chatID, "Precio inválido.")
		return
	}
	oldPrice := getCurrentPrice()
	orderDB.SetPrice(price)
	orderDB.SetAdminValue("price_mode", "manual")

	sendMessage(client, chatID, fmt.Sprintf("💰 Precio manual: %.4f€ → <b>%.4f€</b>\n⚠️ Modo automático desactivado. Usa /autoprecio para reactivar.", oldPrice, price))
	notifyAdmins(client, fmt.Sprintf("💰 Precio cambiado manualmente a %.4f€", price), chatID)
}

func handleAutoPrecio(client *http.Client, chatID int64) {
	orderDB.SetAdminValue("price_mode", "auto")
	price := getCurrentPrice()
	sendMessage(client, chatID, fmt.Sprintf("📈 Precio automático activado. Precio actual: <b>%.4f€</b>\nEl precio sube automáticamente según los SPC vendidos.", price))
}

func handleSetBank(client *http.Client, chatID int64, text string) {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 || len(parts[1]) < 5 {
		sendMessage(client, chatID, "Uso: /setbank ES12 3456 7890 1234 5678 9012")
		return
	}
	bankInfo = strings.TrimSpace(parts[1])
	orderDB.SetAdminValue("bank_info", bankInfo)
	sendMessage(client, chatID, fmt.Sprintf("🏦 Banco actualizado:\n<code>%s</code>", bankInfo))
}

func handleAddAdmin(client *http.Client, chatID int64, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, "Uso: /addadmin 123456789\nEl usuario debe escribir /myid al bot para obtener su ID.")
		return
	}
	newID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		sendMessage(client, chatID, "ID inválido.")
		return
	}
	adminIDs[newID] = "admin"
	sendMessage(client, chatID, fmt.Sprintf("✅ Admin añadido: %d\n⚠️ Para hacerlo permanente, añade el ID a SPC_ADMIN_IDS en el .env", newID))
}

func handleListAdmins(client *http.Client, chatID int64) {
	msg := "👥 <b>Admins actuales:</b>\n"
	for id, role := range adminIDs {
		emoji := "🔑"
		if role == "super" {
			emoji = "👑"
		}
		msg += fmt.Sprintf("%s %d — %s\n", emoji, id, role)
	}
	sendMessage(client, chatID, msg)
}

func handleShowTiers(client *http.Client, chatID int64) {
	totalSPC, _, _, _ := orderDB.GetStats()
	msg := "📊 <b>Escalones de precio:</b>\n\n"
	for i, t := range priceTiers {
		marker := ""
		if totalSPC < t.SoldUpTo && (i == 0 || totalSPC >= priceTiers[i-1].SoldUpTo) {
			marker = " ← AQUÍ"
		}
		msg += fmt.Sprintf("%.0f SPC → %.2f€%s\n", t.SoldUpTo, t.Price, marker)
	}
	msg += fmt.Sprintf("\nVendidos: %.2f SPC", totalSPC)
	sendMessage(client, chatID, msg)
}

func handleSetVendidos(client *http.Client, chatID int64, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, "Uso: /setvendidos 600")
		return
	}
	amount, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || amount < 0 {
		sendMessage(client, chatID, "Cantidad inválida.")
		return
	}

	// Create a synthetic confirmed order to adjust the total
	currentSPC, _, _, _ := orderDB.GetStats()
	diff := amount - currentSPC
	if diff <= 0 {
		sendMessage(client, chatID, fmt.Sprintf("Ya hay %.2f SPC vendidos. No se puede reducir.", currentSPC))
		return
	}

	price := getCurrentPrice()
	order := &database.Order{
		Type:      "buy",
		UserName:  "ajuste-manual",
		AmountEUR: diff * price,
		AmountSPC: diff,
		PriceEUR:  price,
		Status:    database.OrderConfirmed,
	}
	orderDB.CreateOrder(order)
	orderDB.ConfirmOrder(order.ID)

	newPrice := getCurrentPrice()
	sendMessage(client, chatID, fmt.Sprintf("✅ Vendidos ajustados: %.2f → %.2f SPC\n💰 Precio actual: %.4f€", currentSPC, amount, newPrice))
}

// sendSPCToWallet sends SPC from the hot wallet to a recipient address.
func sendSPCToWallet(client *http.Client, toAddress string, amountSPC float64) (string, error) {
	if hotWalletKey == nil {
		return "", fmt.Errorf("hot wallet not configured")
	}

	// Get nonce
	resp, err := client.Get(fmt.Sprintf("%s/address/%s/balance", nodeRPCURL, hotWalletAddr.String()))
	if err != nil {
		return "", fmt.Errorf("nonce fetch failed: %v", err)
	}
	defer resp.Body.Close()
	var nonceResult struct {
		Nonce uint64 `json:"nonce"`
	}
	json.NewDecoder(resp.Body).Decode(&nonceResult)

	// Build transaction
	amountPesetas := uint64(amountSPC * 1_000_000_000_000)
	toAddr, err := crypto.AddressFromHex(toAddress)
	if err != nil {
		return "", fmt.Errorf("invalid address: %v", err)
	}

	tx := block.NewTransaction(hotWalletAddr, toAddr, amountPesetas, nonceResult.Nonce, 1000)
	if err := tx.Sign(hotWalletKey); err != nil {
		return "", fmt.Errorf("sign failed: %v", err)
	}

	sigR := tx.Signature.R.Text(16)
	sigS := tx.Signature.S.Text(16)
	if len(sigR)%2 != 0 {
		sigR = "0" + sigR
	}
	if len(sigS)%2 != 0 {
		sigS = "0" + sigS
	}

	body := map[string]interface{}{
		"from":   hotWalletAddr.String(),
		"to":     toAddress,
		"amount": amountPesetas,
		"nonce":  nonceResult.Nonce,
		"fee":    1000,
		"sig_r":  sigR,
		"sig_s":  sigS,
	}
	jsonBody, _ := json.Marshal(body)

	txResp, err := http.Post(nodeRPCURL+"/tx/send", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("send failed: %v", err)
	}
	defer txResp.Body.Close()

	var txResult struct {
		TxID  string `json:"tx_id"`
		Error string `json:"error"`
	}
	json.NewDecoder(txResp.Body).Decode(&txResult)

	if txResult.Error != "" {
		return "", fmt.Errorf(txResult.Error)
	}

	log.Printf("[HOT-WALLET] Sent %.4f SPC to %s — tx: %s", amountSPC, toAddress, txResult.TxID)
	return txResult.TxID, nil
}

// sendDailyReport sends the daily status to a chat, deleting the previous report.
func sendDailyReport(client *http.Client, chatID int64) {
	price := getCurrentPrice()
	totalSPC, totalEUR, count, _ := orderDB.GetStats()
	walletCount := getWalletCount(client)

	nextTierInfo := ""
	for _, t := range priceTiers {
		if totalSPC < t.SoldUpTo {
			remaining := t.SoldUpTo - totalSPC
			nextTierInfo = fmt.Sprintf("\n⏰ Faltan %.0f SPC para subir de precio", remaining)
			break
		}
	}

	msg := fmt.Sprintf(`🐂 <b>SpainCoin — Informe diario 🇪🇸</b>

💰 Precio: <b>%.4f€</b> por SPC
📈 SPC vendidos: %.2f
💶 Recaudado: %.2f€
🔄 Operaciones: %d
🇪🇸 Comunidad: <b>%d wallets</b>%s

¿Todavía no tienes SPC? Los primeros en comprar son los que más ganan.
Para comprar, vender o crear tu wallet habla conmigo por privado 👇`, price, totalSPC, totalEUR, count, walletCount, nextTierInfo)

	// Delete previous report in this chat
	prevKey := fmt.Sprintf("last_report_%d", chatID)
	prevMsgStr, _ := orderDB.GetAdminValue(prevKey)
	if prevMsgStr != "" {
		prevMsgID, _ := strconv.ParseInt(prevMsgStr, 10, 64)
		if prevMsgID > 0 {
			deleteMessage(client, chatID, prevMsgID)
		}
	}

	// Send new report and save message ID
	newMsgID := sendMessageAndTrack(client, chatID, msg, [][]InlineButton{
		{{Text: "🐂 Comprar SPC", URL: "https://t.me/spaincoin_bot?start=go"}},
		{{Text: "🌐 Web", URL: "https://spaincoin.es"}},
	})
	if newMsgID > 0 {
		orderDB.SetAdminValue(prevKey, fmt.Sprintf("%d", newMsgID))
	}

	log.Printf("[DAILY] Report sent to chat %d (msg %d)", chatID, newMsgID)
}
