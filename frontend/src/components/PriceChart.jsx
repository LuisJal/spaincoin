import { useMemo } from 'react'

// Minimal SVG line chart — no external dependencies
export default function PriceChart({ data, width = 600, height = 200, color }) {
  const chartColor = color || 'var(--green)'

  const { points, polyline, polygon, minPrice, maxPrice, priceRange } = useMemo(() => {
    if (!data || data.length < 2) return { points: [], polyline: '', polygon: '', minPrice: 0, maxPrice: 0, priceRange: 0 }

    const prices = data.map(d => d.price)
    const min = Math.min(...prices)
    const max = Math.max(...prices)
    const range = max - min || 0.001

    const padding = { top: 10, bottom: 10, left: 0, right: 0 }
    const chartW = width - padding.left - padding.right
    const chartH = height - padding.top - padding.bottom

    const pts = data.map((d, i) => ({
      x: padding.left + (i / (data.length - 1)) * chartW,
      y: padding.top + chartH - ((d.price - min) / range) * chartH,
      price: d.price,
      height: d.height,
    }))

    const line = pts.map(p => `${p.x},${p.y}`).join(' ')
    const fill = `${pts[0].x},${height} ${line} ${pts[pts.length - 1].x},${height}`

    return { points: pts, polyline: line, polygon: fill, minPrice: min, maxPrice: max, priceRange: range }
  }, [data, width, height])

  if (!data || data.length < 2) {
    return (
      <div style={{ width, height, display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>
        Cargando datos...
      </div>
    )
  }

  return (
    <svg width="100%" height={height} viewBox={`0 0 ${width} ${height}`} preserveAspectRatio="none" style={{ display: 'block' }}>
      <defs>
        <linearGradient id="chartFill" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor={chartColor} stopOpacity="0.3" />
          <stop offset="100%" stopColor={chartColor} stopOpacity="0.02" />
        </linearGradient>
      </defs>

      {/* Fill under the line */}
      <polygon points={polygon} fill="url(#chartFill)" />

      {/* The line itself */}
      <polyline
        points={polyline}
        fill="none"
        stroke={chartColor}
        strokeWidth="2"
        strokeLinejoin="round"
        strokeLinecap="round"
        vectorEffect="non-scaling-stroke"
      />

      {/* Price labels */}
      <text x="4" y="16" fill="var(--text-secondary)" fontSize="10" fontFamily="Inter, sans-serif">
        {maxPrice.toFixed(4)}€
      </text>
      <text x="4" y={height - 4} fill="var(--text-secondary)" fontSize="10" fontFamily="Inter, sans-serif">
        {minPrice.toFixed(4)}€
      </text>
    </svg>
  )
}
