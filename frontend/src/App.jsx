import { useState, useEffect } from 'react'
import Navbar from './components/Navbar.jsx'
import Footer from './components/Footer.jsx'
import CookieBanner from './components/CookieBanner.jsx'
import Onboarding from './pages/Onboarding.jsx'
import Landing from './pages/Landing.jsx'
import Explorer from './pages/Explorer.jsx'
import BlockDetail from './pages/BlockDetail.jsx'
import MarketInfo from './pages/MarketInfo.jsx'
import WalletDownload from './pages/WalletDownload.jsx'
import Validators from './pages/Validators.jsx'
import WhitePaper from './pages/WhitePaper.jsx'
import HowToSell from './pages/HowToSell.jsx'
import Terms from './pages/legal/Terms.jsx'
import Privacy from './pages/legal/Privacy.jsx'
import Risk from './pages/legal/Risk.jsx'
import Cookies from './pages/legal/Cookies.jsx'

function getPageFromHash() {
  const hash = window.location.hash || '#/'
  const path = hash.startsWith('#') ? hash.slice(1) : hash

  if (!path || path === '/') return { page: '/' }

  const blockMatch = path.match(/^\/block\/(\d+)$/)
  if (blockMatch) return { page: '/block', param: Number(blockMatch[1]) }

  if (path === '/explorer') return { page: '/explorer' }
  if (path === '/market') return { page: '/market' }
  if (path === '/wallet') return { page: '/wallet' }
  if (path === '/validators') return { page: '/validators' }
  if (path === '/whitepaper') return { page: '/whitepaper' }
  if (path === '/como-vender') return { page: '/como-vender' }
  if (path === '/legal/terms') return { page: '/legal/terms' }
  if (path === '/legal/privacy') return { page: '/legal/privacy' }
  if (path === '/legal/risk') return { page: '/legal/risk' }
  if (path === '/legal/cookies') return { page: '/legal/cookies' }

  return { page: '/' }
}

function navigate(path) {
  window.location.hash = '#' + path
}

export default function App() {
  const [route, setRoute] = useState(getPageFromHash)
  const [history, setHistory] = useState([])
  const [showOnboarding, setShowOnboarding] = useState(() => !localStorage.getItem('spc_onboarded'))

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
        return <Landing onNavigate={handleNavigate} />
      case '/explorer':
        return <Explorer onNavigate={handleNavigate} />
      case '/block':
        return <BlockDetail height={route.param} onNavigate={handleNavigate} />
      case '/market':
        return <MarketInfo onNavigate={handleNavigate} />
      case '/wallet':
        return <WalletDownload onNavigate={handleNavigate} />
      case '/validators':
        return <Validators onNavigate={handleNavigate} />
      case '/whitepaper':
        return <WhitePaper onNavigate={handleNavigate} />
      case '/como-vender':
        return <HowToSell onNavigate={handleNavigate} />
      case '/legal/terms':
        return <Terms onNavigate={handleNavigate} />
      case '/legal/privacy':
        return <Privacy onNavigate={handleNavigate} />
      case '/legal/risk':
        return <Risk onNavigate={handleNavigate} />
      case '/legal/cookies':
        return <Cookies onNavigate={handleNavigate} />
      default:
        return <Landing onNavigate={handleNavigate} />
    }
  }

  function handleOnboardingComplete() {
    localStorage.setItem('spc_onboarded', '1')
    setShowOnboarding(false)
  }

  if (showOnboarding) {
    return (
      <div style={{ minHeight: '100vh', background: 'var(--bg-primary)' }}>
        <Onboarding onComplete={handleOnboardingComplete} />
      </div>
    )
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
