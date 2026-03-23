import { useState, useEffect, useCallback } from 'react'
import { getStatus, getExplorer, getPrice } from '../api/client.js'
import { formatNumber } from '../utils/format.js'

export default function Dashboard({ onNavigate }) {
  const [status, setStatus] = useState(null)
  const [explorer, setExplorer] = useState(null)
  const [price, setPrice] = useState(null)
  const [walletInput, setWalletInput] = useState('')
  const [walletResult, setWalletResult] = useState(null)
  const [walletLoading, setWalletLoading] = useState(false)
  const [walletError, setWalletError] = useState(null)
  const [error, setError] = useState(null)

  const fetchData = useCallback(async () => {
    try {
      const [s, e, p] = await Promise.allSettled([getStatus(), getExplorer(), getPrice()])
      if (s.status === 'fulfilled') setStatus(s.value)
      if (e.status === 'fulfilled') setExplorer(e.value)
      if (p.status === 'fulfilled') setPrice(p.value)
      if (s.status === 'rejected' && e.status === 'rejected') setError('Sin conexión al nodo')
      else setError(null)
    } catch (err) {
      setError(err.message)
    }
  }, [])

  useEffect(() => {
    fetchData()
    const interval = setInterval(fetchData, 10000)
    return () => clearInterval(interval)
  }, [fetchData])

  async function handleCheckBalance(e) {
    e.preventDefault()
    if (!walletInput.trim()) return
    setWalletLoading(true)
    setWalletError(null)
    setWalletResult(null)
    try {
      const res = await fetch(`/api/wallet/${walletInput.trim()}`)
      const data = await res.json()
      if (!res.ok) throw new Error(data.error || 'Error consultando balance')
      setWalletResult(data)
    } catch (err) {
      setWalletError(err.message)
    } finally {
      setWalletLoading(false)
    }
  }

  const priceEur = price?.price_eur ?? 0.09
  const height = status?.node?.height ?? 0
  const totalSupply = explorer?.total_supply_spc ?? 0
  const networkOk = !error && status !== null

  return (
    <div style={{ maxWidth: '1100px', margin: '0 auto', padding: '2.5rem 1.5rem 4rem' }}>

      {/* ===== HERO — Precio y CTA ===== */}
      <div style={{
        background: 'linear-gradient(135deg, #0f172a 0%, #1e1b4b 50%, #0f172a 100%)',
        border: '1px solid rgba(99, 102, 241, 0.3)',
        borderRadius: '20px',
        padding: '3rem 2.5rem',
        marginBottom: '2rem',
        textAlign: 'center',
        position: 'relative',
        overflow: 'hidden',
      }}>
        {/* Fondo decorativo */}
        <div style={{
          position: 'absolute', top: '-60px', right: '-60px',
          width: '200px', height: '200px',
          background: 'radial-gradient(circle, rgba(59,130,246,0.15) 0%, transparent 70%)',
          borderRadius: '50%',
        }} />

        {/* Badge red */}
        <div style={{
          display: 'inline-flex', alignItems: 'center', gap: '0.4rem',
          background: 'rgba(16,185,129,0.15)', border: '1px solid rgba(16,185,129,0.3)',
          color: '#10b981', fontSize: '0.75rem', fontWeight: '600',
          padding: '0.3rem 0.8rem', borderRadius: '20px', marginBottom: '1.5rem',
        }}>
          <span style={{ animation: 'pulse 2s infinite', display: 'inline-block' }}>●</span>
          {networkOk ? `Red activa · Bloque #${formatNumber(height)}` : 'Conectando...'}
        </div>

        {/* Precio principal */}
        <div style={{ marginBottom: '0.5rem' }}>
          <span style={{ fontSize: '0.9rem', color: '#9ca3af', fontWeight: '500' }}>Precio $SPC</span>
        </div>
        <div style={{
          fontSize: 'clamp(3rem, 8vw, 5rem)', fontWeight: '800',
          color: '#f9fafb', letterSpacing: '-0.03em', lineHeight: 1,
          marginBottom: '0.5rem',
        }}>
          €{priceEur}
        </div>
        <div style={{ fontSize: '0.9rem', color: '#9ca3af', marginBottom: '2rem' }}>
          ≈ ${(priceEur * 1.09).toFixed(2)} USD · Testnet
        </div>

        {/* CTAs */}
        <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center', flexWrap: 'wrap' }}>
          <button
            onClick={() => onNavigate('/trade')}
            style={{
              background: 'linear-gradient(135deg, #3b82f6, #1d4ed8)',
              color: '#fff', border: 'none', borderRadius: '12px',
              padding: '0.9rem 2rem', fontSize: '1rem', fontWeight: '600',
              cursor: 'pointer', transition: 'opacity 0.15s',
            }}
            onMouseEnter={e => e.target.style.opacity = '0.85'}
            onMouseLeave={e => e.target.style.opacity = '1'}
          >
            Comprar $SPC
          </button>
          <button
            onClick={() => onNavigate('/wallet')}
            style={{
              background: 'transparent',
              color: '#f9fafb', border: '1px solid rgba(255,255,255,0.2)',
              borderRadius: '12px', padding: '0.9rem 2rem',
              fontSize: '1rem', fontWeight: '500',
              cursor: 'pointer', transition: 'border-color 0.15s',
            }}
            onMouseEnter={e => e.target.style.borderColor = 'rgba(255,255,255,0.5)'}
            onMouseLeave={e => e.target.style.borderColor = 'rgba(255,255,255,0.2)'}
          >
            Ver mi saldo
          </button>
        </div>
      </div>

      {/* ===== STATS — Lo que le importa al usuario ===== */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
        gap: '1rem',
        marginBottom: '2rem',
      }}>
        {[
          { label: 'Precio actual', value: `€${priceEur}`, sub: 'por $SPC', color: '#3b82f6' },
          { label: 'Supply en circulación', value: `${formatNumber(Math.round(totalSupply))}`, sub: 'de 21.000.000 máximo', color: '#10b981' },
          { label: 'Bloques producidos', value: formatNumber(height), sub: 'cada 5 segundos', color: '#8b5cf6' },
          { label: 'Estado de la red', value: networkOk ? '✓ Activa' : '⟳ Cargando', sub: 'Proof of Stake', color: networkOk ? '#10b981' : '#9ca3af' },
        ].map(({ label, value, sub, color }) => (
          <div key={label} style={{
            background: 'var(--bg-card)',
            border: '1px solid var(--border)',
            borderRadius: '12px',
            padding: '1.25rem 1.5rem',
          }}>
            <div style={{ fontSize: '0.8rem', color: '#9ca3af', marginBottom: '0.4rem' }}>{label}</div>
            <div style={{ fontSize: '1.5rem', fontWeight: '700', color, marginBottom: '0.2rem' }}>{value}</div>
            <div style={{ fontSize: '0.75rem', color: '#6b7280' }}>{sub}</div>
          </div>
        ))}
      </div>

      {/* ===== CONSULTA DE SALDO ===== */}
      <div style={{
        background: 'var(--bg-card)',
        border: '1px solid var(--border)',
        borderRadius: '16px',
        padding: '2rem',
        marginBottom: '2rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '700', color: '#f9fafb', marginBottom: '0.4rem' }}>
          Consulta tu saldo
        </h2>
        <p style={{ fontSize: '0.85rem', color: '#9ca3af', marginBottom: '1.5rem' }}>
          Introduce tu dirección SpainCoin para ver cuántos $SPC tienes.
        </p>

        <form onSubmit={handleCheckBalance} style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
          <input
            type="text"
            value={walletInput}
            onChange={e => setWalletInput(e.target.value)}
            placeholder="SPCtu_dirección_aquí..."
            style={{
              flex: '1', minWidth: '260px',
              background: '#111827', border: '1px solid #374151',
              borderRadius: '10px', padding: '0.75rem 1rem',
              color: '#f9fafb', fontSize: '0.9rem', fontFamily: 'monospace',
              outline: 'none',
            }}
          />
          <button
            type="submit"
            disabled={walletLoading}
            style={{
              background: '#3b82f6', color: '#fff', border: 'none',
              borderRadius: '10px', padding: '0.75rem 1.5rem',
              fontSize: '0.9rem', fontWeight: '600', cursor: 'pointer',
              opacity: walletLoading ? 0.6 : 1,
            }}
          >
            {walletLoading ? 'Consultando...' : 'Ver saldo'}
          </button>
        </form>

        {walletError && (
          <div style={{ marginTop: '1rem', color: '#ef4444', fontSize: '0.85rem' }}>
            ⚠️ {walletError}
          </div>
        )}

        {walletResult && (
          <div style={{
            marginTop: '1.25rem',
            background: 'rgba(16,185,129,0.08)',
            border: '1px solid rgba(16,185,129,0.2)',
            borderRadius: '10px', padding: '1.25rem',
          }}>
            <div style={{ fontSize: '0.8rem', color: '#9ca3af', marginBottom: '0.3rem' }}>Saldo</div>
            <div style={{ fontSize: '2rem', fontWeight: '800', color: '#10b981', marginBottom: '0.25rem' }}>
              {Number(walletResult.balance_spc).toFixed(3)} SPC
            </div>
            <div style={{ fontSize: '0.8rem', color: '#6b7280', fontFamily: 'monospace' }}>
              {walletResult.address}
            </div>
          </div>
        )}
      </div>

      {/* ===== CÓMO FUNCIONA ===== */}
      <div style={{
        background: 'var(--bg-card)',
        border: '1px solid var(--border)',
        borderRadius: '16px',
        padding: '2rem',
        marginBottom: '2rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '700', color: '#f9fafb', marginBottom: '1.5rem' }}>
          ¿Cómo funciona SpainCoin?
        </h2>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1.5rem' }}>
          {[
            { icon: '🔐', title: 'Tu wallet, tu dinero', desc: 'Nadie puede acceder a tus $SPC sin tu clave privada. Ni nosotros.' },
            { icon: '⚡', title: 'Transacciones rápidas', desc: 'Un bloque cada 5 segundos. Tu transacción confirmada en segundos.' },
            { icon: '🌐', title: 'Red descentralizada', desc: 'No hay un servidor central. La red funciona aunque apagues cualquier nodo.' },
            { icon: '🔒', title: 'Supply limitado', desc: 'Solo existirán 21.000.000 $SPC. Nunca se crearán más.' },
          ].map(({ icon, title, desc }) => (
            <div key={title}>
              <div style={{ fontSize: '1.75rem', marginBottom: '0.6rem' }}>{icon}</div>
              <div style={{ fontSize: '0.9rem', fontWeight: '600', color: '#f9fafb', marginBottom: '0.35rem' }}>{title}</div>
              <div style={{ fontSize: '0.8rem', color: '#9ca3af', lineHeight: '1.5' }}>{desc}</div>
            </div>
          ))}
        </div>
      </div>

      {/* ===== LINK AL EXPLORER (para los técnicos) ===== */}
      <div style={{ textAlign: 'center' }}>
        <button
          onClick={() => onNavigate('/explorer')}
          style={{
            background: 'transparent', border: '1px solid #374151',
            color: '#9ca3af', borderRadius: '8px',
            padding: '0.6rem 1.25rem', fontSize: '0.8rem',
            cursor: 'pointer',
          }}
        >
          Ver explorador de bloques →
        </button>
      </div>

    </div>
  )
}
