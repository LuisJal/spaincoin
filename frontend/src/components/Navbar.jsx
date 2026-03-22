import { useState, useEffect } from 'react'

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

  function getActivePage(hash) {
    if (hash === '#/' || hash === '' || hash === '#') return '/'
    return hash.replace('#', '')
  }

  const active = getActivePage(currentPage)

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

      {/* Live indicator */}
      <div style={styles.liveIndicator}>
        <span className="pulse-dot" style={{ fontSize: '0.6rem' }}>●</span>
        LIVE
      </div>
    </nav>
  )
}
