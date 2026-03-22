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
   │  SpainCoin Node │◄───│  Exchange App         │
   │  (blockchain)   │RPC │  nginx + React + API  │
   │                 │    │                       │
   │  :8545 (RPC)    │    │  :80/:443 (HTTP/S)    │
   │  :30303 (P2P)   │    │  spaincoin.es         │
   └────────┬────────┘    └───────────────────────┘
            │ P2P
   ┌────────▼────────┐
   │  Otros nodos    │
   │  (comunidad)    │
   └─────────────────┘
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

### Setup inicial (ya realizado)
```bash
scp deploy/setup-node.sh root@204.168.176.40:/tmp/
ssh root@204.168.176.40 "bash /tmp/setup-node.sh"
```

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

## VPS 2 — Exchange App (spaincoin.es)

| Parámetro | Valor |
|-----------|-------|
| IP | `46.62.201.94` |
| Dominio | `spaincoin.es` |
| Tipo | Hetzner CX22 |
| SO | Ubuntu 22.04 |
| Puertos | 22 (SSH), 80 (HTTP), 443 (HTTPS) |
| Estado | ⏳ Pendiente de setup |

### Guía completa
Ver: [docs/setup-vps2-paso-a-paso.md](setup-vps2-paso-a-paso.md)

### Setup rápido
```bash
# 1. Conectar
ssh root@46.62.201.94

# 2. Clonar repo
git clone https://github.com/TU_USUARIO/spaincoin /opt/spaincoin-exchange

# 3. Setup base
bash /opt/spaincoin-exchange/deploy/setup-exchange.sh

# 4. Configurar .env
cat > /var/spaincoin-exchange/.env << 'EOF'
SPC_NODE_URL=http://204.168.176.40:8545
PORT=3001
SPC_JWT_SECRET=$(openssl rand -hex 32)
SPC_ALLOWED_ORIGIN=https://spaincoin.es
EOF

# 5. Deploy
bash /opt/spaincoin-exchange/deploy/deploy-exchange.sh

# 6. SSL
certbot --nginx -d spaincoin.es -d www.spaincoin.es
```

---

## Coste mensual

| Concepto | Coste |
|----------|-------|
| VPS 1 — Nodo blockchain | ~7.85€/mes |
| VPS 2 — Exchange App | ~7.85€/mes |
| Dominio spaincoin.es | ~0.65€/mes (~8€/año) |
| SSL (Let's Encrypt) | Gratis |
| **Total** | **~16.35€/mes** |

---

## Checklist lanzamiento testnet

- [x] VPS 1 configurado — nodo corriendo y produciendo bloques
- [x] Exchange app (React + Go API) desarrollada
- [x] Auth sistema login/registro con JWT
- [x] Páginas legales (Términos, Privacidad, Riesgos, Cookies)
- [ ] Dominio `spaincoin.es` comprado y DNS configurado
- [ ] VPS 2 configurado con nginx
- [ ] SSL activo (Let's Encrypt)
- [ ] Exchange app deployed en producción
- [ ] Test end-to-end en producción
- [ ] Firewall configurado en VPS 2

---

## Escalado futuro (cuando crezcas)

Cuando superes los 10.000 usuarios activos:
- Añadir validadores de la comunidad (más nodos P2P)
- VPS 2 → múltiples servidores + load balancer (nginx upstream)
- Base de datos → PostgreSQL en servidor dedicado
- CDN para el frontend (Cloudflare — gratis tier)
- Monitoring con Grafana + Prometheus
