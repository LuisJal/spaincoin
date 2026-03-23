import { useState, useEffect } from 'react'
import PriceChart from '../components/PriceChart.jsx'
import { getMarketTable, getPriceHistory, getMarketStats } from '../api/client.js'

const formatEUR = (n) => new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(n)
const formatNum = (n) => n >= 1_000_000 ? (n / 1_000_000).toFixed(2) + 'M' : n >= 1_000 ? (n / 1_000).toFixed(1) + 'K' : n.toFixed(2)

export default function Market({ onNavigate }) {
  const [tokens, setTokens] = useState([])
  const [stats, setStats] = useState(null)
  const [miniChart, setMiniChart] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      try {
        const [table, s, chart] = await Promise.all([
          getMarketTable(),
          getMarketStats(),
          getPriceHistory(50, '24h'),
        ])
        setTokens(table)
        setStats(s)
        setMiniChart(chart)
      } catch (e) {
        console.error('Market load error:', e)
      } finally {
        setLoading(false)
      }
    }
    load()
    const interval = setInterval(load, 15000)
    return () => clearInterval(interval)
  }, [])

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '4rem' }}>
        <div className="spinner" />
      </div>
    )
  }

  return (
    <div className="page-enter" style={{ maxWidth: '1000px', margin: '0 auto', padding: '1.5rem 1rem' }}>

      {/* Header */}
      <div style={{ marginBottom: '2rem' }}>
        <h1 style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.25rem' }}>
          Mercado
        </h1>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
          Precios en tiempo real del ecosistema SpainCoin
        </p>
      </div>

      {/* Stats cards */}
      {stats && (
        <div style={{
          display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))',
          gap: '0.75rem', marginBottom: '2rem',
        }}>
          {[
            { label: 'Market Cap', value: formatEUR(stats.market_cap) },
            { label: 'Supply', value: `${formatNum(stats.circulating_supply)} / 21M SPC` },
            { label: 'Bloque', value: `#${stats.block_height?.toLocaleString('es-ES')}` },
            { label: 'Peers', value: stats.peer_count },
          ].map((s, i) => (
            <div key={i} style={{
              background: 'var(--bg-card)', borderRadius: '10px', border: '1px solid var(--border)',
              padding: '1rem',
            }}>
              <div style={{ fontSize: '0.7rem', color: 'var(--text-secondary)', marginBottom: '0.3rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
                {s.label}
              </div>
              <div style={{ fontSize: '1.1rem', fontWeight: '700', color: 'var(--text-primary)' }}>
                {s.value}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Token table */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        overflow: 'hidden',
      }}>
        <div style={{ padding: '1rem 1.25rem', borderBottom: '1px solid var(--border)' }}>
          <h2 style={{ fontSize: '1rem', fontWeight: '600', color: 'var(--text-primary)', margin: 0 }}>
            Tokens disponibles
          </h2>
        </div>

        {/* Table header */}
        <div style={{
          display: 'grid', gridTemplateColumns: '2fr 1.5fr 1fr 1.5fr 2fr 1fr',
          padding: '0.6rem 1.25rem', fontSize: '0.7rem', color: 'var(--text-secondary)',
          textTransform: 'uppercase', letterSpacing: '0.05em', borderBottom: '1px solid var(--border)',
        }}>
          <span>Token</span>
          <span style={{ textAlign: 'right' }}>Precio</span>
          <span style={{ textAlign: 'right' }}>24h</span>
          <span style={{ textAlign: 'right' }}>Volumen</span>
          <span style={{ textAlign: 'center' }}>Grafico 24h</span>
          <span style={{ textAlign: 'right' }}>Operar</span>
        </div>

        {/* Token rows */}
        {tokens.map((t, i) => {
          const changeColor = t.change_24h >= 0 ? 'var(--green)' : 'var(--red)'
          return (
            <div
              key={t.symbol}
              style={{
                display: 'grid', gridTemplateColumns: '2fr 1.5fr 1fr 1.5fr 2fr 1fr',
                padding: '1rem 1.25rem', alignItems: 'center',
                borderBottom: i < tokens.length - 1 ? '1px solid var(--border)' : 'none',
                cursor: 'pointer', transition: 'background 0.15s',
              }}
              onMouseEnter={e => e.currentTarget.style.background = 'var(--bg-secondary)'}
              onMouseLeave={e => e.currentTarget.style.background = 'transparent'}
            >
              {/* Token name */}
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem' }}>
                <div style={{
                  width: '36px', height: '36px', borderRadius: '50%',
                  background: 'linear-gradient(135deg, #ffc400, #e6a800)',
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  fontWeight: '700', fontSize: '0.7rem', color: '#c60b1e', flexShrink: 0,
                }}>
                  S
                </div>
                <div>
                  <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)' }}>{t.symbol}</div>
                  <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)' }}>{t.name}</div>
                </div>
              </div>

              {/* Price */}
              <div style={{ textAlign: 'right', fontWeight: '600', fontSize: '0.95rem', color: 'var(--text-primary)' }}>
                {formatEUR(t.price)}
              </div>

              {/* 24h change */}
              <div style={{ textAlign: 'right', fontWeight: '600', fontSize: '0.85rem', color: changeColor }}>
                {t.change_24h >= 0 ? '+' : ''}{t.change_24h}%
              </div>

              {/* Volume */}
              <div style={{ textAlign: 'right', fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
                {formatEUR(t.volume)}
              </div>

              {/* Mini chart */}
              <div style={{ display: 'flex', justifyContent: 'center' }}>
                <PriceChart data={miniChart} width={120} height={40} color={changeColor} />
              </div>

              {/* Trade button */}
              <div style={{ textAlign: 'right' }}>
                <button
                  onClick={(e) => { e.stopPropagation(); onNavigate('/trade') }}
                  style={{
                    padding: '0.35rem 0.75rem', borderRadius: '6px', border: 'none',
                    background: 'var(--accent)', color: '#fff', fontSize: '0.75rem',
                    fontWeight: '600', cursor: 'pointer',
                  }}
                >
                  Operar
                </button>
              </div>
            </div>
          )
        })}
      </div>

      {/* Info */}
      <div style={{
        marginTop: '2rem', padding: '1.25rem', borderRadius: '12px',
        background: 'var(--bg-card)', border: '1px solid var(--border)',
      }}>
        <h3 style={{ fontSize: '0.9rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
          Sobre el mercado SpainCoin
        </h3>
        <ul style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', lineHeight: 1.7, paddingLeft: '1.25rem' }}>
          <li>Los precios se actualizan con cada bloque (cada 5 segundos)</li>
          <li>El par principal es <strong style={{ color: 'var(--text-primary)' }}>SPC/EUR</strong></li>
          <li>Cada usuario nuevo recibe <strong style={{ color: 'var(--green)' }}>1.000€ virtuales</strong> para practicar</li>
          <li>Supply maximo: 21.000.000 SPC</li>
        </ul>
      </div>

      <div style={{ textAlign: 'center', marginTop: '1rem', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>
        Testnet — datos simulados, sin valor real
      </div>
    </div>
  )
}
