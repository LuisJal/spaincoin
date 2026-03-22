import { useState, useEffect, useCallback } from 'react'
import { getExplorer } from '../api/client.js'
import { truncateHash, formatTimeFull, formatNumber } from '../utils/format.js'

const styles = {
  page: {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '2rem 1.5rem 4rem',
    width: '100%',
  },
  header: {
    display: 'flex',
    alignItems: 'flex-start',
    justifyContent: 'space-between',
    flexWrap: 'wrap',
    gap: '1rem',
    marginBottom: '2rem',
  },
  titleGroup: {},
  pageTitle: {
    fontSize: '1.5rem',
    fontWeight: '700',
    color: 'var(--text-primary)',
    marginBottom: '0.35rem',
    letterSpacing: '-0.02em',
  },
  pageSub: {
    fontSize: '0.875rem',
    color: 'var(--text-secondary)',
  },
  heightBadge: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.4rem',
    fontSize: '0.8rem',
    fontWeight: '600',
    color: 'var(--green)',
    background: 'rgba(16, 185, 129, 0.1)',
    border: '1px solid rgba(16, 185, 129, 0.2)',
    padding: '0.4rem 0.8rem',
    borderRadius: '8px',
    flexShrink: 0,
  },
  tableWrap: {
    background: 'var(--bg-card)',
    border: '1px solid var(--border)',
    borderRadius: '12px',
    overflow: 'hidden',
    overflowX: 'auto',
  },
  tableRow: (hovered) => ({
    cursor: 'pointer',
    background: hovered ? 'rgba(59, 130, 246, 0.05)' : 'transparent',
    transition: 'background 0.1s ease',
  }),
  error: {
    background: 'rgba(239, 68, 68, 0.1)',
    border: '1px solid rgba(239, 68, 68, 0.25)',
    borderRadius: '8px',
    padding: '0.75rem 1rem',
    fontSize: '0.875rem',
    color: '#ef4444',
    marginBottom: '1.5rem',
  },
  empty: {
    padding: '3rem',
    textAlign: 'center',
    color: 'var(--text-secondary)',
  },
  footer: {
    marginTop: '1rem',
    fontSize: '0.8rem',
    color: 'var(--text-secondary)',
    textAlign: 'center',
  },
}

function BlockTableRow({ block, onNavigate }) {
  const [hovered, setHovered] = useState(false)

  return (
    <tr
      style={styles.tableRow(hovered)}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
      onClick={() => onNavigate(`/block/${block.height}`)}
    >
      <td>
        <span style={{ color: 'var(--accent)', fontWeight: '600' }}>
          #{formatNumber(block.height)}
        </span>
      </td>
      <td>
        <span className="mono" style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
          {truncateHash(block.hash, 10, 6)}
        </span>
      </td>
      <td style={{ textAlign: 'center' }}>
        <span style={{
          background: (block.tx_count ?? 0) > 0 ? 'rgba(59, 130, 246, 0.1)' : 'rgba(156, 163, 175, 0.1)',
          color: (block.tx_count ?? 0) > 0 ? 'var(--accent)' : 'var(--text-secondary)',
          padding: '0.15rem 0.6rem',
          borderRadius: '4px',
          fontSize: '0.8rem',
          fontWeight: '500',
        }}>
          {block.tx_count ?? 0}
        </span>
      </td>
      <td>
        <span className="mono" style={{ fontSize: '0.78rem', color: 'var(--text-secondary)' }}>
          {truncateHash(block.validator, 8, 6)}
        </span>
      </td>
      <td style={{ color: 'var(--text-secondary)', fontSize: '0.8rem', whiteSpace: 'nowrap' }}>
        {formatTimeFull(block.timestamp)}
      </td>
    </tr>
  )
}

export default function Explorer({ onNavigate }) {
  const [data, setData] = useState(null)
  const [error, setError] = useState(null)
  const [loading, setLoading] = useState(true)

  const fetchData = useCallback(async () => {
    try {
      const result = await getExplorer()
      setData(result)
      setError(null)
    } catch (err) {
      setError('Cannot connect to SpainCoin API: ' + err.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchData()
    const interval = setInterval(fetchData, 10000)
    return () => clearInterval(interval)
  }, [fetchData])

  const blocks = data?.blocks ?? []

  return (
    <div className="page-enter" style={styles.page}>
      <div style={styles.header}>
        <div style={styles.titleGroup}>
          <h1 style={styles.pageTitle}>Block Explorer</h1>
          <p style={styles.pageSub}>Browse all blocks on the SpainCoin network</p>
        </div>
        {data?.height != null && (
          <div style={styles.heightBadge}>
            <span className="pulse-dot" style={{ fontSize: '0.55rem' }}>●</span>
            Height #{formatNumber(data.height)}
          </div>
        )}
      </div>

      {error && <div style={styles.error}>{error}</div>}

      <div style={styles.tableWrap}>
        {loading && blocks.length === 0 ? (
          <div style={styles.empty}>
            <span className="spinner" />
            <div style={{ marginTop: '0.75rem' }}>Loading blocks...</div>
          </div>
        ) : blocks.length === 0 ? (
          <div style={styles.empty}>No blocks found</div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Height</th>
                <th>Hash</th>
                <th style={{ textAlign: 'center' }}>Txs</th>
                <th>Validator</th>
                <th>Timestamp</th>
              </tr>
            </thead>
            <tbody>
              {blocks.map((block) => (
                <BlockTableRow
                  key={block.height}
                  block={block}
                  onNavigate={onNavigate}
                />
              ))}
            </tbody>
          </table>
        )}
      </div>

      {blocks.length > 0 && (
        <p style={styles.footer}>
          Showing {blocks.length} block{blocks.length !== 1 ? 's' : ''} · Click a row to view block detail
        </p>
      )}
    </div>
  )
}
