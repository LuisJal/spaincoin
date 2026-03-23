import { useState, useEffect } from 'react'
import { useAuth } from '../auth/useAuth.jsx'
import { getMe, getTradeBalance } from '../api/client.js'

export default function Account({ onNavigate }) {
  const { user, token, loading, logout } = useAuth()

  const [accountData, setAccountData] = useState(null)
  const [fetchError, setFetchError] = useState(null)

  // Export key modal state
  const [showExportModal, setShowExportModal] = useState(false)
  const [exportStep, setExportStep] = useState('warning') // 'warning' | 'password' | 'result'
  const [exportPassword, setExportPassword] = useState('')

  // Copy address state
  const [copied, setCopied] = useState(false)

  // Redirect if not authenticated
  useEffect(() => {
    if (!loading && user === null) {
      onNavigate('/login')
    }
  }, [user, loading, onNavigate])

  const [eurBalance, setEurBalance] = useState(null)
  const [exchangeSPC, setExchangeSPC] = useState(null)

  // Load full account data (balance, nonce) + exchange balances
  useEffect(() => {
    if (!token || !user) return
    getMe(token)
      .then(data => setAccountData(data))
      .catch(err => setFetchError(err.message))
    getTradeBalance(token)
      .then(data => { setEurBalance(data.eur); setExchangeSPC(data.spc) })
      .catch(() => {})
  }, [token, user])

  function copyAddress() {
    if (!accountData?.address && !user?.address) return
    const addr = accountData?.address || user?.address
    navigator.clipboard.writeText(addr).then(() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    })
  }

  function handleLogout() {
    logout()
    onNavigate('/')
  }

  if (loading) {
    return (
      <div style={{
        minHeight: 'calc(100vh - 60px)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}>
        <span className="spinner" />
      </div>
    )
  }

  if (!user) return null // Will redirect via useEffect

  const address = accountData?.address || user?.address || '—'
  const balance = accountData?.balance_spc ?? user?.balance_spc ?? 0
  const nonce = accountData?.nonce ?? '—'

  const sectionCard = {
    background: 'var(--bg-card)',
    border: '1px solid var(--border)',
    borderRadius: '14px',
    padding: '1.5rem',
    marginBottom: '1.25rem',
  }

  const sectionTitle = {
    fontSize: '0.75rem',
    fontWeight: '600',
    color: 'var(--text-secondary)',
    textTransform: 'uppercase',
    letterSpacing: '0.07em',
    marginBottom: '1.1rem',
  }

  const btnOutline = {
    padding: '0.55rem 1.1rem',
    background: 'transparent',
    border: '1px solid var(--border)',
    borderRadius: '8px',
    color: 'var(--text-secondary)',
    fontSize: '0.8rem',
    cursor: 'pointer',
    transition: 'all 0.15s ease',
  }

  const btnDanger = {
    padding: '0.55rem 1.1rem',
    background: 'rgba(239, 68, 68, 0.1)',
    border: '1px solid rgba(239, 68, 68, 0.35)',
    borderRadius: '8px',
    color: 'var(--red)',
    fontSize: '0.8rem',
    cursor: 'pointer',
    transition: 'all 0.15s ease',
  }

  return (
    <div className="page-enter" style={{
      maxWidth: '680px',
      margin: '0 auto',
      padding: '2rem 1.25rem',
    }}>
      {/* Header */}
      <div style={{ marginBottom: '2rem', display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '1rem', flexWrap: 'wrap' }}>
        <div>
          <h1 style={{ fontSize: '1.4rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.2rem' }}>
            Mi Cuenta
          </h1>
          <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>{user.email}</p>
        </div>
        <button onClick={handleLogout} style={btnOutline}>
          Cerrar sesión
        </button>
      </div>

      {fetchError && (
        <div style={{
          background: 'rgba(239, 68, 68, 0.1)',
          border: '1px solid rgba(239, 68, 68, 0.3)',
          borderRadius: '8px',
          padding: '0.75rem 1rem',
          marginBottom: '1.25rem',
          fontSize: '0.85rem',
          color: 'var(--red)',
        }}>
          {fetchError}
        </div>
      )}

      {/* 1. Mi Wallet */}
      <div style={sectionCard}>
        <div style={sectionTitle}>Mi Wallet</div>

        <div style={{ marginBottom: '1rem' }}>
          <div style={{ fontSize: '0.78rem', color: 'var(--text-secondary)', marginBottom: '0.35rem' }}>Dirección</div>
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '0.75rem',
            background: 'var(--bg-secondary)',
            borderRadius: '8px',
            padding: '0.65rem 0.9rem',
            border: '1px solid var(--border)',
          }}>
            <span style={{
              fontFamily: 'monospace',
              fontSize: '0.8rem',
              color: 'var(--text-primary)',
              flex: 1,
              wordBreak: 'break-all',
            }}>
              {address}
            </span>
            <button
              onClick={copyAddress}
              title="Copiar dirección"
              style={{
                background: copied ? 'rgba(16,185,129,0.15)' : 'var(--border)',
                border: 'none',
                borderRadius: '6px',
                padding: '0.35rem 0.7rem',
                color: copied ? 'var(--green)' : 'var(--text-secondary)',
                fontSize: '0.75rem',
                cursor: 'pointer',
                flexShrink: 0,
                transition: 'all 0.15s ease',
              }}
            >
              {copied ? 'Copiado' : 'Copiar'}
            </button>
          </div>
        </div>

        <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
          <div style={{
            flex: 1,
            minWidth: '130px',
            background: 'var(--bg-secondary)',
            borderRadius: '10px',
            padding: '1rem',
            border: '1px solid var(--border)',
          }}>
            <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginBottom: '0.4rem' }}>Balance SPC</div>
            <div style={{ fontSize: '1.3rem', fontWeight: '700', color: 'var(--text-primary)' }}>
              {Number(exchangeSPC ?? balance).toLocaleString('es-ES', { maximumFractionDigits: 4 })}
              <span style={{ fontSize: '0.8rem', color: 'var(--accent)', marginLeft: '0.3rem', fontWeight: '500' }}>SPC</span>
            </div>
          </div>
          <div style={{
            flex: 1,
            minWidth: '130px',
            background: 'var(--bg-secondary)',
            borderRadius: '10px',
            padding: '1rem',
            border: '1px solid var(--border)',
          }}>
            <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginBottom: '0.4rem' }}>Saldo EUR</div>
            <div style={{ fontSize: '1.3rem', fontWeight: '700', color: 'var(--text-primary)' }}>
              {eurBalance !== null ? new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(eurBalance) : '—'}
            </div>
          </div>
          <div style={{
            flex: 1,
            minWidth: '130px',
            background: 'var(--bg-secondary)',
            borderRadius: '10px',
            padding: '1rem',
            border: '1px solid var(--border)',
          }}>
            <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginBottom: '0.4rem' }}>Nonce</div>
            <div style={{ fontSize: '1.3rem', fontWeight: '700', color: 'var(--text-primary)', fontFamily: 'monospace' }}>
              {nonce}
            </div>
          </div>
        </div>

        {/* Quick links */}
        <div style={{ display: 'flex', gap: '0.5rem', marginTop: '1rem' }}>
          <button onClick={() => onNavigate('/wallet')} style={{
            flex: 1, padding: '0.5rem', borderRadius: '8px', border: '1px solid var(--border)',
            background: 'var(--bg-secondary)', color: 'var(--accent)', fontSize: '0.8rem',
            fontWeight: '600', cursor: 'pointer',
          }}>Ver cartera</button>
          <button onClick={() => onNavigate('/trade/SPC')} style={{
            flex: 1, padding: '0.5rem', borderRadius: '8px', border: 'none',
            background: 'var(--accent)', color: '#fff', fontSize: '0.8rem',
            fontWeight: '600', cursor: 'pointer',
          }}>Operar</button>
        </div>
      </div>

      {/* 2. Seguridad */}
      <div style={sectionCard}>
        <div style={sectionTitle}>Seguridad</div>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '1.25rem', lineHeight: '1.6' }}>
          Tu clave privada está cifrada con tu contraseña. Nunca la compartimos.
        </p>
        <button
          onClick={() => { setShowExportModal(true); setExportStep('warning'); setExportPassword('') }}
          style={btnDanger}
        >
          Exportar clave privada
        </button>
      </div>

      {/* 3. Enviar $SPC */}
      <div style={sectionCard}>
        <div style={sectionTitle}>Enviar $SPC</div>
        <div style={{ marginBottom: '1rem' }}>
          <label style={{ display: 'block', fontSize: '0.78rem', color: 'var(--text-secondary)', marginBottom: '0.35rem' }}>
            Dirección destino
          </label>
          <input
            type="text"
            placeholder="SPC1abc..."
            disabled
            style={{
              width: '100%',
              padding: '0.7rem 1rem',
              background: 'var(--bg-secondary)',
              border: '1px solid var(--border)',
              borderRadius: '8px',
              color: 'var(--text-secondary)',
              fontSize: '0.875rem',
              opacity: 0.6,
              cursor: 'not-allowed',
            }}
          />
        </div>
        <div style={{ marginBottom: '1rem' }}>
          <label style={{ display: 'block', fontSize: '0.78rem', color: 'var(--text-secondary)', marginBottom: '0.35rem' }}>
            Cantidad (SPC)
          </label>
          <input
            type="number"
            placeholder="0.00"
            disabled
            style={{
              width: '100%',
              padding: '0.7rem 1rem',
              background: 'var(--bg-secondary)',
              border: '1px solid var(--border)',
              borderRadius: '8px',
              color: 'var(--text-secondary)',
              fontSize: '0.875rem',
              opacity: 0.6,
              cursor: 'not-allowed',
            }}
          />
        </div>
        <div style={{
          background: 'rgba(59, 130, 246, 0.08)',
          border: '1px solid rgba(59, 130, 246, 0.2)',
          borderRadius: '8px',
          padding: '0.75rem 1rem',
          fontSize: '0.82rem',
          color: 'var(--text-secondary)',
          marginBottom: '1rem',
        }}>
          Proximamente — usa el CLI para enviar transacciones:
          <code style={{ display: 'block', marginTop: '0.4rem', color: 'var(--accent)', fontSize: '0.78rem' }}>
            ./spaincoin send --to SPC1... --amount 10
          </code>
        </div>
        <button
          disabled
          style={{
            padding: '0.7rem 1.5rem',
            background: 'var(--border)',
            border: 'none',
            borderRadius: '8px',
            color: 'var(--text-secondary)',
            fontSize: '0.875rem',
            cursor: 'not-allowed',
            opacity: 0.6,
          }}
        >
          Enviar
        </button>
      </div>

      {/* Export Key Modal */}
      {showExportModal && (
        <div style={{
          position: 'fixed',
          inset: 0,
          background: 'rgba(0,0,0,0.7)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 1000,
          padding: '1rem',
        }}>
          <div style={{
            background: 'var(--bg-card)',
            border: '1px solid var(--border)',
            borderRadius: '16px',
            padding: '2rem',
            maxWidth: '420px',
            width: '100%',
          }}>
            {exportStep === 'warning' && (
              <>
                <div style={{ fontSize: '1.1rem', fontWeight: '700', color: 'var(--red)', marginBottom: '1rem' }}>
                  Advertencia
                </div>
                <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: '1.65', marginBottom: '1.5rem' }}>
                  Tu clave privada da acceso total a tus fondos. <strong style={{ color: 'var(--text-primary)' }}>Nunca la compartas</strong> con nadie. Guárdala en un lugar seguro offline.
                </p>
                <div style={{ display: 'flex', gap: '0.75rem' }}>
                  <button
                    onClick={() => setShowExportModal(false)}
                    style={{ ...btnOutline, flex: 1 }}
                  >
                    Cancelar
                  </button>
                  <button
                    onClick={() => setExportStep('password')}
                    style={{
                      flex: 1,
                      padding: '0.55rem 1.1rem',
                      background: 'rgba(239,68,68,0.15)',
                      border: '1px solid rgba(239,68,68,0.4)',
                      borderRadius: '8px',
                      color: 'var(--red)',
                      fontSize: '0.8rem',
                      cursor: 'pointer',
                    }}
                  >
                    Continuar
                  </button>
                </div>
              </>
            )}

            {exportStep === 'password' && (
              <>
                <div style={{ fontSize: '1rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
                  Confirma tu contraseña
                </div>
                <input
                  type="password"
                  value={exportPassword}
                  onChange={(e) => setExportPassword(e.target.value)}
                  placeholder="Tu contraseña"
                  style={{
                    width: '100%',
                    padding: '0.7rem 1rem',
                    background: 'var(--bg-secondary)',
                    border: '1px solid var(--border)',
                    borderRadius: '8px',
                    color: 'var(--text-primary)',
                    fontSize: '0.875rem',
                    marginBottom: '1.25rem',
                  }}
                />
                <div style={{ display: 'flex', gap: '0.75rem' }}>
                  <button
                    onClick={() => setShowExportModal(false)}
                    style={{ ...btnOutline, flex: 1 }}
                  >
                    Cancelar
                  </button>
                  <button
                    onClick={() => setExportStep('result')}
                    style={{
                      flex: 1,
                      padding: '0.55rem 1.1rem',
                      background: 'var(--accent)',
                      border: 'none',
                      borderRadius: '8px',
                      color: '#fff',
                      fontSize: '0.8rem',
                      cursor: 'pointer',
                    }}
                  >
                    Exportar
                  </button>
                </div>
              </>
            )}

            {exportStep === 'result' && (
              <>
                <div style={{ fontSize: '1rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
                  Exportar clave privada
                </div>
                <div style={{
                  background: 'rgba(59,130,246,0.08)',
                  border: '1px solid rgba(59,130,246,0.2)',
                  borderRadius: '8px',
                  padding: '1rem',
                  fontSize: '0.85rem',
                  color: 'var(--text-secondary)',
                  marginBottom: '1.25rem',
                  lineHeight: '1.6',
                }}>
                  Funcion disponible proximamente
                </div>
                <button
                  onClick={() => setShowExportModal(false)}
                  style={{ ...btnOutline, width: '100%' }}
                >
                  Cerrar
                </button>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
