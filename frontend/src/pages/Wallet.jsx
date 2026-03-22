import { useState } from 'react'
import { getWallet, sendTx } from '../api/client.js'
import { formatSPC, formatNumber } from '../utils/format.js'

const styles = {
  page: {
    maxWidth: '900px',
    margin: '0 auto',
    padding: '2rem 1.5rem 4rem',
    width: '100%',
  },
  pageTitle: {
    fontSize: '1.5rem',
    fontWeight: '700',
    color: 'var(--text-primary)',
    marginBottom: '0.35rem',
    letterSpacing: '-0.02em',
  },
  pageSub: {
    fontSize: '0.875rem',
    color: 'var(--text-secondary)',
    marginBottom: '2.5rem',
  },
  card: {
    background: 'var(--bg-card)',
    border: '1px solid var(--border)',
    borderRadius: '12px',
    padding: '1.75rem',
    marginBottom: '1.5rem',
  },
  cardTitle: {
    fontSize: '1rem',
    fontWeight: '600',
    color: 'var(--text-primary)',
    marginBottom: '0.35rem',
  },
  cardSub: {
    fontSize: '0.8rem',
    color: 'var(--text-secondary)',
    marginBottom: '1.5rem',
  },
  label: {
    display: 'block',
    fontSize: '0.78rem',
    fontWeight: '500',
    color: 'var(--text-secondary)',
    textTransform: 'uppercase',
    letterSpacing: '0.05em',
    marginBottom: '0.4rem',
  },
  input: {
    width: '100%',
    background: 'var(--bg-secondary)',
    border: '1px solid var(--border)',
    borderRadius: '8px',
    padding: '0.65rem 0.875rem',
    color: 'var(--text-primary)',
    fontSize: '0.875rem',
    outline: 'none',
    transition: 'border-color 0.15s ease',
    fontFamily: 'inherit',
  },
  inputFocus: {
    borderColor: 'var(--accent)',
  },
  inputRow: {
    display: 'flex',
    gap: '0.75rem',
    alignItems: 'flex-end',
  },
  inputWrap: {
    flex: 1,
  },
  btn: {
    background: 'var(--accent)',
    color: '#fff',
    border: 'none',
    borderRadius: '8px',
    padding: '0.65rem 1.25rem',
    fontSize: '0.875rem',
    fontWeight: '600',
    cursor: 'pointer',
    transition: 'background 0.15s ease, opacity 0.15s ease',
    whiteSpace: 'nowrap',
    flexShrink: 0,
  },
  btnDisabled: {
    opacity: 0.5,
    cursor: 'not-allowed',
  },
  btnDanger: {
    background: 'var(--red)',
  },
  resultCard: {
    background: 'rgba(59, 130, 246, 0.08)',
    border: '1px solid rgba(59, 130, 246, 0.2)',
    borderRadius: '8px',
    padding: '1rem 1.25rem',
    marginTop: '1rem',
  },
  resultRow: {
    display: 'flex',
    justifyContent: 'space-between',
    padding: '0.35rem 0',
  },
  resultLabel: {
    fontSize: '0.8rem',
    color: 'var(--text-secondary)',
  },
  resultValue: {
    fontSize: '0.875rem',
    color: 'var(--text-primary)',
    fontWeight: '500',
  },
  error: {
    background: 'rgba(239, 68, 68, 0.08)',
    border: '1px solid rgba(239, 68, 68, 0.2)',
    borderRadius: '8px',
    padding: '0.75rem 1rem',
    color: '#ef4444',
    fontSize: '0.8rem',
    marginTop: '1rem',
  },
  success: {
    background: 'rgba(16, 185, 129, 0.08)',
    border: '1px solid rgba(16, 185, 129, 0.2)',
    borderRadius: '8px',
    padding: '0.75rem 1rem',
    color: '#10b981',
    fontSize: '0.8rem',
    marginTop: '1rem',
  },
  warning: {
    background: 'rgba(245, 158, 11, 0.08)',
    border: '1px solid rgba(245, 158, 11, 0.2)',
    borderRadius: '8px',
    padding: '0.875rem 1rem',
    marginBottom: '1.25rem',
    fontSize: '0.82rem',
    color: '#f59e0b',
    lineHeight: 1.5,
  },
  // Prominent red security warning box
  securityAlert: {
    background: 'rgba(239, 68, 68, 0.12)',
    border: '2px solid rgba(239, 68, 68, 0.5)',
    borderRadius: '10px',
    padding: '1rem 1.25rem',
    marginBottom: '1.5rem',
    fontSize: '0.875rem',
    color: '#ef4444',
    lineHeight: 1.6,
  },
  securityAlertTitle: {
    fontWeight: '700',
    fontSize: '0.95rem',
    marginBottom: '0.4rem',
    display: 'block',
  },
  balanceNote: {
    fontSize: '0.78rem',
    color: 'var(--text-secondary)',
    marginTop: '0.5rem',
    fontStyle: 'italic',
  },
  divider: {
    border: 'none',
    borderTop: '1px solid var(--border)',
    margin: '1.5rem 0',
  },
  formGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
    gap: '1rem',
    marginBottom: '1rem',
  },
  fieldWrap: {},
  advancedBadge: {
    display: 'inline-block',
    fontSize: '0.7rem',
    fontWeight: '600',
    color: '#f59e0b',
    background: 'rgba(245, 158, 11, 0.1)',
    border: '1px solid rgba(245, 158, 11, 0.2)',
    padding: '0.15rem 0.5rem',
    borderRadius: '4px',
    marginLeft: '0.5rem',
    verticalAlign: 'middle',
  },
  cliBox: {
    background: 'var(--bg-secondary)',
    border: '1px solid var(--border)',
    borderRadius: '8px',
    padding: '0.875rem 1rem',
    marginTop: '1.25rem',
  },
  cliLabel: {
    fontSize: '0.7rem',
    fontWeight: '600',
    color: 'var(--text-secondary)',
    textTransform: 'uppercase',
    letterSpacing: '0.06em',
    marginBottom: '0.5rem',
    display: 'block',
  },
  cliCode: {
    fontFamily: '"JetBrains Mono", "Fira Code", Consolas, monospace',
    fontSize: '0.8rem',
    color: 'var(--green)',
    wordBreak: 'break-all',
    lineHeight: 1.6,
    display: 'block',
  },
}

function InputField({ label, id, value, onChange, placeholder, mono }) {
  const [focused, setFocused] = useState(false)
  return (
    <div style={styles.fieldWrap}>
      <label htmlFor={id} style={styles.label}>{label}</label>
      <input
        id={id}
        style={{
          ...styles.input,
          ...(focused ? styles.inputFocus : {}),
          ...(mono ? { fontFamily: '"JetBrains Mono", Consolas, monospace', fontSize: '0.8rem' } : {}),
        }}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder || ''}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
        autoComplete="off"
        spellCheck={false}
      />
    </div>
  )
}

export default function Wallet() {
  // Check balance state
  const [address, setAddress] = useState('')
  const [walletData, setWalletData] = useState(null)
  const [balanceError, setBalanceError] = useState(null)
  const [balanceLoading, setBalanceLoading] = useState(false)

  // Send tx state
  const [txForm, setTxForm] = useState({
    from: '', to: '', amount: '', fee: '', nonce: '', sig_r: '', sig_s: ''
  })
  const [sendError, setSendError] = useState(null)
  const [sendSuccess, setSendSuccess] = useState(null)
  const [sendLoading, setSendLoading] = useState(false)

  async function handleCheckBalance(e) {
    e.preventDefault()
    if (!address.trim()) return
    setBalanceLoading(true)
    setBalanceError(null)
    setWalletData(null)
    try {
      const data = await getWallet(address.trim())
      setWalletData(data)
    } catch (err) {
      setBalanceError('Could not fetch wallet: ' + err.message)
    } finally {
      setBalanceLoading(false)
    }
  }

  function updateTxForm(field) {
    return (val) => setTxForm((prev) => ({ ...prev, [field]: val }))
  }

  async function handleSendTx(e) {
    e.preventDefault()

    // Confirmation dialog before broadcasting
    const confirmed = window.confirm(
      '¿Confirmas el envio de esta transaccion?\n\n' +
      'De: ' + (txForm.from || '(sin rellenar)') + '\n' +
      'Para: ' + (txForm.to || '(sin rellenar)') + '\n' +
      'Cantidad: ' + (txForm.amount || '0') + ' pesetas\n\n' +
      'Esta accion es IRREVERSIBLE una vez incluida en un bloque.'
    )
    if (!confirmed) return

    setSendLoading(true)
    setSendError(null)
    setSendSuccess(null)
    try {
      const payload = {
        from: txForm.from,
        to: txForm.to,
        amount: Number(txForm.amount),
        fee: Number(txForm.fee),
        nonce: Number(txForm.nonce),
        sig_r: txForm.sig_r,
        sig_s: txForm.sig_s,
      }
      const result = await sendTx(payload)
      setSendSuccess(`Transaction broadcast: ${result.tx_id || result.id || JSON.stringify(result)}`)
    } catch (err) {
      setSendError('Broadcast failed: ' + err.message)
    } finally {
      setSendLoading(false)
    }
  }

  // CLI sign-offline command (does NOT include private key — user fills that in locally)
  const cliSignExample = `spc tx sign \\
  --to ${txForm.to || '<RECIPIENT_ADDRESS>'} \\
  --amount ${txForm.amount || '<AMOUNT_PESETAS>'} \\
  --nonce ${txForm.nonce || '<NONCE>'} \\
  --fee ${txForm.fee || '<FEE>'} \\
  --key /path/to/keyfile.json`

  const cliBroadcastExample = `spc tx broadcast --signed tx_signed.json`

  return (
    <div className="page-enter" style={styles.page}>
      <h1 style={styles.pageTitle}>Wallet</h1>
      <p style={styles.pageSub}>Check balances and broadcast transactions on the SpainCoin network</p>

      {/* Check Balance */}
      <div style={styles.card}>
        <div style={styles.cardTitle}>Check Balance</div>
        <div style={styles.cardSub}>Enter a SpainCoin address to view its balance and nonce</div>

        <form onSubmit={handleCheckBalance}>
          <div style={styles.inputRow}>
            <div style={styles.inputWrap}>
              <label htmlFor="check-addr" style={styles.label}>Address</label>
              <input
                id="check-addr"
                style={styles.input}
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                placeholder="SPC..."
                autoComplete="off"
                spellCheck={false}
                onFocus={(e) => (e.target.style.borderColor = 'var(--accent)')}
                onBlur={(e) => (e.target.style.borderColor = 'var(--border)')}
              />
            </div>
            <button
              type="submit"
              style={{ ...styles.btn, ...(balanceLoading ? styles.btnDisabled : {}) }}
              disabled={balanceLoading || !address.trim()}
            >
              {balanceLoading ? 'Loading...' : 'Check Balance'}
            </button>
          </div>
        </form>

        <p style={styles.balanceNote}>
          Las consultas de balance son publicas y no requieren ninguna clave.
        </p>

        {balanceError && <div style={styles.error}>{balanceError}</div>}

        {walletData && (
          <div style={styles.resultCard}>
            <div style={styles.resultRow}>
              <span style={styles.resultLabel}>Address</span>
              <span className="mono" style={{ ...styles.resultValue, fontSize: '0.8rem' }}>
                {walletData.address}
              </span>
            </div>
            <div style={styles.resultRow}>
              <span style={styles.resultLabel}>Balance</span>
              <span style={{ ...styles.resultValue, color: 'var(--green)', fontSize: '1rem' }}>
                {walletData.balance_spc != null
                  ? `${formatNumber(Number(walletData.balance_spc).toFixed(4))} SPC`
                  : formatSPC(walletData.balance)}
              </span>
            </div>
            <div style={styles.resultRow}>
              <span style={styles.resultLabel}>Nonce</span>
              <span style={styles.resultValue}>{walletData.nonce ?? 0}</span>
            </div>
          </div>
        )}
      </div>

      {/* Send SPC */}
      <div style={styles.card}>
        <div style={styles.cardTitle}>
          Send SPC
          <span style={styles.advancedBadge}>Advanced Users</span>
        </div>
        <div style={styles.cardSub}>Broadcast a signed transaction directly to the network</div>

        {/* Prominent red security warning */}
        <div style={styles.securityAlert}>
          <span style={styles.securityAlertTitle}>
            NUNCA introduzcas tu clave privada en ningun sitio web.
          </span>
          Usa siempre el CLI para firmar transacciones offline. Tu clave privada nunca debe salir
          de tu dispositivo. Si alguien te pide tu clave privada, es una estafa.
          <br /><br />
          <strong>Flujo correcto:</strong> firma offline con el CLI &rarr; copia los valores R y S &rarr; pegalos aqui.
        </div>

        <div style={styles.warning}>
          <strong>Security notice:</strong> Never enter your private key on a website you do not fully trust and control.
          Use the CLI tool to sign transactions locally, then broadcast the signed payload here.
          Transactions require pre-computed ECDSA signatures (R, S values).
        </div>

        <form onSubmit={handleSendTx}>
          <div style={styles.formGrid}>
            <InputField label="From Address" id="tx-from" value={txForm.from} onChange={updateTxForm('from')} placeholder="SPC..." mono />
            <InputField label="To Address" id="tx-to" value={txForm.to} onChange={updateTxForm('to')} placeholder="SPC..." mono />
            <InputField label="Amount (pesetas)" id="tx-amount" value={txForm.amount} onChange={updateTxForm('amount')} placeholder="e.g. 1000000000000000000" />
            <InputField label="Fee (pesetas)" id="tx-fee" value={txForm.fee} onChange={updateTxForm('fee')} placeholder="e.g. 21000" />
            <InputField label="Nonce" id="tx-nonce" value={txForm.nonce} onChange={updateTxForm('nonce')} placeholder="e.g. 1" />
          </div>
          <div style={styles.formGrid}>
            <InputField label="Signature R (hex)" id="tx-sigr" value={txForm.sig_r} onChange={updateTxForm('sig_r')} placeholder="0x..." mono />
            <InputField label="Signature S (hex)" id="tx-sigs" value={txForm.sig_s} onChange={updateTxForm('sig_s')} placeholder="0x..." mono />
          </div>

          <button
            type="submit"
            style={{
              ...styles.btn,
              marginTop: '0.5rem',
              ...(sendLoading ? styles.btnDisabled : {}),
            }}
            disabled={sendLoading}
          >
            {sendLoading ? 'Broadcasting...' : 'Broadcast Transaction'}
          </button>
        </form>

        {sendError && <div style={styles.error}>{sendError}</div>}
        {sendSuccess && <div style={styles.success}>{sendSuccess}</div>}

        <div style={styles.cliBox}>
          <span style={styles.cliLabel}>Paso 1 — Firmar offline (recomendado)</span>
          <code style={styles.cliCode}>{cliSignExample}</code>
          <hr style={{ border: 'none', borderTop: '1px solid var(--border)', margin: '0.75rem 0' }} />
          <span style={styles.cliLabel}>Paso 2 — Broadcast del fichero firmado</span>
          <code style={styles.cliCode}>{cliBroadcastExample}</code>
        </div>
      </div>
    </div>
  )
}
