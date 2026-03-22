# Requisitos Legales — SpainCoin Exchange

> Guía completa de lo que necesitas para operar legalmente en España y la UE.
> Actualizada: marzo 2026.

---

## Estado actual: TESTNET

En fase testnet no hay dinero real. Las obligaciones legales son mínimas:
- ✅ Términos y Condiciones publicados
- ✅ Política de Privacidad (RGPD) publicada
- ✅ Aviso de Riesgos publicado
- ✅ Política de Cookies + banner de consentimiento
- ✅ Checkboxes de consentimiento en el registro

---

## Antes de lanzar MAINNET — Checklist completo

### 1. Estructura societaria

**Qué necesitas:** Una sociedad mercantil española (SL o SA).

| Opción | Coste aprox. | Tiempo |
|--------|-------------|--------|
| SL (Sociedad Limitada) | 3.000€ (capital mínimo) + ~1.000€ notaría/registro | 2-4 semanas |
| SL online (CIRCE) | 3.000€ + ~300€ | 48-72 horas |

**Por qué:** No puedes operar un exchange como persona física. La SL limita tu responsabilidad personal.

**Pasos:**
1. Elegir nombre → consulta en Registro Mercantil Central (rmercantil.es)
2. Abrir cuenta bancaria de empresa (solo para el capital)
3. Escritura notarial
4. Registro en el Registro Mercantil
5. Alta en Hacienda (IAE, epígrafe 651.9 o similar)

---

### 2. Registro en el Banco de España (PSAV)

**Qué es:** Desde la Ley 11/2021, todos los exchanges de criptomonedas en España deben registrarse como **Proveedor de Servicios de Activos Virtuales (PSAV)** ante el Banco de España.

**Sin este registro no puedes:**
- Ofrecer compra/venta de criptomonedas con dinero real (€)
- Publicitar el servicio en España
- Operar legalmente

**Documentación requerida:**
- CIF de la sociedad
- Estatutos sociales
- Plan de negocio detallado
- Programa de actividades
- Procedimientos AML/KYC documentados
- Identidad y curriculum de los administradores
- Medidas de ciberseguridad implementadas

**Coste:** Gratuito el registro, pero necesitarás asesor legal (~2.000-5.000€)

**Tiempo:** 3-6 meses de tramitación

**Referencia:** Circular 1/2022 del Banco de España

---

### 3. Cumplimiento MiCA (EU)

**Qué es:** El **Reglamento de Mercados en Criptoactivos (MiCA)** es la nueva ley europea que regula los exchanges y emisores de criptomonedas. Aplicable desde diciembre 2024.

**Lo que afecta a SpainCoin:**

| Tipo | Aplica | Qué requiere |
|------|--------|-------------|
| Exchange (compra/venta) | ✅ Sí | Licencia CASP (Crypto Asset Service Provider) |
| Emisión de $SPC | ✅ Sí | Whitepaper publicado + notificación regulador |
| Staking | ✅ Sí | Información clara de riesgos |

**Whitepaper obligatorio** — debe incluir:
- Descripción del proyecto y tecnología
- Distribución del supply (tokenomics)
- Riesgos
- Derechos de los holders
- Información de los fundadores

**Nota:** Al ser un token de utilidad (no stablecoin ni asset-referenced), los requisitos son más ligeros.

---

### 4. KYC/AML (Conoce a tu Cliente / Anti-Blanqueo)

**Qué es:** Obligatorio verificar la identidad de los usuarios antes de permitirles operar con dinero real.

**Niveles de verificación:**

| Nivel | Límite operación | Documentos |
|-------|-----------------|------------|
| Básico | Hasta 1.000€/mes | Email verificado |
| KYC 1 | Hasta 15.000€/mes | DNI/Pasaporte + selfie |
| KYC 2 | Sin límite | DNI + justificante domicilio + origen fondos |

**Proveedores de KYC recomendados:**

| Proveedor | Precio aprox. | Notas |
|-----------|--------------|-------|
| **Sumsub** | 0,50-1€/verificación | El más usado en crypto |
| **Veriff** | 1-2€/verificación | Muy buena UX |
| **Onfido** | 1-3€/verificación | Muy preciso |

**Implementación:** API REST que se integra con el backend de la exchange.

---

### 5. Cuenta bancaria de empresa

**El problema:** La mayoría de bancos tradicionales rechazan empresas cripto.

**Opciones que funcionan:**

| Banco/Neobanco | Para cripto | Notas |
|---------------|-------------|-------|
| **Revolut Business** | ✅ Sí | El más fácil, apertura online |
| **BBVA** | ⚠️ A veces | Depende del gestor |
| **Bankinter** | ⚠️ A veces | Mejor que Santander/CaixaBank |
| **Wise Business** | ✅ Sí | Para cobros internacionales |
| **Santander/CaixaBank** | ❌ No | Suelen cerrar cuentas cripto |

**Recomendación:** Revolut Business como principal + Wise para cobros internacionales.

---

### 6. Fiscalidad

**Para la empresa (SL):**
- Impuesto de Sociedades: 25% sobre beneficios (15% los 2 primeros años)
- IVA: los servicios de exchange están exentos de IVA (art. 20 LIVA)
- Los ingresos por fees del exchange son rendimientos de actividad

**Para los usuarios:**
- Las ganancias con criptos tributan como **ganancias patrimoniales** en el IRPF (19-28%)
- El exchange debe informar a Hacienda de operaciones (Modelo 721 desde 2024)
- Obligación de enviar datos a AEAT si el usuario opera más de 3.000€/año

**Lo que debes implementar:**
- Historial de transacciones exportable (para que usuarios hagan su declaración)
- Registro de operaciones (para informar a Hacienda si se requiere)

---

### 7. Protección de datos (RGPD)

**Registro de actividades de tratamiento:**
- Obligatorio tener un registro interno de qué datos tratas y por qué
- No necesitas notificar a la AEPD si no tratas datos "especialmente sensibles"

**DPD (Delegado de Protección de Datos):**
- No es obligatorio para empresas pequeñas (< 250 empleados) salvo que el tratamiento sea a gran escala

**Transferencias internacionales:**
- Si usas servidores fuera de la UE (AWS us-east, etc.) debes indicarlo en la política de privacidad
- Hetzner (Alemania/Finlandia) está dentro de la UE ✅

**Brecha de seguridad:**
- Si sufres un hackeo con datos personales → debes notificar a la AEPD en 72 horas
- Si afecta a usuarios → debes notificarles también

---

### 8. Seguros

| Seguro | Obligatorio | Recomendado |
|--------|------------|-------------|
| Responsabilidad Civil Profesional | No | Sí — cubre reclamaciones de usuarios |
| Ciberriesgo | No | Sí — cubre hackeos y brechas de datos |
| D&O (Administradores) | No | Para cuando el proyecto crezca |

---

## Timeline recomendado

```
AHORA (testnet)
├── ✅ Páginas legales en la web
├── ✅ Aviso de riesgos
└── ✅ Política de privacidad RGPD

3-6 MESES ANTES DE MAINNET
├── Constituir SL
├── Abrir cuenta Revolut Business
├── Contratar asesor legal especializado en cripto
├── Redactar Whitepaper (MiCA)
└── Iniciar proceso registro Banco de España

AL LANZAR MAINNET
├── Registro PSAV activo
├── KYC/AML integrado
├── Modelo 721 preparado
└── Seguro RC contratado
```

---

## Contactos y recursos útiles

| Recurso | URL |
|---------|-----|
| Registro PSAV Banco de España | bde.es |
| AEPD (privacidad) | aepd.es |
| MiCA texto completo | eur-lex.europa.eu |
| Registro Mercantil Central | rmercantil.es |
| Consultas tributarias AEAT | aeat.es |

---

## Coste legal estimado hasta mainnet

| Concepto | Coste estimado |
|----------|---------------|
| Constitución SL | 1.500-2.000€ |
| Asesor legal registro PSAV | 2.000-5.000€ |
| Whitepaper legal (MiCA) | 1.000-3.000€ |
| KYC provider (setup) | 500-1.000€ |
| Seguro RC (anual) | 1.000-2.000€ |
| **Total estimado** | **6.000-13.000€** |

---

*Este documento es orientativo y no constituye asesoramiento jurídico. Consulta con un abogado especializado en derecho digital y fintech antes de lanzar.*
