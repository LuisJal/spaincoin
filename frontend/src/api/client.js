// En desarrollo: VITE_API_URL=http://localhost:3001 (vite proxy)
// En producción: vacío → rutas relativas /api/... que nginx proxya
const API_BASE = import.meta.env.VITE_API_URL || ''

async function apiFetch(path) {
  const res = await fetch(`${API_BASE}${path}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`)
  return res.json()
}

async function apiAuthFetch(path, token) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  if (!res.ok) {
    const data = await res.json().catch(() => ({}))
    throw new Error(data.error || `HTTP ${res.status}`)
  }
  return res.json()
}

async function apiAuthPost(path, body, token) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(body)
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || `HTTP ${res.status}`)
  return data
}

// Blockchain
export async function getStatus() {
  return apiFetch('/api/status')
}

export async function getExplorer() {
  return apiFetch('/api/explorer')
}

export async function getBlock(height) {
  return apiFetch(`/api/blocks/${height}`)
}

export async function getWallet(address) {
  return apiFetch(`/api/wallet/${address}`)
}

// Market
export async function getPrice() {
  return apiFetch('/api/market/price')
}

export async function getMarketStats() {
  return apiFetch('/api/market/stats')
}

export async function getTicker() {
  return apiFetch('/api/market/ticker')
}

export async function getPriceHistory(points = 100, range_ = '24h') {
  return apiFetch(`/api/market/history?points=${points}&range=${range_}`)
}

export async function getMarketTable() {
  return apiFetch('/api/market/table')
}

// Transactions
export async function sendTx(txData) {
  const res = await fetch(`${API_BASE}/api/wallet/send`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(txData)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`)
  return res.json()
}

// Auth
export async function register(email, password, importKey) {
  const body = { email, password }
  if (importKey) body.import_key = importKey
  const res = await fetch(`${API_BASE}/api/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Error en el registro')
  return data
}

export async function login(email, password) {
  const res = await fetch(`${API_BASE}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Email o contraseña incorrectos')
  return data
}

export async function getMe(token) {
  const res = await fetch(`${API_BASE}/api/auth/me`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Sesión expirada')
  return data
}

// Trading
export async function buySPC(token, amountSPC, amountEUR) {
  const body = {}
  if (amountSPC) body.amount_spc = amountSPC
  else if (amountEUR) body.amount_eur = amountEUR
  return apiAuthPost('/api/trade/buy', body, token)
}

export async function sellSPC(token, amountSPC) {
  return apiAuthPost('/api/trade/sell', { amount_spc: amountSPC }, token)
}

export async function getTradeHistory(token) {
  return apiAuthFetch('/api/trade/history', token)
}

export async function getTradeBalance(token) {
  return apiAuthFetch('/api/trade/balance', token)
}
