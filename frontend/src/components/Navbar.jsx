import { useState } from 'react'
import { useAuth } from '../auth/useAuth.jsx'

const navItems = [
  { label: 'Dashboard', page: '/' },
  { label: 'Trading', page: '/trade' },
  { label: 'Mercado', page: '/market' },
  { label: 'Explorer', page: '/explorer' },
  { label: 'Cartera', page: '/wallet' },
]

function LogoSVG({ size = 32 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="cg" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0%" stopColor="#ffc400"/><stop offset="100%" stopColor="#e6a800"/>
        </linearGradient>
      </defs>
      <circle cx="24" cy="24" r="22" fill="url(#cg)" stroke="#b8860b" strokeWidth="1.5"/>
      <circle cx="24" cy="24" r="18" fill="none" stroke="#b8860b" strokeWidth="0.8" opacity="0.5"/>
      <rect x="10" y="10" width="28" height="5" rx="1" fill="#c60b1e" opacity="0.85"/>
      <rect x="10" y="33" width="28" height="5" rx="1" fill="#c60b1e" opacity="0.85"/>
      <text x="24" y="29" textAnchor="middle" fontFamily="Georgia, serif" fontWeight="700" fontSize="18" fill="#c60b1e" opacity="0.9">S</text>
    </svg>
  )
}

export default function Navbar({ currentPage, onNavigate }) {
  const [menuOpen, setMenuOpen] = useState(false)
  const { user, logout } = useAuth()

  function getActivePage(hash) {
    if (!hash || hash === '#/' || hash === '#') return '/'
    return hash.replace('#', '')
  }

  const active = getActivePage(currentPage)

  function handleLogout() {
    logout()
    onNavigate('/')
    setMenuOpen(false)
  }

  function handleNavigate(page) {
    onNavigate(page)
    setMenuOpen(false)
  }

  function truncateEmail(email) {
    if (!email) return ''
    const [local, domain] = email.split('@')
    if (local.length <= 8) return email
    return local.slice(0, 6) + '...' + (domain ? '@' + domain : '')
  }

  const navLinkStyle = (isActive) => ({
    padding: '0.4rem 0.85rem',
    borderRadius: '6px',
    fontSize: '0.875rem',
    fontWeight: isActive ? '600' : '400',
    color: isActive ? '#f9fafb' : '#9ca3af',
    background: isActive ? '#1f2937' : 'transparent',
    cursor: 'pointer',
    border: 'none',
    textDecoration: 'none',
    display: 'block',
    transition: 'all 0.15s ease',
    whiteSpace: 'nowrap',
  })

  return (
    <nav style={{
      position: 'sticky',
      top: 0,
      zIndex: 100,
      background: 'rgba(10, 14, 26, 0.95)',
      backdropFilter: 'blur(12px)',
      borderBottom: '1px solid #1f2937',
      padding: '0 1.25rem',
      height: '60px',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
    }}>

      {/* Logo */}
      <a
        href="#/"
        style={{ display: 'flex', alignItems: 'center', gap: '0.6rem', textDecoration: 'none', flexShrink: 0 }}
        onClick={(e) => { e.preventDefault(); handleNavigate('/') }}
      >
        <LogoSVG size={32} />
        <div>
          <span style={{ fontWeight: '600', fontSize: '1rem', color: '#f9fafb', letterSpacing: '-0.01em' }}>SpainCoin</span>
          <span style={{ fontSize: '0.7rem', color: '#9ca3af', display: 'block', lineHeight: 1 }}>Exchange</span>
        </div>
      </a>

      {/* Desktop nav links */}
      <ul style={{
        display: 'flex', alignItems: 'center', gap: '0.25rem',
        listStyle: 'none', margin: 0, padding: 0,
        // Ocultar en móvil via media query no disponible en inline styles
        // Lo manejamos con el menú hamburguesa
      }} className="desktop-nav">
        {navItems.map((item) => (
          <li key={item.label}>
            <button
              style={navLinkStyle(active === item.page)}
              onClick={() => handleNavigate(item.page)}
            >
              {item.label}
            </button>
          </li>
        ))}
      </ul>

      {/* Desktop auth — oculto en móvil */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', flexShrink: 0 }} className="desktop-auth">
        {user ? (
          <>
            <button
              onClick={() => handleNavigate('/account')}
              style={{
                display: 'flex', alignItems: 'center', gap: '0.4rem',
                padding: '0.35rem 0.75rem',
                background: 'var(--bg-secondary)', border: '1px solid var(--border)',
                borderRadius: '8px', color: 'var(--text-secondary)',
                fontSize: '0.8rem', cursor: 'pointer', maxWidth: '160px',
              }}
              title={user.email}
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="2" y="5" width="20" height="14" rx="2"/>
                <path d="M16 12h.01"/><path d="M2 10h20"/>
              </svg>
              <span style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                {truncateEmail(user.email)}
              </span>
            </button>
            <button
              onClick={handleLogout}
              style={{
                padding: '0.35rem 0.75rem', background: 'transparent',
                border: '1px solid var(--border)', borderRadius: '8px',
                color: 'var(--text-secondary)', fontSize: '0.8rem', cursor: 'pointer',
              }}
            >
              Salir
            </button>
          </>
        ) : (
          <>
            <div style={{
              display: 'flex', alignItems: 'center', gap: '0.4rem',
              fontSize: '0.75rem', fontWeight: '500', color: '#10b981',
              background: 'rgba(16, 185, 129, 0.1)', border: '1px solid rgba(16, 185, 129, 0.25)',
              padding: '0.3rem 0.7rem', borderRadius: '20px',
            }}>
              <span style={{ fontSize: '0.6rem' }}>●</span>
              LIVE
            </div>
            <button
              onClick={() => handleNavigate('/login')}
              style={{
                padding: '0.35rem 0.9rem', background: 'transparent',
                border: '1px solid var(--accent)', borderRadius: '8px',
                color: 'var(--accent)', fontSize: '0.85rem', fontWeight: '500', cursor: 'pointer',
              }}
            >
              Entrar
            </button>
          </>
        )}
      </div>

      {/* Hamburger button — solo visible en móvil */}
      <button
        className="hamburger"
        onClick={() => setMenuOpen(!menuOpen)}
        style={{
          display: 'none', // se muestra via CSS
          background: 'none', border: 'none',
          color: '#9ca3af', cursor: 'pointer',
          padding: '0.4rem', borderRadius: '6px',
          flexDirection: 'column', gap: '4px',
          alignItems: 'center', justifyContent: 'center',
        }}
        aria-label="Menú"
      >
        <span style={{ display: 'block', width: '20px', height: '2px', background: menuOpen ? '#f9fafb' : '#9ca3af', transition: 'all 0.2s', transform: menuOpen ? 'rotate(45deg) translate(4px, 4px)' : 'none' }} />
        <span style={{ display: 'block', width: '20px', height: '2px', background: menuOpen ? 'transparent' : '#9ca3af', transition: 'all 0.2s' }} />
        <span style={{ display: 'block', width: '20px', height: '2px', background: menuOpen ? '#f9fafb' : '#9ca3af', transition: 'all 0.2s', transform: menuOpen ? 'rotate(-45deg) translate(4px, -4px)' : 'none' }} />
      </button>

      {/* Mobile menu dropdown */}
      {menuOpen && (
        <div style={{
          position: 'absolute', top: '60px', left: 0, right: 0,
          background: 'rgba(10, 14, 26, 0.98)', backdropFilter: 'blur(12px)',
          borderBottom: '1px solid #1f2937',
          padding: '0.75rem 1.25rem 1rem',
          display: 'flex', flexDirection: 'column', gap: '0.25rem',
          zIndex: 99,
        }}>
          {navItems.map((item) => (
            <button
              key={item.label}
              style={{
                ...navLinkStyle(active === item.page),
                textAlign: 'left', width: '100%', padding: '0.65rem 0.85rem',
              }}
              onClick={() => handleNavigate(item.page)}
            >
              {item.label}
            </button>
          ))}

          <div style={{ borderTop: '1px solid #1f2937', marginTop: '0.5rem', paddingTop: '0.75rem' }}>
            {user ? (
              <>
                <button
                  style={{ ...navLinkStyle(active === '/account'), textAlign: 'left', width: '100%', padding: '0.65rem 0.85rem', marginBottom: '0.25rem' }}
                  onClick={() => handleNavigate('/account')}
                >
                  Mi cuenta · {truncateEmail(user.email)}
                </button>
                <button
                  style={{ ...navLinkStyle(false), textAlign: 'left', width: '100%', padding: '0.65rem 0.85rem', color: '#ef4444' }}
                  onClick={handleLogout}
                >
                  Cerrar sesión
                </button>
              </>
            ) : (
              <button
                style={{
                  width: '100%', padding: '0.65rem', background: 'var(--accent)',
                  border: 'none', borderRadius: '8px', color: '#fff',
                  fontSize: '0.9rem', fontWeight: '600', cursor: 'pointer',
                }}
                onClick={() => handleNavigate('/login')}
              >
                Entrar
              </button>
            )}
          </div>
        </div>
      )}
    </nav>
  )
}
