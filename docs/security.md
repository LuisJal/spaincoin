# Seguridad — SpainCoin

## Principios

1. **Non-custodial**: no custodiamos fondos de usuarios. Cada usuario controla sus claves.
2. Las claves privadas se generan client-side (en el navegador del usuario). NUNCA salen del dispositivo.
3. El nodo RPC solo acepta peticiones desde VPS 2 (firewall).
4. Rate limiting en todas las APIs públicas.
5. HTTPS obligatorio en producción.
6. El precio de SPC NUNCA se modifica vía Telegram. Solo auto-tiers o SSH.

## Capas de seguridad

### Nodo blockchain (VPS 1 — 204.168.176.40)
- **Firewall (UFW)**: solo puertos 22, 30303, 8545
- **RPC (8545)**: restringido SOLO a IP de VPS 2 (46.62.201.94)
- **Rate limiting**: 60 req/min por IP
- **Clave validador**: solo en `/var/spaincoin/.env` (chmod 600)
- **SSH**: key-only (password deshabilitado)
- **Backups**: cron diario en `/var/backups/spaincoin/`

### Web + API + Bot (VPS 2 — 46.62.201.94)
- **Firewall (UFW)**: solo puertos 22, 80, 443
- **Rate limiting**: 100 req/min por IP
- **HTTPS**: Let's Encrypt (renovación automática)
- **fail2ban**: activo
- **SSH**: key-only (password deshabilitado)
- **Backups**: cron diario en `/var/backups/spaincoin/`
- **Headers seguridad**: CSP, X-Frame-Options, etc. (via nginx)
- **CORS**: solo spaincoin.es permitido

### Web wallet (client-side)
- Generación de claves 100% en el navegador
- La clave privada NUNCA se envía al servidor
- Self-custody: el usuario es responsable de guardar su clave
- Sin login, sin registro, sin sesiones

### Bot de Telegram
- Admin chat ID verificado para operaciones sensibles
- Precio auto-escalado (no modificable vía Telegram)
- Hot wallet con cantidad limitada de SPC (50,000 SPC)
- Logs de todas las operaciones
- Estructura admin planificada: 1 super admin + 2 admins

### Hot wallet
- Dirección: SPCc119f94ab074c970dc129884163fc00106d65481
- Cantidad limitada: 50,000 SPC
- Propósito: operaciones diarias del bot (auto-envío a compradores)
- Se recarga periódicamente desde la wallet del fundador
- Si se comprometiera, la pérdida máxima está acotada

---

## Hardening SSH (ya aplicado)

En ambos VPS, login por contraseña deshabilitado (solo SSH key):
```bash
# Verificar configuración
grep PasswordAuthentication /etc/ssh/sshd_config
# Debe mostrar: PasswordAuthentication no
```

Ver intentos de acceso fallidos:
```bash
journalctl -u ssh --since "1 hour ago" | grep "Failed"
```

---

## Secretos del sistema

| Secreto | Ubicación | Riesgo si se filtra |
|---------|-----------|-------------------|
| SSH key | Tu Mac ~/.ssh/ | Acceso total a servidores |
| Validator key | VPS 1 .env | Control de los fondos del validador |
| Hot wallet key | VPS 2 .env | Pérdida de hasta 50,000 SPC |
| Bot token | VPS 2 .env | Control del bot de Telegram |
| IBAN | VPS 2 .env | Bajo riesgo (es público en transferencias) |

Ver [rotacion-claves.md](rotacion-claves.md) para guía de rotación de cada secreto.

---

## Checklist de seguridad

- [x] SSH key-only en VPS 1 y VPS 2
- [x] Firewall VPS 1: RPC restringido a IP de VPS 2
- [x] Firewall VPS 2: solo 22, 80, 443
- [x] HTTPS activo en spaincoin.es
- [x] fail2ban activo en VPS 2
- [x] Backups diarios en ambos VPS
- [x] Hot wallet con cantidad limitada
- [x] Precio no modificable vía Telegram
- [ ] Auditoría externa del código Go (antes de mainnet)
- [ ] Test de penetración (antes de mainnet)
- [ ] Monitoring y alertas (uptime + errores)
- [ ] Bug bounty program (cuando crezca la comunidad)

## Amenazas conocidas y mitigaciones

| Amenaza | Mitigación |
|---------|-----------|
| DDoS | Rate limiting + Cloudflare (gratis) |
| Robo de clave validador | .env chmod 600 + acceso SSH solo con key |
| Robo de hot wallet | Cantidad limitada (50,000 SPC), recarga manual |
| XSS en frontend | CSP headers + wallet client-side |
| Phishing via Telegram | Admin chat ID verificado, precio no modificable por bot |
| 51% attack | Solo viable en mainnet con muchos validadores |
| Bot comprometido | Hot wallet limita la exposición |

---

*Última actualización: 2026-03-24 — Protocolo non-custodial, infraestructura asegurada*
