#!/bin/bash
set -e

echo "=== SpainCoin Exchange App Setup ==="
echo "Ubuntu 22.04 / Hetzner CX22"
echo ""

# 1. Update
apt-get update && apt-get upgrade -y

# 2. Install Node.js 20 (for React app)
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs git curl wget ufw fail2ban nginx certbot python3-certbot-nginx

# 3. Install Go (for exchange API)
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
export PATH=$PATH:/usr/local/go/bin

# 4. nginx config (reverse proxy)
cat > /etc/nginx/sites-available/spaincoin << 'EOF'
server {
    listen 80;
    server_name exchange.spaincoin.com;

    # React app (static files)
    location / {
        root /var/www/spaincoin;
        try_files $uri $uri/ /index.html;
    }

    # Exchange API
    location /api/ {
        proxy_pass http://localhost:3001/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
EOF

ln -sf /etc/nginx/sites-available/spaincoin /etc/nginx/sites-enabled/
nginx -t && systemctl reload nginx

# 5. Firewall
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80/tcp   # HTTP
ufw allow 443/tcp  # HTTPS
ufw --force enable

echo ""
echo "=== Setup completo ==="
echo "Siguiente: desplegar la app React en /var/www/spaincoin"
echo "SSL: certbot --nginx -d exchange.spaincoin.com"
