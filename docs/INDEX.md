# SpainCoin — Índice de Documentación

> Toda la documentación del proyecto organizada por tema.
> Actualizada: marzo 2026

---

## Para empezar

| Documento | Descripción |
|-----------|-------------|
| [guia-fundador.md](guia-fundador.md) | Cómo funciona SpainCoin por dentro (sin tecnicismos) — para entender qué estás construyendo |
| [tokenomics.md](tokenomics.md) | La matemática del $SPC: supply, precio, distribución, riesgos |

---

## Infraestructura

| Documento | Descripción |
|-----------|-------------|
| [deployment.md](deployment.md) | Arquitectura completa: 2 VPS, qué corre en cada uno, costes |
| [setup-vps-paso-a-paso.md](setup-vps-paso-a-paso.md) | Guía completa para configurar VPS 1 (nodo blockchain) |
| [setup-vps2-paso-a-paso.md](setup-vps2-paso-a-paso.md) | Guía completa para configurar VPS 2 (exchange app + spaincoin.es) |

**Estado actual:**
- VPS 1: `204.168.176.40` — nodo produciendo bloques ✅
- VPS 2: `46.62.201.94` — exchange app online, pendiente SSL ⏳
- Dominio: `spaincoin.es` — DNS configurado, propagando ⏳

---

## Seguridad y Claves

| Documento | Descripción |
|-----------|-------------|
| [security.md](security.md) | Arquitectura de seguridad, amenazas, checklist |
| [key-rotation.md](key-rotation.md) | **Cómo rotar las claves antes de mainnet** — proceso offline sin depender de nadie |

**Regla de oro:** Las claves privadas solo existen en papel y en el `.env` del servidor. En ningún otro sitio.

---

## Legal y Cumplimiento

| Documento | Descripción |
|-----------|-------------|
| [legal-requirements.md](legal-requirements.md) | Todo lo necesario para operar legalmente en España/UE: PSAV, MiCA, KYC, SL, fiscalidad |

**Estado actual (testnet):**
- ✅ Términos y Condiciones publicados
- ✅ Política de Privacidad (RGPD) publicada
- ✅ Aviso de Riesgos publicado
- ✅ Política de Cookies + banner de consentimiento
- ✅ Checkboxes de consentimiento en el registro

**Antes de mainnet:**
- Constituir SL (~1.500€)
- Registro PSAV en Banco de España (~3-6 meses)
- Integrar KYC/AML (Sumsub, ~0.50€/verificación)
- Whitepaper MiCA

---

## Roadmap del Proyecto

| Fase | Estado | Descripción |
|------|--------|-------------|
| Fase 1 — Core blockchain | ✅ | Bloques, transacciones, Merkle, ECDSA |
| Fase 2 — Red P2P | ✅ | libp2p, gossipsub, mDNS |
| Fase 3 — Wallet + CLI | ✅ | CLI `spc`, wallets, persistencia BoltDB |
| Fase 4 — Testnet infra | ✅ | VPS 1 en Hetzner, nodo 24/7 |
| Fase 5 — Exchange App | ✅ | React + Go API, auth JWT, páginas legales |
| Fase 6 — Deploy Exchange | ⏳ | VPS 2 online, falta SSL con dominio propio |
| Fase 7 — Mainnet | 📋 | Rotación de claves, KYC, registro PSAV |

---

## Comandos rápidos

```bash
# Ver estado del nodo (desde cualquier sitio)
curl http://204.168.176.40:8545/api/status

# Ver estado del exchange
curl https://spaincoin.es/api/status

# Logs del nodo (en VPS 1)
ssh root@204.168.176.40 "journalctl -u spaincoin -f"

# Logs del exchange (en VPS 2)
ssh root@46.62.201.94 "journalctl -u spaincoin-exchange -f"

# Generar nueva wallet
./spc wallet new

# Actualizar exchange tras nuevo commit
ssh root@46.62.201.94 "cd /opt/spaincoin-exchange && git pull && bash deploy/deploy-exchange.sh"
```

---

## Tests

```bash
# Todos los tests
go test ./...

# Tests con detalle
go test ./... -v

# Tests de un módulo
go test ./core/crypto/...
```

**85 tests / 85 PASS** (marzo 2026)
