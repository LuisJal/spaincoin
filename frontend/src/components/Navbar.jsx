import { useState } from 'react'

const navItems = [
  { label: 'Inicio', page: '/' },
  { label: 'Explorer', page: '/explorer' },
  { label: 'Mercado', page: '/market' },
  { label: 'Wallet', page: '/wallet' },
  { label: 'Validadores', page: '/validators' },
]

function Logo({ size = 32 }) {
  return <img src="/logo.jpeg" alt="SpainCoin" style={{ width: size, height: size, borderRadius: '6px', objectFit: 'contain' }} />
}

export default function Navbar({ currentPage, onNavigate }) {
  const [menuOpen, setMenuOpen] = useState(false)

  function getActivePage(hash) {
    if (!hash || hash === '#/' || hash === '#') return '/'
    return hash.replace('#', '')
  }

  const active = getActivePage(currentPage)

  function handleNavigate(page) {
    onNavigate(page)
    setMenuOpen(false)
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
      position: 'sticky', top: 0, zIndex: 100,
      background: 'rgba(10, 14, 26, 0.95)', backdropFilter: 'blur(12px)',
      borderBottom: '1px solid #1f2937', padding: '0 1.25rem',
      height: '60px', display: 'flex', alignItems: 'center', justifyContent: 'space-between',
    }}>
      <a href="#/" style={{ display: 'flex', alignItems: 'center', gap: '0.6rem', textDecoration: 'none', flexShrink: 0 }}
        onClick={(e) => { e.preventDefault(); handleNavigate('/') }}>
        <Logo size={36} />
        <div>
          <span style={{ fontWeight: '600', fontSize: '1rem', color: '#f9fafb' }}>SpainCoin</span>
          <span style={{ fontSize: '0.7rem', color: '#9ca3af', display: 'block', lineHeight: 1 }}>Blockchain</span>
        </div>
      </a>

      <ul style={{ display: 'flex', alignItems: 'center', gap: '0.25rem', listStyle: 'none', margin: 0, padding: 0 }} className="desktop-nav">
        {navItems.map((item) => (
          <li key={item.label}>
            <button style={navLinkStyle(active === item.page)} onClick={() => handleNavigate(item.page)}>
              {item.label}
            </button>
          </li>
        ))}
      </ul>

      {/* LIVE badge */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', flexShrink: 0 }} className="desktop-auth">
        <div style={{
          display: 'flex', alignItems: 'center', gap: '0.4rem',
          fontSize: '0.75rem', fontWeight: '500', color: '#10b981',
          background: 'rgba(16, 185, 129, 0.1)', border: '1px solid rgba(16, 185, 129, 0.25)',
          padding: '0.3rem 0.7rem', borderRadius: '20px',
        }}>
          <span style={{ fontSize: '0.6rem' }}>●</span>LIVE
        </div>
      </div>

      {/* Hamburger */}
      <button className="hamburger" onClick={() => setMenuOpen(!menuOpen)}
        style={{ display: 'none', background: 'none', border: 'none', color: '#9ca3af', cursor: 'pointer', padding: '0.4rem', borderRadius: '6px', flexDirection: 'column', gap: '4px', alignItems: 'center', justifyContent: 'center' }}
        aria-label="Menu">
        <span style={{ display: 'block', width: '20px', height: '2px', background: menuOpen ? '#f9fafb' : '#9ca3af', transition: 'all 0.2s', transform: menuOpen ? 'rotate(45deg) translate(4px, 4px)' : 'none' }} />
        <span style={{ display: 'block', width: '20px', height: '2px', background: menuOpen ? 'transparent' : '#9ca3af', transition: 'all 0.2s' }} />
        <span style={{ display: 'block', width: '20px', height: '2px', background: menuOpen ? '#f9fafb' : '#9ca3af', transition: 'all 0.2s', transform: menuOpen ? 'rotate(-45deg) translate(4px, -4px)' : 'none' }} />
      </button>

      {menuOpen && (
        <div style={{
          position: 'absolute', top: '60px', left: 0, right: 0,
          background: 'rgba(10, 14, 26, 0.98)', backdropFilter: 'blur(12px)',
          borderBottom: '1px solid #1f2937', padding: '0.75rem 1.25rem 1rem',
          display: 'flex', flexDirection: 'column', gap: '0.25rem', zIndex: 99,
        }}>
          {navItems.map((item) => (
            <button key={item.label}
              style={{ ...navLinkStyle(active === item.page), textAlign: 'left', width: '100%', padding: '0.65rem 0.85rem' }}
              onClick={() => handleNavigate(item.page)}>
              {item.label}
            </button>
          ))}
        </div>
      )}
    </nav>
  )
}
