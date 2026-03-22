// En desarrollo: VITE_API_URL=http://localhost:3001 (vite proxy)
// En producción: vacío → rutas relativas /api/... que nginx proxya
const API_BASE = import.meta.env.VITE_API_URL || ''

async function apiFetch(path) {
  const res = await fetch(`${API_BASE}${path}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`)
  return res.json()
}

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

export async function getPrice() {
  return apiFetch('/api/market/price')
}

export async function getMarketStats() {
  return apiFetch('/api/market/stats')
}

export async function sendTx(txData) {
  const res = await fetch(`${API_BASE}/api/wallet/send`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(txData)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`)
  return res.json()
}

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
  return data // {token, address, email}
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
