# Setup VPS Paso a Paso — SpainCoin

> Guía completa para desplegar el nodo SpainCoin en Hetzner desde cero.
> Probada el 2026-03-22.

---

## Requisitos previos

- Cuenta en [hetzner.com](https://hetzner.com) con tarjeta verificada
- Terminal en tu Mac
- Repositorio en GitHub: `https://github.com/LuisJal/spaincoin`
- Tus 3 claves privadas guardadas en papel (Wallet #1, #2, #3)

---

## VPS 1 — Nodo SpainCoin

### 1. Crear el servidor en Hetzner

1. Dashboard → **"Add Server"**
2. Configuración:

| Campo | Valor |
|-------|-------|
| Location | Nuremberg o Helsinki |
| Image | **Ubuntu 22.04** |
| Type | Shared CPU → **CX22** |
| IPv4 | Activado |
| Name | `spaincoin-node` |

3. Antes de crear — añadir SSH key. En tu Mac:
```bash
cat ~/.ssh/id_ed25519.pub 2>/dev/null || ssh-keygen -t ed25519 -C "spaincoin" && cat ~/.ssh/id_ed25519.pub
```
Copia el resultado (`ssh-ed25519 AAAA...`) → Hetzner → **"Add SSH Key"** → pegar → nombre: `mi-mac`.

4. Click **"Create & Buy Now"**
5. Anotar la IP pública que aparece (ej: `65.21.xxx.xxx`)

---

### 2. Conectar al servidor

```bash
ssh root@IP_DEL_SERVIDOR
```
Primera vez pregunta "Are you sure?" → escribir `yes` + Enter.

---

### 3. Actualizar el sistema e instalar dependencias

```bash
apt-get update && apt-get upgrade -y && apt-get install -y git curl wget ufw fail2ban
```
Tarda 1-2 minutos.

---

### 4. Instalar Go

```bash
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz && \
tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz && \
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && \
export PATH=$PATH:/usr/local/go/bin && \
go version
```
Debe mostrar: `go version go1.21.6 linux/amd64`

---

### 5. Clonar el repositorio y compilar

```bash
git clone https://github.com/LuisJal/spaincoin.git /opt/spaincoin && \
cd /opt/spaincoin && \
CGO_ENABLED=0 go build -o /usr/local/bin/spaincoin ./node/cmd/ && \
echo "BUILD OK"
```
Tarda 2-3 minutos (descarga dependencias + compila). Debe terminar con `BUILD OK`.

---

### 6. Crear directorios y configurar claves

```bash
mkdir -p /var/spaincoin/data
chmod 700 /var/spaincoin/data
```

Crear el archivo de configuración con tus claves (**nunca subir esto a GitHub**):
```bash
cat > /var/spaincoin/.env << 'EOF'
SPC_VALIDATOR_KEY=TU_CLAVE_PRIVADA_WALLET_1
SPC_VALIDATOR_ADDRESS=TU_ADDRESS_WALLET_1
SPC_DATA_DIR=/var/spaincoin/data
SPC_RPC_PORT=8545
SPC_P2P_PORT=30303
SPC_BLOCK_TIME=5
SPC_LOG_LEVEL=info
EOF
chmod 600 /var/spaincoin/.env
```

Sustituir:
- `TU_CLAVE_PRIVADA_WALLET_1` → la private key de la Wallet #1 (la que tienes en papel)
- `TU_ADDRESS_WALLET_1` → `SPCeefd4d724f8c0262395aa144a0a6bb624609a25b`

---

### 7. Instalar como servicio del sistema (arranque automático)

```bash
cat > /etc/systemd/system/spaincoin.service << 'EOF'
[Unit]
Description=SpainCoin Blockchain Node
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=/var/spaincoin
EnvironmentFile=/var/spaincoin/.env
ExecStart=/usr/local/bin/spaincoin
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=spaincoin

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable spaincoin
systemctl start spaincoin
```

---

### 8. Verificar que está corriendo

```bash
# Estado del servicio
systemctl status spaincoin

# Ver logs en tiempo real
journalctl -u spaincoin -f

# Comprobar que la API responde en local
curl http://localhost:8545/status

# Comprobar desde fuera (en tu Mac)
curl http://IP_DEL_SERVIDOR:8545/status
```

La API debe responder algo como:
```json
{"status":"ok","height":19,"latest_hash":"...","total_supply":1019000000000000}
```

✅ Si ves esto, el nodo está corriendo y accesible desde internet.

---

### 9. Configurar firewall

```bash
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 30303/tcp
ufw allow 8545/tcp
ufw --force enable
ufw status
```

---

### 10. Actualizar el nodo cuando haya nueva versión

```bash
cd /opt/spaincoin && \
git pull origin main && \
CGO_ENABLED=0 go build -o /usr/local/bin/spaincoin ./node/cmd/ && \
systemctl restart spaincoin && \
echo "Actualizado OK"
```

---

## Comandos de gestión del nodo

```bash
# Ver estado
systemctl status spaincoin

# Ver logs (tiempo real)
journalctl -u spaincoin -f

# Ver últimas 50 líneas de logs
journalctl -u spaincoin -n 50

# Reiniciar nodo
systemctl restart spaincoin

# Parar nodo
systemctl stop spaincoin

# Arrancar nodo
systemctl start spaincoin

# Consultar API
curl http://localhost:8545/status
curl http://localhost:8545/block/0
curl http://localhost:8545/block/latest
```

---

## VPS 2 — Exchange App

*(Se configurará en Fase 5 cuando construyamos la app)*

Resumen del proceso:
1. Crear segundo servidor CX22 en Hetzner (mismo proceso que VPS 1)
2. `bash /opt/spaincoin/deploy/setup-exchange.sh`
3. Apuntar dominio a IP del VPS 2
4. Desplegar app React en `/var/www/spaincoin`
5. SSL con `certbot --nginx -d tudominio.com`
6. La app se conecta al nodo via `http://IP_VPS1:8545`

---

## Costes mensuales

| Concepto | Coste |
|----------|-------|
| VPS 1 — Nodo | ~8€/mes |
| VPS 2 — Exchange (Fase 5) | ~8€/mes |
| Dominio .com | ~1€/mes |
| SSL Let's Encrypt | Gratis |
| **Total actual** | **~8€/mes** |
| **Total en Fase 5** | **~17€/mes** |

---

## Solución de problemas comunes

**El nodo no arranca:**
```bash
journalctl -u spaincoin -n 50
# Busca líneas con "error" o "fatal"
```

**No puedo conectar por SSH:**
- Comprueba que la IP es correcta en el dashboard de Hetzner
- Verifica que la SSH key está añadida en Hetzner

**La API no responde:**
```bash
systemctl status spaincoin  # ¿está corriendo?
curl http://localhost:8545/status  # ¿responde en local?
ufw status  # ¿está el puerto 8545 abierto?
```

**Actualizar Go en el futuro:**
```bash
rm -rf /usr/local/go
wget https://go.dev/dl/go1.XX.X.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.XX.X.linux-amd64.tar.gz
```
