import { useState, useEffect } from 'react'

const STORAGE_KEY = 'spc_cookies_accepted'

export default function CookieBanner({ onNavigate }) {
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    const accepted = localStorage.getItem(STORAGE_KEY)
    if (!accepted) {
      // Small delay so it doesn't pop in immediately with the page
      const timer = setTimeout(() => setVisible(true), 600)
      return () => clearTimeout(timer)
    }
  }, [])

  function handleAccept() {
    localStorage.setItem(STORAGE_KEY, 'true')
    setVisible(false)
  }

  if (!visible) return null

  return (
    <div style={{
      position: 'fixed',
      bottom: 0,
      left: 0,
      right: 0,
      zIndex: 999,
      padding: '0 1rem 1rem',
      pointerEvents: 'none',
    }}>
      <div style={{
        maxWidth: '720px',
        margin: '0 auto',
        background: 'var(--bg-card)',
        border: '1px solid var(--border-accent)',
        borderRadius: '12px',
        padding: '1rem 1.25rem',
        display: 'flex',
        alignItems: 'center',
        gap: '1rem',
        flexWrap: 'wrap',
        boxShadow: '0 -4px 32px rgba(0,0,0,0.4)',
        pointerEvents: 'all',
      }}>
        {/* Cookie icon */}
        <div style={{
          flexShrink: 0,
          width: '36px',
          height: '36px',
          background: 'rgba(59,130,246,0.1)',
          border: '1px solid rgba(59,130,246,0.2)',
          borderRadius: '8px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '1.1rem',
        }}>
          🍪
        </div>

        {/* Text */}
        <div style={{ flex: 1, minWidth: '200px' }}>
          <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: '1.5', margin: 0 }}>
            Usamos cookies técnicas necesarias para el funcionamiento del servicio. No usamos
            cookies publicitarias.{' '}
            <button
              onClick={() => onNavigate && onNavigate('/legal/cookies')}
              style={{
                background: 'none',
                border: 'none',
                color: 'var(--accent)',
                cursor: 'pointer',
                fontSize: 'inherit',
                padding: 0,
                textDecoration: 'underline',
              }}
            >
              Más información
            </button>
          </p>
        </div>

        {/* Accept button */}
        <button
          onClick={handleAccept}
          style={{
            flexShrink: 0,
            padding: '0.5rem 1.25rem',
            background: 'var(--accent)',
            border: 'none',
            borderRadius: '8px',
            color: '#fff',
            fontSize: '0.875rem',
            fontWeight: '600',
            cursor: 'pointer',
            transition: 'background 0.15s ease',
            whiteSpace: 'nowrap',
          }}
          onMouseEnter={(e) => { e.target.style.background = 'var(--accent-hover)' }}
          onMouseLeave={(e) => { e.target.style.background = 'var(--accent)' }}
        >
          Entendido
        </button>
      </div>
    </div>
  )
}
