import { useState, useEffect } from 'react'
import PriceChart from '../components/PriceChart.jsx'
import { getMarketTable, getMarketStats } from '../api/client.js'

const formatEUR = (n) => n >= 1000 ? new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR', maximumFractionDigits: 0 }).format(n) : new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(n)
const formatNum = (n) => n >= 1_000_000_000 ? (n / 1_000_000_000).toFixed(2) + 'B' : n >= 1_000_000 ? (n / 1_000_000).toFixed(2) + 'M' : n >= 1_000 ? (n / 1_000).toFixed(1) + 'K' : n.toFixed(2)

const coinColors = {
  SPC: { bg: 'linear-gradient(135deg, #ffc400, #e6a800)', text: '#c60b1e' },
  BTC: { bg: 'linear-gradient(135deg, #f7931a, #e2820e)', text: '#fff' },
  ETH: { bg: 'linear-gradient(135deg, #627eea, #4a67d6)', text: '#fff' },
  BNB: { bg: 'linear-gradient(135deg, #f3ba2f, #d4a017)', text: '#000' },
  SOL: { bg: 'linear-gradient(135deg, #9945ff, #14f195)', text: '#fff' },
  XRP: { bg: 'linear-gradient(135deg, #23292f, #4a4a4a)', text: '#fff' },
  ADA: { bg: 'linear-gradient(135deg, #0033ad, #0052ff)', text: '#fff' },
  DOGE: { bg: 'linear-gradient(135deg, #c2a633, #ba9f33)', text: '#fff' },
  DOT: { bg: 'linear-gradient(135deg, #e6007a, #c40068)', text: '#fff' },
  AVAX: { bg: 'linear-gradient(135deg, #e84142, #d03031)', text: '#fff' },
  MATIC: { bg: 'linear-gradient(135deg, #8247e5, #6b30d0)', text: '#fff' },
}

// Generate synthetic mini chart data per coin
function generateMiniData(currentPrice, change, seed) {
  const points = 30
  const data = []
  for (let i = 0; i < points; i++) {
    const t = i / points
    const wave = Math.sin(i / 3 + seed * 7) * 0.015 + Math.sin(i / 7 + seed * 3) * 0.01
    const trend = (change / 100) * t
    const price = currentPrice * (1 + wave + trend - (change / 100) * 0.5)
    data.push({ price, height: i })
  }
  return data
}

export default function Market({ onNavigate }) {
  const [tokens, setTokens] = useState([])
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      try {
        const [table, s] = await Promise.all([
          getMarketTable(),
          getMarketStats(),
        ])
        setTokens(table)
        setStats(s)
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
        overflowX: 'auto',
      }}>
        <div style={{ padding: '1rem 1.25rem', borderBottom: '1px solid var(--border)' }}>
          <h2 style={{ fontSize: '1rem', fontWeight: '600', color: 'var(--text-primary)', margin: 0 }}>
            Tokens disponibles
          </h2>
        </div>

        {/* Token rows — mobile-friendly card layout */}
        {tokens.map((t, i) => {
          const changeColor = t.change_24h >= 0 ? 'var(--green)' : 'var(--red)'
          return (
            <div
              key={t.symbol}
              onClick={() => onNavigate(`/trade/${t.symbol}`)}
              style={{
                display: 'flex', alignItems: 'center', justifyContent: 'space-between',
                padding: '0.9rem 1.25rem', gap: '0.75rem',
                borderBottom: i < tokens.length - 1 ? '1px solid var(--border)' : 'none',
                cursor: 'pointer', transition: 'background 0.15s',
              }}
              onMouseEnter={e => e.currentTarget.style.background = 'var(--bg-secondary)'}
              onMouseLeave={e => e.currentTarget.style.background = 'transparent'}
            >
              {/* Left: icon + name */}
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem', minWidth: 0 }}>
                <div style={{
                  width: '36px', height: '36px', borderRadius: '50%',
                  background: (coinColors[t.symbol] || coinColors.SPC).bg,
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  fontWeight: '700', fontSize: '0.7rem',
                  color: (coinColors[t.symbol] || coinColors.SPC).text, flexShrink: 0,
                }}>
                  {t.symbol.slice(0, 1)}
                </div>
                <div style={{ minWidth: 0 }}>
                  <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)' }}>{t.symbol}</div>
                  <div style={{ fontSize: '0.7rem', color: 'var(--text-secondary)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{t.name}</div>
                </div>
              </div>

              {/* Center: mini chart — hidden on small screens via className */}
              <div className="market-chart" style={{ flexShrink: 0 }}>
                <PriceChart data={generateMiniData(t.price, t.change_24h, i)} width={100} height={36} color={changeColor} />
              </div>

              {/* Right: price + change */}
              <div style={{ textAlign: 'right', flexShrink: 0 }}>
                <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)', whiteSpace: 'nowrap' }}>
                  {formatEUR(t.price)}
                </div>
                <div style={{ fontWeight: '600', fontSize: '0.75rem', color: changeColor, whiteSpace: 'nowrap' }}>
                  {t.change_24h >= 0 ? '+' : ''}{t.change_24h}%
                </div>
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
          <li>Supply máximo: 21.000.000 SPC</li>
        </ul>
      </div>

      <div style={{ textAlign: 'center', marginTop: '1rem', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>
        Testnet — datos simulados, sin valor real
      </div>
    </div>
  )
}
