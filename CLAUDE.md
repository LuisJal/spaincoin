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
│   ├── network/           # 🔄 P2P networking (libp2p) — pendiente
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

## Estado Actual — Fase 1 en curso

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
| core/network (P2P) | 🔄 Pendiente | — |
| tests/ integración | 🔄 Pendiente | — |

**Total: 51 tests / 51 PASS**

## Fases del Proyecto

- [x] **Fase 1** - Core blockchain (bloques, transacciones, consenso básico, nodo, CLI)
- [ ] **Fase 2** - Red P2P (múltiples nodos comunicándose con libp2p)
- [ ] **Fase 3** - Wallet avanzada + persistencia en disco
- [ ] **Fase 4** - Testnet pública
- [ ] **Fase 5** - Exchange app (React + Go API)
- [ ] **Fase 6** - Mainnet

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

## Variables de Entorno (nodo)
```
SPC_VALIDATOR_KEY=<hex>    # Clave privada del validador
SPC_VALIDATOR_ADDRESS=<addr>
SPC_RPC_PORT=8545
SPC_P2P_PORT=30303
SPC_BLOCK_TIME=5
SPC_DATA_DIR=./data
SPC_LOG_LEVEL=info
```
