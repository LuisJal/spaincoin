export default function Risk({ onNavigate }) {
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
    marginBottom: '2rem',
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
        <h1 style={titleStyle}>Advertencia de Riesgos</h1>
        <p style={metaStyle}>
          SpainCoin Exchange · Última actualización: marzo 2026
        </p>

        {/* Main risk warning box */}
        <div style={{
          background: 'rgba(239, 68, 68, 0.08)',
          border: '2px solid rgba(239, 68, 68, 0.5)',
          borderRadius: '12px',
          padding: '1.5rem',
          marginBottom: '2.5rem',
        }}>
          <div style={{
            fontSize: '1rem',
            fontWeight: '700',
            color: '#ef4444',
            marginBottom: '0.75rem',
            display: 'flex',
            alignItems: 'center',
            gap: '0.5rem',
          }}>
            ADVERTENCIA DE RIESGO ELEVADO
          </div>
          <p style={{ fontSize: '0.9rem', color: '#fca5a5', lineHeight: '1.7', marginBottom: '0.75rem' }}>
            Los activos digitales, incluido $SPC (SpainCoin), son instrumentos de alto riesgo y
            alta volatilidad. <strong>Puede perder la totalidad del capital invertido.</strong> Solo
            opere con fondos que pueda permitirse perder completamente.
          </p>
          <p style={{ fontSize: '0.875rem', color: '#ef4444', lineHeight: '1.6', margin: 0 }}>
            Esta plataforma se encuentra en fase TESTNET. Los activos actuales no tienen valor
            económico real. El lanzamiento en mainnet no garantiza que $SPC adquiera valor
            de mercado.
          </p>
        </div>

        {/* Section 1 */}
        <h2 style={sectionTitleStyle}>1. Volatilidad extrema de los activos digitales</h2>
        <p style={textStyle}>
          Los activos digitales y criptomonedas están sujetos a fluctuaciones de precio de una
          magnitud muy superior a la de los activos financieros tradicionales. Es posible que el
          valor de $SPC o cualquier otro activo digital experimente caídas del 50%, 80% o incluso
          del 100% en períodos muy cortos de tiempo, sin necesidad de que exista una causa
          técnica o fundamental que lo justifique.
        </p>
        <p style={textStyle}>
          Los mercados de criptomonedas operan 24 horas al día, 7 días a la semana, sin
          mecanismos de cierre de mercado ni cortocircuitos (circuit breakers) como los que
          existen en los mercados regulados de valores.
        </p>

        {/* Section 2 */}
        <h2 style={sectionTitleStyle}>2. Pérdida total del capital</h2>
        <p style={textStyle}>
          A diferencia de los depósitos bancarios cubiertos por el Fondo de Garantía de
          Depósitos (FGD), los activos digitales en SpainCoin Exchange <strong style={{ color: 'var(--text-primary)' }}>
          no están protegidos por ningún fondo de garantía ni seguro de depósitos</strong>. En
          caso de quiebra del operador, hackeo, fallo técnico irreversible o cierre de la
          plataforma, podría perder la totalidad de sus fondos sin derecho a compensación.
        </p>

        {/* Section 3 */}
        <h2 style={sectionTitleStyle}>3. Naturaleza experimental de SpainCoin ($SPC)</h2>
        <p style={textStyle}>
          SpainCoin ($SPC) es un proyecto blockchain en fase de desarrollo activo. Esto implica:
        </p>
        <ul style={{ ...textStyle, paddingLeft: '1.5rem', listStyleType: 'disc' }}>
          <li style={{ marginBottom: '0.4rem' }}>
            El protocolo puede contener errores de software (bugs) o vulnerabilidades de seguridad
            que podrían explotarse en detrimento de los usuarios.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            El equipo de desarrollo puede introducir cambios fundamentales en el protocolo
            (hard forks) que alteren el valor, la distribución o las características técnicas de $SPC.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            El proyecto podría no completar su hoja de ruta o ser abandonado, resultando en
            la pérdida total del valor del activo.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            La red puede sufrir ataques del 51%, ataques Sybil u otros vectores de ataque
            propios de redes blockchain con escasa descentralización inicial.
          </li>
        </ul>

        {/* Section 4 */}
        <h2 style={sectionTitleStyle}>4. $SPC no es dinero de curso legal</h2>
        <p style={textStyle}>
          $SPC no tiene la condición de moneda de curso legal en España ni en ningún otro
          país. No está emitido ni respaldado por ningún banco central, gobierno ni institución
          financiera. Su valor, si lo tuviera, derivaría exclusivamente de la oferta y demanda
          del mercado, sin ningún activo subyacente que lo respalde.
        </p>
        <p style={textStyle}>
          Ningún comerciante u operador está obligado a aceptar $SPC como medio de pago de
          bienes o servicios.
        </p>

        {/* Section 5 */}
        <h2 style={sectionTitleStyle}>5. No es asesoramiento financiero</h2>
        <p style={textStyle}>
          Ninguna información, análisis, gráfica, precio, estimación ni otro contenido
          disponible en SpainCoin Exchange constituye asesoramiento financiero, de inversión,
          fiscal o legal. SpainCoin Exchange no está autorizado como Empresa de Servicios de
          Inversión (ESI) conforme a la normativa MiFID II.
        </p>
        <p style={textStyle}>
          Antes de operar con activos digitales, el usuario debería considerar consultar con un
          asesor financiero independiente, debidamente regulado por la CNMV (Comisión Nacional
          del Mercado de Valores) u organismo equivalente de su país de residencia.
        </p>

        {/* Section 6 */}
        <h2 style={sectionTitleStyle}>6. Marco regulatorio MiCA y legislación aplicable</h2>
        <p style={textStyle}>
          El Reglamento (UE) 2023/1114, relativo a los Mercados de Criptoactivos (MiCA), entró
          en vigor progresivamente a partir de 2024. SpainCoin Exchange es consciente de este
          marco regulatorio y se compromete a analizar su aplicabilidad al proyecto conforme
          se avance hacia mainnet.
        </p>
        <p style={textStyle}>
          En la fase Testnet actual, la plataforma opera exclusivamente con activos ficticios
          sin valor económico, lo que limita el alcance inmediato de la regulación MiCA. No
          obstante, el Usuario debe ser consciente de que la evolución regulatoria podría
          afectar a la disponibilidad del servicio en determinados territorios.
        </p>
        <p style={textStyle}>
          Las ganancias derivadas de operaciones con criptomonedas pueden estar sujetas a
          tributación. Consulte la normativa fiscal de su país y, en caso de duda, asesórese
          con un profesional tributario.
        </p>

        {/* Section 7 */}
        <h2 style={sectionTitleStyle}>7. Riesgos operativos y tecnológicos</h2>
        <ul style={{ ...textStyle, paddingLeft: '1.5rem', listStyleType: 'disc' }}>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Pérdida de clave privada:</strong> Si
            pierde su clave privada, perderá acceso permanente e irrecuperable a sus fondos.
            No existe mecanismo de recuperación.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Irreversibilidad de transacciones:</strong> Las
            transacciones confirmadas en la blockchain son irreversibles. No es posible
            cancelar ni revertir una transacción una vez ejecutada.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Phishing y fraudes:</strong> Existen
            sitios web fraudulentos que imitan plataformas legítimas. Verifique siempre la URL
            del sitio antes de introducir sus credenciales.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Interrupción del servicio:</strong> La
            plataforma puede experimentar interrupciones por mantenimiento, fallos técnicos o
            circunstancias fuera del control del operador.
          </li>
        </ul>

        {/* Section 8 */}
        <h2 style={sectionTitleStyle}>8. Perfil de riesgo adecuado</h2>
        <div style={{
          background: 'rgba(245, 158, 11, 0.08)',
          border: '1px solid rgba(245, 158, 11, 0.3)',
          borderRadius: '10px',
          padding: '1rem 1.25rem',
          marginBottom: '1rem',
          fontSize: '0.875rem',
          color: '#f59e0b',
          lineHeight: '1.65',
        }}>
          SpainCoin Exchange no es adecuado para personas con baja tolerancia al riesgo,
          personas que no puedan permitirse perder la totalidad de su inversión, menores de
          18 años, ni personas que no comprendan el funcionamiento de las tecnologías blockchain
          y los activos digitales.
        </div>
        <p style={textStyle}>
          Al usar SpainCoin Exchange, el Usuario declara haber leído y comprendido todos los
          riesgos aquí descritos y acepta operar bajo su propia responsabilidad.
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
          Para consultas: <span style={{ color: 'var(--accent)' }}>legal@spaincoin.com</span>
          <br />
          SpainCoin Exchange · Proyecto en desarrollo · Testnet v0.1 · Marzo 2026
        </div>
      </div>
    </div>
  )
}
