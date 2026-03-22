# Rotación de Claves — Guía para Mainnet
> Cómo sustituir las claves del validador y del fundador sin depender de nadie.
> Este proceso lo haces TÚ solo, sin que nadie vea tus claves privadas.

---

## Por qué rotar las claves antes de mainnet

Las claves actuales de testnet:
- Fueron generadas durante el desarrollo (pasaron por una conversación de IA)
- Solo tienen valor en testnet — no hay dinero real en juego
- **Antes de mainnet hay que sustituirlas por claves generadas offline por ti**

Con dinero real, la única persona que debe generar y ver las claves eres tú.

---

## Paso 1 — Genera las nuevas claves TÚ SOLO

**En tu Mac, desconectado de internet si quieres ser muy estricto:**

```bash
cd /Users/luison/Documents/spaincoin

# Asegúrate de tener el CLI compilado
go build -o spc ./cli/

# Genera wallet del validador (la que produce bloques)
./spc wallet new
```

Verás algo así:
```
Dirección:     SPCabcdef1234...
Clave privada: a1b2c3d4e5f6...
```

**Anótalo en papel inmediatamente. No lo guardes en el ordenador.**

Repite para las wallets de fundador (dev, marketing) si quieres renovarlas también.

---

## Paso 2 — Actualiza las claves en VPS 1 (nodo)

```bash
ssh root@204.168.176.40
nano /var/spaincoin/.env
```

Cambia estas líneas con tus nuevas claves:
```
SPC_VALIDATOR_KEY=TU_NUEVA_CLAVE_PRIVADA_HEX
SPC_VALIDATOR_ADDRESS=TU_NUEVA_DIRECCION_SPC
```

Guarda (`Ctrl+O`, `Enter`, `Ctrl+X`) y reinicia el nodo:

```bash
systemctl restart spaincoin
systemctl status spaincoin
```

Verifica que sigue produciendo bloques:
```bash
curl http://localhost:8545/api/status
```

---

## Paso 3 — Actualiza la wallet genesis (si cambias el validador)

Si cambias la dirección del validador, también hay que actualizar el bloque génesis para que la pre-mina inicial vaya a la nueva dirección.

> ⚠️ **En mainnet esto implica un reinicio completo de la cadena** — todos los bloques anteriores se borran. Solo se hace UNA VEZ antes del lanzamiento definitivo. Nunca después.

```bash
# En VPS 1: borrar datos antiguos
systemctl stop spaincoin
rm -rf /var/spaincoin/data/*
systemctl start spaincoin
```

La cadena arranca desde bloque 0 con la nueva dirección del validador.

---

## Paso 4 — Rota el JWT Secret del exchange

El JWT secret cifra los tokens de sesión de los usuarios. Rotarlo invalida todas las sesiones activas (los usuarios tendrán que volver a loguearse — no pierden fondos).

```bash
# Genera nuevo secret
openssl rand -hex 32

# Actualiza en VPS 2
ssh root@46.62.201.94
nano /var/spaincoin-exchange/.env
# Cambia SPC_JWT_SECRET por el nuevo valor

systemctl restart spaincoin-exchange
```

---

## Checklist rotación completa (antes de mainnet)

```
[ ] Generar nuevas claves en local (offline)
[ ] Anotar claves en papel — 2 copias en lugares distintos
[ ] Actualizar SPC_VALIDATOR_KEY en VPS 1 .env
[ ] Parar nodo + borrar datos + reiniciar (nueva cadena limpia)
[ ] Verificar que el nodo produce bloques con la nueva clave
[ ] Rotar JWT secret en VPS 2
[ ] Verificar que el exchange API responde
[ ] Borrar cualquier rastro de claves antiguas en ordenadores
[ ] Test end-to-end: crear cuenta, ver saldo, enviar tx
```

---

## Seguridad de las claves en producción

| Clave | Dónde vive | Quién la ve |
|-------|-----------|-------------|
| Validador privkey | VPS 1 `.env` (chmod 600) | Solo tú via SSH |
| Wallets fundador | Papel físico | Solo tú |
| JWT Secret | VPS 2 `.env` (chmod 600) | Solo tú via SSH |
| Contraseñas usuarios | BoltDB (bcrypt hash) | Nadie — solo el hash |
| Claves privadas usuarios | BoltDB (AES-256-GCM cifrado) | Nadie — solo cifrado con contraseña del usuario |

**Nunca:**
- Compartas claves privadas por WhatsApp, email, Telegram
- Las guardes en Google Drive, iCloud, Notion
- Las metas en el repositorio de GitHub
- Las pases por ningún chat (incluido este)

---

## Si crees que una clave está comprometida

1. **Mueve los fondos inmediatamente** a una wallet nueva (antes de que el atacante lo haga)
2. Genera nueva clave con `./spc wallet new`
3. Envía todos los fondos de la clave comprometida a la nueva
4. Actualiza el validador siguiendo los pasos anteriores
5. Notifica a los usuarios si el exchange estuvo en riesgo
