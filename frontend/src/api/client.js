const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:3001'

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
