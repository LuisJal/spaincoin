import { useState } from 'react'
import { useAuth } from '../auth/useAuth.jsx'
import { login, register } from '../api/client.js'

export default function Login({ onNavigate }) {
  const { saveAuth } = useAuth()
  const [tab, setTab] = useState('login') // 'login' | 'register'

  // Form fields
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [acceptTerms, setAcceptTerms] = useState(false)
  const [acceptPrivacy, setAcceptPrivacy] = useState(false)
  const [acceptRisk, setAcceptRisk] = useState(false)
  const [importKey, setImportKey] = useState('')
  const [showImport, setShowImport] = useState(false)

  // UI state
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [successData, setSuccessData] = useState(null) // {address} after register

  function validateEmail(e) {
    return e.includes('@')
  }

  function validatePassword(p) {
    return p.length >= 8
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError(null)

    // Validation
    if (!validateEmail(email)) {
      setError('El email debe contener @')
      return
    }
    if (!validatePassword(password)) {
      setError('La contraseña debe tener al menos 8 caracteres')
      return
    }
    if (tab === 'register') {
      if (password !== confirmPassword) {
        setError('Las contraseñas no coinciden')
        return
      }
      if (!acceptTerms || !acceptPrivacy || !acceptRisk) {
        setError('Debes aceptar todos los términos para continuar')
        return
      }
    }

    setLoading(true)
    try {
      let data
      if (tab === 'login') {
        data = await login(email, password)
        saveAuth(data)
        onNavigate('/account')
      } else {
        data = await register(email, password, importKey || undefined)
        setSuccessData({ address: data.address })
        saveAuth(data)
        setTimeout(() => {
          onNavigate('/account')
        }, 3000)
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const inputStyle = {
    width: '100%',
    padding: '0.75rem 1rem',
    background: 'var(--bg-secondary)',
    border: '1px solid var(--border)',
    borderRadius: '8px',
    color: 'var(--text-primary)',
    fontSize: '0.875rem',
    outline: 'none',
    transition: 'border-color 0.15s ease',
  }

  const labelStyle = {
    display: 'block',
    fontSize: '0.8rem',
    fontWeight: '500',
    color: 'var(--text-secondary)',
    marginBottom: '0.4rem',
    textTransform: 'uppercase',
    letterSpacing: '0.05em',
  }

  const btnPrimary = {
    width: '100%',
    padding: '0.8rem',
    background: loading ? 'var(--border-accent)' : 'var(--accent)',
    border: 'none',
    borderRadius: '8px',
    color: '#fff',
    fontSize: '0.9rem',
    fontWeight: '600',
    cursor: loading ? 'not-allowed' : 'pointer',
    transition: 'background 0.15s ease',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    gap: '0.5rem',
  }

  return (
    <div className="page-enter" style={{
      minHeight: 'calc(100vh - 60px)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '2rem 1rem',
    }}>
      <div style={{
        width: '100%',
        maxWidth: '420px',
      }}>
        {/* Logo */}
        <div style={{ textAlign: 'center', marginBottom: '2rem' }}>
          <div style={{
            width: '56px',
            height: '56px',
            background: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)',
            borderRadius: '14px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontWeight: '700',
            fontSize: '1.1rem',
            color: '#fff',
            letterSpacing: '-0.02em',
            margin: '0 auto 1rem',
          }}>SPC</div>
          <h1 style={{ fontSize: '1.4rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.25rem' }}>
            SpainCoin Exchange
          </h1>
          <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
            Tu acceso a la red $SPC
          </p>
        </div>

        {/* Card */}
        <div style={{
          background: 'var(--bg-card)',
          border: '1px solid var(--border)',
          borderRadius: '16px',
          padding: '2rem',
        }}>
          {/* Tabs */}
          <div style={{
            display: 'flex',
            background: 'var(--bg-secondary)',
            borderRadius: '10px',
            padding: '4px',
            marginBottom: '1.75rem',
          }}>
            {['login', 'register'].map((t) => (
              <button
                key={t}
                onClick={() => { setTab(t); setError(null); setSuccessData(null) }}
                style={{
                  flex: 1,
                  padding: '0.55rem',
                  borderRadius: '8px',
                  border: 'none',
                  background: tab === t ? 'var(--bg-card)' : 'transparent',
                  color: tab === t ? 'var(--text-primary)' : 'var(--text-secondary)',
                  fontSize: '0.875rem',
                  fontWeight: tab === t ? '600' : '400',
                  cursor: 'pointer',
                  transition: 'all 0.15s ease',
                  boxShadow: tab === t ? '0 1px 4px rgba(0,0,0,0.3)' : 'none',
                }}
              >
                {t === 'login' ? 'Iniciar sesión' : 'Crear cuenta'}
              </button>
            ))}
          </div>

          {/* Success box (register only) */}
          {successData && (
            <div style={{
              background: 'rgba(16, 185, 129, 0.1)',
              border: '1px solid rgba(16, 185, 129, 0.35)',
              borderRadius: '10px',
              padding: '1rem 1.25rem',
              marginBottom: '1.25rem',
            }}>
              <div style={{ color: 'var(--green)', fontWeight: '700', fontSize: '0.95rem', marginBottom: '0.5rem' }}>
                Cuenta creada
              </div>
              <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', marginBottom: '0.4rem' }}>
                Tu dirección SpainCoin:
              </div>
              <div style={{
                fontFamily: 'monospace',
                fontSize: '0.8rem',
                color: 'var(--text-primary)',
                background: 'var(--bg-secondary)',
                borderRadius: '6px',
                padding: '0.5rem 0.75rem',
                wordBreak: 'break-all',
                marginBottom: '0.5rem',
              }}>
                {successData.address}
              </div>
              <div style={{ fontSize: '0.78rem', color: 'var(--text-secondary)' }}>
                Guarda esta dirección — es tu identidad en la red.
              </div>
              <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginTop: '0.5rem', opacity: 0.7 }}>
                Redirigiendo a tu cuenta en 3 segundos...
              </div>
            </div>
          )}

          {/* Form */}
          {!successData && (
            <form onSubmit={handleSubmit} noValidate>
              {/* Email */}
              <div style={{ marginBottom: '1.25rem' }}>
                <label style={labelStyle}>Email</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="tu@email.com"
                  required
                  style={inputStyle}
                  onFocus={(e) => e.target.style.borderColor = 'var(--accent)'}
                  onBlur={(e) => e.target.style.borderColor = 'var(--border)'}
                />
              </div>

              {/* Password */}
              <div style={{ marginBottom: tab === 'register' ? '1.25rem' : '1rem' }}>
                <label style={labelStyle}>Contraseña</label>
                <div style={{ position: 'relative' }}>
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Mínimo 8 caracteres"
                    required
                    style={{ ...inputStyle, paddingRight: '3rem' }}
                    onFocus={(e) => e.target.style.borderColor = 'var(--accent)'}
                    onBlur={(e) => e.target.style.borderColor = 'var(--border)'}
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    style={{
                      position: 'absolute',
                      right: '0.75rem',
                      top: '50%',
                      transform: 'translateY(-50%)',
                      background: 'none',
                      border: 'none',
                      color: 'var(--text-secondary)',
                      cursor: 'pointer',
                      fontSize: '0.85rem',
                      padding: '0.2rem',
                    }}
                  >
                    {showPassword ? 'Ocultar' : 'Ver'}
                  </button>
                </div>
              </div>

              {/* Confirm password (register only) */}
              {tab === 'register' && (
                <>
                  <div style={{ marginBottom: '1.25rem' }}>
                    <label style={labelStyle}>Confirmar contraseña</label>
                    <div style={{ position: 'relative' }}>
                      <input
                        type={showConfirmPassword ? 'text' : 'password'}
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        placeholder="Repite tu contraseña"
                        required
                        style={{ ...inputStyle, paddingRight: '3rem' }}
                        onFocus={(e) => e.target.style.borderColor = 'var(--accent)'}
                        onBlur={(e) => e.target.style.borderColor = 'var(--border)'}
                      />
                      <button
                        type="button"
                        onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                        style={{
                          position: 'absolute',
                          right: '0.75rem',
                          top: '50%',
                          transform: 'translateY(-50%)',
                          background: 'none',
                          border: 'none',
                          color: 'var(--text-secondary)',
                          cursor: 'pointer',
                          fontSize: '0.85rem',
                          padding: '0.2rem',
                        }}
                      >
                        {showConfirmPassword ? 'Ocultar' : 'Ver'}
                      </button>
                    </div>
                  </div>

                  {/* Legal checkboxes */}
                  <div style={{ marginBottom: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                    {/* Checkbox 1: Terms */}
                    <label style={{
                      display: 'flex',
                      alignItems: 'flex-start',
                      gap: '0.6rem',
                      cursor: 'pointer',
                      fontSize: '0.8rem',
                      color: 'var(--text-secondary)',
                      lineHeight: '1.5',
                    }}>
                      <input
                        type="checkbox"
                        checked={acceptTerms}
                        onChange={(e) => setAcceptTerms(e.target.checked)}
                        style={{ marginTop: '2px', accentColor: 'var(--accent)', flexShrink: 0 }}
                      />
                      <span>
                        He leído y acepto los{' '}
                        <button
                          type="button"
                          onClick={() => onNavigate && onNavigate('/legal/terms')}
                          style={{ background: 'none', border: 'none', color: 'var(--accent)', cursor: 'pointer', fontSize: 'inherit', padding: 0, textDecoration: 'underline' }}
                        >
                          Términos y Condiciones
                        </button>
                      </span>
                    </label>

                    {/* Checkbox 2: Privacy */}
                    <label style={{
                      display: 'flex',
                      alignItems: 'flex-start',
                      gap: '0.6rem',
                      cursor: 'pointer',
                      fontSize: '0.8rem',
                      color: 'var(--text-secondary)',
                      lineHeight: '1.5',
                    }}>
                      <input
                        type="checkbox"
                        checked={acceptPrivacy}
                        onChange={(e) => setAcceptPrivacy(e.target.checked)}
                        style={{ marginTop: '2px', accentColor: 'var(--accent)', flexShrink: 0 }}
                      />
                      <span>
                        He leído la{' '}
                        <button
                          type="button"
                          onClick={() => onNavigate && onNavigate('/legal/privacy')}
                          style={{ background: 'none', border: 'none', color: 'var(--accent)', cursor: 'pointer', fontSize: 'inherit', padding: 0, textDecoration: 'underline' }}
                        >
                          Política de Privacidad
                        </button>
                        {' '}y consiento el tratamiento de mis datos
                      </span>
                    </label>

                    {/* Checkbox 3: Risk */}
                    <label style={{
                      display: 'flex',
                      alignItems: 'flex-start',
                      gap: '0.6rem',
                      cursor: 'pointer',
                      fontSize: '0.8rem',
                      color: 'var(--text-secondary)',
                      lineHeight: '1.5',
                    }}>
                      <input
                        type="checkbox"
                        checked={acceptRisk}
                        onChange={(e) => setAcceptRisk(e.target.checked)}
                        style={{ marginTop: '2px', accentColor: 'var(--accent)', flexShrink: 0 }}
                      />
                      <span>
                        Entiendo y acepto los{' '}
                        <button
                          type="button"
                          onClick={() => onNavigate && onNavigate('/legal/risk')}
                          style={{ background: 'none', border: 'none', color: 'var(--accent)', cursor: 'pointer', fontSize: 'inherit', padding: 0, textDecoration: 'underline' }}
                        >
                          riesgos de los activos digitales
                        </button>
                      </span>
                    </label>
                  </div>

                  {/* Importar wallet existente */}
                  <div style={{ marginBottom: '1.5rem' }}>
                    <button
                      type="button"
                      onClick={() => setShowImport(!showImport)}
                      style={{
                        background: 'none', border: 'none', padding: 0,
                        color: 'var(--accent)', fontSize: '0.8rem',
                        cursor: 'pointer', textDecoration: 'underline',
                      }}
                    >
                      {showImport ? '▾ Ocultar importación de wallet' : '▸ ¿Tienes una wallet existente? Impórtala'}
                    </button>
                    {showImport && (
                      <div style={{ marginTop: '0.75rem' }}>
                        <div style={{ fontSize: '0.75rem', color: '#f59e0b', background: 'rgba(245,158,11,0.1)', border: '1px solid rgba(245,158,11,0.25)', borderRadius: '8px', padding: '0.6rem 0.75rem', marginBottom: '0.5rem' }}>
                          ⚠️ Solo introduce tu clave privada en dispositivos de confianza. Se cifrará con tu contraseña antes de guardarse.
                        </div>
                        <input
                          type="password"
                          value={importKey}
                          onChange={e => setImportKey(e.target.value)}
                          placeholder="Clave privada (hex)..."
                          style={{ ...inputStyle, fontFamily: 'monospace', fontSize: '0.8rem' }}
                        />
                      </div>
                    )}
                  </div>
                </>
              )}

              {/* Forgot password (login only) */}
              {tab === 'login' && (
                <div style={{ textAlign: 'right', marginBottom: '1.25rem' }}>
                  <span
                    style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', cursor: 'default' }}
                    title="Contacta soporte"
                  >
                    ¿Olvidaste tu contraseña? <span style={{ color: 'var(--accent)' }}>Contacta soporte</span>
                  </span>
                </div>
              )}

              {/* Error */}
              {error && (
                <div style={{
                  background: 'rgba(239, 68, 68, 0.1)',
                  border: '1px solid rgba(239, 68, 68, 0.35)',
                  borderRadius: '8px',
                  padding: '0.75rem 1rem',
                  marginBottom: '1.25rem',
                  fontSize: '0.85rem',
                  color: 'var(--red)',
                }}>
                  {error}
                </div>
              )}

              {/* Submit */}
              <button type="submit" style={btnPrimary} disabled={loading}>
                {loading ? (
                  <>
                    <span className="spinner" style={{ width: '16px', height: '16px' }} />
                    {tab === 'login' ? 'Iniciando sesión...' : 'Creando cuenta...'}
                  </>
                ) : (
                  tab === 'login' ? 'Iniciar sesión' : 'Crear cuenta'
                )}
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  )
}
