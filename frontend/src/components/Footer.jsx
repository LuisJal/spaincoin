export default function Footer({ onNavigate }) {
  const linkStyle = {
    color: 'var(--text-secondary)',
    fontSize: '0.8rem',
    cursor: 'pointer',
    background: 'none',
    border: 'none',
    padding: 0,
    textDecoration: 'none',
    transition: 'color 0.15s ease',
  }

  function NavLink({ children, path }) {
    return (
      <button
        style={linkStyle}
        onClick={() => onNavigate && onNavigate(path)}
        onMouseEnter={(e) => { e.target.style.color = 'var(--accent)' }}
        onMouseLeave={(e) => { e.target.style.color = 'var(--text-secondary)' }}
      >
        {children}
      </button>
    )
  }

  return (
    <footer style={{
      borderTop: '1px solid var(--border)',
      padding: '1.25rem 1.5rem',
      background: 'var(--bg-primary)',
    }}>
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '0.65rem',
      }}>
        {/* Brand + badge */}
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem' }}>
          <span style={{ fontSize: '0.85rem', fontWeight: '600', color: 'var(--text-primary)' }}>
            SpainCoin Exchange
          </span>
          <span style={{
            fontSize: '0.7rem',
            fontWeight: '600',
            color: '#f59e0b',
            background: 'rgba(245, 158, 11, 0.1)',
            border: '1px solid rgba(245, 158, 11, 0.25)',
            padding: '0.15rem 0.5rem',
            borderRadius: '20px',
            letterSpacing: '0.02em',
          }}>
            Testnet v0.1
          </span>
        </div>

        {/* Legal links */}
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: '0.35rem',
          flexWrap: 'wrap',
          justifyContent: 'center',
        }}>
          <NavLink path="/legal/terms">Términos</NavLink>
          <span style={{ color: 'var(--border-accent)', fontSize: '0.75rem' }}>·</span>
          <NavLink path="/legal/privacy">Privacidad</NavLink>
          <span style={{ color: 'var(--border-accent)', fontSize: '0.75rem' }}>·</span>
          <NavLink path="/legal/risk">Riesgos</NavLink>
          <span style={{ color: 'var(--border-accent)', fontSize: '0.75rem' }}>·</span>
          <NavLink path="/legal/cookies">Cookies</NavLink>
        </div>

        {/* Copyright + disclaimer */}
        <p style={{
          fontSize: '0.75rem',
          color: 'var(--text-secondary)',
          textAlign: 'center',
          lineHeight: '1.5',
          opacity: 0.7,
        }}>
          © 2026 SpainCoin. Proyecto en desarrollo — no es asesoramiento financiero.
        </p>
      </div>
    </footer>
  )
}
