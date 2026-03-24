# SpainCoin - Blockchain Project

## Proyecto
SpainCoin es una blockchain propia (Layer 1) con su criptomoneda nativa $SPC.
Protocolo **NON-CUSTODIAL**: no custodiamos fondos de usuarios, no somos un exchange.
Los usuarios generan sus propias wallets (client-side) y controlan sus claves privadas.
Trading P2P vía bot de Telegram (admin confirma pagos, bot auto-envía SPC).

> El código de exchange custodial está guardado en la rama `exchange-v1` para cuando se obtenga la licencia CASP.

## Stack Tecnológico
- **Blockchain core**: Go 1.26+ (estándar en blockchain — Ethereum, Cosmos usan Go)
- **Consenso**: Proof of Stake simplificado
- **Networking**: P2P con libp2p
- **Bot Telegram**: Go (~1300 líneas), inline buttons, auto-envío SPC
- **Web (spaincoin.es)**: React + Vite (informacional + wallet + explorer)
- **Exchange API**: Go (backend para la web, conecta con nodo blockchain)
- **CLI**: `spc` — gestión de wallets y transacciones

## Arquitectura del Proyecto

```
spaincoin/
├── .claude/               # Configuración Claude Code
│   ├── commands/          # Slash commands: /status, /new-module
│   └── settings.json      # Permisos + hooks (go fmt/vet automático)
├── bot/                   # ✅ Bot de Telegram (P2P trading, onboarding, admin)
│   └── main.go            # Bot completo (~1300 líneas)
├── core/                  # Blockchain core (Go)
│   ├── block/             # ✅ Transacciones, bloques, Merkle tree
│   ├── chain/             # ✅ Lógica de la cadena (AddBlock, IsValid)
│   ├── consensus/         # ✅ PoS: selección validadores, slashing
│   ├── crypto/            # ✅ ECDSA keys, SHA-256, Address SPC...
│   ├── mempool/           # ✅ Pool de transacciones pendientes
│   ├── network/           # ✅ P2P networking (libp2p, gossipsub, mDNS)
│   ├── state/             # ✅ Estado global (balances, cuentas)
│   ├── storage/           # ✅ Persistencia con BoltDB
│   └── wallet/            # ✅ Gestión de wallets
├── cli/                   # ✅ CLI `spc` (wallet, tx, chain)
├── deploy/                # Scripts de despliegue para VPS 1 y VPS 2
├── docs/                  # Documentación
├── exchange/              # ✅ API backend Go (handlers, auth, market, database)
├── frontend/              # ✅ Web React+Vite (informacional, wallet, explorer)
│   └── src/pages/         # Landing, Wallet, Explorer, Onboarding, WhitePaper...
├── node/                  # ✅ Nodo ejecutable (produce bloques cada 5s)
│   ├── cmd/main.go        # Entrypoint: go run ./node/cmd/
│   ├── node.go            # Lógica del nodo
│   └── rpc/               # HTTP API del nodo
├── releases/              # Binarios cross-platform del CLI (linux, mac, windows)
└── tests/                 # Tests de integración
```

## Estado Actual — Protocolo NON-CUSTODIAL en producción

| Módulo | Estado | Tests |
|--------|--------|-------|
| core/crypto | ✅ Completo | 5 |
| core/block | ✅ Completo | 12 |
| core/state | ✅ Completo | 8 |
| core/mempool | ✅ Completo | 6 |
| core/chain | ✅ Completo | 9 |
| core/consensus | ✅ Completo | 10 |
| core/network (P2P) | ✅ Completo | 5 |
| core/storage (BoltDB) | ✅ Completo | 6 |
| core/wallet | ✅ Completo | 4 |
| node/ | ✅ En producción | — |
| node/rpc (HTTP API) | ✅ Completo | 12 |
| cli/ | ✅ Completo | — |
| tests/ integración | ✅ Completo | 3 |
| bot/ (Telegram) | ✅ En producción | — |
| exchange/ (Go API) | ✅ En producción | — |
| frontend/ (React) | ✅ En producción | — |

**Total: 80 tests PASS**

## Infraestructura de producción

| Componente | Ubicación | Estado |
|-----------|-----------|--------|
| Nodo blockchain | VPS 1: 204.168.176.40 | ✅ Produciendo bloques 24/7 (cada 5s) |
| Web + API + Bot | VPS 2: 46.62.201.94 | ✅ Corriendo 24/7 |
| Dominio | spaincoin.es | ✅ HTTPS activo (Let's Encrypt) |
| Bot Telegram | @SpainCoinBot | ✅ Operativo (P2P trading) |
| Backups | Cron diario en ambos VPS | ✅ Configurados |
| Firewall | UFW en ambos VPS | ✅ RPC restringido |
| SSH | Solo claves, sin password | ✅ Configurado |

## Wallets principales

| Wallet | Dirección | Uso |
|--------|-----------|-----|
| Fundador | SPC5e2ac672147ea748ba1d0c27aed781995ea7349f | 5,000,000 SPC génesis |
| Hot wallet | SPCc119f94ab074c970dc129884163fc00106d65481 | 50,000 SPC para bot auto-envío |

## Bot de Telegram — Funcionalidades

- **Onboarding**: flujo interactivo de 5 pasos con logo del toro
- **Comprar SPC**: P2P con transferencia bancaria (IBAN Revolut)
- **Vender SPC**: P2P, admin confirma recepción
- **Auto-envío**: bot envía SPC automáticamente al confirmar admin
- **Precio auto-escalado**: tiers según SPC vendidos (500->0.05EUR, 1000->0.08EUR, etc.)
- **Reporte diario**: 9AM al grupo de Telegram
- **Precios crypto**: Binance API en tiempo real
- **Estructura admin**: super admin + 2 admins planificado
- **Seguridad**: precio NUNCA se cambia por Telegram, solo auto-tiers o SSH

## Web (spaincoin.es) — Funcionalidades

- **Landing**: informacional sobre el proyecto
- **Wallet**: generación client-side (self-custody, claves nunca salen del navegador)
- **Explorer**: bloques, transacciones, cuentas
- **Onboarding**: guía interactiva
- **White Paper**: documento del proyecto
- **Market Info**: precios en tiempo real
- **Validadores**: información sobre validadores PoS
- **Legal**: Términos, Privacidad (RGPD), Riesgos (MiCA), Cookies

> La web NO tiene trading. El trading es P2P vía Telegram.

## Fases del Proyecto

- [x] **Fase 1** - Core blockchain (bloques, transacciones, consenso básico)
- [x] **Fase 2** - Red P2P (múltiples nodos comunicándose con libp2p)
- [x] **Fase 3** - Wallet + CLI + persistencia
- [x] **Fase 4** - Testnet infra (VPS Hetzner, nodo en producción 24/7)
- [x] **Fase 5** - Web + Bot Telegram (protocolo non-custodial, P2P trading)
- [x] **Fase 6** - Deploy (spaincoin.es + HTTPS + bot en producción)
- [ ] **Fase 7** - Comunidad (crecer base de usuarios, marketing, validadores externos)
- [ ] **Fase 8** - Mainnet (SL + licencia CASP + KYC + claves nuevas)

> Exchange custodial (rama `exchange-v1`) se retomará cuando se obtenga la licencia CASP.

## Convenciones de Código
- Go: `gofmt` + `go vet` (se ejecutan automáticamente via hooks tras cada edición)
- Commits: `feat:`, `fix:`, `docs:`, `test:` prefijos
- Cada módulo tiene sus propios tests unitarios en `*_test.go`
- Documentar todas las funciones públicas con comentarios Go

## Comandos Útiles

```bash
# Tests (core + node/rpc + integración)
go test ./core/... ./node/rpc/ ./tests/

# Build nodo
go build -o spaincoin ./node/cmd/

# Correr el nodo
go run ./node/cmd/main.go

# Build CLI
go build -o spc ./cli/

# Generar wallet
./spc wallet new

# Build bot
go build -o spaincoin-bot ./bot/

# Build frontend
cd frontend && npm run build

# Ver estado del proyecto
/status
```

## Slash Commands Disponibles
- `/status` — muestra estado de tests, build y estructura
- `/new-module [nombre]` — crea estructura base para nuevo módulo del core

## Tokenomics $SPC
- **Supply máximo**: 21,000,000 SPC
- **Génesis**: 5,000,000 SPC al fundador (SPC5e2ac672...ea7349f)
- **Hot wallet**: 50,000 SPC (SPCc119f94a...d65481) para operaciones diarias del bot
- **Validadores**: Recompensa por bloque (se acumula con el tiempo)
- **Decimales**: 18 (1 SPC = 10^18 pesetas — unidad mínima)
- **Stake mínimo validador**: 1 SPC = 1_000_000_000_000_000_000 pesetas
- **Precio**: auto-escalado por tiers de SPC vendidos (no fijado manualmente)
- **Pago**: transferencia bancaria a IBAN Revolut (ES87 1583 0001 1890 5361 0687)

## Variables de Entorno

### Nodo (VPS 1)
```
SPC_VALIDATOR_KEY=<hex>    # Clave privada del validador
SPC_VALIDATOR_ADDRESS=<addr>
SPC_RPC_PORT=8545
SPC_P2P_PORT=30303
SPC_BLOCK_TIME=5
SPC_DATA_DIR=./data
SPC_LOG_LEVEL=info
```

### Web + API + Bot (VPS 2)
```
SPC_NODE_URL=http://204.168.176.40:8545
PORT=3001
SPC_BOT_TOKEN=<telegram bot token>
SPC_ADMIN_CHAT_ID=<admin telegram chat id>
SPC_HOT_WALLET_KEY=<hex — clave hot wallet>
SPC_IBAN=ES87 1583 0001 1890 5361 0687
SPC_ALLOWED_ORIGIN=https://spaincoin.es
```

## Deploy

### VPS 1 (Nodo)
```bash
ssh root@204.168.176.40
systemctl status spaincoin
journalctl -u spaincoin -f
curl http://204.168.176.40:8545/api/status
```

### VPS 2 (Web + Bot)
```bash
ssh root@46.62.201.94
systemctl status spaincoin-exchange
systemctl status spaincoin-bot
journalctl -u spaincoin-bot -f
```
