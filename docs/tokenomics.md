# Tokenomics $SPC — Cómo funciona el dinero de SpainCoin

> En lenguaje claro. Sin tecnicismos. Para que cualquiera lo entienda.

---

## ¿Qué es un token y por qué tiene valor?

Piensa en $SPC como las **fichas de un casino**. Cuando entras al casino (el exchange de SpainCoin), necesitas cambiar tu dinero real por fichas ($SPC) para poder jugar. Esas fichas tienen valor porque:

1. Hay una cantidad limitada de ellas
2. Las necesitas para usar el casino
3. Cuanta más gente quiere entrar al casino, más valen las fichas

Cuando el casino se llena de gente, todos quieren fichas → las fichas escasean → el precio sube.

---

## Los números de SpainCoin

### Supply — cuántos $SPC existen

| Concepto | Cantidad | Explicación |
|----------|---------|-------------|
| **Supply máximo** | 21.000.000 SPC | No se pueden crear más nunca. Igual que Bitcoin. |
| **Génesis (día 1)** | 1.000.000 SPC | Los que existen desde el principio |
| **Recompensa por bloque** | ~0,000001 SPC | Lo que gana el validador por cada bloque |
| **Tiempo para llegar al máximo** | ~décadas | Los bloques se producen despacio |

El número máximo (21 millones) es un **guiño a Bitcoin** y tiene un propósito: la escasez. Si hubiera billones de $SPC, cada uno valdría una miseria. Al haber solo 21 millones, cada $SPC es más escaso y por tanto más valioso si hay demanda.

---

## Tú como fundador — qué tienes y qué significa

### Distribución del génesis (1.000.000 SPC)

```
┌─────────────────────────────────────────┐
│           1.000.000 SPC totales          │
├──────────┬──────────┬────────┬──────────┤
│ Fundador │Desarrollo│Marketing│ Liquidez │
│ 400.000  │ 300.000  │ 100.000 │ 200.000  │
│   40%    │   30%    │   10%   │   20%    │
└──────────┴──────────┴────────┴──────────┘
```

**¿Qué significa cada parte?**

- **Fundador (40% = 400.000 SPC)** — Es tuyo. Recomendamos bloquearlo 2 años (no puedes venderlo antes). Esto genera CONFIANZA: la gente sabe que no vas a salir corriendo con el dinero al día siguiente.

- **Desarrollo (30% = 300.000 SPC)** — Para pagar servidores, mejoras, contratar ayuda si crece. Es el "presupuesto de la empresa".

- **Marketing (10% = 100.000 SPC)** — Para campañas, influencers, comunidad. Aquí entra tu trabajo de publicidad.

- **Liquidez (20% = 200.000 SPC)** — Van al exchange para que la gente pueda comprar y vender desde el primer día. Sin liquidez, nadie puede comprar aunque quiera.

---

## Cómo sube el precio — la matemática simple

### La fórmula

```
Precio = Cuánto dinero quiere entrar ÷ Cuántos SPC hay disponibles
```

### Ejemplo real

**Día del lanzamiento:**
- Pones 200.000 SPC en el exchange a 0,10€ cada uno
- Market cap = 1.000.000 SPC × 0,10€ = **100.000€**

**Un mes después (si el marketing va bien):**
- 1.000 personas compran 100€ de $SPC cada una
- Entra 100.000€ de dinero nuevo
- Los SPC disponibles siguen siendo los mismos
- El precio sube a ~0,20€

**Resultado para ti:**
- Tu wallet (400.000 SPC) × 0,20€ = **80.000€**
- Eso sin vender nada — solo por que más gente quiere $SPC

**Si sube a 1€:**
- Tu wallet vale **400.000€**

**Si sube a 10€** (como le pasó a muchas cryptos):
- Tu wallet vale **4.000.000€**

Esto no es una promesa — es la matemática. Puede subir, puede bajar, puede quedarse igual. Depende de cuánta gente quiera usar SpainCoin.

---

## Por qué la gente querría comprar $SPC

Esta es la pregunta más importante. Un token sin utilidad real no vale nada. $SPC tiene (o tendrá) utilidad real porque:

| Utilidad | Explicación |
|----------|-------------|
| **Pagar fees del exchange** | Para operar en el exchange, pagas una pequeña comisión en $SPC |
| **Hacer staking** | Puedes bloquear $SPC y ganar recompensas por validar la red |
| **Especulación** | Gente que cree que subirá y quiere comprar antes |
| **Comunidad** | Proyectos que se construyen encima de SpainCoin |

---

## Lo peligroso que nunca debes hacer

### El "rug pull" — el mayor enemigo de la confianza

Un rug pull es cuando el creador vende todo de golpe y se lleva el dinero. Es la razón por la que la gente desconfía de las cryptos nuevas. **Nunca lo hagas.**

La solución: el **vesting** — comprometerte públicamente a no vender tu 40% durante 2 años. Esto se puede hacer en el código (el smart contract bloquea los fondos). Cuando lleguemos a esa fase, lo implementamos.

### Lo que destruye el valor

| Acción | Por qué es malo |
|--------|----------------|
| Vender todo tu 40% de golpe | El precio se desploma, pierdes toda credibilidad |
| Prometer rendimientos fijos ("te doy el 20% mensual") | Es ilegal en España y en la UE — es un esquema Ponzi |
| Supply infinito | Si puedes crear más tokens cuando quieras, nadie los quiere |
| No tener producto real | Sin exchange funcionando, $SPC no tiene utilidad |

---

## Las recompensas de bloque — cómo se crean nuevos $SPC

Cada 5 segundos se produce un bloque nuevo. El validador que lo produce recibe una pequeña recompensa en $SPC. Así es como se "minan" los SPC que faltan hasta llegar a los 21 millones.

```
Bloque nuevo cada 5 segundos
       ↓
Validador recibe 0,000001 SPC de recompensa
       ↓
Esos SPC nuevos entran en circulación
       ↓
En décadas → llegamos a 21.000.000 SPC
```

Esto se llama **emisión controlada** — la cantidad de nuevos SPC que entran al mercado es predecible y pequeña. No hay sorpresas.

---

## Resumen en una frase

> SpainCoin tiene 21 millones de tokens máximo, tú controlas el 40% del génesis, el valor sube cuando más gente quiere usarlo, y el exchange que estamos construyendo es la razón por la que lo querrán usar.

---

## Fases del valor

| Fase | Dónde estamos | Precio estimado |
|------|--------------|----------------|
| Desarrollo | Ahora | Sin precio (no hay mercado) |
| Testnet | Próximamente | Sin precio real |
| Lanzamiento exchange | Fase 5 | Tú fijas el precio inicial |
| Crecimiento | Depende del marketing | El mercado lo decide |
| Mainnet establecida | Fase 6 | El mercado lo decide |

---

---

## La matemática real — Exchange y dinero del fundador

### Estado actual de la blockchain

El nodo está produciendo bloques cada 5 segundos. Cada bloque genera 1 SPC de recompensa (1.000.000.000.000.000 pesetas). El suministro inicial fue de 1.000 SPC repartido en 3 wallets:

| Wallet | Concepto | SPC iniciales |
|--------|----------|---------------|
| Wallet A (validador) | Fundador — la que valida bloques | ~334 SPC |
| Wallet B | Fundador — reserva | ~333 SPC |
| Wallet C | Fundador — reserva | ~333 SPC |

El validador (Wallet A) además acumula +1 SPC por bloque. Si lleva 7.000 bloques, tiene ~7.334 SPC.

### ¿Cuánto dinero real necesitas meter?

**Respuesta corta: 0€ al principio.**

En testnet no hay dinero real. Cuando pases a mainnet:

| Concepto | Coste | Cuándo |
|----------|-------|--------|
| Servidores (2x VPS) | ~16€/mes | Ya lo estás pagando |
| Dominio spaincoin.es | ~8€/año | Ya lo has pagado |
| Constitución SL | ~1.500-2.000€ | Antes de mainnet |
| Registro PSAV | ~600€ tasa | Antes de mainnet |
| KYC/AML (Sumsub) | ~0.50€/verificación | Al empezar a operar |
| **Total para arrancar** | **~2.500-3.000€** | |

### ¿Cómo funciona el exchange por dentro?

```
USUARIO quiere comprar 100 SPC
         │
         ▼
1. El exchange mira el precio actual: 0,09€/SPC
2. Calcula: 100 × 0,09 = 9,00€
3. El usuario paga 9€ (con su balance EUR en la plataforma)
4. El exchange transfiere 100 SPC al wallet del usuario
5. El exchange recibe los 9€
```

### ¿De dónde salen los SPC que compra la gente?

De **tu wallet de liquidez**. Tú pones SPC disponibles en el exchange, y la gente te compra a ti. Cuando alguien compra 100 SPC a 0,09€:

- Tú pierdes 100 SPC
- Tú ganas 9€
- El usuario gana 100 SPC

Cuando alguien vende 100 SPC a 0,09€:

- Tú ganas 100 SPC
- Tú pierdes 9€
- El usuario gana 9€

**Tú eres el "market maker"** — el que pone liquidez en ambos lados.

### Ejemplo con números reales

**Arranque:**
- Pones 200.000 SPC a precio de 0,10€
- Eso representa 20.000€ en "valor" de SPC

**Escenario: 500 personas compran 20€ cada una en el primer mes:**
- Se venden 100.000 SPC (500 × 200 SPC a 0,10€)
- Tú recibes 10.000€ reales
- Te quedan 100.000 SPC de liquidez
- Como hay menos SPC disponibles y más demanda → el precio sube

**El precio ahora es 0,15€:**
- Tus 100.000 SPC restantes valen 15.000€
- Más los 10.000€ que ya recibiste
- Total: 25.000€ (empezaste con 0€ de inversión en SPC)

**¿Suena demasiado bien? Es la misma matemática que Bitcoin en 2010.**
Pero también puede ir al revés: si nadie quiere comprar, tu SPC vale 0.

### El escenario peligroso (sé honesto contigo mismo)

Si el exchange no genera tracción:
- Los SPC no valen nada
- Has gastado ~3.000€ en la SL + servidores
- Has invertido tu tiempo

Es un riesgo real. Por eso el marketing es TU trabajo más importante.

### Resumen para el fundador

```
LO QUE TIENES:                    LO QUE NECESITAS:
─────────────────                  ──────────────────
~1.000 SPC (génesis)               Marketing (tu trabajo)
+ ~7.000 SPC (recompensas)         ~3.000€ para arrancar legal
= ~8.000 SPC total                 Paciencia (6-12 meses)

Si precio = 0,10€ → valen 800€
Si precio = 1€    → valen 8.000€
Si precio = 10€   → valen 80.000€
Si precio = 100€  → valen 800.000€
```

Los 200.000 SPC de liquidez los crearás en la génesis del mainnet. En testnet son simulados.

---

*Última actualización: 2026-03-23 — Exchange en producción, Trading funcional en testnet*
