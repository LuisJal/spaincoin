# Setup VPS 2 — Exchange App (spaincoin.es)
> Guía completa paso a paso para desplegar la app del exchange.
> Actualizada: marzo 2026

---

## Datos del servidor

| Concepto | Valor |
|----------|-------|
| Proveedor | Hetzner Cloud |
| Servidor | CX22 (2 vCPU, 4GB RAM, 40GB SSD) |
| IP pública | `46.62.201.94` |
| SO | Ubuntu 22.04 LTS |
| Dominio | `spaincoin.es` |
| Nodo blockchain (VPS 1) | `204.168.176.40:8545` |

---

## Paso 1 — DNS del dominio

En el panel de Dondominio (o donde compraste el dominio), crea estos registros DNS:

| Tipo | Nombre | Valor | TTL |
|------|--------|-------|-----|
| A | `@` | `46.62.201.94` | 3600 |
| A | `www` | `46.62.201.94` | 3600 |

> ⏳ Los DNS pueden tardar hasta 24h en propagarse (normalmente 5-15 minutos).
> Verifica con: `nslookup spaincoin.es` — debe devolver `46.62.201.94`

---

## Paso 2 — Conectar al servidor

```bash
ssh root@46.62.201.94
```

---

## Paso 3 — Setup base del servidor

```bash
# Descargar y ejecutar el script de setup
curl -fsSL https://raw.githubusercontent.com/TU_USUARIO/spaincoin/main/deploy/setup-exchange.sh | bash
```

O manualmente:
```bash
git clone https://github.com/TU_USUARIO/spaincoin /opt/spaincoin-exchange
bash /opt/spaincoin-exchange/deploy/setup-exchange.sh
```

Esto instala: Node.js 20, Go 1.21, nginx, certbot, ufw, fail2ban.

---

## Paso 4 — Clonar el repositorio

```bash
git clone https://github.com/TU_USUARIO/spaincoin /opt/spaincoin-exchange
```

---

## Paso 5 — Configurar variables de entorno

```bash
cat > /var/spaincoin-exchange/.env << 'EOF'
# URL del nodo blockchain (VPS 1)
SPC_NODE_URL=http://204.168.176.40:8545

# Puerto donde escucha la API
PORT=3001

# Secreto JWT (cámbialo por uno aleatorio seguro)
SPC_JWT_SECRET=CAMBIA_ESTO_POR_UNA_CADENA_ALEATORIA_LARGA

# Origen permitido para CORS
SPC_ALLOWED_ORIGIN=https://spaincoin.es
EOF

chmod 600 /var/spaincoin-exchange/.env
```

> ⚠️ **Genera un JWT secret seguro:**
> ```bash
> openssl rand -hex 32
> ```
> Copia el resultado y reemplaza `CAMBIA_ESTO_POR_UNA_CADENA_ALEATORIA_LARGA`

---

## Paso 6 — Deploy de la app

```bash
bash /opt/spaincoin-exchange/deploy/deploy-exchange.sh
```

Este script:
1. Hace `git pull` para tener la última versión
2. Build del frontend React (`npm run build`)
3. Copia los archivos a `/var/www/spaincoin`
4. Build del binario Go de la API
5. Crea y arranca el servicio systemd

---

## Paso 7 — SSL con Let's Encrypt (HTTPS gratuito)

> ⚠️ Primero asegúrate de que el DNS ya apunta a la IP correcta.

```bash
certbot --nginx -d spaincoin.es -d www.spaincoin.es
```

Sigue las instrucciones:
- Email para notificaciones: tu email
- Acepta los términos: `Y`
- ¿Compartir email?: `N`
- Redirigir HTTP a HTTPS: **`2`** (recomendado)

Let's Encrypt se renueva automáticamente cada 90 días.

---

## Paso 8 — Verificar que todo funciona

```bash
# API del exchange
curl https://spaincoin.es/api/status

# Debería devolver algo como:
# {"node":{"height":...,"status":"ok"},"exchange":{"version":"0.1.0"}}
```

También abre el navegador en `https://spaincoin.es` — deberías ver la web.

---

## Comandos útiles

```bash
# Ver estado del servicio
systemctl status spaincoin-exchange

# Ver logs en tiempo real
journalctl -u spaincoin-exchange -f

# Reiniciar el servicio
systemctl restart spaincoin-exchange

# Actualizar cuando haya nueva versión
cd /opt/spaincoin-exchange && git pull
bash /opt/spaincoin-exchange/deploy/deploy-exchange.sh

# Ver logs de nginx
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
```

---

## Troubleshooting

**La web no carga:**
```bash
systemctl status nginx
nginx -t  # comprueba la config de nginx
```

**La API no responde:**
```bash
systemctl status spaincoin-exchange
journalctl -u spaincoin-exchange -n 50
curl http://localhost:3001/api/status  # sin nginx, directo
```

**El nodo no conecta:**
```bash
curl http://204.168.176.40:8545/api/status
# Si falla: el nodo está caído o el firewall del VPS 1 bloquea
```

**SSL no funciona:**
```bash
certbot renew --dry-run  # test de renovación
```

---

## Arquitectura final

```
Usuario (navegador)
        │
        ▼ HTTPS 443
┌─────────────────────┐
│  VPS 2 (46.62.201.94)│
│  spaincoin.es        │
│                     │
│  nginx              │
│  ├── / → React app  │
│  └── /api → :3001   │
│                     │
│  spaincoin-exchange │
│  (Go API :3001)     │
└──────────┬──────────┘
           │ HTTP :8545
           ▼
┌─────────────────────┐
│  VPS 1 (204.168.176.40)│
│  Nodo SpainCoin     │
│  (blockchain)       │
└─────────────────────┘
```
