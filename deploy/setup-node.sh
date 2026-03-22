#!/bin/bash
set -e

echo "=== SpainCoin Node Setup ==="
echo "Ubuntu 22.04 / Hetzner CX22"
echo ""

# 1. System update
apt-get update && apt-get upgrade -y

# 2. Install dependencies
apt-get install -y git curl wget ufw fail2ban

# 3. Install Go 1.21
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
export PATH=$PATH:/usr/local/go/bin

# 4. Clone repo
git clone https://github.com/LuisJal/spaincoin.git /opt/spaincoin
cd /opt/spaincoin

# 5. Build
CGO_ENABLED=0 go build -o /usr/local/bin/spaincoin ./node/cmd/

# 6. Create data directory
mkdir -p /var/spaincoin/data
chmod 700 /var/spaincoin/data

# 7. Create systemd service
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

# 8. Firewall
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 30303/tcp  # P2P
ufw allow 8545/tcp   # RPC API (solo si quieres acceso externo)
ufw --force enable

echo ""
echo "=== Setup completo ==="
echo ""
echo "Siguiente paso: configura /var/spaincoin/.env con tus claves:"
echo "  SPC_VALIDATOR_KEY=tu_clave_privada"
echo "  SPC_VALIDATOR_ADDRESS=tu_address"
echo "  SPC_DATA_DIR=/var/spaincoin/data"
echo "  SPC_RPC_PORT=8545"
echo "  SPC_P2P_PORT=30303"
echo "  SPC_BLOCK_TIME=5"
echo ""
echo "Luego: systemctl enable spaincoin && systemctl start spaincoin"
echo "Logs:  journalctl -u spaincoin -f"
