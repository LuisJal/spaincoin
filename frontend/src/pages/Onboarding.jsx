import { useState, useEffect } from 'react'

// ==========================================
// CRYPTO — same as WalletDownload (client-side)
// ==========================================
async function generateWallet() {
  const keyPair = await crypto.subtle.generateKey(
    { name: 'ECDSA', namedCurve: 'P-256' }, true, ['sign', 'verify']
  )
  const privJwk = await crypto.subtle.exportKey('jwk', keyPair.privateKey)
  const pubJwk = await crypto.subtle.exportKey('jwk', keyPair.publicKey)
  const privHex = b64toHex(privJwk.d)
  const xBytes = b64toBytes(pubJwk.x)
  const yBytes = b64toBytes(pubJwk.y)
  const pubBytes = new Uint8Array([...xBytes, ...yBytes])
  const hashBuf = await crypto.subtle.digest('SHA-256', pubBytes)
  const addrBytes = new Uint8Array(hashBuf).slice(12, 32)
  const address = 'SPC' + [...addrBytes].map(b => b.toString(16).padStart(2, '0')).join('')
  return { privateKey: privHex, address }
}
function b64toBytes(b64) {
  const p = b64.replace(/-/g, '+').replace(/_/g, '/') + '=='.slice(0, (4 - b64.length % 4) % 4)
  return new Uint8Array([...atob(p)].map(c => c.charCodeAt(0)))
}
function b64toHex(b64) { return [...b64toBytes(b64)].map(b => b.toString(16).padStart(2, '0')).join('') }

function saveWallet(address, privKey) {
  const wallets = JSON.parse(localStorage.getItem('spc_wallets') || '[]')
  if (!wallets.find(w => w.address === address)) {
    wallets.push({ address, key: privKey, created: Date.now() })
    localStorage.setItem('spc_wallets', JSON.stringify(wallets))
  }
  // Register with backend for wallet count
  fetch('/api/wallets/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ address }),
  }).catch(() => {})
}

function LogoImage({ size = 140 }) {
  return (
    <div className="torito-bounce" style={{ margin: '0 auto', width: size }}>
      <img src="/logo.jpeg" alt="SpainCoin" style={{ width: size, height: 'auto', borderRadius: '12px' }} />
    </div>
  )
}

// ==========================================
// STEP COMPONENTS
// ==========================================

function Step1Welcome({ onNext }) {
  const [visible, setVisible] = useState(false)
  const [walletCount, setWalletCount] = useState(0)
  useEffect(() => {
    setTimeout(() => setVisible(true), 100)
    fetch('/api/wallets/count').then(r => r.json()).then(d => setWalletCount(d.total)).catch(() => {})
  }, [])

  return (
    <div style={{ opacity: visible ? 1 : 0, transition: 'opacity 0.6s', textAlign: 'center', padding: '2rem 1.5rem' }}>
      <LogoImage size={160} />

      <h1 style={{ fontSize: '1.8rem', fontWeight: '800', color: 'var(--text-primary)', marginTop: '1.5rem', marginBottom: '0.5rem' }}>
        Bienvenido a <span style={{ color: '#ffc400' }}>SpainCoin</span>
      </h1>
      <p style={{ fontSize: '1rem', color: 'var(--text-secondary)', marginBottom: '1rem', lineHeight: 1.7 }}>
        La primera criptomoneda española.
        Vamos a explicarte todo paso a paso.
      </p>

      {walletCount > 0 && (
        <div style={{
          display: 'inline-block', padding: '0.4rem 1rem', borderRadius: '20px',
          background: 'rgba(16, 185, 129, 0.1)', border: '1px solid rgba(16, 185, 129, 0.25)',
          fontSize: '0.82rem', color: 'var(--green)', fontWeight: '600', marginBottom: '1.5rem',
        }}>
          🇪🇸 Ya somos {walletCount.toLocaleString('es-ES')} wallets
        </div>
      )}

      <div style={{
        background: 'var(--bg-card)', borderRadius: '16px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '2rem', textAlign: 'left',
      }}>
        <h3 style={{ fontSize: '1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
          ¿Qué es una blockchain?
        </h3>
        <p style={{ fontSize: '0.88rem', color: 'var(--text-secondary)', lineHeight: 1.8 }}>
          Imagina un <strong style={{ color: 'var(--text-primary)' }}>libro de cuentas</strong> que nadie puede borrar ni modificar.
          Cada página es un <strong style={{ color: '#ffc400' }}>bloque</strong>, y todas las páginas están encadenadas.
          Si alguien intenta cambiar una página, la cadena se rompe y todos lo ven.
        </p>
        <p style={{ fontSize: '0.88rem', color: 'var(--text-secondary)', lineHeight: 1.8, marginTop: '0.75rem' }}>
          SpainCoin es ese libro. Cada 5 segundos se escribe una nueva página.
          Y tú puedes tener <strong style={{ color: 'var(--green)' }}>monedas ($SPC)</strong> dentro de ese libro
          que son tuyas y <strong style={{ color: 'var(--text-primary)' }}>de nadie más</strong>.
        </p>
      </div>

      <div style={{
        background: 'rgba(16, 185, 129, 0.08)', border: '1px solid rgba(16, 185, 129, 0.25)',
        borderRadius: '12px', padding: '1rem', marginBottom: '2rem',
        fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.6,
      }}>
        Solo existirán <strong style={{ color: 'var(--green)' }}>21 millones de SPC</strong>.
        Cuantos menos queden, más valen. Los primeros en comprar son los que más ganan.
      </div>

      <button onClick={onNext} style={{
        width: '100%', padding: '1rem', borderRadius: '12px', border: 'none',
        background: 'linear-gradient(135deg, #ffc400, #e6a800)', color: '#000',
        fontSize: '1.1rem', fontWeight: '700', cursor: 'pointer',
      }}>
        Entendido, ¡vamos! →
      </button>

      <div style={{ marginTop: '1rem', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>
        Paso 1 de 5
      </div>
    </div>
  )
}

function Step2CreateWallet({ onNext, onWalletCreated }) {
  const [creating, setCreating] = useState(false)
  const [wallet, setWallet] = useState(null)
  const [copiedAddr, setCopiedAddr] = useState(false)
  const [copiedKey, setCopiedKey] = useState(false)
  const [confirmedAddr, setConfirmedAddr] = useState(false)
  const [confirmedKey, setConfirmedKey] = useState(false)

  async function handleCreate() {
    setCreating(true)
    try {
      const w = await generateWallet()
      saveWallet(w.address, w.privateKey)
      setWallet(w)
      onWalletCreated(w)
    } catch (e) {
      alert('Error: ' + e.message)
    }
    setCreating(false)
  }

  function copyAddr() {
    navigator.clipboard.writeText(wallet.address)
    setCopiedAddr(true)
  }

  function copyKey() {
    navigator.clipboard.writeText(wallet.privateKey)
    setCopiedKey(true)
  }

  const bothConfirmed = confirmedAddr && confirmedKey

  return (
    <div className="page-enter" style={{ textAlign: 'center', padding: '2rem 1.5rem' }}>
      <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>🔐</div>

      <h2 style={{ fontSize: '1.4rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
        {!wallet ? 'Crea tu wallet' : '¡Wallet creada!'}
      </h2>

      {!wallet ? (
        <>
          <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '1.5rem', lineHeight: 1.7 }}>
            Tu wallet es como tu <strong style={{ color: 'var(--text-primary)' }}>cuenta bancaria</strong> en SpainCoin.
            Se genera aquí, en tu dispositivo. Nadie más la ve.
          </p>

          <button onClick={handleCreate} disabled={creating} style={{
            width: '100%', padding: '1rem', borderRadius: '12px', border: 'none',
            background: 'linear-gradient(135deg, #ffc400, #e6a800)', color: '#000',
            fontSize: '1.1rem', fontWeight: '700', cursor: 'pointer',
            opacity: creating ? 0.6 : 1,
          }}>
            {creating ? 'Generando...' : 'Crear mi wallet 🔐'}
          </button>
        </>
      ) : (
        <>
          {/* ADDRESS */}
          <div style={{
            background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
            padding: '1.25rem', marginBottom: '1rem', textAlign: 'left',
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.5rem' }}>
              <span style={{ fontSize: '1.2rem' }}>📬</span>
              <strong style={{ fontSize: '0.9rem', color: 'var(--text-primary)' }}>Tu dirección (pública)</strong>
            </div>
            <p style={{ fontSize: '0.78rem', color: 'var(--text-secondary)', marginBottom: '0.6rem' }}>
              Es como tu número de cuenta. La compartes para que te envíen SPC.
            </p>
            <div style={{
              background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.65rem 0.75rem',
              fontFamily: 'monospace', fontSize: '0.72rem', color: 'var(--accent)',
              wordBreak: 'break-all', marginBottom: '0.5rem',
            }}>
              {wallet.address}
            </div>
            <button onClick={copyAddr} style={{
              width: '100%', padding: '0.5rem', borderRadius: '8px',
              border: '1px solid var(--border)',
              background: copiedAddr ? 'rgba(16,185,129,0.15)' : 'var(--bg-secondary)',
              color: copiedAddr ? 'var(--green)' : 'var(--text-secondary)',
              fontSize: '0.82rem', fontWeight: '600', cursor: 'pointer',
            }}>{copiedAddr ? '✓ Dirección copiada' : 'Copiar dirección'}</button>
          </div>

          {/* PRIVATE KEY */}
          <div style={{
            background: 'rgba(239, 68, 68, 0.06)', border: '1px solid rgba(239, 68, 68, 0.25)',
            borderRadius: '12px', padding: '1.25rem', marginBottom: '1rem', textAlign: 'left',
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.5rem' }}>
              <span style={{ fontSize: '1.2rem' }}>🔑</span>
              <strong style={{ fontSize: '0.9rem', color: 'var(--red)' }}>Tu clave privada (SECRETA)</strong>
            </div>
            <p style={{ fontSize: '0.78rem', color: 'var(--text-secondary)', marginBottom: '0.6rem' }}>
              Con esto firmas tus transacciones. Es como la contraseña de tu banco pero <strong style={{ color: 'var(--red)' }}>NO se puede recuperar</strong>.
              Si la pierdes, pierdes tus fondos para siempre.
            </p>
            <div style={{
              background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.65rem 0.75rem',
              fontFamily: 'monospace', fontSize: '0.65rem', color: 'var(--red)',
              wordBreak: 'break-all', marginBottom: '0.5rem',
            }}>
              {wallet.privateKey}
            </div>
            <button onClick={copyKey} style={{
              width: '100%', padding: '0.5rem', borderRadius: '8px',
              border: '1px solid rgba(239,68,68,0.3)',
              background: copiedKey ? 'rgba(16,185,129,0.15)' : 'rgba(239,68,68,0.08)',
              color: copiedKey ? 'var(--green)' : 'var(--red)',
              fontSize: '0.82rem', fontWeight: '600', cursor: 'pointer',
            }}>{copiedKey ? '✓ Clave copiada' : 'Copiar clave privada'}</button>
          </div>

          {/* WHERE TO SAVE */}
          <div style={{
            background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
            padding: '1rem', marginBottom: '1.25rem', textAlign: 'left',
          }}>
            <div style={{ fontSize: '0.85rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
              ¿Dónde guardar las dos?
            </div>
            <ul style={{ fontSize: '0.78rem', color: 'var(--text-secondary)', lineHeight: 1.8, paddingLeft: '1.25rem', margin: 0 }}>
              <li>Apúntalas <strong style={{ color: 'var(--green)' }}>en papel</strong> y guárdalas en un lugar seguro</li>
              <li>O usa un gestor de contraseñas (1Password, Bitwarden)</li>
              <li><strong style={{ color: 'var(--red)' }}>NUNCA</strong> las mandes por WhatsApp, email o Telegram</li>
            </ul>
          </div>

          {/* Confirmations */}
          <div style={{ textAlign: 'left', marginBottom: '1.25rem' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.75rem', cursor: 'pointer' }}>
              <input type="checkbox" checked={confirmedAddr} onChange={e => setConfirmedAddr(e.target.checked)}
                style={{ width: '20px', height: '20px', accentColor: 'var(--green)', flexShrink: 0 }} />
              <span style={{ fontSize: '0.85rem', color: 'var(--text-primary)' }}>
                He guardado mi <strong>dirección</strong> (SPCxxx...)
              </span>
            </label>
            <label style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', cursor: 'pointer' }}>
              <input type="checkbox" checked={confirmedKey} onChange={e => setConfirmedKey(e.target.checked)}
                style={{ width: '20px', height: '20px', accentColor: 'var(--green)', flexShrink: 0 }} />
              <span style={{ fontSize: '0.85rem', color: 'var(--text-primary)' }}>
                He guardado mi <strong>clave privada</strong> en un lugar seguro
              </span>
            </label>
          </div>

          <button onClick={onNext} disabled={!bothConfirmed} style={{
            width: '100%', padding: '1rem', borderRadius: '12px', border: 'none',
            background: bothConfirmed ? 'linear-gradient(135deg, #ffc400, #e6a800)' : 'var(--border)',
            color: bothConfirmed ? '#000' : 'var(--text-secondary)',
            fontSize: '1.1rem', fontWeight: '700',
            cursor: bothConfirmed ? 'pointer' : 'not-allowed',
          }}>
            Las he guardado, ¡vamos! →
          </button>
        </>
      )}

      <div style={{ marginTop: '1rem', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Paso 2 de 5</div>
    </div>
  )
}

function Step3SaveKey({ wallet, onNext }) {
  const [copied, setCopied] = useState(false)
  const [confirmed, setConfirmed] = useState(false)

  function handleCopy() {
    navigator.clipboard.writeText(wallet.privateKey)
    setCopied(true)
  }

  return (
    <div className="page-enter" style={{ textAlign: 'center', padding: '2rem 1.5rem' }}>
      <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>⚠️</div>

      <h2 style={{ fontSize: '1.4rem', fontWeight: '700', color: 'var(--red)', marginBottom: '0.5rem' }}>
        Guarda tu clave privada
      </h2>
      <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '1.5rem', lineHeight: 1.7 }}>
        Esta es la <strong style={{ color: 'var(--text-primary)' }}>única forma</strong> de acceder a tus fondos.
        Si la pierdes, <strong style={{ color: 'var(--red)' }}>pierdes todo para siempre</strong>.
        Nadie puede recuperarla, ni siquiera nosotros.
      </p>

      <div style={{
        background: 'rgba(239, 68, 68, 0.08)', border: '1px solid rgba(239, 68, 68, 0.3)',
        borderRadius: '12px', padding: '1.25rem', marginBottom: '1rem', textAlign: 'left',
      }}>
        <div style={{ fontSize: '0.7rem', color: 'var(--red)', fontWeight: '600', marginBottom: '0.5rem' }}>
          TU CLAVE PRIVADA (SECRETA)
        </div>
        <div style={{
          fontFamily: 'monospace', fontSize: '0.7rem', color: 'var(--text-primary)',
          wordBreak: 'break-all', lineHeight: 1.6, marginBottom: '0.75rem',
        }}>
          {wallet.privateKey}
        </div>
        <button onClick={handleCopy} style={{
          width: '100%', padding: '0.6rem', borderRadius: '8px',
          border: '1px solid var(--border)',
          background: copied ? 'rgba(16,185,129,0.15)' : 'var(--bg-secondary)',
          color: copied ? 'var(--green)' : 'var(--text-secondary)',
          fontSize: '0.85rem', fontWeight: '600', cursor: 'pointer',
        }}>
          {copied ? '✓ Copiada' : 'Copiar clave'}
        </button>
      </div>

      <div style={{
        background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
        padding: '1rem', marginBottom: '1.5rem', textAlign: 'left',
      }}>
        <div style={{ fontSize: '0.85rem', color: 'var(--text-primary)', fontWeight: '600', marginBottom: '0.5rem' }}>
          ¿Dónde guardarla?
        </div>
        <ul style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', lineHeight: 1.8, paddingLeft: '1.25rem' }}>
          <li>Apúntala <strong style={{ color: 'var(--green)' }}>en papel</strong> y guárdala en un lugar seguro</li>
          <li>O usa un gestor de contraseñas (1Password, Bitwarden)</li>
          <li><strong style={{ color: 'var(--red)' }}>NUNCA</strong> la mandes por WhatsApp, email o Telegram</li>
        </ul>
      </div>

      <label style={{
        display: 'flex', alignItems: 'center', gap: '0.75rem',
        marginBottom: '1.25rem', cursor: 'pointer', textAlign: 'left',
      }}>
        <input type="checkbox" checked={confirmed} onChange={e => setConfirmed(e.target.checked)}
          style={{ width: '20px', height: '20px', accentColor: 'var(--green)', flexShrink: 0 }} />
        <span style={{ fontSize: '0.85rem', color: 'var(--text-primary)' }}>
          He guardado mi clave privada en un lugar seguro
        </span>
      </label>

      <button onClick={onNext} disabled={!confirmed} style={{
        width: '100%', padding: '1rem', borderRadius: '12px', border: 'none',
        background: confirmed ? 'linear-gradient(135deg, #ffc400, #e6a800)' : 'var(--border)',
        color: confirmed ? '#000' : 'var(--text-secondary)',
        fontSize: '1.1rem', fontWeight: '700', cursor: confirmed ? 'pointer' : 'not-allowed',
      }}>
        Ya la guardé, siguiente →
      </button>

      <div style={{ marginTop: '1rem', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Paso 3 de 5</div>
    </div>
  )
}

function Step4FirstBuy({ wallet, onNext }) {
  const [price, setPrice] = useState(0.05)

  useEffect(() => {
    fetch('/api/market/price').then(r => r.json()).then(d => setPrice(d.price_eur)).catch(() => {})
  }, [])

  return (
    <div className="page-enter" style={{ textAlign: 'center', padding: '2rem 1.5rem' }}>
      <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>💰</div>

      <h2 style={{ fontSize: '1.4rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
        Tu primera compra
      </h2>
      <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '1.5rem', lineHeight: 1.7 }}>
        Ahora mismo 1 SPC vale <strong style={{ color: '#ffc400' }}>{price.toFixed(4)}€</strong>.
        Con 10€ conseguirías <strong style={{ color: 'var(--green)' }}>{Math.round(10/price)} SPC</strong>.
      </p>

      {/* Price projection */}
      <div style={{
        background: 'var(--bg-card)', borderRadius: '16px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '1.5rem', textAlign: 'left',
      }}>
        <div style={{ fontSize: '0.9rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '1rem' }}>
          Si compras hoy {Math.round(10/price)} SPC por 10€...
        </div>
        <div style={{ display: 'grid', gap: '0.5rem' }}>
          {[
            { when: 'Compras hoy', price: price, color: 'var(--text-secondary)' },
            { when: 'A 0.50€', price: 0.50, color: 'var(--accent)' },
            { when: 'A 1€', price: 1.00, color: 'var(--green)' },
            { when: 'A 5€', price: 5.00, color: '#ffc400' },
          ].map((p, i) => {
            const spc = Math.round(10/price)
            const value = (spc * p.price).toFixed(2)
            const gain = ((p.price / price - 1) * 100).toFixed(0)
            return (
              <div key={i} style={{
                display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                padding: '0.5rem 0.75rem', borderRadius: '8px',
                background: i === 0 ? 'var(--bg-secondary)' : 'transparent',
              }}>
                <span style={{ fontSize: '0.82rem', color: 'var(--text-secondary)' }}>{p.when}</span>
                <span style={{ fontSize: '0.9rem', fontWeight: '700', color: p.color }}>
                  {value}€ {i > 0 && <span style={{ fontSize: '0.7rem', fontWeight: '400' }}>(+{gain}%)</span>}
                </span>
              </div>
            )
          })}
        </div>
      </div>

      <div style={{
        background: 'rgba(255, 196, 0, 0.08)', border: '1px solid rgba(255, 196, 0, 0.25)',
        borderRadius: '12px', padding: '1rem', marginBottom: '1.5rem',
        fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.6,
      }}>
        🐂 Los <strong style={{ color: '#ffc400' }}>primeros 33.000 SPC</strong> son para la comunidad fundadora.
        Cuanto antes compres, más barato. El precio sube automáticamente con cada venta.
      </div>

      <a href={`https://t.me/spaincoin_bot?start=buy_${wallet?.address || ''}`} target="_blank" rel="noopener noreferrer" style={{
        display: 'block', width: '100%', padding: '1rem', borderRadius: '12px', border: 'none',
        background: '#0088cc', color: '#fff', textAlign: 'center',
        fontSize: '1.1rem', fontWeight: '700', textDecoration: 'none', marginBottom: '0.75rem',
      }}>
        Comprar en Telegram →
      </a>

      <button onClick={onNext} style={{
        width: '100%', padding: '0.75rem', borderRadius: '10px',
        border: '1px solid var(--border)', background: 'transparent',
        color: 'var(--text-secondary)', fontSize: '0.9rem', cursor: 'pointer',
      }}>
        Ya he comprado / Lo haré luego →
      </button>

      <div style={{ marginTop: '1rem', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Paso 4 de 5</div>
    </div>
  )
}

function Step5Share({ wallet, onFinish }) {
  const shareText = '🐂 Acabo de entrar en SpainCoin — la primera crypto española. Los primeros en comprar son los que más ganan. Yo ya tengo mis SPC. 🇪🇸 spaincoin.es'

  function shareWhatsApp() {
    window.open(`https://wa.me/?text=${encodeURIComponent(shareText)}`, '_blank')
  }
  function shareTelegram() {
    window.open(`https://t.me/share/url?url=${encodeURIComponent('https://spaincoin.es')}&text=${encodeURIComponent(shareText)}`, '_blank')
  }
  function shareTwitter() {
    window.open(`https://twitter.com/intent/tweet?text=${encodeURIComponent(shareText)}`, '_blank')
  }
  function copyLink() {
    navigator.clipboard.writeText('https://spaincoin.es')
  }

  return (
    <div className="page-enter" style={{ textAlign: 'center', padding: '2rem 1.5rem' }}>
      <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>🎉</div>

      <h2 style={{ fontSize: '1.4rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
        ¡Ya eres parte de SpainCoin!
      </h2>
      <p style={{ fontSize: '0.9rem', color: 'var(--text-secondary)', marginBottom: '1.5rem', lineHeight: 1.7 }}>
        Tienes tu wallet lista. Ahora viene lo más importante...
      </p>

      <div style={{
        background: 'var(--bg-card)', borderRadius: '16px', border: '1px solid var(--border)',
        padding: '1.5rem', marginBottom: '1.5rem', textAlign: 'left',
      }}>
        <h3 style={{ fontSize: '1rem', fontWeight: '600', color: '#ffc400', marginBottom: '0.75rem' }}>
          ¿Cómo sube el precio?
        </h3>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.8 }}>
          <strong style={{ color: 'var(--text-primary)' }}>Cada persona que compra SPC hace que el precio suba.</strong>
          Es matemática pura: hay una cantidad limitada (21 millones) y cada vez quedan menos disponibles.
        </p>
        <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.8, marginTop: '0.75rem' }}>
          Si compartes SpainCoin con un amigo y él compra → tu SPC vale más.
          Si lo compartes con 10 amigos → vale mucho más.
          <strong style={{ color: 'var(--green)' }}> Los de la primera comunidad son los que más ganan.</strong>
        </p>
      </div>

      <div style={{
        background: 'rgba(255, 196, 0, 0.08)', border: '1px solid rgba(255, 196, 0, 0.25)',
        borderRadius: '16px', padding: '1.5rem', marginBottom: '1.5rem', textAlign: 'left',
      }}>
        <h3 style={{ fontSize: '1rem', fontWeight: '600', color: '#ffc400', marginBottom: '0.5rem' }}>
          El plan SpainCoin
        </h3>
        <ul style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', lineHeight: 1.8, paddingLeft: '1.25rem' }}>
          <li>Primeros 33.000 SPC → <strong style={{ color: 'var(--green)' }}>comunidad fundadora</strong> (los que más ganan)</li>
          <li>Después → salimos a exchanges internacionales</li>
          <li>Gente de fuera de España compra → el precio sube más</li>
          <li>Los españoles que compraron primero → <strong style={{ color: '#ffc400' }}>los más beneficiados</strong></li>
        </ul>
      </div>

      <div style={{ fontSize: '0.95rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '1rem' }}>
        Comparte con tus amigos:
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem', marginBottom: '1rem' }}>
        <button onClick={shareWhatsApp} style={{
          padding: '0.75rem', borderRadius: '10px', border: 'none',
          background: '#25D366', color: '#fff', fontSize: '0.9rem', fontWeight: '600', cursor: 'pointer',
        }}>WhatsApp</button>
        <button onClick={shareTelegram} style={{
          padding: '0.75rem', borderRadius: '10px', border: 'none',
          background: '#0088cc', color: '#fff', fontSize: '0.9rem', fontWeight: '600', cursor: 'pointer',
        }}>Telegram</button>
        <button onClick={shareTwitter} style={{
          padding: '0.75rem', borderRadius: '10px', border: 'none',
          background: '#1DA1F2', color: '#fff', fontSize: '0.9rem', fontWeight: '600', cursor: 'pointer',
        }}>Twitter / X</button>
        <button onClick={copyLink} style={{
          padding: '0.75rem', borderRadius: '10px', border: '1px solid var(--border)',
          background: 'var(--bg-secondary)', color: 'var(--text-secondary)', fontSize: '0.9rem', cursor: 'pointer',
        }}>Copiar link</button>
      </div>

      <div style={{
        background: 'rgba(239, 68, 68, 0.06)', border: '1px solid rgba(239, 68, 68, 0.2)',
        borderRadius: '10px', padding: '0.85rem', marginBottom: '1.5rem',
        fontSize: '0.75rem', color: 'var(--text-secondary)', lineHeight: 1.5,
      }}>
        ⚠️ Invertir en criptomonedas tiene riesgos. Invierte solo lo que puedas permitirte perder.
        SPC no tiene valor garantizado. Lee los riesgos en spaincoin.es
      </div>

      <button onClick={onFinish} style={{
        width: '100%', padding: '1rem', borderRadius: '12px', border: 'none',
        background: 'linear-gradient(135deg, #ffc400, #e6a800)', color: '#000',
        fontSize: '1.1rem', fontWeight: '700', cursor: 'pointer',
      }}>
        Entrar a SpainCoin 🇪🇸
      </button>

      <div style={{ marginTop: '1rem', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Paso 5 de 5</div>
    </div>
  )
}

// ==========================================
// MAIN ONBOARDING COMPONENT
// ==========================================
export default function Onboarding({ onComplete }) {
  const [step, setStep] = useState(1)
  const [wallet, setWallet] = useState(null)

  return (
    <div style={{ maxWidth: '500px', margin: '0 auto', minHeight: '100vh', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
      {step === 1 && <Step1Welcome onNext={() => setStep(2)} />}
      {step === 2 && <Step2CreateWallet onNext={() => setStep(4)} onWalletCreated={setWallet} />}
      {step === 4 && <Step4FirstBuy wallet={wallet} onNext={() => setStep(5)} />}
      {step === 5 && <Step5Share wallet={wallet} onFinish={onComplete} />}
    </div>
  )
}
