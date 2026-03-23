# Rotación de Claves — Guía Completa

> Cuándo y cómo cambiar cada secreto del sistema.

---

## 1. Clave del Bot de Telegram

**Dónde está:** VPS 2 → `/var/spaincoin-exchange/.env` → `SPC_BOT_TOKEN`

**Cuándo cambiar:** Si sospechas que alguien la tiene, o antes de mainnet.

**Cómo:**
1. Abre Telegram → @BotFather → `/revoke` → selecciona tu bot
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

**Dónde está:** Tu Mac → `~/.ssh/id_ed25519`

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

**Dónde está:** VPS 1 → `/var/spaincoin/.env` → `SPC_VALIDATOR_KEY`

**Cuándo cambiar:** Antes de mainnet (obligatorio). La clave actual pasó por esta conversación.

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

## 4. JWT Secret (SPC_JWT_SECRET)

**Dónde está:** VPS 2 → `/var/spaincoin-exchange/.env` → `SPC_JWT_SECRET`

**Cuándo cambiar:** Si sospechas filtración. Al cambiarlo, TODOS los tokens de sesión existentes se invalidan.

**Cómo:**
```bash
ssh root@46.62.201.94
# Genera nuevo secret
NEW_SECRET=$(openssl rand -hex 32)
echo $NEW_SECRET

# Edita el .env
nano /var/spaincoin-exchange/.env
# Pega el nuevo SPC_JWT_SECRET

systemctl restart spaincoin-exchange
```

---

## 5. Admin Chat ID de Telegram

**Dónde está:** VPS 2 → `/var/spaincoin-exchange/.env` → `SPC_ADMIN_CHAT_ID`

**Cómo obtener tu chat ID:**
1. Escribe `/myadmin` al bot
2. Te responde con tu chat ID
3. Ponlo en el .env

---

## 6. Bizum (teléfono para pagos)

**Dónde está:** VPS 2 → `/var/spaincoin-exchange/.env` → `SPC_BIZUM_PHONE`

**Cuándo cambiar:** Si cambias de número de teléfono.

---

## Resumen rápido

| Secreto | Ubicación | Riesgo si se filtra |
|---------|-----------|-------------------|
| SSH key | Tu Mac ~/.ssh/ | Acceso total a servidores |
| Validator key | VPS 1 .env | Control de los fondos del validador |
| JWT secret | VPS 2 .env | Suplantación de sesiones |
| Bot token | VPS 2 .env | Control del bot de Telegram |
| Bizum phone | VPS 2 .env | Spam, bajo riesgo |

**Regla de oro:** si dudas, cambia. Es gratis y tarda 2 minutos.
