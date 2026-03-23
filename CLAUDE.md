# SpainCoin - Blockchain Project

## Proyecto
SpainCoin es una blockchain propia (Layer 1) con su criptomoneda nativa $SPC.
Objetivo final: exchange descentralizado tipo crypto.com centrado en $SPC.

## Stack Tecnológico
- **Blockchain core**: Go 1.26+ (estándar en blockchain — Ethereum, Cosmos usan Go)
- **Consenso**: Proof of Stake simplificado
- **Networking**: P2P con libp2p (Fase 2)
- **Smart contracts** (Fase 3): EVM-compatible o VM propia
- **Exchange app** (Fase 5): React + Go API

## Arquitectura del Proyecto

```
spaincoin/
├── .claude/               # Configuración Claude Code
│   ├── commands/          # Slash commands: /status, /new-module
│   └── settings.json      # Permisos + hooks (go fmt/vet automático)
├── core/                  # Blockchain core (Go)
│   ├── block/             # ✅ Transacciones, bloques, Merkle tree
│   ├── chain/             # ✅ Lógica de la cadena (AddBlock, IsValid)
│   ├── consensus/         # ✅ PoS: selección validadores, slashing
│   ├── crypto/            # ✅ ECDSA keys, SHA-256, Address SPC...
│   ├── mempool/           # ✅ Pool de transacciones pendientes
│   ├── network/           # ✅ P2P networking (libp2p, gossipsub, mDNS)
│   ├── state/             # ✅ Estado global (balances, cuentas)
│   └── wallet/            # 📋 Gestión de wallets — pendiente
├── node/                  # ✅ Nodo ejecutable (produce bloques)
│   └── cmd/main.go        # Entrypoint: go run ./node/cmd/
├── cli/                   # ✅ CLI `spc` (wallet, tx, chain)
├── docs/                  # Documentación
│   ├── guia-fundador.md   # ✅ Guía completa para el fundador
│   ├── architecture.md    # 📋 Pendiente
│   └── tokenomics.md      # 📋 Pendiente
└── tests/                 # 📋 Tests de integración — pendiente
```

## Estado Actual — Fase 6 completada

| Módulo | Estado | Tests |
|--------|--------|-------|
| core/crypto | ✅ Completo | 5 |
| core/block | ✅ Completo | 12 |
| core/state | ✅ Completo | 8 |
| core/mempool | ✅ Completo | 6 |
| core/chain | ✅ Completo | 10 |
| core/consensus | ✅ Completo | 10 |
| node/ | ✅ Compila | — |
| cli/ | ✅ Compila | — |
| core/network (P2P) | ✅ Completo | 5 |
| core/storage (BoltDB) | ✅ Completo | 6 |
| core/wallet | ✅ Completo | 4 |
| node/rpc (HTTP API) | ✅ Completo | 12 |
| tests/ integración | ✅ Completo | 7 |
| exchange/ (Go API) | ✅ Completo | — |
| exchange/market (Binance) | ✅ Completo | — |
| exchange/auth (JWT+bcrypt) | ✅ Completo | — |
| exchange/database (trades) | ✅ Completo | — |
| frontend/ (React) | ✅ Completo | — |

**Total: 85 tests / 85 PASS**

## Infraestructura de producción

| Componente | Ubicación | Estado |
|-----------|-----------|--------|
| Nodo blockchain | VPS 1: 204.168.176.40 | ✅ Corriendo 24/7 |
| Exchange API + Web | VPS 2: 46.62.201.94 | ✅ Corriendo 24/7 |
| Dominio | spaincoin.es | ✅ HTTPS activo |
| Backups | Cron diario en ambos VPS | ✅ Configurados |
| Firewall | UFW en ambos VPS | ✅ RPC restringido |
| SSH | Solo claves, sin password | ✅ Configurado |

## Fases del Proyecto

- [x] **Fase 1** - Core blockchain (bloques, transacciones, consenso básico)
- [x] **Fase 2** - Red P2P (múltiples nodos comunicándose con libp2p)
- [x] **Fase 3** - Wallet + CLI + persistencia
- [x] **Fase 4** - Testnet infra (VPS Hetzner, nodo en producción 24/7)
- [x] **Fase 5** - Exchange App (React + Go API) con trading multi-par
- [x] **Fase 6** - Deploy Exchange (spaincoin.es + HTTPS + seguridad)
- [ ] **Fase 7** - Mainnet (SL + PSAV + KYC + claves nuevas)

## Convenciones de Código
- Go: `gofmt` + `go vet` (se ejecutan automáticamente via hooks tras cada edición)
- Commits: `feat:`, `fix:`, `docs:`, `test:` prefijos
- Cada módulo tiene sus propios tests unitarios en `*_test.go`
- Documentar todas las funciones públicas con comentarios Go

## Comandos Útiles

```bash
# Tests
go test ./...

# Build nodo
go build -o spaincoin ./node/cmd/

# Correr el nodo
go run ./node/cmd/main.go

# Build CLI
go build -o spc ./cli/

# Generar wallet
./spc wallet new

# Ver estado del proyecto
/status
```

## Slash Commands Disponibles
- `/status` — muestra estado de tests, build y estructura
- `/new-module [nombre]` — crea estructura base para nuevo módulo del core

## Tokenomics $SPC
- **Supply máximo**: 21,000,000 SPC (guiño a Bitcoin)
- **Génesis**: 1,000,000 SPC pre-minados (fundadores, desarrollo, liquidez)
- **Validadores**: Recompensa por bloque (actualmente 0.001 SPC = 1_000_000_000_000 pesetas)
- **Decimales**: 18 (1 SPC = 10^18 pesetas — unidad mínima)
- **Stake mínimo validador**: 1 SPC = 1_000_000_000_000_000_000 pesetas

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

### Exchange (VPS 2)
```
SPC_NODE_URL=http://204.168.176.40:8545
PORT=3001
SPC_JWT_SECRET=<hex>       # openssl rand -hex 32
SPC_ALLOWED_ORIGIN=https://spaincoin.es
```

## Exchange Features
- **Auth**: JWT (HS256, 7 días) + bcrypt (cost 12) + AES-256-GCM para claves
- **Trading**: Compra/venta cualquier par (SPC/EUR, BTC/EUR, ETH/EUR...)
- **Precios**: Binance API (gratis, actualización cada 30s) + simulador para SPC
- **Portfolio**: Holdings, valor total, depósito EUR (testnet: 1000€ iniciales)
- **Legal**: Términos, Privacidad (RGPD), Riesgos (MiCA), Cookies
- **Seguridad**: Rate limiting, CORS, security headers, audit logs

## Deploy Exchange (VPS 2)
```bash
# Build y deploy completo
ssh root@46.62.201.94
cd /opt/spaincoin-exchange && git pull
CGO_ENABLED=0 go build -o /usr/local/bin/spaincoin-exchange ./exchange/
cd frontend && npm run build && cp -r dist/* /var/www/spaincoin/
systemctl restart spaincoin-exchange
```
