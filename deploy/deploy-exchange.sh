#!/bin/bash
# Deploy completo del Exchange App (API Go + Frontend React)
# Ejecutar desde el VPS 2 después de setup-exchange.sh
set -e

REPO_DIR="/opt/spaincoin-exchange"
WWW_DIR="/var/www/spaincoin"
ENV_FILE="/var/spaincoin-exchange/.env"
BINARY="/usr/local/bin/spaincoin-exchange"

echo "=== Deploy SpainCoin Exchange ==="
echo ""

# Verificar que existe el .env
if [ ! -f "$ENV_FILE" ]; then
  echo "ERROR: Falta $ENV_FILE — créalo antes de continuar"
  echo "Ver: docs/setup-vps2-paso-a-paso.md"
  exit 1
fi

# 1. Pull última versión
echo "[1/5] Actualizando código..."
cd "$REPO_DIR"
git pull origin main

# 2. Build React frontend
echo "[2/5] Build del frontend React..."
cd "$REPO_DIR/frontend"
npm install --silent
npm run build
cp -r dist/* "$WWW_DIR/"
echo "      → Frontend copiado a $WWW_DIR"

# 3. Build Exchange API (Go)
echo "[3/5] Build de la Exchange API..."
cd "$REPO_DIR"
CGO_ENABLED=0 GOOS=linux go build -o "$BINARY" ./exchange/
echo "      → Binario en $BINARY"

# 4. Crear/actualizar servicio systemd
echo "[4/5] Configurando servicio systemd..."
cat > /etc/systemd/system/spaincoin-exchange.service << SVCEOF
[Unit]
Description=SpainCoin Exchange API
After=network.target

[Service]
Type=simple
User=root
EnvironmentFile=$ENV_FILE
ExecStart=$BINARY
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
SVCEOF

systemctl daemon-reload
systemctl enable spaincoin-exchange

# 5. Reiniciar servicio
echo "[5/5] Reiniciando servicio..."
systemctl restart spaincoin-exchange
sleep 2
systemctl status spaincoin-exchange --no-pager

echo ""
echo "=== Deploy completado ==="
echo "Exchange API: http://localhost:3001/api/status"
echo "Web:          https://spaincoin.es"
