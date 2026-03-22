# Despliegue — SpainCoin Infrastructure

## Arquitectura de producción

```
                    Internet
                        │
          ┌─────────────┼─────────────┐
          │                           │
   ┌──────▼──────┐            ┌───────▼──────┐
   │  VPS 1 (4€) │            │  VPS 2 (4€) │
   │  Hetzner    │            │  Hetzner    │
   │  CX22       │            │  CX22       │
   │             │            │             │
   │ SpainCoin   │◄──RPC──────│ Exchange    │
   │ Node        │            │ App         │
   │             │            │ (React+API) │
   │ Puerto 30303│            │ Puerto 80   │
   │ (P2P)       │            │ Puerto 443  │
   │ Puerto 8545 │            │ (HTTPS)     │
   │ (RPC)       │            │             │
   └──────┬──────┘            └─────────────┘
          │ P2P
   ┌──────▼──────┐
   │ Otros nodos │
   │ comunidad   │
   └─────────────┘
```

## VPS 1 — Nodo SpainCoin

**Specs mínimas:** CX22 (2 vCPU, 4GB RAM, 40GB SSD)
**SO:** Ubuntu 22.04 LTS
**Puertos abiertos:** 22 (SSH), 30303 (P2P), 8545 (RPC — solo para VPS 2)

### Setup inicial
```bash
scp deploy/setup-node.sh root@IP_VPS1:/tmp/
ssh root@IP_VPS1 "bash /tmp/setup-node.sh"
```

### Configurar claves (IMPORTANTE: nunca en el repositorio)
```bash
ssh root@IP_VPS1
cat > /var/spaincoin/.env << 'EOF'
SPC_VALIDATOR_KEY=tu_clave_privada_aqui
SPC_VALIDATOR_ADDRESS=SPCtu_address_aqui
SPC_DATA_DIR=/var/spaincoin/data
SPC_RPC_PORT=8545
SPC_P2P_PORT=30303
SPC_BLOCK_TIME=5
SPC_LOG_LEVEL=info
EOF
chmod 600 /var/spaincoin/.env
```

### Arrancar el nodo
```bash
systemctl enable spaincoin
systemctl start spaincoin
systemctl status spaincoin
```

### Ver logs en tiempo real
```bash
journalctl -u spaincoin -f
```

### Actualizar cuando haya nueva versión
```bash
bash /opt/spaincoin/deploy/update-node.sh
```

---

## VPS 2 — Exchange App

**Specs mínimas:** CX22 (2 vCPU, 4GB RAM, 40GB SSD)
**SO:** Ubuntu 22.04 LTS
**Puertos abiertos:** 22 (SSH), 80 (HTTP), 443 (HTTPS)

### Setup inicial
```bash
scp deploy/setup-exchange.sh root@IP_VPS2:/tmp/
ssh root@IP_VPS2 "bash /tmp/setup-exchange.sh"
```

### Conectar con el nodo
La app exchange se comunica con el nodo via RPC:
```
http://IP_VPS1:8545
```

### SSL con Let's Encrypt (gratis)
```bash
certbot --nginx -d exchange.spaincoin.com
```
Requiere que el dominio apunte a la IP del VPS 2.

---

## Coste mensual estimado

| Concepto | Coste |
|----------|-------|
| VPS 1 — Nodo | ~4€/mes |
| VPS 2 — Exchange | ~4€/mes |
| Dominio (.com) | ~1€/mes |
| SSL (Let's Encrypt) | Gratis |
| **Total** | **~9€/mes** |

---

## Checklist antes de lanzar testnet

- [ ] VPS 1 configurado y nodo corriendo
- [ ] VPS 2 configurado con nginx
- [ ] Dominio apuntando a VPS 2
- [ ] SSL activo
- [ ] Firewall configurado en ambos VPS
- [ ] Clave privada del validador en `.env` (chmod 600, nunca en git)
- [ ] Backups automáticos de `/var/spaincoin/data/`
- [ ] Monitoring básico (uptime check)

---

## Escalado futuro

Cuando superes los 10.000 usuarios activos:
- Añadir más nodos validadores (comunidad)
- VPS 2 → múltiples servidores + load balancer
- Base de datos → PostgreSQL en servidor dedicado
- CDN para el frontend (Cloudflare gratis)
