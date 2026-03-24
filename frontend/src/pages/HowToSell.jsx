export default function HowToSell({ onNavigate }) {
  return (
    <div className="page-enter" style={{ maxWidth: '700px', margin: '0 auto', padding: '2rem 1.5rem' }}>
      <h1 style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
        Cómo vender $SPC
      </h1>
      <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '2rem', lineHeight: 1.6 }}>
        Vender tus SPC y recibir euros en tu cuenta bancaria. Paso a paso.
      </p>

      {/* Overview */}
      <div style={{
        background: 'rgba(16, 185, 129, 0.08)', border: '1px solid rgba(16, 185, 129, 0.25)',
        borderRadius: '12px', padding: '1.25rem', marginBottom: '1.5rem',
      }}>
        <div style={{ fontWeight: '600', color: 'var(--green)', marginBottom: '0.5rem' }}>Resumen rápido</div>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.7 }}>
          Le dices al bot cuántos SPC quieres vender → envías los SPC a la dirección que te indica →
          verificamos la recepción → te hacemos transferencia bancaria.
        </p>
      </div>

      {/* Steps */}
      {[
        {
          step: '1',
          title: 'Abre el bot de Telegram',
          content: 'Habla con @spaincoin_bot en Telegram y pulsa "Vender SPC" o escribe /vender seguido de la cantidad.',
          example: '/vender 100',
          note: 'El bot te dirá cuántos euros recibirás al precio actual.',
        },
        {
          step: '2',
          title: 'Envía los SPC',
          content: 'El bot te dará una dirección SPC donde enviar tus monedas. Desde tu ordenador, abre una terminal y ejecuta:',
          code: './spc send --to SPCc119f94a...d65481 --amount 100 --node http://204.168.176.40:8545',
          note: 'Necesitas tu clave privada para firmar la transacción. Solo se usa en tu ordenador, nunca la compartas.',
        },
        {
          step: '3',
          title: 'Espera la confirmación',
          content: 'Un admin verificará que los SPC han llegado. Normalmente en menos de 10 minutos.',
          note: null,
        },
        {
          step: '4',
          title: 'Recibe tus euros',
          content: 'Te haremos una transferencia bancaria por el importe correspondiente. El bot te notificará cuando esté hecho.',
          note: 'Para recibir la transferencia, el admin te pedirá tu IBAN por mensaje privado.',
        },
      ].map((s, i) => (
        <div key={i} style={{
          background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
          padding: '1.25rem', marginBottom: '1rem',
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem', marginBottom: '0.75rem' }}>
            <div style={{
              width: '28px', height: '28px', borderRadius: '50%',
              background: 'var(--accent)', color: '#fff', display: 'flex',
              alignItems: 'center', justifyContent: 'center', fontWeight: '700',
              fontSize: '0.8rem', flexShrink: 0,
            }}>{s.step}</div>
            <span style={{ fontWeight: '600', fontSize: '0.95rem', color: 'var(--text-primary)' }}>{s.title}</span>
          </div>
          <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.7, marginBottom: s.example || s.code ? '0.75rem' : '0' }}>
            {s.content}
          </p>
          {s.example && (
            <div style={{
              background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.6rem 0.85rem',
              fontFamily: 'monospace', fontSize: '0.82rem', color: 'var(--accent)', marginBottom: '0.5rem',
            }}>{s.example}</div>
          )}
          {s.code && (
            <div style={{
              background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.75rem 0.85rem',
              fontFamily: 'monospace', fontSize: '0.72rem', color: 'var(--accent)',
              wordBreak: 'break-all', lineHeight: 1.6, marginBottom: '0.5rem',
            }}>{s.code}</div>
          )}
          {s.note && (
            <p style={{ fontSize: '0.78rem', color: 'var(--text-secondary)', fontStyle: 'italic' }}>
              {s.note}
            </p>
          )}
        </div>
      ))}

      {/* FAQ */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1.25rem', marginBottom: '1.5rem',
      }}>
        <div style={{ fontWeight: '600', color: 'var(--text-primary)', marginBottom: '1rem' }}>Preguntas frecuentes</div>
        {[
          { q: '¿Puedo vender desde el móvil?', a: 'De momento necesitas un ordenador para enviar SPC (firmar la transacción). Estamos trabajando para que sea posible desde el móvil en el futuro.' },
          { q: '¿Cuánto tarda?', a: 'Una vez que envías los SPC, la verificación y transferencia suelen tardar menos de 10 minutos.' },
          { q: '¿Hay comisión?', a: 'No cobramos comisión por vender. El precio es el mismo que para comprar.' },
          { q: '¿Cuánto puedo vender?', a: 'Todo lo que tengas en tu wallet. No hay límite mínimo ni máximo.' },
          { q: '¿Qué pasa si envío los SPC y no recibo el pago?', a: 'Todas las transacciones quedan registradas en la blockchain. Si hay cualquier problema, contacta con un admin en Telegram.' },
        ].map((faq, i) => (
          <div key={i} style={{ marginBottom: i < 4 ? '1rem' : 0 }}>
            <div style={{ fontSize: '0.85rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.25rem' }}>{faq.q}</div>
            <div style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', lineHeight: 1.6 }}>{faq.a}</div>
          </div>
        ))}
      </div>

      {/* CTA */}
      <div style={{ display: 'flex', gap: '0.5rem' }}>
        <a href="https://t.me/spaincoin_bot" target="_blank" rel="noopener noreferrer" style={{
          flex: 1, padding: '0.75rem', borderRadius: '10px', border: 'none',
          background: 'var(--accent)', color: '#fff', fontSize: '0.9rem',
          fontWeight: '700', textAlign: 'center', textDecoration: 'none',
        }}>Abrir bot de Telegram</a>
        <button onClick={() => onNavigate('/wallet')} style={{
          flex: 1, padding: '0.75rem', borderRadius: '10px',
          border: '1px solid var(--border)', background: 'transparent',
          color: 'var(--text-secondary)', fontSize: '0.9rem', cursor: 'pointer',
        }}>Mi wallet</button>
      </div>
    </div>
  )
}
