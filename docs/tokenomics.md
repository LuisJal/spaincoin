# Tokenomics $SPC — Cómo funciona el dinero de SpainCoin

> En lenguaje claro. Sin tecnicismos. Para que cualquiera lo entienda.

---

## Qué es SpainCoin

SpainCoin es una blockchain propia (Layer 1) con su criptomoneda nativa $SPC.
Es un protocolo **non-custodial**: nadie custodia tus fondos. Tú generas tu wallet,
tú controlas tus claves privadas, tú decides qué hacer con tu $SPC.

---

## Los números de SpainCoin

### Supply — cuántos $SPC existen

| Concepto | Cantidad | Explicación |
|----------|---------|-------------|
| **Supply máximo** | 21.000.000 SPC | No se pueden crear más nunca. Igual que Bitcoin. |
| **Génesis (día 1)** | 5.000.000 SPC | Los que existen desde el principio |
| **Hot wallet** | 50.000 SPC | Reserva para operaciones diarias del bot |
| **Recompensa por bloque** | Variable | Lo que gana el validador por cada bloque |
| **Tiempo para llegar al máximo** | Décadas | Los bloques se producen despacio |

El número máximo (21 millones) es un **guiño a Bitcoin** y tiene un propósito: la escasez. Si hubiera billones de $SPC, cada uno valdría una miseria. Al haber solo 21 millones, cada $SPC es más escaso y por tanto más valioso si hay demanda.

---

## Distribución del génesis

```
┌──────────────────────────────────────────┐
│           5.000.000 SPC totales          │
├──────────────────────────────────────────┤
│ Fundador: SPC5e2ac672...ea7349f          │
│ 5.000.000 SPC (100% génesis)             │
├──────────────────────────────────────────┤
│ Hot wallet: SPCc119f94a...d65481         │
│ 50.000 SPC (para bot auto-envío)         │
└──────────────────────────────────────────┘
```

**Hot wallet**: contiene una cantidad limitada de SPC para las operaciones diarias del bot de Telegram. Cuando un usuario compra SPC y el admin confirma el pago, el bot envía automáticamente desde esta wallet. Se recarga periódicamente desde la wallet del fundador.

---

## Cómo se compra y vende $SPC — Modelo P2P

SpainCoin opera como protocolo **non-custodial** con trading P2P vía Telegram:

```
USUARIO quiere comprar SPC
         │
         ▼
1. Abre el bot de Telegram (@SpainCoinBot)
2. Pulsa "Comprar SPC" → ve el precio actual
3. Transfiere EUR al IBAN del proyecto (Revolut)
4. Admin confirma el pago en Telegram
5. Bot auto-envía SPC a la wallet del usuario
6. Transacción registrada en la blockchain
```

### Precio auto-escalado

El precio NO se fija manualmente. Se calcula automáticamente por tiers según la cantidad de SPC vendidos:

| SPC vendidos | Precio por SPC |
|-------------|---------------|
| 0 - 500 | 0.05 EUR |
| 500 - 1.000 | 0.08 EUR |
| 1.000 - 2.500 | 0.12 EUR |
| 2.500 - 5.000 | 0.18 EUR |
| 5.000 - 10.000 | 0.25 EUR |
| 10.000+ | El mercado decide |

Cuanto más SPC se vende, más sube el precio. Esto incentiva comprar pronto.

> El precio NUNCA se cambia vía Telegram. Solo se modifica por auto-tiers o directamente por SSH en el servidor.

---

## Pago

- **IBAN**: ES87 1583 0001 1890 5361 0687 (Revolut)
- **Método**: Transferencia bancaria
- **Confirmación**: Admin verifica el pago y confirma en Telegram
- **Envío**: Bot auto-envía SPC a la wallet del comprador

---

## Por qué la gente querría comprar $SPC

| Utilidad | Explicación |
|----------|-------------|
| **Usar la blockchain** | Para enviar transacciones en la red SpainCoin |
| **Hacer staking** | Bloquear $SPC y ganar recompensas por validar la red |
| **Especulación** | Gente que cree que subirá y quiere comprar antes |
| **Comunidad** | Formar parte del proyecto desde el principio |
| **Futuro exchange** | Cuando se obtenga la licencia CASP, $SPC será el token nativo |

---

## Las recompensas de bloque — cómo se crean nuevos $SPC

Cada 5 segundos se produce un bloque nuevo. El validador que lo produce recibe una recompensa en $SPC. Así es como se "minan" los SPC que faltan hasta llegar a los 21 millones.

```
Bloque nuevo cada 5 segundos
       ↓
Validador recibe recompensa en SPC
       ↓
Esos SPC nuevos entran en circulación
       ↓
En décadas → llegamos a 21.000.000 SPC
```

Esto se llama **emisión controlada** — la cantidad de nuevos SPC que entran al mercado es predecible y pequeña.

---

## Wallets del proyecto

| Wallet | Dirección | SPC | Uso |
|--------|-----------|-----|-----|
| Fundador | SPC5e2ac672147ea748ba1d0c27aed781995ea7349f | 5.000.000 (génesis) | Reserva principal |
| Hot wallet | SPCc119f94ab074c970dc129884163fc00106d65481 | 50.000 | Operaciones diarias bot |

---

## Costes operativos

| Concepto | Coste | Estado |
|----------|-------|--------|
| VPS 1 (nodo blockchain) | ~7.85 EUR/mes | Pagando |
| VPS 2 (web + bot) | ~7.85 EUR/mes | Pagando |
| Dominio spaincoin.es | ~8 EUR/año | Pagando |
| SSL (Let's Encrypt) | Gratis | Activo |
| **Total** | **~16.35 EUR/mes** | |

---

## Visión a futuro

1. **Ahora**: Protocolo non-custodial + P2P trading vía Telegram
2. **Próximo**: Crecer comunidad, más validadores, marketing
3. **Futuro**: Licencia CASP + exchange custodial (código en rama `exchange-v1`)
4. **Largo plazo**: Mainnet establecida con ecosistema de validadores

---

## Lo peligroso que nunca debes hacer

### El "rug pull" — el mayor enemigo de la confianza

Un rug pull es cuando el creador vende todo de golpe y se lleva el dinero. Es la razón por la que la gente desconfía de las cryptos nuevas. **Nunca lo hagas.**

### Lo que destruye el valor

| Acción | Por qué es malo |
|--------|----------------|
| Vender todos los SPC de golpe | El precio se desploma, pierdes toda credibilidad |
| Prometer rendimientos fijos | Es ilegal en España y en la UE — es un esquema Ponzi |
| Supply infinito | Si puedes crear más tokens cuando quieras, nadie los quiere |
| No tener producto real | Sin blockchain funcionando, $SPC no tiene utilidad |

---

*Última actualización: 2026-03-24 — Protocolo non-custodial en producción, P2P trading vía Telegram*
