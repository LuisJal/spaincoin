export default function WhitePaper({ onNavigate }) {
  return (
    <div className="page-enter" style={{ maxWidth: '800px', margin: '0 auto', padding: '2rem 1.5rem' }}>
      <h1 style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
        White Paper — SpainCoin ($SPC)
      </h1>
      <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '2rem' }}>
        Version 1.0 — Marzo 2026
      </p>

      {[
        {
          title: '1. Introduccion',
          content: `SpainCoin es una blockchain Layer 1 construida desde cero en Go con consenso Proof of Stake (PoS). Su objetivo es crear una red descentralizada, eficiente y accesible, con foco en la comunidad hispanohablante.

La moneda nativa $SPC tiene un supply maximo de 21.000.000 unidades — escasez programada, inspirada en Bitcoin. SpainCoin no depende de ninguna otra blockchain: tiene su propia red P2P, su propio consenso y su propio explorador de bloques.`
        },
        {
          title: '2. Arquitectura',
          content: `La blockchain esta escrita integramente en Go, utilizando:

• Consenso: Proof of Stake con seleccion ponderada de validadores
• Bloques: cada 5 segundos, con Merkle Tree para integridad
• Criptografia: ECDSA P256, SHA-256, direcciones con prefijo SPC
• Red: libp2p con gossipsub para propagacion de bloques y transacciones
• Almacenamiento: BoltDB (puro Go, sin dependencias C)
• Wallet: CLI self-custody, claves generadas localmente`
        },
        {
          title: '3. Tokenomics',
          content: `Supply maximo: 21.000.000 SPC
Decimales: 18 (1 SPC = 10^18 pesetas)
Genesis: 1.000 SPC (distribucion inicial)
Recompensa por bloque: 1 SPC

Distribucion del genesis:
• Fundadores y desarrollo: reserva para sostenibilidad del proyecto
• Liquidez futura: reservada para cuando se lance el exchange descentralizado
• Comunidad: recompensas de validacion

La emision es predecible y decreciente. No existe mecanismo para crear SPC fuera del protocolo de consenso. El supply maximo es inmutable.`
        },
        {
          title: '4. Consenso Proof of Stake',
          content: `SpainCoin utiliza un sistema de Proof of Stake simplificado:

1. Los validadores depositan SPC como stake
2. Cada bloque, un validador es seleccionado aleatoriamente (ponderado por stake)
3. El validador seleccionado produce el bloque y recibe la recompensa
4. Si un validador actua maliciosamente (double-sign), pierde el 50% de su stake (slashing)
5. Si un validador esta offline, pierde el 1% de su stake

Este modelo es energeticamente eficiente (a diferencia del Proof of Work) y alinea los incentivos de los validadores con la salud de la red.`
        },
        {
          title: '5. Red P2P',
          content: `La red utiliza libp2p para comunicacion entre nodos:

• Protocolo de descubrimiento: mDNS (red local) + bootstrap nodes
• Propagacion: gossipsub para bloques y transacciones
• Puertos: 30303 (P2P), 8545 (RPC, solo para servicios autorizados)

Cualquier persona puede correr un nodo y unirse a la red. No se requiere permiso.`
        },
        {
          title: '6. Self-Custody',
          content: `SpainCoin es 100% self-custody. Esto significa:

• Las claves privadas se generan y almacenan en el dispositivo del usuario
• Ningun servidor almacena claves privadas
• Las transacciones se firman localmente antes de enviarse a la red
• No existe mecanismo de recuperacion de claves — el usuario es responsable de su seguridad

Esta filosofia garantiza que ningun tercero (incluido el equipo de SpainCoin) puede acceder a los fondos de los usuarios.`
        },
        {
          title: '7. Roadmap',
          content: `Fase 1 ✅ — Core blockchain (bloques, transacciones, consenso)
Fase 2 ✅ — Red P2P (multiples nodos, gossipsub)
Fase 3 ✅ — Wallet CLI + persistencia
Fase 4 ✅ — Testnet en produccion (nodo 24/7)
Fase 5 ✅ — Web informativa + explorer
Fase 6 — DEX (exchange descentralizado en la propia red)
Fase 7 — Smart contracts / tokens sobre SpainCoin
Fase 8 — Mainnet con comunidad de validadores`
        },
        {
          title: '8. Riesgos',
          content: `• SPC no tiene valor garantizado. Su precio depende de la oferta y la demanda
• La tecnologia es experimental. Pueden existir vulnerabilidades no descubiertas
• Invertir en criptomonedas conlleva riesgo de perdida total del capital
• SpainCoin no es un producto financiero regulado
• Este documento no constituye una oferta de inversion ni asesoramiento financiero`
        },
        {
          title: '9. Contacto',
          content: `Web: spaincoin.es
GitHub: github.com/spaincoin/spaincoin
Telegram: t.me/spaincoin`
        },
      ].map((section, i) => (
        <div key={i} style={{
          background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
          padding: '1.5rem', marginBottom: '1rem',
        }}>
          <h2 style={{ fontSize: '1.05rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
            {section.title}
          </h2>
          <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.75, whiteSpace: 'pre-line' }}>
            {section.content}
          </p>
        </div>
      ))}
    </div>
  )
}
