# Guía del Fundador — SpainCoin ($SPC)

> Para Luis. Sin tecnicismos innecesarios. Lo que necesitas entender, lo que debes proteger, y lo que nadie más puede saber.

---

## 1. Cómo funciona SpainCoin por dentro

### El bloque — la unidad básica

Imagina un bloque como una página de un libro contable. Cada página contiene:
- **Quién envió dinero a quién** (transacciones)
- **La firma de la página anterior** (así nadie puede cambiar el pasado)
- **La firma del validador** que "selló" esa página

Si alguien intenta cambiar una transacción del pasado, la firma de esa página cambia, lo que rompe la cadena entera. Por eso la blockchain es inmutable.

```
Bloque 1  →  Bloque 2  →  Bloque 3  →  ...
(genesis)     (txs)         (txs)
```

### Las claves — tu identidad en la red

Cada wallet tiene dos claves:

| Clave | Qué es | ¿Se comparte? |
|-------|--------|---------------|
| **Clave Pública** | Tu dirección `SPC...` | Sí — es tu "número de cuenta" |
| **Clave Privada** | La contraseña de tu cuenta | **NUNCA** — quien la tenga, tiene tus fondos |

Funciona como un buzón: cualquiera puede meter cartas (enviarte $SPC), pero solo tú puedes abrirlo (gastar los fondos) con tu clave privada.

### Las firmas — por qué nadie puede robarte

Cuando envías $SPC, tu wallet:
1. Crea la transacción (de: ti, para: X, cantidad: Y)
2. La **firma** con tu clave privada (como firmar un cheque)
3. Envía la transacción firmada a la red

La red verifica que la firma es tuya usando tu clave pública. Si no coincide → transacción rechazada. Nadie puede gastar tus fondos sin tu clave privada. Es matemáticamente imposible falsificar la firma.

### El Merkle Tree — por qué no se pueden colar transacciones falsas

Si un bloque tiene 1.000 transacciones, verificar todas sería lento. El Merkle Tree es un sistema de verificación en árbol:

```
        [Root Hash]
       /            \
   [Hash A-B]    [Hash C-D]
   /      \      /       \
[Hash A] [Hash B] [Hash C] [Hash D]
  tx1      tx2     tx3      tx4
```

Si cambias cualquier transacción, su hash cambia → el hash padre cambia → el Root Hash cambia → el bloque es inválido. Con solo el Root Hash puedes verificar que ninguna de las 1.000 transacciones fue alterada.

### El Estado — quién tiene qué

El "estado" es el balance actual de todas las cuentas. Cuando se aplica un bloque:
1. Se verifican todas las transacciones
2. Se actualizan los balances
3. Se actualiza el State Root (hash del estado completo)

El State Root en el bloque garantiza que todos los nodos tienen exactamente el mismo estado.

### El Mempool — la sala de espera

Cuando envías una transacción, no va directa a un bloque. Primero entra al **mempool** (memory pool) — una cola de espera en cada nodo. El validador selecciona las transacciones del mempool y las incluye en el siguiente bloque, priorizando las que pagan mayor **fee**.

---

## 2. Consenso Proof of Stake — quién crea los bloques

### El problema que resuelve

En una red descentralizada, ¿quién decide cuál es el siguiente bloque válido? Necesitas un sistema donde:
- Nadie pueda controlar la red solo
- Los que intentan hacer trampa pierdan dinero
- Sea sostenible (no queme electricidad como Bitcoin)

### Cómo funciona nuestro PoS

1. Los **validadores** bloquean ("hacen stake") una cantidad de $SPC como garantía
2. Para cada bloque, se selecciona un validador pseudo-aleatoriamente, **con más probabilidad cuanto más stake tenga**
3. El validador propone el bloque, los demás lo verifican
4. Si el validador intentó hacer trampa → pierde su stake (**slashing**)
5. Si todo va bien → el validador gana la **recompensa de bloque**

```
Validador A: 1.000 SPC stake → 10% probabilidad de ser elegido
Validador B: 4.000 SPC stake → 40% probabilidad
Validador C: 5.000 SPC stake → 50% probabilidad
```

### Por qué es seguro

Para atacar la red necesitarías el 51% del stake total. Eso significa comprar la mayoría de los $SPC en circulación. El ataque costaría más de lo que ganarías, y además derrumbaría el valor de los $SPC que tienes.

---

## 3. Tokenomics — los números de $SPC

| Concepto | Valor | Por qué |
|----------|-------|---------|
| **Supply máximo** | 21.000.000 SPC | Escasez (como Bitcoin) |
| **Pre-minado génesis** | 1.000.000 SPC | Fundadores + desarrollo + liquidez inicial |
| **Unidad mínima** | 1 "peseta" = 0.000000000000000001 SPC | Precisión en micropagos |
| **Recompensa por bloque** | A definir en Fase 1 | Incentivo para validadores |

### Distribución del génesis (1.000.000 SPC)
Esto lo decides tú. Recomendación típica:
- 40% — Fundadores (bloqueado 2 años, liberación gradual)
- 30% — Desarrollo y operaciones
- 20% — Liquidez inicial en el exchange
- 10% — Comunidad / marketing

**Esto es importante**: la distribución del génesis es una decisión que tienes que tomar tú. Una vez generado el bloque génesis, no se puede cambiar.

---

## 4. Lo que SOLO tú puedes saber

### Las claves privadas del génesis

Cuando generemos el bloque génesis, se crearán las wallets fundadoras. Las **claves privadas** de esas wallets son lo más crítico del proyecto:

- **Si las pierdes**: los fondos del génesis desaparecen para siempre. Nadie puede recuperarlos. Ni yo.
- **Si alguien las roba**: puede vaciar todas las wallets fundadoras al instante.
- **No existe un "recuperar contraseña"**: la blockchain no tiene soporte técnico. Las claves son las claves.

### Cómo protegerlas (cuando llegue el momento)

1. **Escríbelas en papel** — sí, papel físico. Dos copias en sitios distintos.
2. **Hardware wallet** — dispositivo físico (Ledger/Trezor) para los fondos grandes
3. **Nunca en email, cloud, fotos, WhatsApp, Notion, etc.**
4. **Seed phrase** — la clave privada se puede representar como 12-24 palabras (mnemónico). Es lo mismo, igual de peligroso.

### El nodo validador principal

Cuando lancemos la testnet, el nodo principal necesitará:
- La clave privada del validador fundador
- Estar en un servidor seguro (no tu portátil personal)
- Acceso SSH solo desde tu IP

Esto lo configuraremos cuando lleguemos a esa fase.

---

## 5. Qué es peligroso y qué no

### PELIGROSO ⚠️

| Situación | Consecuencia |
|-----------|-------------|
| Perder la clave privada del génesis | Fondos fundadores perdidos para siempre |
| Exponer clave privada en GitHub/código | Fondos vaciados en minutos |
| Lanzar mainnet sin auditoría de seguridad | Exploit → pérdida de fondos de usuarios |
| Prometer rendimientos fijos | Problema legal en muchos países |
| Supply inflacionario sin control | Devaluación de $SPC |

### NECESARIO ✅

| Qué | Por qué |
|-----|---------|
| Código abierto (open source) | Los usuarios necesitan verificar que no hay trampa |
| Auditoría externa antes de mainnet | Encontrar bugs antes que los atacantes |
| Plan de contingencia para bugs críticos | Un exploit en mainnet puede matar el proyecto |
| Whitepaper público | Credibilidad y transparencia |
| Multi-sig para fondos del proyecto | Que no dependa solo de una persona/clave |

### NO CRÍTICO (pero importante)

- El dominio del exchange
- El diseño de la app
- Los contratos de marketing

---

## 6. Arquitectura del sistema — vista general

```
                    ┌─────────────┐
                    │   APP Web   │  ← React (lo que ve el usuario)
                    │  Exchange   │
                    └──────┬──────┘
                           │ API REST / WebSocket
                    ┌──────┴──────┐
                    │  Go API     │  ← Capa intermedia
                    └──────┬──────┘
                           │ RPC
              ┌────────────┴────────────┐
              │        NODO SPC         │
              │  ┌─────────────────┐    │
              │  │   Blockchain    │    │
              │  │  core/chain/    │    │
              │  ├─────────────────┤    │
              │  │   Consenso PoS  │    │
              │  │  core/consensus/│    │
              │  ├─────────────────┤    │
              │  │   Mempool       │    │
              │  │  core/mempool/  │    │
              │  ├─────────────────┤    │
              │  │   Estado        │    │
              │  │  core/state/    │    │
              │  └─────────────────┘    │
              └────────────┬────────────┘
                           │ P2P (libp2p)
              ┌────────────┴────────────┐
         [Nodo 2]                   [Nodo 3]
              └────────────┬────────────┘
                      [Nodo 4...]
```

---

## 7. Fases y qué significa cada una

| Fase | Qué construimos | Lo que consigues |
|------|----------------|-----------------|
| **1 - Core** | Bloques, txs, consenso, estado | Blockchain funcional en local |
| **2 - P2P** | Múltiples nodos comunicándose | Red descentralizada real |
| **3 - Wallet + CLI** | Herramientas de usuario | Puedes crear wallets y enviar $SPC |
| **4 - Testnet** | Red pública de pruebas | Pruebas reales con comunidad |
| **5 - Exchange app** | App tipo crypto.com | Usuarios pueden comprar/vender $SPC |
| **6 - Mainnet** | Red principal con $SPC real | El proyecto está live |

**Dónde estamos ahora:** Fase 1, construyendo los módulos del core.

---

*Última actualización: Fase 1 en curso — core/crypto ✅ core/block ✅*
