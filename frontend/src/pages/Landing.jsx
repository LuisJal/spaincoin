import { useState, useEffect } from 'react'
import { getStatus, getPrice } from '../api/client.js'

export default function Landing({ onNavigate }) {
  const [status, setStatus] = useState(null)
  const [price, setPrice] = useState(null)

  useEffect(() => {
    async function load() {
      try {
        const [s, p] = await Promise.all([getStatus(), getPrice()])
        setStatus(s)
        setPrice(p)
      } catch (e) { console.error(e) }
    }
    load()
    const i = setInterval(load, 15000)
    return () => clearInterval(i)
  }, [])

  const blockHeight = status?.node?.height || 0
  const supply = status?.node?.total_supply ? (status.node.total_supply / 1_000_000_000_000_000).toFixed(2) : '—'

  return (
    <div className="page-enter">
      {/* Hero */}
      <div style={{
        textAlign: 'center', padding: '4rem 1.5rem 3rem',
        background: 'linear-gradient(180deg, rgba(255,196,0,0.06) 0%, transparent 60%)',
      }}>
        <div style={{ fontSize: '3rem', fontWeight: '800', color: 'var(--text-primary)', lineHeight: 1.1, marginBottom: '1rem' }}>
          La blockchain <span style={{ color: '#ffc400' }}>de Espana</span>
        </div>
        <p style={{ fontSize: '1.15rem', color: 'var(--text-secondary)', maxWidth: '600px', margin: '0 auto 2rem', lineHeight: 1.6 }}>
          SpainCoin es una blockchain Layer 1 con consenso Proof of Stake.
          Codigo abierto, descentralizada, construida desde cero en Go.
        </p>
        <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'center', flexWrap: 'wrap' }}>
          <button onClick={() => onNavigate('/wallet')} style={{
            padding: '0.85rem 2rem', background: 'linear-gradient(135deg, #ffc400, #e6a800)',
            border: 'none', borderRadius: '10px', color: '#000', fontSize: '1rem',
            fontWeight: '700', cursor: 'pointer',
          }}>Descargar Wallet</button>
          <button onClick={() => onNavigate('/whitepaper')} style={{
            padding: '0.85rem 2rem', background: 'transparent',
            border: '1px solid rgba(255,255,255,0.2)', borderRadius: '10px',
            color: '#f9fafb', fontSize: '1rem', fontWeight: '500', cursor: 'pointer',
          }}>White Paper</button>
        </div>
      </div>

      {/* Stats */}
      <div style={{
        display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(160px, 1fr))',
        gap: '0.75rem', maxWidth: '800px', margin: '0 auto 3rem', padding: '0 1.5rem',
      }}>
        {[
          { label: 'Precio SPC', value: price ? `${price.price_eur.toFixed(4)} EUR` : '—' },
          { label: 'Bloque', value: blockHeight ? `#${blockHeight.toLocaleString('es-ES')}` : '—' },
          { label: 'Supply', value: `${supply} SPC` },
          { label: 'Max Supply', value: '21M SPC' },
        ].map((s, i) => (
          <div key={i} style={{
            background: 'var(--bg-card)', borderRadius: '10px', border: '1px solid var(--border)',
            padding: '1rem', textAlign: 'center',
          }}>
            <div style={{ fontSize: '0.7rem', color: 'var(--text-secondary)', textTransform: 'uppercase', marginBottom: '0.3rem' }}>{s.label}</div>
            <div style={{ fontSize: '1.1rem', fontWeight: '700', color: 'var(--text-primary)' }}>{s.value}</div>
          </div>
        ))}
      </div>

      {/* Features */}
      <div style={{ maxWidth: '800px', margin: '0 auto 3rem', padding: '0 1.5rem' }}>
        <h2 style={{ fontSize: '1.4rem', fontWeight: '700', color: 'var(--text-primary)', textAlign: 'center', marginBottom: '2rem' }}>
          Por que SpainCoin
        </h2>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '1rem' }}>
          {[
            { title: 'Proof of Stake', desc: 'Consenso eficiente sin mineria. Valida bloques con tu stake y gana recompensas.' },
            { title: 'Self-Custody', desc: 'Tus claves, tus fondos. Nunca almacenamos claves privadas. Tu wallet, tu control.' },
            { title: 'Codigo Abierto', desc: 'Todo el codigo es publico y auditable. Transparencia total, sin puertas traseras.' },
            { title: 'Red P2P', desc: 'Nodos conectados con libp2p. Descubrimiento automatico. Resistente a censura.' },
            { title: 'Supply Limitado', desc: '21 millones de SPC maximo. Escasez programada, como Bitcoin.' },
            { title: 'Hecha en Espana', desc: 'Disenada y construida desde cero. La primera blockchain espanola.' },
          ].map((f, i) => (
            <div key={i} style={{
              background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
              padding: '1.25rem',
            }}>
              <div style={{ fontWeight: '600', fontSize: '0.95rem', color: 'var(--text-primary)', marginBottom: '0.4rem' }}>{f.title}</div>
              <div style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', lineHeight: 1.6 }}>{f.desc}</div>
            </div>
          ))}
        </div>
      </div>

      {/* CTA */}
      <div style={{
        textAlign: 'center', padding: '3rem 1.5rem',
        background: 'var(--bg-card)', borderTop: '1px solid var(--border)',
      }}>
        <h2 style={{ fontSize: '1.3rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
          Unete a la red
        </h2>
        <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '1.5rem', maxWidth: '500px', margin: '0 auto 1.5rem' }}>
          Descarga el wallet, corre un nodo validador y forma parte de la primera blockchain espanola.
        </p>
        <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'center', flexWrap: 'wrap' }}>
          <button onClick={() => onNavigate('/validators')} style={{
            padding: '0.75rem 1.5rem', background: 'var(--accent)', border: 'none',
            borderRadius: '8px', color: '#fff', fontSize: '0.9rem', fontWeight: '600', cursor: 'pointer',
          }}>Ser validador</button>
          <button onClick={() => onNavigate('/explorer')} style={{
            padding: '0.75rem 1.5rem', background: 'transparent',
            border: '1px solid var(--border)', borderRadius: '8px',
            color: 'var(--text-secondary)', fontSize: '0.9rem', cursor: 'pointer',
          }}>Ver explorer</button>
        </div>
      </div>
    </div>
  )
}
