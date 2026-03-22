import { useState, useEffect, useCallback } from 'react'
import { getBlock } from '../api/client.js'
import { truncateHash, formatTimeFull, formatSPC, formatNumber } from '../utils/format.js'

const COINBASE_ADDR = 'SPC0000000000000000000000000000000000000000'

const styles = {
  page: {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '2rem 1.5rem 4rem',
    width: '100%',
  },
  backLink: {
    display: 'inline-flex',
    alignItems: 'center',
    gap: '0.4rem',
    fontSize: '0.875rem',
    color: 'var(--text-secondary)',
    marginBottom: '1.5rem',
    cursor: 'pointer',
    background: 'none',
    border: 'none',
    padding: 0,
    transition: 'color 0.15s ease',
  },
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
    marginBottom: '2rem',
  },
  card: {
    background: 'var(--bg-card)',
    border: '1px solid var(--border)',
    borderRadius: '12px',
    padding: '1.5rem',
    marginBottom: '1.5rem',
  },
  cardTitle: {
    fontSize: '0.875rem',
    fontWeight: '600',
    color: 'var(--text-primary)',
    marginBottom: '1rem',
    paddingBottom: '0.75rem',
    borderBottom: '1px solid var(--border)',
    display: 'flex',
    alignItems: 'center',
    gap: '0.5rem',
  },
  fieldGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
    gap: '0.75rem 2rem',
  },
  field: {
    display: 'flex',
    flexDirection: 'column',
    gap: '0.2rem',
  },
  fieldLabel: {
    fontSize: '0.7rem',
    fontWeight: '500',
    color: 'var(--text-secondary)',
    textTransform: 'uppercase',
    letterSpacing: '0.06em',
  },
  fieldValue: {
    fontSize: '0.875rem',
    color: 'var(--text-primary)',
    wordBreak: 'break-all',
  },
  tableWrap: {
    background: 'var(--bg-card)',
    border: '1px solid var(--border)',
    borderRadius: '12px',
    overflow: 'hidden',
    overflowX: 'auto',
  },
  typeBadge: (isCoinbase) => ({
    display: 'inline-block',
    fontSize: '0.7rem',
    fontWeight: '600',
    padding: '0.15rem 0.5rem',
    borderRadius: '4px',
    background: isCoinbase ? 'rgba(16, 185, 129, 0.1)' : 'rgba(59, 130, 246, 0.1)',
    color: isCoinbase ? 'var(--green)' : 'var(--accent)',
    border: `1px solid ${isCoinbase ? 'rgba(16, 185, 129, 0.2)' : 'rgba(59, 130, 246, 0.2)'}`,
  }),
  error: {
    background: 'rgba(239, 68, 68, 0.1)',
    border: '1px solid rgba(239, 68, 68, 0.25)',
    borderRadius: '8px',
    padding: '1.5rem',
    fontSize: '0.875rem',
    color: '#ef4444',
    textAlign: 'center',
  },
  empty: {
    padding: '2rem',
    textAlign: 'center',
    color: 'var(--text-secondary)',
    fontSize: '0.875rem',
  },
}

function isCoinbaseTx(tx) {
  return (
    !tx.from ||
    tx.from === '' ||
    tx.from === COINBASE_ADDR ||
    tx.from.replace(/0/g, '').replace('SPC', '') === '' ||
    tx.type === 'COINBASE'
  )
}

export default function BlockDetail({ height, onNavigate }) {
  const [block, setBlock] = useState(null)
  const [error, setError] = useState(null)
  const [loading, setLoading] = useState(true)

  const fetchBlock = useCallback(async () => {
    if (height == null) return
    setLoading(true)
    try {
      const data = await getBlock(height)
      setBlock(data)
      setError(null)
    } catch (err) {
      setError(`Block #${height} not found or API error: ${err.message}`)
    } finally {
      setLoading(false)
    }
  }, [height])

  useEffect(() => {
    fetchBlock()
  }, [fetchBlock])

  const transactions = block?.transactions ?? []

  return (
    <div className="page-enter" style={styles.page}>
      <button
        style={styles.backLink}
        onClick={() => onNavigate('/explorer')}
        onMouseEnter={(e) => (e.target.style.color = 'var(--text-primary)')}
        onMouseLeave={(e) => (e.target.style.color = 'var(--text-secondary)')}
      >
        ← Back to Explorer
      </button>

      <h1 style={styles.pageTitle}>
        Block #{height != null ? formatNumber(height) : '—'}
      </h1>
      <p style={styles.pageSub}>Full block details and transactions</p>

      {loading && !block && (
        <div style={{ textAlign: 'center', padding: '3rem', color: 'var(--text-secondary)' }}>
          <span className="spinner" />
          <div style={{ marginTop: '0.75rem' }}>Loading block...</div>
        </div>
      )}

      {error && !loading && (
        <div style={styles.error}>{error}</div>
      )}

      {block && (
        <>
          {/* Block Header */}
          <div style={styles.card}>
            <div style={styles.cardTitle}>
              <span>⬡</span> Block Header
            </div>
            <div style={styles.fieldGrid}>
              <div style={styles.field}>
                <span style={styles.fieldLabel}>Height</span>
                <span style={{ ...styles.fieldValue, color: 'var(--accent)', fontWeight: '600' }}>
                  #{formatNumber(block.height)}
                </span>
              </div>
              <div style={styles.field}>
                <span style={styles.fieldLabel}>Timestamp</span>
                <span style={styles.fieldValue}>{formatTimeFull(block.timestamp)}</span>
              </div>
              <div style={styles.field}>
                <span style={styles.fieldLabel}>Hash</span>
                <span className="mono" style={{ ...styles.fieldValue, fontSize: '0.8rem' }}>
                  {block.hash}
                </span>
              </div>
              <div style={styles.field}>
                <span style={styles.fieldLabel}>Previous Hash</span>
                <span className="mono" style={{ ...styles.fieldValue, fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
                  {block.prev_hash}
                </span>
              </div>
              <div style={styles.field}>
                <span style={styles.fieldLabel}>Validator</span>
                <span className="mono" style={{ ...styles.fieldValue, fontSize: '0.82rem' }}>
                  {block.validator}
                </span>
              </div>
              {block.merkle_root && (
                <div style={styles.field}>
                  <span style={styles.fieldLabel}>Merkle Root</span>
                  <span className="mono" style={{ ...styles.fieldValue, fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
                    {block.merkle_root}
                  </span>
                </div>
              )}
              <div style={styles.field}>
                <span style={styles.fieldLabel}>Transactions</span>
                <span style={styles.fieldValue}>{block.tx_count ?? transactions.length}</span>
              </div>
            </div>
          </div>

          {/* Transactions */}
          <div>
            <div style={{ marginBottom: '1rem', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <span style={{ fontWeight: '600', color: 'var(--text-primary)' }}>
                Transactions ({transactions.length})
              </span>
            </div>
            <div style={styles.tableWrap}>
              {transactions.length === 0 ? (
                <div style={styles.empty}>No transactions in this block</div>
              ) : (
                <table>
                  <thead>
                    <tr>
                      <th>ID</th>
                      <th>From</th>
                      <th>To</th>
                      <th>Amount</th>
                      <th>Fee</th>
                      <th>Type</th>
                    </tr>
                  </thead>
                  <tbody>
                    {transactions.map((tx, i) => {
                      const coinbase = isCoinbaseTx(tx)
                      return (
                        <tr key={tx.id || tx.hash || i}>
                          <td>
                            <span className="mono" style={{ fontSize: '0.78rem', color: 'var(--text-secondary)' }}>
                              {truncateHash(tx.id || tx.hash || String(i), 8, 4)}
                            </span>
                          </td>
                          <td>
                            {coinbase ? (
                              <span style={{ color: 'var(--green)', fontSize: '0.8rem', fontStyle: 'italic' }}>
                                coinbase
                              </span>
                            ) : (
                              <span className="mono" style={{ fontSize: '0.78rem', color: 'var(--text-secondary)' }}>
                                {truncateHash(tx.from, 6, 4)}
                              </span>
                            )}
                          </td>
                          <td>
                            <span className="mono" style={{ fontSize: '0.78rem', color: 'var(--text-secondary)' }}>
                              {truncateHash(tx.to, 6, 4)}
                            </span>
                          </td>
                          <td style={{ fontWeight: '500', color: 'var(--text-primary)' }}>
                            {formatSPC(tx.amount)}
                          </td>
                          <td style={{ color: 'var(--text-secondary)' }}>
                            {tx.fee != null ? formatSPC(tx.fee) : '—'}
                          </td>
                          <td>
                            <span style={styles.typeBadge(coinbase)}>
                              {coinbase ? 'COINBASE' : (tx.type || 'TRANSFER')}
                            </span>
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              )}
            </div>
          </div>
        </>
      )}
    </div>
  )
}
