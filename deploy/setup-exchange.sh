#!/bin/bash
set -e

echo "=== SpainCoin Exchange App Setup ==="
echo "VPS 2 — Ubuntu 22.04 / Hetzner CX22"
echo "Dominio: spaincoin.es"
echo ""

# 1. Update del sistema
apt-get update && apt-get upgrade -y

# 2. Instalar dependencias
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs git curl wget ufw fail2ban nginx certbot python3-certbot-nginx

# 3. Instalar Go 1.21 (exchange API)
wget -q https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
rm go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
export PATH=$PATH:/usr/local/go/bin
go version

# 4. Crear usuario y directorio de la app
useradd -r -s /bin/false spaincoin-exchange 2>/dev/null || true
mkdir -p /opt/spaincoin-exchange
mkdir -p /var/www/spaincoin
mkdir -p /var/spaincoin-exchange

# 5. Configurar nginx
cat > /etc/nginx/sites-available/spaincoin << 'NGINXEOF'
server {
    listen 80;
    server_name spaincoin.es www.spaincoin.es;

    # Seguridad básica
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # React app (archivos estáticos)
    location / {
        root /var/www/spaincoin;
        try_files $uri $uri/ /index.html;
        expires 1h;
        add_header Cache-Control "public, no-transform";
    }

    # Exchange API (proxy)
    location /api/ {
        proxy_pass http://127.0.0.1:3001/api/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 30s;
        proxy_connect_timeout 10s;
    }
}
NGINXEOF

ln -sf /etc/nginx/sites-available/spaincoin /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
nginx -t && systemctl reload nginx

# 6. Firewall
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw --force enable

# 7. fail2ban (protección SSH)
systemctl enable fail2ban
systemctl start fail2ban

echo ""
echo "=== Setup base completo ==="
echo ""
echo "Siguientes pasos:"
echo "  1. Subir código: git clone https://github.com/TU_USUARIO/spaincoin /opt/spaincoin-exchange"
echo "  2. Configurar variables: cat > /var/spaincoin-exchange/.env"
echo "  3. Build y deploy: bash /opt/spaincoin-exchange/deploy/deploy-exchange.sh"
echo "  4. SSL: certbot --nginx -d spaincoin.es -d www.spaincoin.es"
