import { useState, useEffect } from 'react'
import Navbar from './components/Navbar.jsx'
import Dashboard from './pages/Dashboard.jsx'
import Explorer from './pages/Explorer.jsx'
import BlockDetail from './pages/BlockDetail.jsx'
import Wallet from './pages/Wallet.jsx'

/**
 * Parse the current window.location.hash into a route object.
 * Supported hashes:
 *   #/            → { page: '/' }
 *   #/explorer    → { page: '/explorer' }
 *   #/wallet      → { page: '/wallet' }
 *   #/block/42    → { page: '/block', param: 42 }
 */
function getPageFromHash() {
  const hash = window.location.hash || '#/'
  // strip leading #
  const path = hash.startsWith('#') ? hash.slice(1) : hash

  if (!path || path === '/') return { page: '/' }

  const blockMatch = path.match(/^\/block\/(\d+)$/)
  if (blockMatch) return { page: '/block', param: Number(blockMatch[1]) }

  if (path === '/explorer') return { page: '/explorer' }
  if (path === '/wallet') return { page: '/wallet' }

  return { page: '/' }
}

function navigate(path) {
  window.location.hash = '#' + path
}

export default function App() {
  const [route, setRoute] = useState(getPageFromHash)

  useEffect(() => {
    function onHashChange() {
      setRoute(getPageFromHash())
    }
    window.addEventListener('hashchange', onHashChange)
    return () => window.removeEventListener('hashchange', onHashChange)
  }, [])

  function handleNavigate(path) {
    navigate(path)
    setRoute(getPageFromHash())
  }

  function renderPage() {
    switch (route.page) {
      case '/':
        return <Dashboard onNavigate={handleNavigate} />
      case '/explorer':
        return <Explorer onNavigate={handleNavigate} />
      case '/wallet':
        return <Wallet />
      case '/block':
        return <BlockDetail height={route.param} onNavigate={handleNavigate} />
      default:
        return <Dashboard onNavigate={handleNavigate} />
    }
  }

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column', background: 'var(--bg-primary)' }}>
      <Navbar currentPage={'#' + route.page} onNavigate={handleNavigate} />
      <main style={{ flex: 1 }}>
        {renderPage()}
      </main>
      <footer style={{
        borderTop: '1px solid var(--border)',
        padding: '1rem 1.5rem',
        textAlign: 'center',
        fontSize: '0.75rem',
        color: 'var(--text-secondary)',
      }}>
        SpainCoin Exchange · $SPC · Layer 1 Blockchain · Built with React + Vite
      </footer>
    </div>
  )
}
