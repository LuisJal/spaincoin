export default function WalletDownload({ onNavigate }) {
  return (
    <div className="page-enter" style={{ maxWidth: '800px', margin: '0 auto', padding: '2rem 1.5rem' }}>
      <h1 style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
        Wallet SpainCoin
      </h1>
      <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '2rem' }}>
        Tu wallet, tus claves, tus fondos. SpainCoin es 100% self-custody.
      </p>

      {/* Download options */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '1.5rem',
      }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '1rem' }}>
          CLI Wallet (recomendado)
        </h2>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.6, marginBottom: '1.25rem' }}>
          Wallet por linea de comandos. Maxima seguridad, sin interfaz web. Tus claves nunca salen de tu ordenador.
        </p>

        <div style={{
          background: 'var(--bg-secondary)', borderRadius: '8px', padding: '1rem',
          fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--accent)',
          overflowX: 'auto', marginBottom: '1rem',
        }}>
          <div style={{ color: 'var(--text-secondary)', marginBottom: '0.5rem' }}># Instalar Go y compilar</div>
          <div>git clone https://github.com/spaincoin/spaincoin</div>
          <div>cd spaincoin</div>
          <div>go build -o spc ./cli/</div>
          <div style={{ marginTop: '0.75rem', color: 'var(--text-secondary)' }}># Crear nueva wallet</div>
          <div>./spc wallet new</div>
          <div style={{ marginTop: '0.75rem', color: 'var(--text-secondary)' }}># Ver balance</div>
          <div>./spc wallet balance</div>
          <div style={{ marginTop: '0.75rem', color: 'var(--text-secondary)' }}># Enviar SPC</div>
          <div>./spc send --to SPCxxx... --amount 10</div>
        </div>

        <a href="https://github.com/spaincoin/spaincoin" target="_blank" rel="noopener noreferrer" style={{
          display: 'inline-block', padding: '0.6rem 1.5rem',
          background: 'var(--accent)', border: 'none', borderRadius: '8px',
          color: '#fff', fontSize: '0.85rem', fontWeight: '600', textDecoration: 'none',
        }}>Ver en GitHub</a>
      </div>

      {/* Security warning */}
      <div style={{
        background: 'rgba(239, 68, 68, 0.08)', border: '1px solid rgba(239, 68, 68, 0.25)',
        borderRadius: '12px', padding: '1.25rem', marginBottom: '1.5rem',
      }}>
        <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--red)', marginBottom: '0.5rem' }}>
          Seguridad
        </div>
        <ul style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', lineHeight: 1.7, paddingLeft: '1.25rem' }}>
          <li><strong style={{ color: 'var(--text-primary)' }}>Nunca compartas tu clave privada</strong> con nadie</li>
          <li>Guarda tu clave en papel o en un gestor de contrasenas cifrado</li>
          <li>No existe recuperacion de claves — si la pierdes, pierdes tus fondos</li>
          <li>SpainCoin nunca te pedira tu clave privada</li>
        </ul>
      </div>

      {/* How it works */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.5rem',
      }}>
        <h2 style={{ fontSize: '1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '1rem' }}>
          Como funciona
        </h2>
        <div style={{ display: 'grid', gap: '1rem' }}>
          {[
            { step: '1', title: 'Genera tu wallet', desc: 'El CLI crea un par de claves (publica + privada) en tu ordenador. La clave privada nunca sale de tu maquina.' },
            { step: '2', title: 'Recibe SPC', desc: 'Comparte tu direccion publica (SPCxxx...) para que te envien SPC. Es como tu numero de cuenta.' },
            { step: '3', title: 'Envia SPC', desc: 'Firma transacciones con tu clave privada desde el CLI. Solo tu puedes mover tus fondos.' },
          ].map((s, i) => (
            <div key={i} style={{ display: 'flex', gap: '1rem', alignItems: 'flex-start' }}>
              <div style={{
                width: '32px', height: '32px', borderRadius: '50%',
                background: 'var(--accent)', color: '#fff', display: 'flex',
                alignItems: 'center', justifyContent: 'center', fontWeight: '700',
                fontSize: '0.85rem', flexShrink: 0,
              }}>{s.step}</div>
              <div>
                <div style={{ fontWeight: '600', fontSize: '0.9rem', color: 'var(--text-primary)', marginBottom: '0.2rem' }}>{s.title}</div>
                <div style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', lineHeight: 1.5 }}>{s.desc}</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
