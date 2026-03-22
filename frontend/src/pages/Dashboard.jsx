import { useState, useEffect, useCallback } from 'react'
import StatCard from '../components/StatCard.jsx'
import { getStatus, getExplorer, getPrice } from '../api/client.js'
import { truncateHash, formatTime, formatNumber } from '../utils/format.js'

const styles = {
  page: {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '2rem 1.5rem 4rem',
    width: '100%',
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
  statsGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
    gap: '1rem',
    marginBottom: '2.5rem',
  },
  section: {
    marginBottom: '2.5rem',
  },
  sectionHeader: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: '1rem',
  },
  sectionTitle: {
    fontSize: '1rem',
    fontWeight: '600',
    color: 'var(--text-primary)',
    letterSpacing: '-0.01em',
  },
  badge: {
    fontSize: '0.7rem',
    fontWeight: '600',
    color: 'var(--accent)',
    background: 'rgba(59, 130, 246, 0.1)',
    border: '1px solid rgba(59, 130, 246, 0.25)',
    padding: '0.2rem 0.6rem',
    borderRadius: '20px',
  },
  tableWrap: {
    background: 'var(--bg-card)',
    border: '1px solid var(--border)',
    borderRadius: '12px',
    overflow: 'hidden',
  },
  tableRow: (clickable, hovered) => ({
    cursor: clickable ? 'pointer' : 'default',
    background: hovered ? 'rgba(59, 130, 246, 0.05)' : 'transparent',
    transition: 'background 0.1s ease',
  }),
  infoGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
    gap: '1rem',
  },
  infoCard: {
    background: 'var(--bg-card)',
    border: '1px solid var(--border)',
    borderRadius: '12px',
    padding: '1.5rem',
  },
  infoCardTitle: {
    fontSize: '0.875rem',
    fontWeight: '600',
    color: 'var(--text-primary)',
    marginBottom: '1rem',
    paddingBottom: '0.75rem',
    borderBottom: '1px solid var(--border)',
  },
  infoRow: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '0.4rem 0',
  },
  infoLabel: {
    fontSize: '0.8125rem',
    color: 'var(--text-secondary)',
  },
  infoValue: {
    fontSize: '0.8125rem',
    color: 'var(--text-primary)',
    fontWeight: '500',
  },
  error: {
    background: 'rgba(239, 68, 68, 0.1)',
    border: '1px solid rgba(239, 68, 68, 0.25)',
    borderRadius: '8px',
    padding: '0.75rem 1rem',
    fontSize: '0.875rem',
    color: '#ef4444',
    marginBottom: '1.5rem',
  },
}

function BlockRow({ block, onNavigate }) {
  const [hovered, setHovered] = useState(false)

  return (
    <tr
      style={styles.tableRow(true, hovered)}
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
        <span className="mono" style={{ color: 'var(--text-secondary)' }}>
          {truncateHash(block.hash, 8, 4)}
        </span>
      </td>
      <td style={{ textAlign: 'center' }}>
        <span style={{
          background: block.tx_count > 0 ? 'rgba(59, 130, 246, 0.1)' : 'rgba(156, 163, 175, 0.1)',
          color: block.tx_count > 0 ? 'var(--accent)' : 'var(--text-secondary)',
          padding: '0.15rem 0.5rem',
          borderRadius: '4px',
          fontSize: '0.8rem',
          fontWeight: '500',
        }}>
          {block.tx_count ?? 0}
        </span>
      </td>
      <td>
        <span className="mono" style={{ color: 'var(--text-secondary)', fontSize: '0.8rem' }}>
          {truncateHash(block.validator, 6, 4)}
        </span>
      </td>
      <td style={{ color: 'var(--text-secondary)', fontSize: '0.8rem' }}>
        {formatTime(block.timestamp)}
      </td>
    </tr>
  )
}

export default function Dashboard({ onNavigate }) {
  const [status, setStatus] = useState(null)
  const [explorer, setExplorer] = useState(null)
  const [price, setPrice] = useState(null)
  const [error, setError] = useState(null)
  const [lastUpdated, setLastUpdated] = useState(null)

  const fetchData = useCallback(async () => {
    try {
      const [s, e, p] = await Promise.allSettled([
        getStatus(),
        getExplorer(),
        getPrice(),
      ])
      if (s.status === 'fulfilled') setStatus(s.value)
      if (e.status === 'fulfilled') setExplorer(e.value)
      if (p.status === 'fulfilled') setPrice(p.value)
      if (s.status === 'rejected' && e.status === 'rejected') {
        setError('Cannot connect to SpainCoin API. Make sure the node is running on port 3001.')
      } else {
        setError(null)
      }
      setLastUpdated(new Date())
    } catch (err) {
      setError(err.message)
    }
  }, [])

  useEffect(() => {
    fetchData()
    const interval = setInterval(fetchData, 10000)
    return () => clearInterval(interval)
  }, [fetchData])

  const height = status?.node?.height
  const totalSupply = explorer?.total_supply_spc
  const mempoolSize = status?.node?.mempool_size
  const priceEur = price?.price_eur
  const change24h = price?.change_24h

  const blocks = explorer?.blocks?.slice(0, 10) ?? []

  return (
    <div className="page-enter" style={styles.page}>
      <h1 style={styles.pageTitle}>Dashboard</h1>
      <p style={styles.pageSub}>
        Real-time SpainCoin network overview
        {lastUpdated && (
          <span style={{ marginLeft: '0.5rem', opacity: 0.6 }}>
            · Updated {formatTime(lastUpdated.getTime() * 1e6)}
          </span>
        )}
      </p>

      {error && <div style={styles.error}>{error}</div>}

      {/* Hero stats */}
      <div style={styles.statsGrid}>
        <StatCard
          title="Block Height"
          value={height != null ? `#${formatNumber(height)}` : null}
          subtitle="latest block"
          icon="⬡"
        />
        <StatCard
          title="Total Supply"
          value={totalSupply != null ? `${formatNumber(Number(totalSupply).toFixed(0))} SPC` : null}
          subtitle={`of 21,000,000 max`}
          icon="◈"
        />
        <StatCard
          title="SPC Price"
          value={priceEur != null ? `€${priceEur}` : null}
          subtitle={
            change24h != null
              ? `${change24h >= 0 ? '+' : ''}${change24h}% 24h`
              : 'market price'
          }
          icon="$"
          accent={true}
        />
        <StatCard
          title="Mempool"
          value={mempoolSize != null ? `${formatNumber(mempoolSize)}` : null}
          subtitle="txs pending"
          icon="⏳"
        />
      </div>

      {/* Latest Blocks */}
      <div style={styles.section}>
        <div style={styles.sectionHeader}>
          <span style={styles.sectionTitle}>Latest Blocks</span>
          {height != null && (
            <span style={styles.badge}>Height #{formatNumber(height)}</span>
          )}
        </div>
        <div style={styles.tableWrap}>
          {blocks.length === 0 && !error ? (
            <div style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-secondary)' }}>
              <span className="spinner" style={{ marginRight: '0.5rem' }} />
              Loading blocks...
            </div>
          ) : blocks.length === 0 ? (
            <div style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-secondary)' }}>
              No blocks available
            </div>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>Height</th>
                  <th>Hash</th>
                  <th style={{ textAlign: 'center' }}>Txs</th>
                  <th>Validator</th>
                  <th>Time</th>
                </tr>
              </thead>
              <tbody>
                {blocks.map((block) => (
                  <BlockRow
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
          <div style={{ marginTop: '0.75rem', textAlign: 'right' }}>
            <a
              href="#/explorer"
              style={{ fontSize: '0.8rem', color: 'var(--accent)' }}
              onClick={(e) => { e.preventDefault(); onNavigate('/explorer') }}
            >
              View full explorer →
            </a>
          </div>
        )}
      </div>

      {/* Network info */}
      <div style={styles.section}>
        <div style={styles.sectionHeader}>
          <span style={styles.sectionTitle}>Network Information</span>
        </div>
        <div style={styles.infoGrid}>
          <div style={styles.infoCard}>
            <div style={styles.infoCardTitle}>SpainCoin Network</div>
            {[
              ['P2P Port', '30303'],
              ['RPC Port', '8545'],
              ['Block Time', `${status?.block_time_seconds ?? 5}s`],
              ['Consensus', 'Proof of Stake'],
              ['Version', status?.version ?? '—'],
            ].map(([label, value]) => (
              <div key={label} style={styles.infoRow}>
                <span style={styles.infoLabel}>{label}</span>
                <span style={styles.infoValue}>{value}</span>
              </div>
            ))}
          </div>
          <div style={styles.infoCard}>
            <div style={styles.infoCardTitle}>About $SPC</div>
            {[
              ['Max Supply', '21,000,000 SPC'],
              ['Genesis Allocation', '1,000,000 SPC'],
              ['Decimals', '18 (pesetas)'],
              ['Symbol', '$SPC'],
              ['Chain', 'SpainCoin L1'],
            ].map(([label, value]) => (
              <div key={label} style={styles.infoRow}>
                <span style={styles.infoLabel}>{label}</span>
                <span style={styles.infoValue}>{value}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
