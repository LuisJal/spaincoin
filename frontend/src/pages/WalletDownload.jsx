export default function WalletDownload({ onNavigate }) {
  const downloads = [
    { os: 'macOS (Apple Silicon)', file: 'spc-macos-arm64', icon: '🍎', note: 'MacBook M1/M2/M3/M4' },
    { os: 'macOS (Intel)', file: 'spc-macos-amd64', icon: '🍎', note: 'MacBook antes de 2020' },
    { os: 'Windows', file: 'spc-windows-amd64.exe', icon: '🪟', note: 'Windows 10/11 (64 bits)' },
    { os: 'Linux', file: 'spc-linux-amd64', icon: '🐧', note: 'Ubuntu, Debian, Fedora...' },
  ]

  const ghRelease = 'https://github.com/spaincoin/spaincoin/releases/latest/download'

  return (
    <div className="page-enter" style={{ maxWidth: '800px', margin: '0 auto', padding: '2rem 1.5rem' }}>
      <h1 style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
        Wallet SpainCoin
      </h1>
      <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '2rem' }}>
        Tu wallet, tus claves, tus fondos. 100% self-custody — nadie más tiene acceso.
      </p>

      {/* Download buttons */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '1.5rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '1rem' }}>
          Paso 1 — Descarga
        </h2>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '1.25rem' }}>
          Elige tu sistema operativo. Es un archivo único, no necesita instalación.
        </p>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '0.75rem' }}>
          {downloads.map((d) => (
            <a key={d.file} href={`${ghRelease}/${d.file}`} target="_blank" rel="noopener noreferrer"
              style={{
                display: 'flex', alignItems: 'center', gap: '0.75rem',
                padding: '0.9rem 1rem', borderRadius: '10px',
                background: 'var(--bg-secondary)', border: '1px solid var(--border)',
                textDecoration: 'none', cursor: 'pointer', transition: 'border-color 0.15s',
              }}
              onMouseEnter={e => e.currentTarget.style.borderColor = 'var(--accent)'}
              onMouseLeave={e => e.currentTarget.style.borderColor = 'var(--border)'}
            >
              <span style={{ fontSize: '1.5rem' }}>{d.icon}</span>
              <div>
                <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)' }}>{d.os}</div>
                <div style={{ fontSize: '0.7rem', color: 'var(--text-secondary)' }}>{d.note}</div>
              </div>
            </a>
          ))}
        </div>
      </div>

      {/* Step 2: Create wallet */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '1.5rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '1rem' }}>
          Paso 2 — Crea tu wallet
        </h2>
        <div style={{ display: 'grid', gap: '1.25rem' }}>
          <div>
            <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)', marginBottom: '0.4rem' }}>
              Abre una terminal y ejecuta:
            </div>
            <div style={{
              background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.85rem 1rem',
              fontFamily: 'monospace', fontSize: '0.85rem', color: 'var(--accent)',
            }}>
              ./spc wallet new
            </div>
          </div>
          <div>
            <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)', marginBottom: '0.4rem' }}>
              Verás algo así:
            </div>
            <div style={{
              background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.85rem 1rem',
              fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--green)', lineHeight: 1.8,
            }}>
              <div>Nueva wallet creada</div>
              <div>Dirección: <span style={{ color: 'var(--accent)' }}>SPCa1b2c3d4e5f6...</span></div>
              <div>Clave privada: <span style={{ color: 'var(--red)' }}>1a2b3c4d5e6f... (GUÁRDALA EN PAPEL)</span></div>
            </div>
          </div>
        </div>
      </div>

      {/* Step 3: Use it */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '1.5rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '1rem' }}>
          Paso 3 — Úsalo
        </h2>
        <div style={{ display: 'grid', gap: '1rem' }}>
          {[
            { cmd: './spc wallet balance', desc: 'Ver tu saldo' },
            { cmd: './spc send --to SPCxxx... --amount 10', desc: 'Enviar SPC a alguien' },
            { cmd: './spc chain status', desc: 'Ver estado de la red' },
          ].map((c, i) => (
            <div key={i} style={{ display: 'flex', gap: '1rem', alignItems: 'center', flexWrap: 'wrap' }}>
              <div style={{
                background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.6rem 0.85rem',
                fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--accent)', flex: '1 1 250px',
              }}>{c.cmd}</div>
              <div style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>{c.desc}</div>
            </div>
          ))}
        </div>
      </div>

      {/* macOS permission note */}
      <div style={{
        background: 'rgba(59, 130, 246, 0.08)', border: '1px solid rgba(59, 130, 246, 0.2)',
        borderRadius: '12px', padding: '1.25rem', marginBottom: '1.5rem',
      }}>
        <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--accent)', marginBottom: '0.5rem' }}>
          Nota para macOS
        </div>
        <p style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', lineHeight: 1.6 }}>
          La primera vez que lo abras, macOS puede bloquearlo. Ve a <strong style={{ color: 'var(--text-primary)' }}>Ajustes del Sistema → Privacidad y Seguridad</strong> y haz clic en "Abrir de todos modos". También puedes ejecutar en terminal:
        </p>
        <div style={{
          background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.6rem 0.85rem',
          fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--accent)', marginTop: '0.5rem',
        }}>
          chmod +x spc-macos-arm64 && xattr -d com.apple.quarantine spc-macos-arm64
        </div>
      </div>

      {/* Security */}
      <div style={{
        background: 'rgba(239, 68, 68, 0.08)', border: '1px solid rgba(239, 68, 68, 0.25)',
        borderRadius: '12px', padding: '1.25rem', marginBottom: '1.5rem',
      }}>
        <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--red)', marginBottom: '0.5rem' }}>
          Seguridad importante
        </div>
        <ul style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', lineHeight: 1.7, paddingLeft: '1.25rem' }}>
          <li><strong style={{ color: 'var(--text-primary)' }}>GUARDA tu clave privada en papel</strong> — si la pierdes, pierdes tus fondos para siempre</li>
          <li><strong style={{ color: 'var(--text-primary)' }}>NUNCA la compartas</strong> con nadie — ni con nosotros, ni por Telegram, ni por email</li>
          <li>Tu dirección pública (SPCxxx...) sí puedes compartirla — es como un número de cuenta</li>
        </ul>
      </div>

      {/* Community link */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem', textAlign: 'center',
      }}>
        <h2 style={{ fontSize: '1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
          ¿Necesitas ayuda?
        </h2>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '1rem' }}>
          Únete al Telegram y te ayudamos paso a paso
        </p>
        <a href="https://t.me/spaincoin" target="_blank" rel="noopener noreferrer" style={{
          display: 'inline-block', padding: '0.65rem 1.5rem',
          background: '#0088cc', border: 'none', borderRadius: '8px',
          color: '#fff', fontSize: '0.85rem', fontWeight: '600', textDecoration: 'none',
        }}>Telegram SpainCoin</a>
      </div>
    </div>
  )
}
