import { useState, useEffect } from 'react'
import { AuthProvider } from './auth/useAuth.jsx'
import Navbar from './components/Navbar.jsx'
import Footer from './components/Footer.jsx'
import CookieBanner from './components/CookieBanner.jsx'
import Dashboard from './pages/Dashboard.jsx'
import Explorer from './pages/Explorer.jsx'
import BlockDetail from './pages/BlockDetail.jsx'
import Wallet from './pages/Wallet.jsx'
import Login from './pages/Login.jsx'
import Account from './pages/Account.jsx'
import Trade from './pages/Trade.jsx'
import Market from './pages/Market.jsx'
import Terms from './pages/legal/Terms.jsx'
import Privacy from './pages/legal/Privacy.jsx'
import Risk from './pages/legal/Risk.jsx'
import Cookies from './pages/legal/Cookies.jsx'

/**
 * Parse the current window.location.hash into a route object.
 * Supported hashes:
 *   #/            → { page: '/' }
 *   #/explorer    → { page: '/explorer' }
 *   #/wallet      → { page: '/wallet' }
 *   #/login       → { page: '/login' }
 *   #/account     → { page: '/account' }
 *   #/block/42    → { page: '/block', param: 42 }
 *   #/legal/terms    → { page: '/legal/terms' }
 *   #/legal/privacy  → { page: '/legal/privacy' }
 *   #/legal/risk     → { page: '/legal/risk' }
 *   #/legal/cookies  → { page: '/legal/cookies' }
 */
function getPageFromHash() {
  const hash = window.location.hash || '#/'
  // strip leading #
  const path = hash.startsWith('#') ? hash.slice(1) : hash

  if (!path || path === '/') return { page: '/' }

  const blockMatch = path.match(/^\/block\/(\d+)$/)
  if (blockMatch) return { page: '/block', param: Number(blockMatch[1]) }

  if (path === '/explorer') return { page: '/explorer' }
  if (path === '/trade') return { page: '/trade' }
  if (path === '/market') return { page: '/market' }
  if (path === '/wallet') return { page: '/wallet' }
  if (path === '/login') return { page: '/login' }
  if (path === '/account') return { page: '/account' }
  if (path === '/legal/terms') return { page: '/legal/terms' }
  if (path === '/legal/privacy') return { page: '/legal/privacy' }
  if (path === '/legal/risk') return { page: '/legal/risk' }
  if (path === '/legal/cookies') return { page: '/legal/cookies' }

  return { page: '/' }
}

function navigate(path) {
  window.location.hash = '#' + path
}

function AppInner() {
  const [route, setRoute] = useState(getPageFromHash)
  // Keep previous route for back navigation from legal pages
  const [history, setHistory] = useState([])

  useEffect(() => {
    function onHashChange() {
      setRoute((prev) => {
        setHistory((h) => [...h, prev.page])
        return getPageFromHash()
      })
    }
    window.addEventListener('hashchange', onHashChange)
    return () => window.removeEventListener('hashchange', onHashChange)
  }, [])

  function handleNavigate(path) {
    // Support numeric -1 for "go back"
    if (path === -1) {
      const prev = history[history.length - 1]
      setHistory((h) => h.slice(0, -1))
      navigate(prev || '/')
      setRoute(getPageFromHash())
      return
    }
    navigate(path)
    setRoute(getPageFromHash())
  }

  function renderPage() {
    switch (route.page) {
      case '/':
        return <Dashboard onNavigate={handleNavigate} />
      case '/explorer':
        return <Explorer onNavigate={handleNavigate} />
      case '/trade':
        return <Trade onNavigate={handleNavigate} />
      case '/market':
        return <Market onNavigate={handleNavigate} />
      case '/wallet':
        return <Wallet />
      case '/block':
        return <BlockDetail height={route.param} onNavigate={handleNavigate} />
      case '/login':
        return <Login onNavigate={handleNavigate} />
      case '/account':
        return <Account onNavigate={handleNavigate} />
      case '/legal/terms':
        return <Terms onNavigate={handleNavigate} />
      case '/legal/privacy':
        return <Privacy onNavigate={handleNavigate} />
      case '/legal/risk':
        return <Risk onNavigate={handleNavigate} />
      case '/legal/cookies':
        return <Cookies onNavigate={handleNavigate} />
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
      <Footer onNavigate={handleNavigate} />
      <CookieBanner onNavigate={handleNavigate} />
    </div>
  )
}

export default function App() {
  return (
    <AuthProvider>
      <AppInner />
    </AuthProvider>
  )
}
