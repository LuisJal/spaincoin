/**
 * Format pesetas (1e18 base units) into SPC with commas.
 * Node uses 1e18 decimals like Ethereum but the prompt says divide by 1e15 — keeping 1e15.
 */
export function formatSPC(pesetas) {
  if (pesetas == null) return '—'
  const spc = Number(pesetas) / 1e15
  return formatNumber(spc.toFixed(3)) + ' SPC'
}

/**
 * Truncate a hash: first `start` chars + "..." + last `end` chars.
 */
export function truncateHash(hash, start = 8, end = 6) {
  if (!hash) return '—'
  if (hash.length <= start + end + 3) return hash
  return `${hash.slice(0, start)}...${hash.slice(-end)}`
}

/**
 * Format a Go UnixNano timestamp into a relative time string.
 * Go time.Now().UnixNano() → divide by 1e6 to get ms for JS Date.
 */
export function formatTime(timestamp) {
  if (!timestamp) return '—'
  const ms = Number(timestamp) / 1e6
  const date = new Date(ms)
  const now = Date.now()
  const diffMs = now - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)

  if (diffSec < 0) return 'just now'
  if (diffSec < 60) return `${diffSec}s ago`
  const diffMin = Math.floor(diffSec / 60)
  if (diffMin < 60) return `${diffMin}m ago`
  const diffHr = Math.floor(diffMin / 60)
  if (diffHr < 24) return `${diffHr}h ago`
  const diffDay = Math.floor(diffHr / 24)
  return `${diffDay}d ago`
}

/**
 * Format a Go UnixNano timestamp into a full human-readable date string.
 */
export function formatTimeFull(timestamp) {
  if (!timestamp) return '—'
  const ms = Number(timestamp) / 1e6
  const date = new Date(ms)
  return date.toLocaleString('en-GB', {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

/**
 * Add thousands separators to a number or string.
 */
export function formatNumber(n) {
  if (n == null) return '—'
  const parts = String(n).split('.')
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',')
  return parts.join('.')
}
