import { useState, useEffect, useCallback } from 'react'
import { useAuth } from '../auth/useAuth.jsx'
import PriceChart from '../components/PriceChart.jsx'
import { getTicker, getPriceHistory, buySPC, sellSPC, getTradeBalance, getTradeHistory } from '../api/client.js'

const formatEUR = (n) => new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(n)
const formatSPC = (n) => n >= 1 ? n.toLocaleString('es-ES', { maximumFractionDigits: 4 }) : n.toFixed(6)

export default function Trade({ onNavigate }) {
  const { user, token } = useAuth()
  const [ticker, setTicker] = useState(null)
  const [chartData, setChartData] = useState([])
  const [chartRange, setChartRange] = useState('24h')
  const [tab, setTab] = useState('buy') // 'buy' | 'sell'
  const [amount, setAmount] = useState('')
  const [balance, setBalance] = useState({ eur: 0, spc: 0 })
  const [trades, setTrades] = useState([])
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)

  const loadData = useCallback(async () => {
    try {
      const t = await getTicker()
      setTicker(t)
      const h = await getPriceHistory(120, chartRange)
      setChartData(h)
    } catch (e) {
      console.error('Failed to load ticker:', e)
    }
  }, [chartRange])

  const loadUserData = useCallback(async () => {
    if (!token) return
    try {
      const [b, t] = await Promise.all([
        getTradeBalance(token),
        getTradeHistory(token),
      ])
      setBalance(b)
      setTrades(t)
    } catch (e) {
      console.error('Failed to load user data:', e)
    }
  }, [token])

  useEffect(() => {
    loadData()
    const interval = setInterval(loadData, 10000)
    return () => clearInterval(interval)
  }, [loadData])

  useEffect(() => {
    loadUserData()
  }, [loadUserData])

  const amountNum = parseFloat(amount) || 0
  const price = ticker?.price || 0
  const totalEUR = tab === 'buy' ? amountNum * price : amountNum * price

  async function handleTrade() {
    if (amountNum <= 0) return
    setLoading(true)
    setError(null)
    setResult(null)

    try {
      let res
      if (tab === 'buy') {
        res = await buySPC(token, amountNum)
      } else {
        res = await sellSPC(token, amountNum)
      }
      setResult(res)
      setAmount('')
      loadUserData()
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  const changeColor = (ticker?.change_24h || 0) >= 0 ? 'var(--green)' : 'var(--red)'

  return (
    <div className="page-enter" style={{ maxWidth: '800px', margin: '0 auto', padding: '1.5rem 1rem' }}>

      {/* Price header */}
      <div style={{ marginBottom: '1.5rem' }}>
        <div style={{ display: 'flex', alignItems: 'baseline', gap: '0.75rem', flexWrap: 'wrap' }}>
          <h1 style={{ fontSize: '2rem', fontWeight: '700', color: 'var(--text-primary)', margin: 0 }}>
            SPC/EUR
          </h1>
          {ticker && (
            <>
              <span style={{ fontSize: '2rem', fontWeight: '600', color: 'var(--text-primary)' }}>
                {formatEUR(ticker.price)}
              </span>
              <span style={{ fontSize: '1rem', fontWeight: '600', color: changeColor }}>
                {ticker.change_24h >= 0 ? '+' : ''}{ticker.change_24h}%
              </span>
            </>
          )}
        </div>
        {ticker && (
          <div style={{ display: 'flex', gap: '1.5rem', marginTop: '0.5rem', fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
            <span>24h High: {formatEUR(ticker.high_24h)}</span>
            <span>24h Low: {formatEUR(ticker.low_24h)}</span>
            <span>Vol: {formatEUR(ticker.volume_24h)}</span>
          </div>
        )}
      </div>

      {/* Chart */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1rem', marginBottom: '1.5rem',
      }}>
        <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.75rem' }}>
          {['1h', '24h', '7d', '30d'].map(r => (
            <button
              key={r}
              onClick={() => setChartRange(r)}
              style={{
                padding: '0.3rem 0.7rem', borderRadius: '6px', border: 'none', cursor: 'pointer',
                fontSize: '0.75rem', fontWeight: chartRange === r ? '600' : '400',
                background: chartRange === r ? 'var(--accent)' : 'var(--bg-secondary)',
                color: chartRange === r ? '#fff' : 'var(--text-secondary)',
              }}
            >
              {r.toUpperCase()}
            </button>
          ))}
        </div>
        <PriceChart data={chartData} width={760} height={220} color={changeColor} />
      </div>

      {/* Buy/Sell panel */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.25rem', marginBottom: '1.5rem',
      }}>
        {!user ? (
          <div style={{ textAlign: 'center', padding: '2rem 1rem' }}>
            <p style={{ color: 'var(--text-secondary)', marginBottom: '1rem' }}>Inicia sesión para operar</p>
            <button
              onClick={() => onNavigate('/login')}
              style={{
                padding: '0.6rem 2rem', background: 'var(--accent)', border: 'none', borderRadius: '8px',
                color: '#fff', fontSize: '0.9rem', fontWeight: '600', cursor: 'pointer',
              }}
            >
              Entrar
            </button>
          </div>
        ) : (
          <>
            {/* Tabs */}
            <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1.25rem' }}>
              <button
                onClick={() => { setTab('buy'); setResult(null); setError(null) }}
                style={{
                  flex: 1, padding: '0.6rem', borderRadius: '8px', border: 'none', cursor: 'pointer',
                  fontSize: '0.95rem', fontWeight: '600',
                  background: tab === 'buy' ? 'var(--green)' : 'var(--bg-secondary)',
                  color: tab === 'buy' ? '#fff' : 'var(--text-secondary)',
                }}
              >
                Comprar
              </button>
              <button
                onClick={() => { setTab('sell'); setResult(null); setError(null) }}
                style={{
                  flex: 1, padding: '0.6rem', borderRadius: '8px', border: 'none', cursor: 'pointer',
                  fontSize: '0.95rem', fontWeight: '600',
                  background: tab === 'sell' ? 'var(--red)' : 'var(--bg-secondary)',
                  color: tab === 'sell' ? '#fff' : 'var(--text-secondary)',
                }}
              >
                Vender
              </button>
            </div>

            {/* Balance */}
            <div style={{
              display: 'flex', justifyContent: 'space-between', marginBottom: '1rem',
              fontSize: '0.8rem', color: 'var(--text-secondary)',
            }}>
              <span>Balance disponible:</span>
              <span style={{ fontWeight: '600', color: 'var(--text-primary)' }}>
                {tab === 'buy' ? formatEUR(balance.eur) : `${formatSPC(balance.spc)} SPC`}
              </span>
            </div>

            {/* Amount input */}
            <div style={{ marginBottom: '0.75rem' }}>
              <label style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', display: 'block', marginBottom: '0.3rem' }}>
                Cantidad (SPC)
              </label>
              <input
                type="number"
                value={amount}
                onChange={e => setAmount(e.target.value)}
                placeholder="0.00"
                min="0"
                step="0.01"
                style={{
                  width: '100%', padding: '0.65rem 0.85rem', borderRadius: '8px',
                  border: '1px solid var(--border)', background: 'var(--bg-secondary)',
                  color: 'var(--text-primary)', fontSize: '1.1rem', fontWeight: '600',
                  outline: 'none',
                }}
              />
            </div>

            {/* EUR equivalent */}
            {amountNum > 0 && (
              <div style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '0.75rem', textAlign: 'right' }}>
                = <span style={{ fontWeight: '600', color: 'var(--text-primary)' }}>{formatEUR(totalEUR)}</span>
              </div>
            )}

            {/* Quick amounts */}
            <div style={{ display: 'flex', gap: '0.4rem', marginBottom: '1rem' }}>
              {[10, 50, 100, 500].map(q => (
                <button
                  key={q}
                  onClick={() => setAmount(String(q))}
                  style={{
                    flex: 1, padding: '0.35rem', borderRadius: '6px', border: '1px solid var(--border)',
                    background: 'var(--bg-secondary)', color: 'var(--text-secondary)',
                    fontSize: '0.75rem', cursor: 'pointer',
                  }}
                >
                  {q} SPC
                </button>
              ))}
              <button
                onClick={() => {
                  if (tab === 'buy' && price > 0) setAmount(String(Math.floor((balance.eur / price) * 100) / 100))
                  else setAmount(String(balance.spc))
                }}
                style={{
                  flex: 1, padding: '0.35rem', borderRadius: '6px', border: '1px solid var(--border)',
                  background: 'var(--bg-secondary)', color: 'var(--accent)',
                  fontSize: '0.75rem', fontWeight: '600', cursor: 'pointer',
                }}
              >
                Max
              </button>
            </div>

            {/* Submit */}
            <button
              onClick={handleTrade}
              disabled={loading || amountNum <= 0}
              style={{
                width: '100%', padding: '0.75rem', borderRadius: '8px', border: 'none',
                cursor: loading || amountNum <= 0 ? 'not-allowed' : 'pointer',
                fontSize: '1rem', fontWeight: '700',
                background: tab === 'buy' ? 'var(--green)' : 'var(--red)',
                color: '#fff',
                opacity: loading || amountNum <= 0 ? 0.5 : 1,
              }}
            >
              {loading ? 'Procesando...' : tab === 'buy' ? `Comprar ${amountNum > 0 ? formatSPC(amountNum) : ''} SPC` : `Vender ${amountNum > 0 ? formatSPC(amountNum) : ''} SPC`}
            </button>

            {/* Result */}
            {result && (
              <div style={{
                marginTop: '1rem', padding: '0.75rem', borderRadius: '8px',
                background: 'rgba(16, 185, 129, 0.1)', border: '1px solid rgba(16, 185, 129, 0.3)',
                fontSize: '0.85rem', color: 'var(--green)',
              }}>
                {result.type === 'buy' ? 'Compra' : 'Venta'} completada: {formatSPC(result.amount_spc)} SPC a {formatEUR(result.price_eur)} = {formatEUR(result.total_eur)}
              </div>
            )}
            {error && (
              <div style={{
                marginTop: '1rem', padding: '0.75rem', borderRadius: '8px',
                background: 'rgba(239, 68, 68, 0.1)', border: '1px solid rgba(239, 68, 68, 0.3)',
                fontSize: '0.85rem', color: 'var(--red)',
              }}>
                {error}
              </div>
            )}
          </>
        )}
      </div>

      {/* Trade history */}
      {user && trades.length > 0 && (
        <div style={{
          background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
          padding: '1.25rem',
        }}>
          <h3 style={{ fontSize: '0.9rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
            Historial de operaciones
          </h3>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.8rem' }}>
            <thead>
              <tr style={{ color: 'var(--text-secondary)', borderBottom: '1px solid var(--border)' }}>
                <th style={{ textAlign: 'left', padding: '0.5rem 0.25rem', fontWeight: '500' }}>Tipo</th>
                <th style={{ textAlign: 'right', padding: '0.5rem 0.25rem', fontWeight: '500' }}>Cantidad</th>
                <th style={{ textAlign: 'right', padding: '0.5rem 0.25rem', fontWeight: '500' }}>Precio</th>
                <th style={{ textAlign: 'right', padding: '0.5rem 0.25rem', fontWeight: '500' }}>Total</th>
              </tr>
            </thead>
            <tbody>
              {trades.slice(0, 20).map((t, i) => (
                <tr key={t.id || i} style={{ borderBottom: '1px solid var(--border)' }}>
                  <td style={{
                    padding: '0.5rem 0.25rem', fontWeight: '600',
                    color: t.type === 'buy' ? 'var(--green)' : 'var(--red)',
                  }}>
                    {t.type === 'buy' ? 'Compra' : 'Venta'}
                  </td>
                  <td style={{ textAlign: 'right', padding: '0.5rem 0.25rem', color: 'var(--text-primary)' }}>
                    {formatSPC(t.amount_spc)} SPC
                  </td>
                  <td style={{ textAlign: 'right', padding: '0.5rem 0.25rem', color: 'var(--text-secondary)' }}>
                    {formatEUR(t.price_eur)}
                  </td>
                  <td style={{ textAlign: 'right', padding: '0.5rem 0.25rem', color: 'var(--text-primary)' }}>
                    {formatEUR(t.total_eur)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Testnet badge */}
      <div style={{
        textAlign: 'center', marginTop: '1.5rem', padding: '0.5rem',
        fontSize: '0.75rem', color: 'var(--text-secondary)',
      }}>
        Testnet — los fondos no tienen valor real. Cada usuario recibe 1.000€ virtuales al registrarse.
      </div>
    </div>
  )
}
