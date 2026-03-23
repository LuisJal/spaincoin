package main

import (
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

	"github.com/spaincoin/spaincoin/exchange/database"
)

// Telegram Bot API base URL.
const telegramAPI = "https://api.telegram.org/bot"

var (
	botToken  string
	adminChat int64 // Your Telegram chat ID (admin)
	orderDB   *database.OrderDB
	bizumInfo string // Bizum phone number for payments
)

// TelegramUpdate represents an incoming update from Telegram.
type TelegramUpdate struct {
	UpdateID int64            `json:"update_id"`
	Message  *TelegramMessage `json:"message"`
}

// TelegramMessage represents a Telegram message.
type TelegramMessage struct {
	MessageID int64         `json:"message_id"`
	From      *TelegramUser `json:"from"`
	Chat      *TelegramChat `json:"chat"`
	Text      string        `json:"text"`
	Date      int64         `json:"date"`
}

// TelegramUser represents a Telegram user.
type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// TelegramChat represents a Telegram chat.
type TelegramChat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

func main() {
	botToken = os.Getenv("SPC_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("SPC_BOT_TOKEN required")
	}

	adminStr := os.Getenv("SPC_ADMIN_CHAT_ID")
	if adminStr != "" {
		adminChat, _ = strconv.ParseInt(adminStr, 10, 64)
	}

	bizumInfo = os.Getenv("SPC_BIZUM_PHONE")
	if bizumInfo == "" {
		bizumInfo = "(no configurado)"
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

	log.Println("SpainCoin Bot starting...")
	log.Printf("Admin chat ID: %d", adminChat)

	// Long polling
	offset := int64(0)
	client := &http.Client{Timeout: 35 * time.Second}

	for {
		updates, err := getUpdates(client, offset)
		if err != nil {
			log.Printf("getUpdates error: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		for _, u := range updates {
			if u.Message != nil {
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

func jsonStr(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func handleMessage(client *http.Client, msg *TelegramMessage) {
	text := strings.TrimSpace(msg.Text)
	chatID := msg.Chat.ID
	isAdmin := chatID == adminChat
	userName := msg.From.FirstName
	if msg.From.Username != "" {
		userName = "@" + msg.From.Username
	}

	log.Printf("[MSG] from=%s chat=%d text=%s", userName, chatID, text)

	switch {
	// ===== PUBLIC COMMANDS =====
	case text == "/start":
		handleStart(client, chatID, userName)

	case text == "/precio" || text == "/price":
		handlePrecio(client, chatID)

	case strings.HasPrefix(text, "/comprar"):
		handleComprar(client, chatID, userName, text)

	case strings.HasPrefix(text, "/vender"):
		handleVender(client, chatID, userName, text)

	case text == "/ayuda" || text == "/help":
		handleAyuda(client, chatID, isAdmin)

	case text == "/miwallet":
		handleMiWallet(client, chatID)

	// ===== ADMIN COMMANDS =====
	case isAdmin && strings.HasPrefix(text, "/setprecio"):
		handleSetPrecio(client, chatID, text)

	case isAdmin && text == "/ventas":
		handleVentas(client, chatID)

	case isAdmin && strings.HasPrefix(text, "/confirmar"):
		handleConfirmar(client, chatID, text)

	case isAdmin && strings.HasPrefix(text, "/cancelar"):
		handleCancelar(client, chatID, text)

	case isAdmin && text == "/stats":
		handleStats(client, chatID)

	case isAdmin && text == "/myadmin":
		sendMessage(client, chatID, fmt.Sprintf("Tu chat ID es: <code>%d</code>", chatID))

	default:
		if strings.HasPrefix(text, "/") {
			sendMessage(client, chatID, "Comando no reconocido. Escribe /ayuda para ver los comandos disponibles.")
		}
	}
}

func handleStart(client *http.Client, chatID int64, name string) {
	price, _ := orderDB.GetPrice()
	msg := fmt.Sprintf(`¡Hola %s! 👋

Bienvenido a <b>SpainCoin ($SPC)</b>
La blockchain de España 🇪🇸

💰 Precio actual: <b>%.4f€</b> por SPC

<b>Comandos:</b>
/precio — Ver precio actual
/comprar 50 — Comprar 50€ de SPC
/vender 100 — Vender 100 SPC
/miwallet — Cómo crear tu wallet
/ayuda — Todos los comandos

🌐 Web: spaincoin.es`, name, price)
	sendMessage(client, chatID, msg)
}

func handlePrecio(client *http.Client, chatID int64) {
	price, _ := orderDB.GetPrice()
	totalSPC, totalEUR, count, _ := orderDB.GetStats()
	msg := fmt.Sprintf(`💰 <b>Precio SPC:</b> %.4f€

📊 Estadísticas:
• Vendidos: %.2f SPC
• Recaudado: %.2f€
• Operaciones: %d

🌐 spaincoin.es`, price, totalSPC, totalEUR, count)
	sendMessage(client, chatID, msg)
}

func handleComprar(client *http.Client, chatID int64, userName, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, "Uso: /comprar 50\n(cantidad en EUR que quieres gastar)")
		return
	}
	amountEUR, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || amountEUR < 1 || amountEUR > 10000 {
		sendMessage(client, chatID, "Cantidad inválida. Ejemplo: /comprar 50")
		return
	}

	price, _ := orderDB.GetPrice()
	amountSPC := math.Round((amountEUR/price)*10000) / 10000

	// Check if they provided wallet address
	walletAddr := ""
	if len(parts) >= 3 && strings.HasPrefix(parts[2], "SPC") {
		walletAddr = parts[2]
	}

	if walletAddr == "" {
		sendMessage(client, chatID, fmt.Sprintf(`Para completar la compra necesito tu dirección SPC.

Envía:
<code>/comprar %.0f SPCtu_direccion_aqui</code>

¿No tienes wallet? Escribe /miwallet`, amountEUR))
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

	// Notify user
	msg := fmt.Sprintf(`✅ <b>Orden de compra #%d creada</b>

💶 Pagas: <b>%.2f€</b>
💰 Recibes: <b>%.4f SPC</b>
📍 Precio: %.4f€/SPC
📬 Wallet: <code>%s</code>

<b>Siguiente paso:</b>
Envía %.2f€ por Bizum al %s
Concepto: "SPC-%d"

Cuando enviemos los SPC recibirás confirmación aquí.`, order.ID, amountEUR, amountSPC, price, walletAddr, amountEUR, bizumInfo, order.ID)
	sendMessage(client, chatID, msg)

	// Notify admin
	if adminChat > 0 {
		adminMsg := fmt.Sprintf(`🔔 <b>NUEVA COMPRA #%d</b>

👤 %s
💶 %.2f€ → %.4f SPC
📍 Precio: %.4f€
📬 <code>%s</code>

Para confirmar: /confirmar %d
Para cancelar: /cancelar %d`, order.ID, userName, amountEUR, amountSPC, price, walletAddr, order.ID, order.ID)
		sendMessage(client, adminChat, adminMsg)
	}
}

func handleVender(client *http.Client, chatID int64, userName, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, "Uso: /vender 100\n(cantidad de SPC que quieres vender)")
		return
	}
	amountSPC, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || amountSPC < 0.01 {
		sendMessage(client, chatID, "Cantidad inválida. Ejemplo: /vender 100")
		return
	}

	price, _ := orderDB.GetPrice()
	amountEUR := math.Round(amountSPC*price*100) / 100

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

	// Admin's SPC address for receiving
	adminAddr, _ := orderDB.GetAdminValue("admin_spc_address")
	if adminAddr == "" {
		adminAddr = "(pendiente de configurar)"
	}

	msg := fmt.Sprintf(`✅ <b>Orden de venta #%d creada</b>

💰 Vendes: <b>%.4f SPC</b>
💶 Recibes: <b>%.2f€</b>
📍 Precio: %.4f€/SPC

<b>Siguiente paso:</b>
Envía %.4f SPC a esta dirección:
<code>%s</code>

Cuando recibamos los SPC te haremos Bizum por %.2f€.`, order.ID, amountSPC, amountEUR, price, amountSPC, adminAddr, amountEUR)
	sendMessage(client, chatID, msg)

	if adminChat > 0 {
		adminMsg := fmt.Sprintf(`🔔 <b>NUEVA VENTA #%d</b>

👤 %s
💰 %.4f SPC → %.2f€

Para confirmar: /confirmar %d
Para cancelar: /cancelar %d`, order.ID, userName, amountSPC, amountEUR, order.ID, order.ID)
		sendMessage(client, adminChat, adminMsg)
	}
}

func handleMiWallet(client *http.Client, chatID int64) {
	msg := `🔐 <b>Cómo crear tu wallet SpainCoin</b>

<b>Opción 1 — Desde el móvil:</b>
Entra en spaincoin.es/wallet y pulsa "Crear Wallet"
(las claves se generan en TU dispositivo, nunca salen de él)

<b>Opción 2 — Descarga el CLI:</b>
Descarga desde spaincoin.es → Wallet → tu sistema operativo

Luego ejecuta:
<code>./spc wallet new</code>

⚠️ <b>IMPORTANTE:</b> Guarda tu clave privada en papel. Si la pierdes, pierdes tus fondos para siempre. NUNCA la compartas con nadie.

🌐 spaincoin.es/#/wallet`
	sendMessage(client, chatID, msg)
}

func handleAyuda(client *http.Client, chatID int64, isAdmin bool) {
	msg := `📖 <b>Comandos SpainCoin Bot</b>

/precio — Ver precio actual de SPC
/comprar 50 SPCxxx — Comprar 50€ de SPC
/vender 100 — Vender 100 SPC
/miwallet — Cómo crear tu wallet
/ayuda — Este mensaje`

	if isAdmin {
		msg += `

🔑 <b>Admin:</b>
/setprecio 0.15 — Cambiar precio
/ventas — Ver órdenes pendientes
/confirmar 23 — Confirmar orden
/cancelar 23 — Cancelar orden
/stats — Estadísticas totales`
	}

	sendMessage(client, chatID, msg)
}

// ===== ADMIN COMMANDS =====

func handleSetPrecio(client *http.Client, chatID int64, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, "Uso: /setprecio 0.15")
		return
	}
	price, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || price <= 0 || price > 1000000 {
		sendMessage(client, chatID, "Precio inválido. Ejemplo: /setprecio 0.15")
		return
	}
	oldPrice, _ := orderDB.GetPrice()
	orderDB.SetPrice(price)

	change := ((price - oldPrice) / oldPrice) * 100
	arrow := "📈"
	if price < oldPrice {
		arrow = "📉"
	}

	msg := fmt.Sprintf(`%s <b>Precio actualizado</b>

Anterior: %.4f€
Nuevo: <b>%.4f€</b>
Cambio: %.2f%%`, arrow, oldPrice, price, change)
	sendMessage(client, chatID, msg)
}

func handleVentas(client *http.Client, chatID int64) {
	orders, _ := orderDB.GetPendingOrders()
	if len(orders) == 0 {
		sendMessage(client, chatID, "No hay órdenes pendientes.")
		return
	}

	msg := fmt.Sprintf("📋 <b>%d órdenes pendientes:</b>\n\n", len(orders))
	for _, o := range orders {
		emoji := "🟢"
		if o.Type == "sell" {
			emoji = "🔴"
		}
		msg += fmt.Sprintf("%s #%d %s — %.2f€ / %.4f SPC — %s\n", emoji, o.ID, o.Type, o.AmountEUR, o.AmountSPC, o.UserName)
	}
	msg += "\nPara confirmar: /confirmar ID"
	sendMessage(client, chatID, msg)
}

func handleConfirmar(client *http.Client, chatID int64, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendMessage(client, chatID, "Uso: /confirmar 23")
		return
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		sendMessage(client, chatID, "ID inválido.")
		return
	}

	order, err := orderDB.GetOrder(id)
	if err != nil {
		sendMessage(client, chatID, fmt.Sprintf("Orden #%d no encontrada.", id))
		return
	}

	if order.Status != database.OrderPending {
		sendMessage(client, chatID, fmt.Sprintf("Orden #%d ya está %s.", id, order.Status))
		return
	}

	orderDB.ConfirmOrder(id)

	sendMessage(client, chatID, fmt.Sprintf("✅ Orden #%d confirmada.", id))

	// Notify user
	if order.UserChatID > 0 {
		var userMsg string
		if order.Type == "buy" {
			userMsg = fmt.Sprintf(`✅ <b>¡Compra completada!</b>

Has recibido <b>%.4f SPC</b> en tu wallet:
<code>%s</code>

Gracias por confiar en SpainCoin 🇪🇸`, order.AmountSPC, order.WalletAddr)
		} else {
			userMsg = fmt.Sprintf(`✅ <b>¡Venta completada!</b>

Te hemos enviado <b>%.2f€</b> por Bizum.

Gracias por confiar en SpainCoin 🇪🇸`, order.AmountEUR)
		}
		sendMessage(client, order.UserChatID, userMsg)
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
		sendMessage(client, order.UserChatID, fmt.Sprintf("❌ Tu orden #%d ha sido cancelada. Si tienes dudas, escribe /ayuda.", id))
	}
}

func handleStats(client *http.Client, chatID int64) {
	totalSPC, totalEUR, count, _ := orderDB.GetStats()
	price, _ := orderDB.GetPrice()
	msg := fmt.Sprintf(`📊 <b>Estadísticas SpainCoin</b>

💰 Precio actual: %.4f€
📈 Total vendido: %.2f SPC
💶 Total recaudado: %.2f€
🔄 Operaciones completadas: %d
📦 Valor SPC vendido hoy: %.2f€`, price, totalSPC, totalEUR, count, totalSPC*price)
	sendMessage(client, chatID, msg)
}
