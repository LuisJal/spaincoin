import { useState } from 'react'
import { useAuth } from '../auth/useAuth.jsx'

const styles = {
  nav: {
    position: 'sticky',
    top: 0,
    zIndex: 100,
    background: 'rgba(10, 14, 26, 0.92)',
    backdropFilter: 'blur(12px)',
    borderBottom: '1px solid #1f2937',
    padding: '0 1.5rem',
    height: '60px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: '1rem',
  },
  logo: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.6rem',
    textDecoration: 'none',
    flexShrink: 0,
  },
  logoIcon: {
    width: '32px',
    height: '32px',
    background: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)',
    borderRadius: '8px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    fontWeight: '700',
    fontSize: '0.8rem',
    color: '#fff',
    letterSpacing: '-0.02em',
    flexShrink: 0,
  },
  logoText: {
    fontWeight: '600',
    fontSize: '1rem',
    color: '#f9fafb',
    letterSpacing: '-0.01em',
  },
  logoSub: {
    fontSize: '0.7rem',
    color: '#9ca3af',
    fontWeight: '400',
    display: 'block',
    lineHeight: 1,
  },
  navLinks: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.25rem',
    listStyle: 'none',
  },
  navLink: (active) => ({
    padding: '0.4rem 0.85rem',
    borderRadius: '6px',
    fontSize: '0.875rem',
    fontWeight: active ? '600' : '400',
    color: active ? '#f9fafb' : '#9ca3af',
    background: active ? '#1f2937' : 'transparent',
    cursor: 'pointer',
    transition: 'all 0.15s ease',
    border: 'none',
    textDecoration: 'none',
    display: 'block',
  }),
  liveIndicator: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.4rem',
    fontSize: '0.75rem',
    fontWeight: '500',
    color: '#10b981',
    background: 'rgba(16, 185, 129, 0.1)',
    border: '1px solid rgba(16, 185, 129, 0.25)',
    padding: '0.3rem 0.7rem',
    borderRadius: '20px',
    flexShrink: 0,
  },
}

const navItems = [
  { label: 'Dashboard', hash: '#/' },
  { label: 'Explorer', hash: '#/explorer' },
  { label: 'Wallet', hash: '#/wallet' },
]

export default function Navbar({ currentPage, onNavigate }) {
  const [hovered, setHovered] = useState(null)
  const { user, logout } = useAuth()

  function getActivePage(hash) {
    if (hash === '#/' || hash === '' || hash === '#') return '/'
    return hash.replace('#', '')
  }

  const active = getActivePage(currentPage)

  function handleLogout() {
    logout()
    onNavigate('/')
  }

  // Truncate email for display
  function truncateEmail(email) {
    if (!email) return ''
    const [local, domain] = email.split('@')
    if (local.length <= 8) return email
    return local.slice(0, 6) + '...' + (domain ? '@' + domain : '')
  }

  return (
    <nav style={styles.nav}>
      {/* Logo */}
      <a
        href="#/"
        style={styles.logo}
        onClick={(e) => { e.preventDefault(); onNavigate('/') }}
      >
        <div style={styles.logoIcon}>SPC</div>
        <div>
          <span style={styles.logoText}>SpainCoin</span>
          <span style={styles.logoSub}>Exchange</span>
        </div>
      </a>

      {/* Nav links */}
      <ul style={styles.navLinks}>
        {navItems.map((item) => {
          const page = item.hash.replace('#', '') || '/'
          const isActive = active === page || (page === '/' && active === '/')
          const isHovered = hovered === item.label

          return (
            <li key={item.label}>
              <a
                href={item.hash}
                style={{
                  ...styles.navLink(isActive),
                  ...(isHovered && !isActive ? { color: '#f9fafb', background: '#111827' } : {}),
                }}
                onClick={(e) => { e.preventDefault(); onNavigate(page) }}
                onMouseEnter={() => setHovered(item.label)}
                onMouseLeave={() => setHovered(null)}
              >
                {item.label}
              </a>
            </li>
          )
        })}
      </ul>

      {/* Right side: auth section */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', flexShrink: 0 }}>
        {user ? (
          // Logged in: show email + wallet icon + logout
          <>
            <button
              onClick={() => onNavigate('/account')}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '0.4rem',
                padding: '0.35rem 0.75rem',
                background: 'var(--bg-secondary)',
                border: '1px solid var(--border)',
                borderRadius: '8px',
                color: 'var(--text-secondary)',
                fontSize: '0.8rem',
                cursor: 'pointer',
                transition: 'all 0.15s ease',
                maxWidth: '180px',
              }}
              title={user.email}
            >
              {/* Wallet icon */}
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="2" y="5" width="20" height="14" rx="2"/>
                <path d="M16 12h.01"/>
                <path d="M2 10h20"/>
              </svg>
              <span style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                {truncateEmail(user.email)}
              </span>
            </button>
            <button
              onClick={handleLogout}
              style={{
                padding: '0.35rem 0.75rem',
                background: 'transparent',
                border: '1px solid var(--border)',
                borderRadius: '8px',
                color: 'var(--text-secondary)',
                fontSize: '0.8rem',
                cursor: 'pointer',
                transition: 'all 0.15s ease',
                flexShrink: 0,
              }}
            >
              Cerrar sesión
            </button>
          </>
        ) : (
          // Not logged in: show live indicator + login button
          <>
            <div style={styles.liveIndicator}>
              <span className="pulse-dot" style={{ fontSize: '0.6rem' }}>●</span>
              LIVE
            </div>
            <button
              onClick={() => onNavigate('/login')}
              style={{
                padding: '0.35rem 0.9rem',
                background: 'transparent',
                border: '1px solid var(--accent)',
                borderRadius: '8px',
                color: 'var(--accent)',
                fontSize: '0.85rem',
                fontWeight: '500',
                cursor: 'pointer',
                transition: 'all 0.15s ease',
                flexShrink: 0,
              }}
            >
              Entrar
            </button>
          </>
        )}
      </div>
    </nav>
  )
}
