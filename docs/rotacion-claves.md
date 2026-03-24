# Rotación de Claves — Guía Completa

> Cuándo y cómo cambiar cada secreto del sistema.

---

## 1. Clave del Bot de Telegram

**Dónde está:** VPS 2 -> `/var/spaincoin-exchange/.env` -> `SPC_BOT_TOKEN`

**Cuándo cambiar:** Si sospechas que alguien la tiene, o antes de mainnet.

**Cómo:**
1. Abre Telegram -> @BotFather -> `/revoke` -> selecciona tu bot
2. Te da un token nuevo
3. En VPS 2:
```bash
ssh root@46.62.201.94
nano /var/spaincoin-exchange/.env
# Cambia SPC_BOT_TOKEN por el nuevo
systemctl restart spaincoin-bot
```

---

## 2. Clave SSH (acceso a servidores)

**Dónde está:** Tu Mac -> `~/.ssh/id_ed25519`

**Cuándo cambiar:** Si pierdes el Mac, o por precaución anual.

**Cómo:**
```bash
# En tu Mac — genera nueva clave
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519_new

# Sube la nueva a ambos VPS (desde sesión existente)
ssh root@204.168.176.40
echo "TU_NUEVA_CLAVE_PUBLICA" >> ~/.ssh/authorized_keys

ssh root@46.62.201.94
echo "TU_NUEVA_CLAVE_PUBLICA" >> ~/.ssh/authorized_keys

# Verifica que puedes entrar con la nueva
ssh -i ~/.ssh/id_ed25519_new root@204.168.176.40

# Si funciona, borra la vieja de authorized_keys en ambos VPS
# Y renombra la nueva en tu Mac:
mv ~/.ssh/id_ed25519_new ~/.ssh/id_ed25519
mv ~/.ssh/id_ed25519_new.pub ~/.ssh/id_ed25519.pub
```

---

## 3. Clave del Validador (SPC_VALIDATOR_KEY)

**Dónde está:** VPS 1 -> `/var/spaincoin/.env` -> `SPC_VALIDATOR_KEY`

**Cuándo cambiar:** Antes de mainnet (obligatorio). La clave actual pasó por conversaciones anteriores.

**Cómo:**
```bash
# En tu Mac — genera nueva clave (SIN que nadie vea)
./spc wallet new
# Apunta la clave privada en PAPEL
# Apunta la dirección SPC...

# En VPS 1:
ssh root@204.168.176.40
nano /var/spaincoin/.env
# Cambia SPC_VALIDATOR_KEY y SPC_VALIDATOR_ADDRESS
systemctl restart spaincoin

# Verifica
journalctl -u spaincoin -f
```

**IMPORTANTE:** Al cambiar la clave del validador, la nueva dirección empieza con 0 SPC. Necesitas transferir los fondos de la wallet antigua a la nueva.

---

## 4. Clave de la Hot Wallet (SPC_HOT_WALLET_KEY)

**Dónde está:** VPS 2 -> `/var/spaincoin-exchange/.env` -> `SPC_HOT_WALLET_KEY`

**Cuándo cambiar:** Si sospechas filtración, o por precaución periódica.

**Cómo:**
```bash
# En tu Mac — genera nueva wallet
./spc wallet new
# Apunta la clave privada y la dirección SPC...

# Transfiere los SPC de la hot wallet vieja a la nueva
./spc tx send --from <vieja> --to <nueva> --amount <balance>

# En VPS 2:
ssh root@46.62.201.94
nano /var/spaincoin-exchange/.env
# Cambia SPC_HOT_WALLET_KEY por la nueva clave privada
systemctl restart spaincoin-bot

# Verifica que el bot funciona
journalctl -u spaincoin-bot -f
```

**IMPORTANTE:** La hot wallet actual es `SPCc119f94ab074c970dc129884163fc00106d65481` con 50,000 SPC. Al rotar, hay que:
1. Crear wallet nueva
2. Transferir SPC de la vieja a la nueva
3. Actualizar .env en VPS 2
4. Reiniciar el bot

---

## 5. Admin Chat ID de Telegram

**Dónde está:** VPS 2 -> `/var/spaincoin-exchange/.env` -> `SPC_ADMIN_CHAT_ID`

**Cuándo cambiar:** Si cambias de cuenta de Telegram o añades admins.

**Estructura planificada:**
- 1 super admin (fundador)
- 2 admins adicionales (cuando crezca el equipo)

---

## 6. IBAN (cuenta para pagos)

**Dónde está:** VPS 2 -> `/var/spaincoin-exchange/.env` -> `SPC_IBAN`

**Valor actual:** ES87 1583 0001 1890 5361 0687 (Revolut)

**Cuándo cambiar:** Si cambias de cuenta bancaria.

---

## Resumen rápido

| Secreto | Ubicación | Riesgo si se filtra | Prioridad rotación |
|---------|-----------|-------------------|--------------------|
| SSH key | Tu Mac ~/.ssh/ | Acceso total a servidores | Alta |
| Validator key | VPS 1 .env | Control de fondos del validador | Alta (antes de mainnet) |
| Hot wallet key | VPS 2 .env | Pérdida de hasta 50,000 SPC | Media |
| Bot token | VPS 2 .env | Control del bot de Telegram | Media |
| IBAN | VPS 2 .env | Bajo riesgo (es público en transferencias) | Baja |

**Regla de oro:** si dudas, cambia. Es gratis y tarda 2 minutos.

---

*Última actualización: 2026-03-24 — Añadida hot wallet, estructura admin*
