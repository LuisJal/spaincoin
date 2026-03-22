export default function Terms({ onNavigate }) {
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

  const warningBoxStyle = {
    background: 'rgba(245, 158, 11, 0.08)',
    border: '1px solid rgba(245, 158, 11, 0.3)',
    borderRadius: '10px',
    padding: '1rem 1.25rem',
    marginBottom: '2rem',
    fontSize: '0.875rem',
    color: '#f59e0b',
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
        <h1 style={titleStyle}>Términos y Condiciones de Uso</h1>
        <p style={metaStyle}>
          SpainCoin Exchange · Última actualización: marzo 2026
        </p>

        {/* Testnet warning */}
        <div style={warningBoxStyle}>
          <strong>AVISO IMPORTANTE:</strong> SpainCoin Exchange se encuentra actualmente en fase
          TESTNET. Se trata de un entorno de pruebas con activos digitales experimentales
          sin valor económico real. No invierta dinero real en esta fase.
        </div>

        {/* Section 1 */}
        <h2 style={sectionTitleStyle}>1. Objeto y naturaleza del servicio</h2>
        <p style={textStyle}>
          Los presentes Términos y Condiciones (en adelante, "Términos") regulan el acceso y uso
          de SpainCoin Exchange (en adelante, "el Servicio" o "la Plataforma"), un exchange
          descentralizado de activos digitales centrado en la criptomoneda $SPC (SpainCoin),
          operado en el marco de un proyecto de blockchain Layer 1 propio.
        </p>
        <p style={textStyle}>
          Al acceder o utilizar el Servicio, el usuario (en adelante, "el Usuario") declara haber
          leído, comprendido y aceptado en su totalidad los presentes Términos. Si no está de
          acuerdo con alguna de las condiciones aquí expuestas, deberá abstenerse de usar la
          Plataforma.
        </p>
        <p style={textStyle}>
          SpainCoin Exchange no es una entidad financiera regulada, ni un banco, ni una empresa de
          servicios de inversión (ESI) conforme a la normativa MiFID II o legislación española
          equivalente. $SPC es un activo digital experimental y no constituye dinero de curso
          legal, ni un instrumento financiero regulado en el sentido de la Directiva 2014/65/UE.
        </p>

        {/* Section 2 */}
        <h2 style={sectionTitleStyle}>2. Estado actual: Fase Testnet</h2>
        <p style={textStyle}>
          La Plataforma se encuentra en fase de pruebas (TESTNET). Esto implica que:
        </p>
        <ul style={{ ...textStyle, paddingLeft: '1.5rem', listStyleType: 'disc' }}>
          <li style={{ marginBottom: '0.4rem' }}>
            Los activos $SPC en circulación durante esta fase son ficticios y no tienen valor económico.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            La red puede reiniciarse total o parcialmente en cualquier momento, lo que podría implicar
            la pérdida de todos los datos, balances y transacciones registradas.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            Las funcionalidades están sujetas a cambios frecuentes sin previo aviso.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            No se garantiza la disponibilidad, estabilidad ni continuidad del servicio durante
            esta fase.
          </li>
        </ul>
        <p style={textStyle}>
          La transición a mainnet requerirá la aceptación de nuevos Términos y el cumplimiento de
          requisitos adicionales, incluyendo verificación KYC/AML.
        </p>

        {/* Section 3 */}
        <h2 style={sectionTitleStyle}>3. Riesgos de los activos digitales</h2>
        <p style={textStyle}>
          El Usuario reconoce y acepta expresamente los siguientes riesgos inherentes al uso de
          activos digitales y plataformas blockchain:
        </p>
        <ul style={{ ...textStyle, paddingLeft: '1.5rem', listStyleType: 'disc' }}>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Volatilidad extrema:</strong> El
            valor de los activos digitales puede fluctuar de manera drástica en períodos muy
            cortos de tiempo, pudiendo llegar a cero.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Pérdida total del capital:</strong> El
            Usuario puede perder la totalidad de los fondos invertidos.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Riesgo tecnológico:</strong> Los
            sistemas blockchain y los contratos inteligentes pueden contener vulnerabilidades,
            errores o fallos que deriven en pérdidas.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Riesgo regulatorio:</strong> La
            normativa aplicable a los activos digitales está en evolución constante. Cambios
            legislativos podrían afectar la operatividad de la Plataforma.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Ausencia de garantías:</strong> SpainCoin
            Exchange no garantiza beneficios, rentabilidades ni el mantenimiento del valor futuro
            de $SPC o cualquier otro activo.
          </li>
        </ul>

        {/* Section 4 */}
        <h2 style={sectionTitleStyle}>4. Responsabilidades del usuario</h2>
        <p style={textStyle}>
          El Usuario se compromete a:
        </p>
        <ul style={{ ...textStyle, paddingLeft: '1.5rem', listStyleType: 'disc' }}>
          <li style={{ marginBottom: '0.4rem' }}>
            Ser mayor de 18 años y tener plena capacidad legal para contratar conforme a la
            legislación de su país de residencia.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            Usar la Plataforma de manera lícita, no fraudulenta y conforme a la normativa vigente.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            <strong style={{ color: 'var(--text-primary)' }}>Custodia de claves privadas:</strong> El
            Usuario es el único responsable de la custodia y seguridad de su clave privada. SpainCoin
            Exchange nunca almacena claves privadas en texto plano y no podrá recuperarlas en caso
            de pérdida. La pérdida de la clave privada implica la pérdida irreversible del acceso
            a los fondos asociados.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            No realizar actividades de lavado de dinero, financiación del terrorismo u otras
            actividades ilícitas.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            Mantener la confidencialidad de sus credenciales de acceso y notificar inmediatamente
            cualquier acceso no autorizado.
          </li>
        </ul>

        {/* Section 5 */}
        <h2 style={sectionTitleStyle}>5. Responsabilidad del operador</h2>
        <p style={textStyle}>
          En la medida en que lo permita la legislación aplicable, SpainCoin Exchange no asume
          responsabilidad alguna por:
        </p>
        <ul style={{ ...textStyle, paddingLeft: '1.5rem', listStyleType: 'disc' }}>
          <li style={{ marginBottom: '0.4rem' }}>
            Pérdidas directas o indirectas derivadas del uso o la imposibilidad de uso de la
            Plataforma.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            Interrupciones, errores, vulnerabilidades o fallos técnicos del sistema blockchain
            subyacente.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            Pérdida de acceso a fondos por extravío, robo o compromiso de claves privadas.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            Decisiones de inversión o pérdidas económicas tomadas por el Usuario basándose en
            información disponible en la Plataforma.
          </li>
          <li style={{ marginBottom: '0.4rem' }}>
            Cambios normativos que afecten a la operatividad de la Plataforma o al valor de $SPC.
          </li>
        </ul>
        <p style={textStyle}>
          La responsabilidad total del operador, en cualquier caso, queda limitada al importe
          máximo de las comisiones efectivamente abonadas por el Usuario en los últimos doce (12)
          meses anteriores al evento que origina la reclamación.
        </p>

        {/* Section 6 */}
        <h2 style={sectionTitleStyle}>6. KYC/AML (próximamente en mainnet)</h2>
        <p style={textStyle}>
          En la fase actual (Testnet), no se exige verificación de identidad (KYC —
          Know Your Customer) ni procedimientos de previsión de blanqueo de capitales (AML —
          Anti-Money Laundering).
        </p>
        <p style={textStyle}>
          Con el lanzamiento de la mainnet, y de conformidad con la normativa europea aplicable —
          incluyendo la Directiva (UE) 2018/843 (5ª Directiva AML), el Reglamento (UE) 2023/1113
          sobre información en las transferencias de fondos y criptoactivos, y el Reglamento MiCA
          (UE) 2023/1114 — se requerirá a todos los usuarios la superación de un proceso de
          verificación de identidad previo a la operación con activos digitales de valor real.
        </p>

        {/* Section 7 */}
        <h2 style={sectionTitleStyle}>7. Propiedad intelectual</h2>
        <p style={textStyle}>
          Todos los derechos de propiedad intelectual e industrial sobre la Plataforma, su código
          fuente, diseño, marca "SpainCoin" y demás elementos son titularidad de SpainCoin
          Exchange o de sus licenciantes. Queda prohibida su reproducción, distribución,
          modificación o uso comercial sin autorización expresa y escrita.
        </p>
        <p style={textStyle}>
          El protocolo blockchain subyacente puede estar sujeto a una licencia de código abierto.
          En tal caso, los términos de dicha licencia prevalecen sobre los presentes Términos en
          lo relativo al uso del código fuente del protocolo.
        </p>

        {/* Section 8 */}
        <h2 style={sectionTitleStyle}>8. Modificaciones</h2>
        <p style={textStyle}>
          SpainCoin Exchange se reserva el derecho de modificar los presentes Términos en
          cualquier momento. Los cambios serán notificados mediante aviso en la Plataforma con
          al menos siete (7) días de antelación, salvo en casos de urgencia por razones legales
          o de seguridad. El uso continuado de la Plataforma tras la entrada en vigor de las
          modificaciones implicará la aceptación de los nuevos Términos.
        </p>

        {/* Section 9 */}
        <h2 style={sectionTitleStyle}>9. Ley aplicable y jurisdicción</h2>
        <p style={textStyle}>
          Los presentes Términos se rigen e interpretan conforme al ordenamiento jurídico español.
          Para la resolución de cualquier controversia derivada de la interpretación o ejecución
          de estos Términos, las partes se someten expresamente a la jurisdicción de los Juzgados
          y Tribunales de España, con renuncia a cualquier otro fuero que pudiera corresponderles.
        </p>
        <p style={textStyle}>
          Lo anterior se entiende sin perjuicio de los derechos que la normativa de consumidores
          y usuarios reconoce a los usuarios que tengan la condición de consumidores conforme al
          Real Decreto Legislativo 1/2007, de 16 de noviembre.
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
          Para cualquier consulta legal: <span style={{ color: 'var(--accent)' }}>legal@spaincoin.com</span>
          <br />
          SpainCoin Exchange · Proyecto en desarrollo · Testnet v0.1 · Marzo 2026
        </div>
      </div>
    </div>
  )
}
