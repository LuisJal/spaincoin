import { useState, useEffect, useCallback } from 'react'
import { useAuth } from '../auth/useAuth.jsx'
import PriceChart from '../components/PriceChart.jsx'
import { getMarketTable, buyAsset, sellAsset, getPortfolio, getTradeHistory } from '../api/client.js'

const formatEUR = (n) => n >= 1000
  ? new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR', maximumFractionDigits: 0 }).format(n)
  : new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(n)
const formatAmount = (n, sym) => {
  if (sym === 'BTC') return n.toFixed(8)
  if (['ETH', 'BNB', 'SOL', 'AVAX'].includes(sym)) return n.toFixed(6)
  return n >= 1 ? n.toLocaleString('es-ES', { maximumFractionDigits: 4 }) : n.toFixed(6)
}

function generateChartData(price, change, seed) {
  const points = 80
  const data = []
  for (let i = 0; i < points; i++) {
    const t = i / points
    const w = Math.sin(i / 4 + seed * 7) * 0.012 + Math.sin(i / 9 + seed * 3) * 0.008
    const trend = (change / 100) * t
    const p = price * (1 + w + trend - (change / 100) * 0.5)
    data.push({ price: p, height: i })
  }
  return data
}

export default function Trade({ symbol = 'SPC', onNavigate }) {
  const { user, token } = useAuth()
  const [coinData, setCoinData] = useState(null)
  const [allCoins, setAllCoins] = useState([])
  const [tab, setTab] = useState('buy')
  const [amount, setAmount] = useState('')
  const [portfolio, setPortfolio] = useState({ eur: 0, holdings: [] })
  const [trades, setTrades] = useState([])
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)

  const loadMarket = useCallback(async () => {
    try {
      const table = await getMarketTable()
      setAllCoins(table)
      const coin = table.find(c => c.symbol === symbol) || table[0]
      setCoinData(coin)
    } catch (e) {
      console.error('Market load error:', e)
    }
  }, [symbol])

  const loadUserData = useCallback(async () => {
    if (!token) return
    try {
      const [p, t] = await Promise.all([getPortfolio(token), getTradeHistory(token)])
      setPortfolio(p)
      setTrades(t.filter(tr => !tr.symbol || tr.symbol === symbol || tr.pair === symbol + '/EUR'))
    } catch (e) {
      console.error('User data error:', e)
    }
  }, [token, symbol])

  useEffect(() => { loadMarket(); const i = setInterval(loadMarket, 10000); return () => clearInterval(i) }, [loadMarket])
  useEffect(() => { loadUserData() }, [loadUserData])
  useEffect(() => { setAmount(''); setResult(null); setError(null) }, [symbol])

  const price = coinData?.price || 0
  const change = coinData?.change_24h || 0
  const amountNum = parseFloat(amount) || 0
  const totalEUR = amountNum * price
  const holding = portfolio.holdings?.find(h => h.symbol === symbol)
  const holdingAmount = holding?.amount || 0
  const changeColor = change >= 0 ? 'var(--green)' : 'var(--red)'
  const chartData = coinData ? generateChartData(price, change, allCoins.findIndex(c => c.symbol === symbol)) : []

  async function handleTrade() {
    if (amountNum <= 0) return
    setLoading(true); setError(null); setResult(null)
    try {
      const res = tab === 'buy'
        ? await buyAsset(token, symbol, amountNum)
        : await sellAsset(token, symbol, amountNum)
      setResult(res)
      setAmount('')
      loadUserData()
    } catch (e) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="page-enter" style={{ maxWidth: '800px', margin: '0 auto', padding: '1.5rem 1rem' }}>

      {/* Symbol selector */}
      <div style={{ display: 'flex', gap: '0.4rem', marginBottom: '1rem', flexWrap: 'wrap' }}>
        {allCoins.map(c => (
          <button
            key={c.symbol}
            onClick={() => onNavigate(`/trade/${c.symbol}`)}
            style={{
              padding: '0.3rem 0.6rem', borderRadius: '6px', border: 'none', cursor: 'pointer',
              fontSize: '0.75rem', fontWeight: symbol === c.symbol ? '700' : '400',
              background: symbol === c.symbol ? 'var(--accent)' : 'var(--bg-secondary)',
              color: symbol === c.symbol ? '#fff' : 'var(--text-secondary)',
            }}
          >
            {c.symbol}
          </button>
        ))}
      </div>

      {/* Price header */}
      <div style={{ marginBottom: '1.5rem' }}>
        <div style={{ display: 'flex', alignItems: 'baseline', gap: '0.75rem', flexWrap: 'wrap' }}>
          <h1 style={{ fontSize: '2rem', fontWeight: '700', color: 'var(--text-primary)', margin: 0 }}>
            {symbol}/EUR
          </h1>
          {coinData && (
            <>
              <span style={{ fontSize: '2rem', fontWeight: '600', color: 'var(--text-primary)' }}>
                {formatEUR(price)}
              </span>
              <span style={{ fontSize: '1rem', fontWeight: '600', color: changeColor }}>
                {change >= 0 ? '+' : ''}{change}%
              </span>
            </>
          )}
        </div>
        {coinData && (
          <div style={{ display: 'flex', gap: '1.5rem', marginTop: '0.5rem', fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
            <span>Vol: {formatEUR(coinData.volume)}</span>
            <span>MCap: {formatEUR(coinData.market_cap)}</span>
          </div>
        )}
      </div>

      {/* Chart */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1rem', marginBottom: '1.5rem',
      }}>
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
            <button onClick={() => onNavigate('/login')} style={{
              padding: '0.6rem 2rem', background: 'var(--accent)', border: 'none', borderRadius: '8px',
              color: '#fff', fontSize: '0.9rem', fontWeight: '600', cursor: 'pointer',
            }}>Entrar</button>
          </div>
        ) : (
          <>
            {/* Tabs */}
            <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1.25rem' }}>
              <button onClick={() => { setTab('buy'); setResult(null); setError(null) }}
                style={{
                  flex: 1, padding: '0.6rem', borderRadius: '8px', border: 'none', cursor: 'pointer',
                  fontSize: '0.95rem', fontWeight: '600',
                  background: tab === 'buy' ? 'var(--green)' : 'var(--bg-secondary)',
                  color: tab === 'buy' ? '#fff' : 'var(--text-secondary)',
                }}>Comprar</button>
              <button onClick={() => { setTab('sell'); setResult(null); setError(null) }}
                style={{
                  flex: 1, padding: '0.6rem', borderRadius: '8px', border: 'none', cursor: 'pointer',
                  fontSize: '0.95rem', fontWeight: '600',
                  background: tab === 'sell' ? 'var(--red)' : 'var(--bg-secondary)',
                  color: tab === 'sell' ? '#fff' : 'var(--text-secondary)',
                }}>Vender</button>
            </div>

            {/* Balance */}
            <div style={{
              display: 'flex', justifyContent: 'space-between', marginBottom: '1rem',
              fontSize: '0.8rem', color: 'var(--text-secondary)',
            }}>
              <span>Balance disponible:</span>
              <span style={{ fontWeight: '600', color: 'var(--text-primary)' }}>
                {tab === 'buy' ? formatEUR(portfolio.eur) : `${formatAmount(holdingAmount, symbol)} ${symbol}`}
              </span>
            </div>

            {/* Amount input */}
            <div style={{ marginBottom: '0.75rem' }}>
              <label style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', display: 'block', marginBottom: '0.3rem' }}>
                Cantidad ({symbol})
              </label>
              <input type="number" value={amount} onChange={e => setAmount(e.target.value)}
                placeholder="0.00" min="0" step="any"
                style={{
                  width: '100%', padding: '0.65rem 0.85rem', borderRadius: '8px',
                  border: '1px solid var(--border)', background: 'var(--bg-secondary)',
                  color: 'var(--text-primary)', fontSize: '1.1rem', fontWeight: '600', outline: 'none',
                }} />
            </div>

            {amountNum > 0 && (
              <div style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '0.75rem', textAlign: 'right' }}>
                = <span style={{ fontWeight: '600', color: 'var(--text-primary)' }}>{formatEUR(totalEUR)}</span>
              </div>
            )}

            {/* Quick amounts in EUR */}
            <div style={{ display: 'flex', gap: '0.4rem', marginBottom: '1rem' }}>
              {[10, 50, 100, 500].map(eur => (
                <button key={eur} onClick={() => price > 0 && setAmount(String(Math.floor((eur / price) * 1000000) / 1000000))}
                  style={{
                    flex: 1, padding: '0.35rem', borderRadius: '6px', border: '1px solid var(--border)',
                    background: 'var(--bg-secondary)', color: 'var(--text-secondary)',
                    fontSize: '0.75rem', cursor: 'pointer',
                  }}>{eur}€</button>
              ))}
              <button onClick={() => {
                if (tab === 'buy' && price > 0) setAmount(String(Math.floor((portfolio.eur / price) * 1000000) / 1000000))
                else setAmount(String(holdingAmount))
              }} style={{
                flex: 1, padding: '0.35rem', borderRadius: '6px', border: '1px solid var(--border)',
                background: 'var(--bg-secondary)', color: 'var(--accent)',
                fontSize: '0.75rem', fontWeight: '600', cursor: 'pointer',
              }}>Max</button>
            </div>

            <button onClick={handleTrade} disabled={loading || amountNum <= 0}
              style={{
                width: '100%', padding: '0.75rem', borderRadius: '8px', border: 'none',
                cursor: loading || amountNum <= 0 ? 'not-allowed' : 'pointer',
                fontSize: '1rem', fontWeight: '700',
                background: tab === 'buy' ? 'var(--green)' : 'var(--red)', color: '#fff',
                opacity: loading || amountNum <= 0 ? 0.5 : 1,
              }}>
              {loading ? 'Procesando...' : tab === 'buy'
                ? `Comprar ${amountNum > 0 ? formatAmount(amountNum, symbol) : ''} ${symbol}`
                : `Vender ${amountNum > 0 ? formatAmount(amountNum, symbol) : ''} ${symbol}`}
            </button>

            {result && (
              <div style={{
                marginTop: '1rem', padding: '0.75rem', borderRadius: '8px',
                background: 'rgba(16, 185, 129, 0.1)', border: '1px solid rgba(16, 185, 129, 0.3)',
                fontSize: '0.85rem', color: 'var(--green)',
              }}>
                {result.type === 'buy' ? 'Compra' : 'Venta'} completada: {formatAmount(result.amount, symbol)} {symbol} a {formatEUR(result.price_eur)} = {formatEUR(result.total_eur)}
              </div>
            )}
            {error && (
              <div style={{
                marginTop: '1rem', padding: '0.75rem', borderRadius: '8px',
                background: 'rgba(239, 68, 68, 0.1)', border: '1px solid rgba(239, 68, 68, 0.3)',
                fontSize: '0.85rem', color: 'var(--red)',
              }}>{error}</div>
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
                <th style={{ textAlign: 'right', padding: '0.5rem 0.25rem', fontWeight: '500' }}>Par</th>
                <th style={{ textAlign: 'right', padding: '0.5rem 0.25rem', fontWeight: '500' }}>Cantidad</th>
                <th style={{ textAlign: 'right', padding: '0.5rem 0.25rem', fontWeight: '500' }}>Total</th>
              </tr>
            </thead>
            <tbody>
              {trades.slice(0, 20).map((t, i) => (
                <tr key={t.id || i} style={{ borderBottom: '1px solid var(--border)' }}>
                  <td style={{ padding: '0.5rem 0.25rem', fontWeight: '600', color: t.type === 'buy' ? 'var(--green)' : 'var(--red)' }}>
                    {t.type === 'buy' ? 'Compra' : 'Venta'}
                  </td>
                  <td style={{ textAlign: 'right', padding: '0.5rem 0.25rem', color: 'var(--text-secondary)' }}>{t.pair}</td>
                  <td style={{ textAlign: 'right', padding: '0.5rem 0.25rem', color: 'var(--text-primary)' }}>
                    {formatAmount(t.amount || t.amount_spc, t.symbol || 'SPC')}
                  </td>
                  <td style={{ textAlign: 'right', padding: '0.5rem 0.25rem', color: 'var(--text-primary)' }}>{formatEUR(t.total_eur)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div style={{ textAlign: 'center', marginTop: '1.5rem', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>
        Testnet — los fondos no tienen valor real
      </div>
    </div>
  )
}
