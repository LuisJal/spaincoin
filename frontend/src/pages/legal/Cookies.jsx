export default function Cookies({ onNavigate }) {
  const containerStyle = {
    minHeight: 'calc(100vh - 60px)',
    padding: '2rem 1rem 4rem',
    display: 'flex',
    justifyContent: 'center',
  }

  const contentStyle = {
    width: '100%',
    maxWidth: '800px',
  }

  const backStyle = {
    display: 'inline-flex',
    alignItems: 'center',
    gap: '0.4rem',
    color: 'var(--accent)',
    fontSize: '0.875rem',
    cursor: 'pointer',
    background: 'none',
    border: 'none',
    padding: 0,
    marginBottom: '2rem',
    textDecoration: 'none',
  }

  const titleStyle = {
    fontSize: '1.75rem',
    fontWeight: '700',
    color: 'var(--text-primary)',
    marginBottom: '0.5rem',
    lineHeight: 1.2,
  }

  const metaStyle = {
    fontSize: '0.85rem',
    color: 'var(--text-secondary)',
    marginBottom: '2.5rem',
    paddingBottom: '1.5rem',
    borderBottom: '1px solid var(--border)',
  }

  const sectionTitleStyle = {
    fontSize: '1.05rem',
    fontWeight: '600',
    color: 'var(--accent)',
    marginTop: '2rem',
    marginBottom: '0.75rem',
  }

  const textStyle = {
    fontSize: '0.9rem',
    color: 'var(--text-secondary)',
    lineHeight: '1.75',
    marginBottom: '0.75rem',
  }

  const tableStyle = {
    width: '100%',
    borderCollapse: 'collapse',
    marginBottom: '1.25rem',
    fontSize: '0.875rem',
  }

  const thStyle = {
    background: 'var(--bg-secondary)',
    color: 'var(--text-primary)',
    padding: '0.65rem 1rem',
    textAlign: 'left',
    fontWeight: '600',
    fontSize: '0.8rem',
    textTransform: 'uppercase',
    letterSpacing: '0.05em',
    borderBottom: '1px solid var(--border)',
  }

  const tdStyle = {
    padding: '0.65rem 1rem',
    color: 'var(--text-secondary)',
    borderBottom: '1px solid var(--border)',
    verticalAlign: 'top',
    lineHeight: '1.5',
  }

  const greenBoxStyle = {
    background: 'rgba(16, 185, 129, 0.08)',
    border: '1px solid rgba(16, 185, 129, 0.25)',
    borderRadius: '10px',
    padding: '1rem 1.25rem',
    marginBottom: '1.5rem',
    fontSize: '0.875rem',
    color: '#10b981',
    lineHeight: '1.65',
  }

  const browserStepStyle = {
    background: 'var(--bg-secondary)',
    borderRadius: '8px',
    padding: '0.75rem 1rem',
    marginBottom: '0.5rem',
    fontSize: '0.85rem',
    color: 'var(--text-secondary)',
    lineHeight: '1.6',
  }

  return (
    <div className="page-enter" style={containerStyle}>
      <div style={contentStyle}>
        {/* Back button */}
        <button
          style={backStyle}
          onClick={() => onNavigate ? onNavigate(-1) : window.history.back()}
        >
          ← Volver
        </button>

        {/* Header */}
        <h1 style={titleStyle}>Política de Cookies</h1>
        <p style={metaStyle}>
          SpainCoin Exchange · Última actualización: marzo 2026
        </p>

        {/* Green summary box */}
        <div style={greenBoxStyle}>
          <strong>Resumen:</strong> SpainCoin Exchange utiliza exclusivamente cookies técnicas
          necesarias para el funcionamiento del servicio. No utilizamos cookies publicitarias,
          de seguimiento ni analíticas. No compartimos datos de cookies con terceros.
        </div>

        {/* Section 1 */}
        <h2 style={sectionTitleStyle}>1. ¿Qué son las cookies?</h2>
        <p style={textStyle}>
          Las cookies son pequeños archivos de texto que los sitios web almacenan en el
          dispositivo del usuario cuando este los visita. Permiten que el sitio web recuerde
          información sobre su visita, como el idioma preferido y otras opciones, lo que
          facilita su próxima visita y hace que el sitio le resulte más útil.
        </p>
        <p style={textStyle}>
          La normativa aplicable en España incluye el art. 22.2 de la Ley de Servicios de la
          Sociedad de la Información (LSSI) y el RGPD (UE) 2016/679.
        </p>

        {/* Section 2 */}
        <h2 style={sectionTitleStyle}>2. Cookies que utilizamos</h2>
        <p style={textStyle}>
          SpainCoin Exchange utiliza únicamente las siguientes cookies técnicas, estrictamente
          necesarias para el funcionamiento de la plataforma:
        </p>
        <table style={tableStyle}>
          <thead>
            <tr>
              <th style={thStyle}>Nombre</th>
              <th style={thStyle}>Tipo</th>
              <th style={thStyle}>Finalidad</th>
              <th style={thStyle}>Duración</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td style={{ ...tdStyle, fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--text-primary)' }}>
                spc_auth_token
              </td>
              <td style={tdStyle}>Sesión / Autenticación</td>
              <td style={tdStyle}>
                Almacena el token JWT cifrado para mantener la sesión iniciada del usuario.
                Sin esta cookie, no es posible acceder a las funciones autenticadas.
              </td>
              <td style={tdStyle}>Sesión (se elimina al cerrar el navegador) o hasta 7 días si el usuario selecciona "Recordarme"</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--text-primary)' }}>
                spc_cookies_accepted
              </td>
              <td style={tdStyle}>Preferencias</td>
              <td style={tdStyle}>
                Almacena si el usuario ha aceptado la política de cookies para evitar mostrar
                el banner repetidamente.
              </td>
              <td style={tdStyle}>1 año</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--text-primary)' }}>
                spc_theme
              </td>
              <td style={tdStyle}>Preferencias</td>
              <td style={tdStyle}>
                Almacena la preferencia de tema del usuario (actualmente solo tema oscuro).
                Reservado para uso futuro.
              </td>
              <td style={tdStyle}>1 año</td>
            </tr>
          </tbody>
        </table>
        <p style={textStyle}>
          Además, utilizamos <strong style={{ color: 'var(--text-primary)' }}>localStorage</strong> del
          navegador para almacenar preferencias de sesión del lado del cliente. El localStorage
          no es técnicamente una cookie, pero tiene funcionalidad similar para guardar datos
          localmente en el dispositivo del usuario.
        </p>

        {/* Section 3 */}
        <h2 style={sectionTitleStyle}>3. Cookies que NO utilizamos</h2>
        <p style={textStyle}>
          SpainCoin Exchange no utiliza ninguno de los siguientes tipos de cookies:
        </p>
        <ul style={{ ...textStyle, paddingLeft: '1.5rem', listStyleType: 'disc' }}>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Cookies publicitarias:</strong> No
            mostramos publicidad personalizada ni utilizamos redes publicitarias.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Cookies de seguimiento/rastreo:</strong> No
            rastreamos el comportamiento del usuario entre distintos sitios web.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Cookies analíticas de terceros:</strong> No
            utilizamos Google Analytics, Hotjar, Mixpanel ni herramientas similares de análisis
            de comportamiento.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Cookies de redes sociales:</strong> No
            integramos botones de compartir ni píxeles de redes sociales (Facebook Pixel, etc.).
          </li>
        </ul>

        {/* Section 4 */}
        <h2 style={sectionTitleStyle}>4. Cómo gestionar y eliminar las cookies</h2>
        <p style={textStyle}>
          Puede gestionar o eliminar las cookies en cualquier momento desde la configuración de
          su navegador. Tenga en cuenta que deshabilitar las cookies técnicas puede afectar al
          funcionamiento correcto de la plataforma, incluyendo la posibilidad de mantener la
          sesión iniciada.
        </p>

        <p style={{ ...textStyle, fontWeight: '600', color: 'var(--text-primary)', marginTop: '1.25rem', marginBottom: '0.5rem' }}>
          Google Chrome
        </p>
        <div style={browserStepStyle}>
          Menú (tres puntos) → <strong>Configuración</strong> → <strong>Privacidad y seguridad</strong> →{' '}
          <strong>Cookies y otros datos de sitios</strong> → <strong>Ver todos los datos y permisos de sitios</strong> →
          Busque "spaincoin" y elimine las cookies asociadas.
        </div>

        <p style={{ ...textStyle, fontWeight: '600', color: 'var(--text-primary)', marginTop: '1rem', marginBottom: '0.5rem' }}>
          Mozilla Firefox
        </p>
        <div style={browserStepStyle}>
          Menú (tres rayas) → <strong>Configuración</strong> → <strong>Privacidad y seguridad</strong> →
          Sección "Cookies y datos del sitio" → <strong>Gestionar datos</strong> → Busque "spaincoin"
          y elimínelos.
        </div>

        <p style={{ ...textStyle, fontWeight: '600', color: 'var(--text-primary)', marginTop: '1rem', marginBottom: '0.5rem' }}>
          Safari (macOS / iOS)
        </p>
        <div style={browserStepStyle}>
          <strong>Preferencias</strong> (macOS) / <strong>Ajustes</strong> (iOS) → <strong>Privacidad</strong> →{' '}
          <strong>Gestionar datos del sitio web</strong> → Busque "spaincoin" y elimínelo.
        </div>

        <p style={{ ...textStyle, fontWeight: '600', color: 'var(--text-primary)', marginTop: '1rem', marginBottom: '0.5rem' }}>
          Microsoft Edge
        </p>
        <div style={browserStepStyle}>
          Menú (tres puntos) → <strong>Configuración</strong> → <strong>Cookies y permisos del sitio</strong> →{' '}
          <strong>Administrar y eliminar cookies y datos del sitio</strong> → Busque "spaincoin".
        </div>

        {/* Section 5 */}
        <h2 style={sectionTitleStyle}>5. Cambios en esta política</h2>
        <p style={textStyle}>
          SpainCoin Exchange se reserva el derecho a actualizar esta Política de Cookies en
          cualquier momento. Los cambios se notificarán a través de la plataforma. Si continúa
          usando la plataforma tras la publicación de los cambios, se considerará que los ha
          aceptado.
        </p>

        {/* Footer note */}
        <div style={{
          marginTop: '3rem',
          paddingTop: '1.5rem',
          borderTop: '1px solid var(--border)',
          fontSize: '0.8rem',
          color: 'var(--text-secondary)',
          lineHeight: '1.6',
        }}>
          Para consultas sobre cookies: <span style={{ color: 'var(--accent)' }}>legal@spaincoin.com</span>
          <br />
          SpainCoin Exchange · Proyecto en desarrollo · Testnet v0.1 · Marzo 2026
        </div>
      </div>
    </div>
  )
}
