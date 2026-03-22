const styles = {
  card: {
    background: 'var(--bg-card)',
    border: '1px solid var(--border)',
    borderRadius: '12px',
    padding: '1.25rem 1.5rem',
    display: 'flex',
    flexDirection: 'column',
    gap: '0.35rem',
    transition: 'border-color 0.15s ease, transform 0.15s ease',
    cursor: 'default',
  },
  cardHover: {
    borderColor: 'var(--border-accent)',
    transform: 'translateY(-1px)',
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  title: {
    fontSize: '0.75rem',
    fontWeight: '500',
    color: 'var(--text-secondary)',
    textTransform: 'uppercase',
    letterSpacing: '0.06em',
  },
  icon: {
    fontSize: '1.2rem',
    opacity: 0.6,
  },
  value: {
    fontSize: '1.6rem',
    fontWeight: '700',
    color: 'var(--text-primary)',
    letterSpacing: '-0.02em',
    lineHeight: 1.2,
  },
  subtitle: {
    fontSize: '0.75rem',
    color: 'var(--text-secondary)',
    fontWeight: '400',
  },
}

import { useState } from 'react'

export default function StatCard({ title, value, subtitle, icon, accent }) {
  const [hovered, setHovered] = useState(false)

  return (
    <div
      style={{
        ...styles.card,
        ...(hovered ? styles.cardHover : {}),
        ...(accent ? { borderColor: 'rgba(59, 130, 246, 0.3)' } : {}),
      }}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <div style={styles.header}>
        <span style={styles.title}>{title}</span>
        {icon && <span style={styles.icon}>{icon}</span>}
      </div>
      <div style={{
        ...styles.value,
        ...(accent ? { color: 'var(--accent)' } : {}),
      }}>
        {value ?? <span className="spinner" style={{ width: 24, height: 24 }} />}
      </div>
      {subtitle && <div style={styles.subtitle}>{subtitle}</div>}
    </div>
  )
}
