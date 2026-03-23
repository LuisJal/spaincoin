import { useState, useEffect, useCallback } from 'react'
import { useAuth } from '../auth/useAuth.jsx'
import { getPortfolio, depositEUR, getTradeHistory } from '../api/client.js'

const formatEUR = (n) => new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(n)
const formatAmount = (n, sym) => {
  if (sym === 'BTC') return n.toFixed(8)
  if (['ETH', 'BNB', 'SOL', 'AVAX'].includes(sym)) return n.toFixed(6)
  return n >= 1 ? n.toLocaleString('es-ES', { maximumFractionDigits: 4 }) : n.toFixed(6)
}

const coinColors = {
  SPC: '#ffc400', BTC: '#f7931a', ETH: '#627eea', BNB: '#f3ba2f',
  SOL: '#9945ff', XRP: '#23292f', ADA: '#0033ad', DOGE: '#c2a633',
  DOT: '#e6007a', AVAX: '#e84142', MATIC: '#8247e5',
}

export default function Wallet({ onNavigate }) {
  const { user, token } = useAuth()
  const [portfolio, setPortfolio] = useState(null)
  const [trades, setTrades] = useState([])
  const [depositAmount, setDepositAmount] = useState('')
  const [showDeposit, setShowDeposit] = useState(false)
  const [depositLoading, setDepositLoading] = useState(false)
  const [depositResult, setDepositResult] = useState(null)

  const loadData = useCallback(async () => {
    if (!token) return
    try {
      const [p, t] = await Promise.all([getPortfolio(token), getTradeHistory(token)])
      setPortfolio(p)
      setTrades(t)
    } catch (e) {
      console.error('Portfolio error:', e)
    }
  }, [token])

  useEffect(() => {
    loadData()
    const i = setInterval(loadData, 15000)
    return () => clearInterval(i)
  }, [loadData])

  // Redirect to login if not authenticated
  if (!user) {
    return (
      <div className="page-enter" style={{ maxWidth: '600px', margin: '0 auto', padding: '4rem 1rem', textAlign: 'center' }}>
        <h2 style={{ color: 'var(--text-primary)', marginBottom: '1rem' }}>Mi Cartera</h2>
        <p style={{ color: 'var(--text-secondary)', marginBottom: '1.5rem' }}>Inicia sesión para ver tu portfolio</p>
        <button onClick={() => onNavigate('/login')} style={{
          padding: '0.7rem 2rem', background: 'var(--accent)', border: 'none', borderRadius: '8px',
          color: '#fff', fontSize: '0.95rem', fontWeight: '600', cursor: 'pointer',
        }}>Entrar</button>
      </div>
    )
  }

  if (!portfolio) {
    return <div style={{ display: 'flex', justifyContent: 'center', padding: '4rem' }}><div className="spinner" /></div>
  }

  async function handleDeposit() {
    const amt = parseFloat(depositAmount)
    if (!amt || amt <= 0) return
    setDepositLoading(true); setDepositResult(null)
    try {
      const res = await depositEUR(token, amt)
      setDepositResult(`+${formatEUR(res.deposited)} depositados. Nuevo saldo: ${formatEUR(res.new_balance)}`)
      setDepositAmount('')
      loadData()
    } catch (e) {
      setDepositResult('Error: ' + e.message)
    } finally {
      setDepositLoading(false)
    }
  }

  return (
    <div className="page-enter" style={{ maxWidth: '800px', margin: '0 auto', padding: '1.5rem 1rem' }}>

      {/* Total value */}
      <div style={{ marginBottom: '1.5rem' }}>
        <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', marginBottom: '0.25rem' }}>Valor total del portfolio</div>
        <div style={{ fontSize: '2.5rem', fontWeight: '700', color: 'var(--text-primary)' }}>
          {formatEUR(portfolio.total_value_eur)}
        </div>
      </div>

      {/* EUR balance card */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.25rem', marginBottom: '1rem', display: 'flex', justifyContent: 'space-between',
        alignItems: 'center', flexWrap: 'wrap', gap: '1rem',
      }}>
        <div>
          <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginBottom: '0.25rem' }}>SALDO EUR</div>
          <div style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--text-primary)' }}>{formatEUR(portfolio.eur)}</div>
        </div>
        <button onClick={() => setShowDeposit(!showDeposit)} style={{
          padding: '0.5rem 1.25rem', borderRadius: '8px', border: 'none',
          background: 'var(--green)', color: '#fff', fontSize: '0.85rem',
          fontWeight: '600', cursor: 'pointer',
        }}>
          {showDeposit ? 'Cerrar' : 'Depositar EUR'}
        </button>
      </div>

      {/* Deposit form */}
      {showDeposit && (
        <div style={{
          background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
          padding: '1.25rem', marginBottom: '1rem',
        }}>
          <div style={{ fontSize: '0.85rem', color: 'var(--text-primary)', fontWeight: '600', marginBottom: '0.75rem' }}>
            Depositar EUR (testnet)
          </div>
          <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.75rem' }}>
            {[100, 500, 1000, 5000].map(amt => (
              <button key={amt} onClick={() => setDepositAmount(String(amt))} style={{
                flex: 1, padding: '0.4rem', borderRadius: '6px', border: '1px solid var(--border)',
                background: depositAmount === String(amt) ? 'var(--accent)' : 'var(--bg-secondary)',
                color: depositAmount === String(amt) ? '#fff' : 'var(--text-secondary)',
                fontSize: '0.8rem', cursor: 'pointer',
              }}>{amt}€</button>
            ))}
          </div>
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <input type="number" value={depositAmount} onChange={e => setDepositAmount(e.target.value)}
              placeholder="Cantidad en EUR" min="0" style={{
                flex: 1, padding: '0.6rem', borderRadius: '8px', border: '1px solid var(--border)',
                background: 'var(--bg-secondary)', color: 'var(--text-primary)', fontSize: '0.9rem', outline: 'none',
              }} />
            <button onClick={handleDeposit} disabled={depositLoading} style={{
              padding: '0.6rem 1.5rem', borderRadius: '8px', border: 'none',
              background: 'var(--green)', color: '#fff', fontWeight: '600', cursor: 'pointer',
              opacity: depositLoading ? 0.5 : 1,
            }}>{depositLoading ? '...' : 'Depositar'}</button>
          </div>
          {depositResult && (
            <div style={{ marginTop: '0.75rem', fontSize: '0.8rem', color: depositResult.startsWith('Error') ? 'var(--red)' : 'var(--green)' }}>
              {depositResult}
            </div>
          )}
        </div>
      )}

      {/* Holdings */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        overflow: 'hidden', marginBottom: '1.5rem',
      }}>
        <div style={{ padding: '1rem 1.25rem', borderBottom: '1px solid var(--border)' }}>
          <h2 style={{ fontSize: '1rem', fontWeight: '600', color: 'var(--text-primary)', margin: 0 }}>
            Mis activos
          </h2>
        </div>

        {(!portfolio.holdings || portfolio.holdings.length === 0) ? (
          <div style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>
            No tienes activos todavía. <button onClick={() => onNavigate('/trade/SPC')} style={{
              background: 'none', border: 'none', color: 'var(--accent)', cursor: 'pointer', fontWeight: '600',
            }}>Compra tu primera crypto</button>
          </div>
        ) : (
          portfolio.holdings.map((h, i) => (
            <div key={h.symbol} style={{
              display: 'grid', gridTemplateColumns: '2fr 1.5fr 1.5fr 1fr',
              padding: '1rem 1.25rem', alignItems: 'center',
              borderBottom: i < portfolio.holdings.length - 1 ? '1px solid var(--border)' : 'none',
            }}>
              {/* Coin */}
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem' }}>
                <div style={{
                  width: '36px', height: '36px', borderRadius: '50%',
                  background: coinColors[h.symbol] || '#666',
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  fontWeight: '700', fontSize: '0.7rem', color: '#fff', flexShrink: 0,
                }}>{h.symbol.slice(0, 1)}</div>
                <div>
                  <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)' }}>{h.symbol}</div>
                  <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)' }}>{h.name}</div>
                </div>
              </div>

              {/* Amount */}
              <div style={{ textAlign: 'right' }}>
                <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)' }}>
                  {formatAmount(h.amount, h.symbol)}
                </div>
                <div style={{ fontSize: '0.7rem', color: 'var(--text-secondary)' }}>
                  @ {formatEUR(h.price)}
                </div>
              </div>

              {/* Value */}
              <div style={{ textAlign: 'right', fontWeight: '600', fontSize: '0.95rem', color: 'var(--text-primary)' }}>
                {formatEUR(h.value_eur)}
              </div>

              {/* Action */}
              <div style={{ textAlign: 'right' }}>
                <button onClick={() => onNavigate(`/trade/${h.symbol}`)} style={{
                  padding: '0.3rem 0.65rem', borderRadius: '6px', border: 'none',
                  background: 'var(--accent)', color: '#fff', fontSize: '0.75rem',
                  fontWeight: '600', cursor: 'pointer',
                }}>Operar</button>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Recent trades */}
      {trades.length > 0 && (
        <div style={{
          background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
          padding: '1.25rem',
        }}>
          <h3 style={{ fontSize: '0.9rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
            Últimas operaciones
          </h3>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.8rem' }}>
            <thead>
              <tr style={{ color: 'var(--text-secondary)', borderBottom: '1px solid var(--border)' }}>
                <th style={{ textAlign: 'left', padding: '0.5rem 0', fontWeight: '500' }}>Tipo</th>
                <th style={{ textAlign: 'left', padding: '0.5rem 0', fontWeight: '500' }}>Par</th>
                <th style={{ textAlign: 'right', padding: '0.5rem 0', fontWeight: '500' }}>Cantidad</th>
                <th style={{ textAlign: 'right', padding: '0.5rem 0', fontWeight: '500' }}>Total</th>
              </tr>
            </thead>
            <tbody>
              {trades.slice(0, 10).map((t, i) => (
                <tr key={t.id || i} style={{ borderBottom: '1px solid var(--border)' }}>
                  <td style={{ padding: '0.5rem 0', fontWeight: '600', color: t.type === 'buy' ? 'var(--green)' : 'var(--red)' }}>
                    {t.type === 'buy' ? 'Compra' : 'Venta'}
                  </td>
                  <td style={{ padding: '0.5rem 0', color: 'var(--text-secondary)' }}>{t.pair}</td>
                  <td style={{ textAlign: 'right', padding: '0.5rem 0', color: 'var(--text-primary)' }}>
                    {formatAmount(t.amount || t.amount_spc, t.symbol || 'SPC')}
                  </td>
                  <td style={{ textAlign: 'right', padding: '0.5rem 0', color: 'var(--text-primary)' }}>{formatEUR(t.total_eur)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div style={{ textAlign: 'center', marginTop: '1.5rem', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>
        Testnet — fondos virtuales sin valor real
      </div>
    </div>
  )
}
