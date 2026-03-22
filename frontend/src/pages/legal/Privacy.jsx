export default function Privacy({ onNavigate }) {
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

  const infoBoxStyle = {
    background: 'rgba(59, 130, 246, 0.08)',
    border: '1px solid rgba(59, 130, 246, 0.25)',
    borderRadius: '10px',
    padding: '1rem 1.25rem',
    marginBottom: '1.5rem',
    fontSize: '0.875rem',
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
        <h1 style={titleStyle}>Política de Privacidad</h1>
        <p style={metaStyle}>
          SpainCoin Exchange · Última actualización: marzo 2026 · Conforme al RGPD (UE) 2016/679
        </p>

        {/* Intro box */}
        <div style={infoBoxStyle}>
          Esta Política de Privacidad describe cómo SpainCoin Exchange recoge, trata y protege
          los datos personales de sus usuarios, en cumplimiento del Reglamento General de
          Protección de Datos (RGPD) y la Ley Orgánica 3/2018, de 5 de diciembre, de
          Protección de Datos Personales y garantía de los derechos digitales (LOPDGDD).
        </div>

        {/* Section 1 */}
        <h2 style={sectionTitleStyle}>1. Responsable del tratamiento</h2>
        <table style={tableStyle}>
          <tbody>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)', width: '35%' }}>Denominación</td>
              <td style={tdStyle}>SpainCoin Exchange</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)' }}>Correo de contacto</td>
              <td style={tdStyle}>legal@spaincoin.com</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)' }}>DPD / DPO</td>
              <td style={tdStyle}>En proceso de designación (próximamente en mainnet)</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)' }}>Marco normativo</td>
              <td style={tdStyle}>RGPD (UE) 2016/679 · LOPDGDD 3/2018</td>
            </tr>
          </tbody>
        </table>

        {/* Section 2 */}
        <h2 style={sectionTitleStyle}>2. Datos que se recogen</h2>
        <p style={textStyle}>
          SpainCoin Exchange recoge únicamente los datos estrictamente necesarios para la
          prestación del servicio:
        </p>
        <table style={tableStyle}>
          <thead>
            <tr>
              <th style={thStyle}>Dato</th>
              <th style={thStyle}>Descripción</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)' }}>Dirección de email</td>
              <td style={tdStyle}>Identificador de cuenta. Necesario para el acceso y notificaciones.</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)' }}>Contraseña</td>
              <td style={tdStyle}>Almacenada únicamente en forma de hash (bcrypt). No es recuperable.</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)' }}>Dirección de wallet</td>
              <td style={tdStyle}>Dirección pública blockchain (0x...). Identificador de cuenta en la red.</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)' }}>Clave privada cifrada</td>
              <td style={tdStyle}>Si el usuario opta por custodia delegada, se almacena cifrada con su contraseña. Nunca en texto plano.</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)' }}>Dirección IP de acceso</td>
              <td style={tdStyle}>Registrada en logs de seguridad para prevención de fraude y accesos no autorizados.</td>
            </tr>
            <tr>
              <td style={{ ...tdStyle, fontWeight: '600', color: 'var(--text-primary)' }}>Historial de transacciones</td>
              <td style={tdStyle}>Transacciones blockchain realizadas a través de la plataforma. Inherentemente públicas en la cadena.</td>
            </tr>
          </tbody>
        </table>
        <p style={textStyle}>
          <strong style={{ color: 'var(--text-primary)' }}>Importante:</strong> SpainCoin Exchange
          no recoge datos de tarjetas de crédito, documentos de identidad (DNI/pasaporte) ni
          información bancaria durante la fase Testnet.
        </p>

        {/* Section 3 */}
        <h2 style={sectionTitleStyle}>3. Finalidad del tratamiento</h2>
        <ul style={{ ...textStyle, paddingLeft: '1.5rem', listStyleType: 'disc' }}>
          <li style={{ marginBottom: '0.4rem' }}>Gestión de cuentas de usuario y autenticación en la plataforma.</li>
          <li style={{ marginBottom: '0.4rem' }}>Ejecución de transacciones blockchain solicitadas por el usuario.</li>
          <li style={{ marginBottom: '0.4rem' }}>Seguridad de la plataforma: detección de accesos no autorizados, fraude y actividad maliciosa.</li>
          <li style={{ marginBottom: '0.4rem' }}>Comunicaciones relacionadas con el servicio (cambios relevantes, avisos de seguridad).</li>
          <li style={{ marginBottom: '0.4rem' }}>Cumplimiento de obligaciones legales aplicables.</li>
        </ul>

        {/* Section 4 */}
        <h2 style={sectionTitleStyle}>4. Base jurídica del tratamiento</h2>
        <table style={tableStyle}>
          <thead>
            <tr>
              <th style={thStyle}>Tratamiento</th>
              <th style={thStyle}>Base jurídica (art. 6 RGPD)</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td style={tdStyle}>Gestión de cuenta y prestación del servicio</td>
              <td style={tdStyle}>Art. 6.1.b — Ejecución de un contrato en el que el interesado es parte</td>
            </tr>
            <tr>
              <td style={tdStyle}>Seguridad y prevención de fraude</td>
              <td style={tdStyle}>Art. 6.1.f — Interés legítimo del responsable</td>
            </tr>
            <tr>
              <td style={tdStyle}>Conservación por obligaciones legales</td>
              <td style={tdStyle}>Art. 6.1.c — Cumplimiento de una obligación legal</td>
            </tr>
            <tr>
              <td style={tdStyle}>Comunicaciones comerciales (si aplica)</td>
              <td style={tdStyle}>Art. 6.1.a — Consentimiento del interesado</td>
            </tr>
          </tbody>
        </table>

        {/* Section 5 */}
        <h2 style={sectionTitleStyle}>5. Derechos RGPD del usuario</h2>
        <p style={textStyle}>
          De conformidad con el RGPD, el Usuario tiene derecho a:
        </p>
        <ul style={{ ...textStyle, paddingLeft: '1.5rem', listStyleType: 'disc' }}>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Acceso (art. 15):</strong> Obtener
            confirmación sobre si tratamos sus datos y acceder a una copia.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Rectificación (art. 16):</strong> Solicitar
            la corrección de datos inexactos o incompletos.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Supresión (art. 17):</strong> Solicitar
            la eliminación de sus datos cuando ya no sean necesarios o retire su consentimiento
            (sujeto a obligaciones legales de conservación).
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Portabilidad (art. 20):</strong> Recibir
            sus datos en formato estructurado y de uso común, y transmitirlos a otro responsable.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Oposición (art. 21):</strong> Oponerse
            al tratamiento basado en interés legítimo, incluyendo la elaboración de perfiles.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Limitación (art. 18):</strong> Solicitar
            la restricción del tratamiento en determinadas circunstancias.
          </li>
        </ul>
        <p style={textStyle}>
          Para ejercer cualquiera de estos derechos, contacte con nosotros en:
          <span style={{ color: 'var(--accent)', marginLeft: '0.3rem' }}>legal@spaincoin.com</span>
        </p>
        <p style={textStyle}>
          Si considera que el tratamiento de sus datos no es conforme al RGPD, tiene derecho a
          presentar una reclamación ante la Agencia Española de Protección de Datos (AEPD):
          <span style={{ color: 'var(--accent)', marginLeft: '0.3rem' }}>www.aepd.es</span>
        </p>

        {/* Section 6 */}
        <h2 style={sectionTitleStyle}>6. Conservación de los datos</h2>
        <table style={tableStyle}>
          <thead>
            <tr>
              <th style={thStyle}>Tipo de dato</th>
              <th style={thStyle}>Período de conservación</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td style={tdStyle}>Datos de cuenta activa</td>
              <td style={tdStyle}>Durante la vigencia de la cuenta de usuario</td>
            </tr>
            <tr>
              <td style={tdStyle}>Datos tras cancelación de cuenta</td>
              <td style={tdStyle}>5 años por obligaciones legales (normativa contable y AML)</td>
            </tr>
            <tr>
              <td style={tdStyle}>Logs de seguridad (IPs)</td>
              <td style={tdStyle}>12 meses</td>
            </tr>
            <tr>
              <td style={tdStyle}>Historial de transacciones blockchain</td>
              <td style={tdStyle}>Permanente (inherente a la naturaleza inmutable de la blockchain)</td>
            </tr>
          </tbody>
        </table>

        {/* Section 7 */}
        <h2 style={sectionTitleStyle}>7. Transferencias internacionales de datos</h2>
        <p style={textStyle}>
          SpainCoin Exchange no realiza transferencias internacionales de datos personales a
          países fuera del Espacio Económico Europeo (EEE). Toda la infraestructura de tratamiento
          de datos se encuentra dentro del territorio de la Unión Europea.
        </p>

        {/* Section 8 */}
        <h2 style={sectionTitleStyle}>8. Cookies</h2>
        <p style={textStyle}>
          SpainCoin Exchange utiliza únicamente cookies técnicas estrictamente necesarias para el
          funcionamiento del servicio (gestión de sesión, autenticación JWT). No se utilizan
          cookies publicitarias, de perfilado ni de rastreo. Para más información, consulte nuestra{' '}
          <button
            onClick={() => onNavigate && onNavigate('/legal/cookies')}
            style={{ background: 'none', border: 'none', color: 'var(--accent)', cursor: 'pointer', fontSize: 'inherit', padding: 0 }}
          >
            Política de Cookies
          </button>.
        </p>

        {/* Section 9 */}
        <h2 style={sectionTitleStyle}>9. Seguridad de los datos</h2>
        <p style={textStyle}>
          Aplicamos medidas técnicas y organizativas apropiadas para proteger sus datos personales
          frente a accesos no autorizados, pérdida, destrucción o divulgación accidental, de
          conformidad con el art. 32 RGPD. Entre estas medidas se incluyen: cifrado de contraseñas
          mediante bcrypt, cifrado de claves privadas, comunicaciones cifradas mediante TLS/HTTPS,
          y acceso restringido a los sistemas de producción.
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
          Contacto para asuntos de privacidad: <span style={{ color: 'var(--accent)' }}>legal@spaincoin.com</span>
          <br />
          SpainCoin Exchange · Proyecto en desarrollo · Testnet v0.1 · Marzo 2026
        </div>
      </div>
    </div>
  )
}
