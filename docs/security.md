# Seguridad — SpainCoin

## Principios
1. Las claves privadas NUNCA salen del dispositivo del usuario
2. El nodo RPC solo acepta peticiones del servidor del exchange
3. Rate limiting en todas las APIs públicas
4. HTTPS obligatorio en producción
5. Auditoría antes de mainnet

## Capas de seguridad

### Nodo blockchain (VPS 1)
- Firewall: solo puertos 22, 30303, 8545
- Rate limiting: 60 req/min por IP
- Clave privada del validador: solo en /var/spaincoin/.env (chmod 600)
- RPC solo accesible desde VPS 2 (configurar firewall para restringir 8545 a IP del VPS 2)

### Exchange API (VPS 2)
- Rate limiting: 100 req/min por IP
- Validación estricta de todas las entradas
- Logs de auditoría para todas las transacciones
- Headers de seguridad en todas las respuestas
- HTTPS con Let's Encrypt (obligatorio antes de lanzar)

### Frontend React
- Sin claves privadas en el navegador nunca
- CSP headers en producción (via nginx)
- Firma de transacciones siempre offline con CLI

## Hardening SSH (pendiente — aplicar antes de mainnet)

En ambos VPS, deshabilitar login por contraseña (solo SSH key):
```bash
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sed -i 's/PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
systemctl restart ssh
```

Ver intentos de acceso fallidos:
```bash
journalctl -u ssh --since "1 hour ago" | grep "Failed"
```

> ⚠️ Solo aplicar si tienes la SSH key configurada y funcionando — si la pierdes, te quedas sin acceso.

---

## Antes del lanzamiento (checklist)
- [ ] Deshabilitar login SSH por contraseña en VPS 1 y VPS 2
- [ ] Auditoría externa del código Go
- [ ] Test de penetración del exchange
- [ ] HTTPS activo en VPS 2
- [ ] Firewall VPS 1: restringir puerto 8545 solo a IP de VPS 2
- [ ] Backup automático de /var/spaincoin/data
- [ ] Monitoring y alertas (uptime + errores)
- [ ] Bug bounty program

## Amenazas conocidas y mitigaciones
| Amenaza | Mitigación |
|---------|-----------|
| DDoS | Rate limiting + Cloudflare (gratis) |
| Robo de clave validador | .env chmod 600 + acceso SSH solo con key |
| XSS en frontend | CSP headers + no innerHTML |
| SQL injection | No hay SQL — blockchain es el estado |
| 51% attack | Solo viable en mainnet con muchos validadores |
| Rug pull | Vesting 2 años en contrato (Fase 3) |
