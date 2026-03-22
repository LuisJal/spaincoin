#!/bin/bash
set -e
cd /opt/spaincoin
git pull origin main
CGO_ENABLED=0 go build -o /usr/local/bin/spaincoin ./node/cmd/
systemctl restart spaincoin
echo "Nodo actualizado y reiniciado"
journalctl -u spaincoin -n 20
