# Despliegue — SpainCoin Infrastructure

## Arquitectura de producción

```
                    Internet
                        │
          ┌─────────────┼─────────────┐
          │                           │
   ┌──────▼──────────┐    ┌───────────▼──────────┐
   │  VPS 1          │    │  VPS 2                │
   │  204.168.176.40 │    │  46.62.201.94         │
   │  Hetzner CX22   │    │  Hetzner CX22         │
   │                 │    │                       │
   │  SpainCoin Node │◄───│  Web + API + Bot TG   │
   │  (blockchain)   │RPC │  nginx + React + Go   │
   │                 │    │                       │
   │  :8545 (RPC)    │    │  :80/:443 (HTTP/S)    │
   │  :30303 (P2P)   │    │  spaincoin.es         │
   └────────┬────────┘    └───────────┬───────────┘
            │ P2P                     │
   ┌────────▼────────┐    ┌───────────▼───────────┐
   │  Otros nodos    │    │  Telegram Bot API      │
   │  (comunidad)    │    │  (@SpainCoinBot)       │
   └─────────────────┘    └───────────────────────┘
```

---

## VPS 1 — Nodo SpainCoin

| Parámetro | Valor |
|-----------|-------|
| IP | `204.168.176.40` |
| Tipo | Hetzner CX22 |
| SO | Ubuntu 22.04 |
| Puertos | 22 (SSH), 30303 (P2P), 8545 (RPC) |
| Estado | ✅ Activo, produciendo bloques cada 5s |

### Variables de entorno (en `/var/spaincoin/.env`)
```
SPC_VALIDATOR_KEY=<hex — nunca en git>
SPC_VALIDATOR_ADDRESS=<address>
SPC_DATA_DIR=/var/spaincoin/data
SPC_RPC_PORT=8545
SPC_P2P_PORT=30303
SPC_BLOCK_TIME=5
SPC_LOG_LEVEL=info
```

### Comandos útiles VPS 1
```bash
# Estado del nodo
systemctl status spaincoin
journalctl -u spaincoin -f

# Verificar que produce bloques
curl http://204.168.176.40:8545/api/status

# Actualizar
bash /opt/spaincoin/deploy/update-node.sh
```

---

## VPS 2 — Web + API + Bot Telegram (spaincoin.es)

| Parámetro | Valor |
|-----------|-------|
| IP | `46.62.201.94` |
| Dominio | `spaincoin.es` |
| Tipo | Hetzner CX22 |
| SO | Ubuntu 22.04 |
| Puertos | 22 (SSH), 80 (HTTP), 443 (HTTPS) |
| Estado | ✅ Activo — Web + API + Bot corriendo 24/7 |

### Servicios en VPS 2

| Servicio | Systemd unit | Puerto | Descripción |
|----------|-------------|--------|-------------|
| Exchange API | spaincoin-exchange | 3001 | Backend Go, conecta con nodo |
| Bot Telegram | spaincoin-bot | — | Bot P2P trading, auto-envío SPC |
| nginx | nginx | 80/443 | Reverse proxy + frontend estático |

### Variables de entorno (en `/var/spaincoin-exchange/.env`)
```
SPC_NODE_URL=http://204.168.176.40:8545
PORT=3001
SPC_BOT_TOKEN=<telegram bot token>
SPC_ADMIN_CHAT_ID=<admin telegram chat id>
SPC_HOT_WALLET_KEY=<hex — clave hot wallet>
SPC_IBAN=ES87 1583 0001 1890 5361 0687
SPC_ALLOWED_ORIGIN=https://spaincoin.es
```

### Comandos útiles VPS 2
```bash
# Estado de los servicios
systemctl status spaincoin-exchange
systemctl status spaincoin-bot
journalctl -u spaincoin-bot -f

# Logs del bot
journalctl -u spaincoin-bot --since "1 hour ago"

# Deploy actualización web
cd /opt/spaincoin-exchange && git pull
cd frontend && npm run build && cp -r dist/* /var/www/spaincoin/
systemctl restart spaincoin-exchange

# Deploy actualización bot
cd /opt/spaincoin-exchange && git pull
CGO_ENABLED=0 go build -o /usr/local/bin/spaincoin-bot ./bot/
systemctl restart spaincoin-bot
```

### Guía detallada
Ver: [docs/setup-vps2-paso-a-paso.md](setup-vps2-paso-a-paso.md)

---

## Coste mensual

| Concepto | Coste |
|----------|-------|
| VPS 1 — Nodo blockchain | ~7.85 EUR/mes |
| VPS 2 — Web + Bot | ~7.85 EUR/mes |
| Dominio spaincoin.es | ~0.65 EUR/mes (~8 EUR/año) |
| SSL (Let's Encrypt) | Gratis |
| **Total** | **~16.35 EUR/mes** |

---

## Checklist de producción

- [x] VPS 1 configurado — nodo corriendo y produciendo bloques
- [x] VPS 2 configurado con nginx + SSL
- [x] Dominio spaincoin.es comprado y DNS configurado
- [x] HTTPS activo (Let's Encrypt)
- [x] Web React desplegada (informacional + wallet + explorer)
- [x] Exchange API en producción
- [x] Bot Telegram operativo (P2P trading, auto-envío SPC)
- [x] Firewall configurado en ambos VPS (RPC restringido)
- [x] SSH key-only en ambos VPS
- [x] Backups diarios en ambos VPS
- [x] fail2ban activo en VPS 2
- [ ] Monitoring con alertas (uptime + errores)
- [ ] CDN para frontend (Cloudflare)

---

## Seguridad

- **SSH**: solo claves, sin password, en ambos VPS
- **Firewall (UFW)**: puertos restringidos, RPC solo accesible desde VPS 2
- **fail2ban**: activo en VPS 2
- **Backups**: cron diario en ambos VPS (`/var/backups/spaincoin/`)
- **SSL**: Let's Encrypt, renovación automática
- **Hot wallet**: cantidad limitada de SPC para operaciones diarias
- **Precio SPC**: nunca se cambia vía Telegram, solo auto-tiers o SSH

---

## Escalado futuro

Cuando crezca la comunidad:
- Añadir validadores externos (más nodos P2P)
- VPS 2 -> múltiples servidores + load balancer (nginx upstream)
- CDN para el frontend (Cloudflare — tier gratis)
- Monitoring con Grafana + Prometheus
- Segundo nodo de respaldo para la blockchain
