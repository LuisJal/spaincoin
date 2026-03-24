export default function Validators({ onNavigate }) {
  return (
    <div className="page-enter" style={{ maxWidth: '800px', margin: '0 auto', padding: '2rem 1.5rem' }}>
      <h1 style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
        Ser Validador
      </h1>
      <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '2rem' }}>
        Corre un nodo, valida bloques y gana recompensas en SPC.
      </p>

      {/* What is a validator */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '1.5rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
          Qué es un validador
        </h2>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.7 }}>
          Un validador es un nodo que participa en el consenso de la red. Cada 5 segundos, un validador
          es seleccionado para producir el siguiente bloque. A cambio, recibe <strong style={{ color: 'var(--green)' }}>1 SPC de recompensa</strong> por bloque.
          Cuanto más stake tengas, más probabilidad de ser seleccionado.
        </p>
      </div>

      {/* Requirements */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '1.5rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
          Requisitos
        </h2>
        <div style={{ display: 'grid', gap: '0.75rem' }}>
          {[
            { label: 'Servidor', value: 'VPS con 2 vCPU, 4GB RAM (~8 EUR/mes)' },
            { label: 'Sistema', value: 'Ubuntu 22.04 o similar' },
            { label: 'Stake mínimo', value: '1 SPC' },
            { label: 'Conocimientos', value: 'Básico de terminal/SSH' },
            { label: 'Disponibilidad', value: '24/7 (el nodo debe estar siempre activo)' },
          ].map((r, i) => (
            <div key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '0.5rem 0', borderBottom: i < 4 ? '1px solid var(--border)' : 'none' }}>
              <span style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>{r.label}</span>
              <span style={{ fontSize: '0.85rem', fontWeight: '600', color: 'var(--text-primary)', textAlign: 'right' }}>{r.value}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Step by step */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '1.5rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '1rem' }}>
          Guía paso a paso
        </h2>

        {[
          { step: '1', title: 'Crea tu wallet', code: 'go build -o spc ./cli/\n./spc wallet new' },
          { step: '2', title: 'Clona el repositorio en tu servidor', code: 'git clone https://github.com/LuisJal/spaincoin\ncd spaincoin' },
          { step: '3', title: 'Compila el nodo', code: 'CGO_ENABLED=0 go build -o spaincoin ./node/cmd/' },
          { step: '4', title: 'Configura las variables de entorno', code: 'export SPC_VALIDATOR_KEY=tu_clave_privada\nexport SPC_VALIDATOR_ADDRESS=SPCtu_address\nexport SPC_RPC_PORT=8545\nexport SPC_P2P_PORT=30303\nexport SPC_BLOCK_TIME=5\nexport SPC_DATA_DIR=./data' },
          { step: '5', title: 'Arranca el nodo', code: './spaincoin' },
        ].map((s, i) => (
          <div key={i} style={{ marginBottom: i < 4 ? '1.25rem' : 0 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem', marginBottom: '0.5rem' }}>
              <div style={{
                width: '24px', height: '24px', borderRadius: '50%',
                background: 'var(--accent)', color: '#fff', display: 'flex',
                alignItems: 'center', justifyContent: 'center', fontWeight: '700', fontSize: '0.75rem', flexShrink: 0,
              }}>{s.step}</div>
              <span style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)' }}>{s.title}</span>
            </div>
            <div style={{
              background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.75rem 1rem',
              fontFamily: 'monospace', fontSize: '0.78rem', color: 'var(--accent)',
              overflowX: 'auto', whiteSpace: 'pre', lineHeight: 1.8,
            }}>{s.code}</div>
          </div>
        ))}
      </div>

      {/* Rewards */}
      <div style={{
        background: 'rgba(16, 185, 129, 0.08)', border: '1px solid rgba(16, 185, 129, 0.25)',
        borderRadius: '12px', padding: '1.5rem', marginBottom: '1.5rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '600', color: 'var(--green)', marginBottom: '0.5rem' }}>
          Recompensas
        </h2>
        <ul style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.8, paddingLeft: '1.25rem' }}>
          <li><strong style={{ color: 'var(--text-primary)' }}>1 SPC por bloque</strong> producido</li>
          <li>Un bloque cada 5 segundos</li>
          <li>~17.280 bloques/día = hasta <strong style={{ color: 'var(--green)' }}>17.280 SPC/día</strong> (repartidos entre validadores)</li>
          <li>Cuanto más stake, más bloques te tocan</li>
        </ul>
      </div>

      {/* Community */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem', textAlign: 'center',
      }}>
        <h2 style={{ fontSize: '1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
          ¿Necesitas ayuda?
        </h2>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '1rem' }}>
          Únete a la comunidad de validadores en Telegram
        </p>
        <a href="https://t.me/spaincoin_comunidad" target="_blank" rel="noopener noreferrer" style={{
          display: 'inline-block', padding: '0.65rem 1.5rem',
          background: '#0088cc', border: 'none', borderRadius: '8px',
          color: '#fff', fontSize: '0.85rem', fontWeight: '600', textDecoration: 'none',
        }}>Telegram SpainCoin</a>
      </div>
    </div>
  )
}
