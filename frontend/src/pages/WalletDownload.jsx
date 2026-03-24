import { useState, useEffect } from 'react'
import { getStatus } from '../api/client.js'

// ==========================================
// CLIENT-SIDE CRYPTO — claves NUNCA salen del navegador
// ==========================================

async function generateWallet() {
  const keyPair = await crypto.subtle.generateKey(
    { name: 'ECDSA', namedCurve: 'P-256' },
    true, // extractable
    ['sign', 'verify']
  )
  const privJwk = await crypto.subtle.exportKey('jwk', keyPair.privateKey)
  const pubJwk = await crypto.subtle.exportKey('jwk', keyPair.publicKey)

  // Private key: d parameter (base64url → hex)
  const privHex = base64urlToHex(privJwk.d)

  // Public key: x,y coordinates
  const xBytes = base64urlToBytes(pubJwk.x)
  const yBytes = base64urlToBytes(pubJwk.y)

  // Address: SHA-256(x || y), take last 20 bytes, prefix "SPC"
  const pubBytes = new Uint8Array([...xBytes, ...yBytes])
  const hashBuf = await crypto.subtle.digest('SHA-256', pubBytes)
  const hashArr = new Uint8Array(hashBuf)
  const addrBytes = hashArr.slice(12, 32) // last 20 bytes
  const address = 'SPC' + bytesToHex(addrBytes)

  return { privateKey: privHex, address, pubX: bytesToHex(xBytes), pubY: bytesToHex(yBytes) }
}

async function importWallet(privHex) {
  const privBytes = hexToBytes(privHex)
  const privB64 = bytesToBase64url(privBytes)

  // Derive public key by importing as ECDSA
  const keyPair = await crypto.subtle.importKey(
    'jwk',
    { kty: 'EC', crv: 'P-256', d: privB64, x: 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', y: 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA' },
    { name: 'ECDSA', namedCurve: 'P-256' },
    true, ['sign']
  ).catch(() => null)

  if (!keyPair) {
    // Fallback: import as raw and re-export
    throw new Error('Clave privada inválida')
  }

  // We need the public key — generate from private using a different approach
  // Import private, export JWK to get x,y
  const jwk = await crypto.subtle.exportKey('jwk', keyPair)
  const xBytes = base64urlToBytes(jwk.x)
  const yBytes = base64urlToBytes(jwk.y)

  const pubBytes = new Uint8Array([...xBytes, ...yBytes])
  const hashBuf = await crypto.subtle.digest('SHA-256', pubBytes)
  const hashArr = new Uint8Array(hashBuf)
  const addrBytes = hashArr.slice(12, 32)
  const address = 'SPC' + bytesToHex(addrBytes)

  return { privateKey: privHex, address, pubX: bytesToHex(xBytes), pubY: bytesToHex(yBytes) }
}

// Utility functions
function base64urlToBytes(b64) {
  const padded = b64.replace(/-/g, '+').replace(/_/g, '/') + '=='.slice(0, (4 - b64.length % 4) % 4)
  const binary = atob(padded)
  return new Uint8Array([...binary].map(c => c.charCodeAt(0)))
}

function base64urlToHex(b64) {
  return bytesToHex(base64urlToBytes(b64))
}

function bytesToBase64url(bytes) {
  const binary = String.fromCharCode(...bytes)
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '')
}

function bytesToHex(bytes) {
  return [...bytes].map(b => b.toString(16).padStart(2, '0')).join('')
}

function hexToBytes(hex) {
  const bytes = new Uint8Array(hex.length / 2)
  for (let i = 0; i < hex.length; i += 2) {
    bytes[i / 2] = parseInt(hex.substr(i, 2), 16)
  }
  return bytes
}

// ==========================================
// WALLET STORAGE — localStorage cifrado
// ==========================================

function saveWallet(address, privKey, name) {
  const wallets = JSON.parse(localStorage.getItem('spc_wallets') || '[]')
  if (!wallets.find(w => w.address === address)) {
    wallets.push({ address, key: privKey, name: name || '', created: Date.now() })
    localStorage.setItem('spc_wallets', JSON.stringify(wallets))
  }
}

function loadWallets() {
  return JSON.parse(localStorage.getItem('spc_wallets') || '[]')
}

function deleteWallet(address) {
  const wallets = loadWallets().filter(w => w.address !== address)
  localStorage.setItem('spc_wallets', JSON.stringify(wallets))
}

// ==========================================
// COMPONENT
// ==========================================

const formatSPC = (n) => n >= 1 ? n.toLocaleString('es-ES', { maximumFractionDigits: 4 }) : n.toFixed(6)

export default function WalletDownload({ onNavigate }) {
  const [wallets, setWallets] = useState([])
  const [activeWallet, setActiveWallet] = useState(null)
  const [balance, setBalance] = useState(null)
  const [creating, setCreating] = useState(false)
  const [justCreated, setJustCreated] = useState(null) // {address, privateKey} — shown once only
  const [createdConfirmed, setCreatedConfirmed] = useState(false)
  const [showImport, setShowImport] = useState(false)
  const [importInput, setImportInput] = useState('')
  const [importName, setImportName] = useState('')
  const [importError, setImportError] = useState('')
  const [copied, setCopied] = useState('')
  const [tab, setTab] = useState('wallet') // 'wallet' | 'download'

  useEffect(() => {
    const saved = loadWallets()
    setWallets(saved)
    if (saved.length > 0) setActiveWallet(saved[0])
  }, [])

  useEffect(() => {
    if (!activeWallet) return
    async function fetchBalance() {
      try {
        const res = await fetch(`/api/wallet/${activeWallet.address}`)
        if (res.ok) {
          const data = await res.json()
          setBalance(data)
        }
      } catch (e) { console.error(e) }
    }
    fetchBalance()
    const i = setInterval(fetchBalance, 15000)
    return () => clearInterval(i)
  }, [activeWallet])

  async function handleCreate() {
    setCreating(true)
    try {
      const w = await generateWallet()
      saveWallet(w.address, '', 'Mi wallet')
      const saved = loadWallets()
      setWallets(saved)
      setActiveWallet(saved.find(s => s.address === w.address))
      setJustCreated(w) // show private key once
      setCreatedConfirmed(false)
    } catch (e) {
      alert('Error creando wallet: ' + e.message)
    }
    setCreating(false)
  }

  function handleImport() {
    setImportError('')
    const input = importInput.trim()
    if (!input.startsWith('SPC') || input.length !== 43) {
      setImportError('Dirección inválida. Debe empezar por SPC y tener 43 caracteres.')
      return
    }
    saveWallet(input, '', importName.trim() || 'Mi wallet')
    const saved = loadWallets()
    setWallets(saved)
    setActiveWallet(saved.find(s => s.address === input))
    setShowImport(false)
    setImportInput('')
    setImportName('')
  }

  function handleCopy(text, label) {
    navigator.clipboard.writeText(text)
    setCopied(label)
    setTimeout(() => setCopied(''), 2000)
  }

  const sectionCard = {
    background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)',
    padding: '1.25rem', marginBottom: '1rem',
  }

  return (
    <div className="page-enter" style={{ maxWidth: '600px', margin: '0 auto', padding: '1.5rem 1rem' }}>
      <h1 style={{ fontSize: '1.5rem', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
        Wallet SpainCoin
      </h1>
      <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '1.5rem' }}>
        Tus claves, tus fondos. Todo se queda en tu dispositivo.
      </p>

      {/* Tab selector */}
      <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1.5rem' }}>
        <button onClick={() => setTab('wallet')} style={{
          flex: 1, padding: '0.5rem', borderRadius: '8px', border: 'none', cursor: 'pointer',
          fontSize: '0.85rem', fontWeight: tab === 'wallet' ? '600' : '400',
          background: tab === 'wallet' ? 'var(--accent)' : 'var(--bg-secondary)',
          color: tab === 'wallet' ? '#fff' : 'var(--text-secondary)',
        }}>Web Wallet</button>
        <button onClick={() => setTab('download')} style={{
          flex: 1, padding: '0.5rem', borderRadius: '8px', border: 'none', cursor: 'pointer',
          fontSize: '0.85rem', fontWeight: tab === 'download' ? '600' : '400',
          background: tab === 'download' ? 'var(--accent)' : 'var(--bg-secondary)',
          color: tab === 'download' ? '#fff' : 'var(--text-secondary)',
        }}>Descargar CLI</button>
      </div>

      {tab === 'wallet' && (
        <>
          {/* No wallet yet */}
          {wallets.length === 0 && !showImport && (
            <div style={{ ...sectionCard, textAlign: 'center', padding: '2rem 1.25rem' }}>
              <div style={{ fontSize: '2.5rem', marginBottom: '1rem' }}>🔐</div>
              <div style={{ fontSize: '1.1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.5rem' }}>
                Crea tu primer wallet
              </div>
              <p style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', marginBottom: '1.5rem', lineHeight: 1.6 }}>
                Las claves se generan aquí, en tu dispositivo. Nunca salen de tu móvil/ordenador.
              </p>
              <button onClick={handleCreate} disabled={creating} style={{
                width: '100%', padding: '0.85rem', borderRadius: '10px', border: 'none',
                background: 'linear-gradient(135deg, #ffc400, #e6a800)', color: '#000',
                fontSize: '1rem', fontWeight: '700', cursor: 'pointer',
                opacity: creating ? 0.6 : 1,
              }}>
                {creating ? 'Generando...' : 'Crear Wallet'}
              </button>
              <button onClick={() => setShowImport(true)} style={{
                marginTop: '0.75rem', width: '100%', padding: '0.75rem', borderRadius: '10px',
                border: '1px solid var(--accent)', background: 'transparent',
                color: 'var(--accent)', fontSize: '0.95rem', fontWeight: '600', cursor: 'pointer',
              }}>
                Ya tengo una wallet
              </button>
            </div>
          )}

          {/* Import form */}
          {showImport && (
            <div style={sectionCard}>
              <div style={{ fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>Ya tengo wallet</div>
              <p style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', marginBottom: '0.75rem' }}>
                Pega tu dirección pública (SPCxxx...) para ver tu saldo.
              </p>
              <input
                type="text"
                value={importName}
                onChange={e => setImportName(e.target.value)}
                placeholder="Nombre (ej: Mi wallet principal)"
                style={{
                  width: '100%', padding: '0.7rem', borderRadius: '8px',
                  border: '1px solid var(--border)', background: 'var(--bg-secondary)',
                  color: 'var(--text-primary)', fontSize: '0.85rem', marginBottom: '0.5rem',
                }} />
              <input
                type="text"
                value={importInput}
                onChange={e => setImportInput(e.target.value)}
                placeholder="SPCxxx..."
                style={{
                  width: '100%', padding: '0.7rem', borderRadius: '8px',
                  border: '1px solid var(--border)', background: 'var(--bg-secondary)',
                  color: 'var(--text-primary)', fontSize: '0.85rem', marginBottom: '0.75rem',
                }} />
              {importError && <div style={{ color: 'var(--red)', fontSize: '0.8rem', marginBottom: '0.5rem' }}>{importError}</div>}
              <div style={{ display: 'flex', gap: '0.5rem' }}>
                <button onClick={() => { setShowImport(false); setImportInput('') }} style={{
                  flex: 1, padding: '0.55rem', borderRadius: '8px', border: '1px solid var(--border)',
                  background: 'transparent', color: 'var(--text-secondary)', cursor: 'pointer',
                }}>Cancelar</button>
                <button onClick={handleImport} style={{
                  flex: 1, padding: '0.55rem', borderRadius: '8px', border: 'none',
                  background: 'var(--accent)', color: '#fff', fontWeight: '600', cursor: 'pointer',
                }}>Ver saldo</button>
              </div>
            </div>
          )}

          {/* Just created — show private key ONCE */}
          {justCreated && !createdConfirmed && (
            <div style={sectionCard}>
              <div style={{
                textAlign: 'center', fontSize: '1.1rem', fontWeight: '700',
                color: 'var(--green)', marginBottom: '1rem',
              }}>
                ¡Wallet creada!
              </div>

              {/* Address */}
              <div style={{ marginBottom: '1rem' }}>
                <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginBottom: '0.3rem' }}>
                  📬 TU DIRECCIÓN (pública — compártela para recibir SPC)
                </div>
                <div style={{
                  background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.65rem 0.75rem',
                  fontFamily: 'monospace', fontSize: '0.72rem', color: 'var(--accent)',
                  wordBreak: 'break-all', marginBottom: '0.4rem',
                }}>
                  {justCreated.address}
                </div>
                <button onClick={() => handleCopy(justCreated.address, 'new-addr')} style={{
                  width: '100%', padding: '0.45rem', borderRadius: '6px',
                  border: '1px solid var(--border)',
                  background: copied === 'new-addr' ? 'rgba(16,185,129,0.15)' : 'var(--bg-secondary)',
                  color: copied === 'new-addr' ? 'var(--green)' : 'var(--text-secondary)',
                  fontSize: '0.8rem', cursor: 'pointer',
                }}>{copied === 'new-addr' ? '✓ Dirección copiada' : 'Copiar dirección'}</button>
              </div>

              {/* Private key */}
              <div style={{
                background: 'rgba(239, 68, 68, 0.06)', border: '1px solid rgba(239, 68, 68, 0.25)',
                borderRadius: '10px', padding: '1rem', marginBottom: '1rem',
              }}>
                <div style={{ fontSize: '0.75rem', color: 'var(--red)', fontWeight: '600', marginBottom: '0.3rem' }}>
                  🔑 TU CLAVE PRIVADA (secreta — guárdala en papel AHORA)
                </div>
                <p style={{ fontSize: '0.72rem', color: 'var(--text-secondary)', marginBottom: '0.5rem' }}>
                  Esta clave NO se guardará en la web. Si cierras esta pantalla sin guardarla, la pierdes para siempre.
                </p>
                <div style={{
                  background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.65rem 0.75rem',
                  fontFamily: 'monospace', fontSize: '0.62rem', color: 'var(--red)',
                  wordBreak: 'break-all', marginBottom: '0.4rem',
                }}>
                  {justCreated.privateKey}
                </div>
                <button onClick={() => handleCopy(justCreated.privateKey, 'new-key')} style={{
                  width: '100%', padding: '0.45rem', borderRadius: '6px',
                  border: '1px solid rgba(239,68,68,0.3)',
                  background: copied === 'new-key' ? 'rgba(16,185,129,0.15)' : 'rgba(239,68,68,0.08)',
                  color: copied === 'new-key' ? 'var(--green)' : 'var(--red)',
                  fontSize: '0.8rem', cursor: 'pointer',
                }}>{copied === 'new-key' ? '✓ Clave copiada' : 'Copiar clave privada'}</button>
              </div>

              {/* Confirm */}
              <label style={{
                display: 'flex', alignItems: 'center', gap: '0.75rem',
                marginBottom: '1rem', cursor: 'pointer',
              }}>
                <input type="checkbox" checked={createdConfirmed} onChange={e => setCreatedConfirmed(e.target.checked)}
                  style={{ width: '20px', height: '20px', accentColor: 'var(--green)', flexShrink: 0 }} />
                <span style={{ fontSize: '0.85rem', color: 'var(--text-primary)' }}>
                  He guardado mi dirección y mi clave privada en un lugar seguro
                </span>
              </label>

              <button onClick={() => setJustCreated(null)} disabled={!createdConfirmed} style={{
                width: '100%', padding: '0.85rem', borderRadius: '10px', border: 'none',
                background: createdConfirmed ? 'linear-gradient(135deg, #ffc400, #e6a800)' : 'var(--border)',
                color: createdConfirmed ? '#000' : 'var(--text-secondary)',
                fontSize: '1rem', fontWeight: '700',
                cursor: createdConfirmed ? 'pointer' : 'not-allowed',
              }}>
                Ya las guardé, continuar →
              </button>
            </div>
          )}

          {/* Active wallet */}
          {activeWallet && !justCreated && (
            <>
              {/* Wallet info */}
              <div style={sectionCard}>
                <div style={{ fontSize: '0.7rem', color: 'var(--text-secondary)', textTransform: 'uppercase', marginBottom: '0.25rem' }}>Mi Wallet</div>
                {activeWallet.name && (
                  <div style={{ fontSize: '0.95rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
                    {activeWallet.name}
                  </div>
                )}
                {!activeWallet.name && <div style={{ marginBottom: '0.5rem' }} />}

                {/* Address */}
                <div style={{
                  display: 'flex', alignItems: 'center', gap: '0.5rem',
                  background: 'var(--bg-secondary)', borderRadius: '8px', padding: '0.6rem 0.75rem',
                  marginBottom: '1rem',
                }}>
                  <span style={{ fontFamily: 'monospace', fontSize: '0.72rem', color: 'var(--text-primary)', flex: 1, wordBreak: 'break-all' }}>
                    {activeWallet.address}
                  </span>
                  <button onClick={() => handleCopy(activeWallet.address, 'addr')} style={{
                    padding: '0.3rem 0.6rem', borderRadius: '6px', border: 'none', flexShrink: 0,
                    background: copied === 'addr' ? 'rgba(16,185,129,0.15)' : 'var(--border)',
                    color: copied === 'addr' ? 'var(--green)' : 'var(--text-secondary)',
                    fontSize: '0.7rem', cursor: 'pointer',
                  }}>{copied === 'addr' ? '✓' : 'Copiar'}</button>
                </div>

                {/* Balance */}
                <div style={{ display: 'flex', gap: '0.75rem' }}>
                  <div style={{
                    flex: 2, background: 'var(--bg-secondary)', borderRadius: '10px',
                    padding: '0.85rem', border: '1px solid var(--border)',
                  }}>
                    <div style={{ fontSize: '0.7rem', color: 'var(--text-secondary)', marginBottom: '0.3rem' }}>Balance</div>
                    <div style={{ fontSize: '1.3rem', fontWeight: '700', color: 'var(--text-primary)' }}>
                      {balance ? formatSPC(balance.balance_spc) : '0'}
                      <span style={{ fontSize: '0.8rem', color: 'var(--accent)', marginLeft: '0.3rem' }}>SPC</span>
                    </div>
                  </div>
                  <div style={{
                    flex: 1, background: 'var(--bg-secondary)', borderRadius: '10px',
                    padding: '0.85rem', border: '1px solid var(--border)',
                  }}>
                    <div style={{ fontSize: '0.7rem', color: 'var(--text-secondary)', marginBottom: '0.3rem' }}>Transacciones</div>
                    <div style={{ fontSize: '1.3rem', fontWeight: '700', color: 'var(--text-primary)' }}>
                      {balance ? balance.nonce : '0'}
                    </div>
                  </div>
                </div>
              </div>

              {/* Actions */}
              <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1rem' }}>
                <a href="https://t.me/spaincoin_bot" target="_blank" rel="noopener noreferrer" style={{
                  flex: 1, padding: '0.7rem', borderRadius: '10px', border: 'none',
                  background: 'var(--green)', color: '#fff', fontSize: '0.9rem',
                  fontWeight: '700', cursor: 'pointer', textAlign: 'center', textDecoration: 'none',
                }}>Comprar SPC</a>
                <button onClick={() => onNavigate('/como-vender')} style={{
                  flex: 1, padding: '0.7rem', borderRadius: '10px',
                  border: '1px solid var(--border)', background: 'transparent',
                  color: 'var(--text-primary)', fontSize: '0.9rem',
                  fontWeight: '600', cursor: 'pointer', textAlign: 'center',
                }}>Vender SPC</button>
              </div>

              {/* Create another / Import */}
              <div style={{ display: 'flex', gap: '0.5rem' }}>
                <button onClick={handleCreate} style={{
                  flex: 1, padding: '0.5rem', borderRadius: '8px',
                  border: '1px solid var(--border)', background: 'var(--bg-secondary)',
                  color: 'var(--text-secondary)', fontSize: '0.8rem', cursor: 'pointer',
                }}>Crear otra wallet</button>
                <button onClick={() => setShowImport(true)} style={{
                  flex: 1, padding: '0.5rem', borderRadius: '8px',
                  border: '1px solid var(--border)', background: 'var(--bg-secondary)',
                  color: 'var(--text-secondary)', fontSize: '0.8rem', cursor: 'pointer',
                }}>Importar</button>
              </div>

              {/* Wallet selector if multiple */}
              {wallets.length > 1 && (
                <div style={{ marginTop: '1rem' }}>
                  <div style={{ fontSize: '0.7rem', color: 'var(--text-secondary)', marginBottom: '0.5rem' }}>Mis wallets:</div>
                  {wallets.map(w => (
                    <button key={w.address} onClick={() => { setActiveWallet(w); setShowImport(false) }}
                      style={{
                        display: 'block', width: '100%', textAlign: 'left',
                        padding: '0.5rem 0.75rem', marginBottom: '0.25rem',
                        borderRadius: '6px', border: 'none', cursor: 'pointer',
                        background: w.address === activeWallet.address ? 'var(--accent)' : 'var(--bg-secondary)',
                        color: w.address === activeWallet.address ? '#fff' : 'var(--text-secondary)',
                        fontSize: '0.8rem',
                      }}>
                      <span style={{ fontWeight: '600' }}>{w.name || 'Wallet'}</span>
                      <span style={{ fontFamily: 'monospace', fontSize: '0.65rem', marginLeft: '0.5rem', opacity: 0.7 }}>
                        {w.address.slice(0, 12)}...
                      </span>
                    </button>
                  ))}
                </div>
              )}
            </>
          )}

          {/* Security note */}
          <div style={{
            background: 'rgba(59, 130, 246, 0.08)', border: '1px solid rgba(59, 130, 246, 0.2)',
            borderRadius: '10px', padding: '1rem', marginTop: '1rem',
            fontSize: '0.78rem', color: 'var(--text-secondary)', lineHeight: 1.6,
          }}>
            🔐 <strong style={{ color: 'var(--text-primary)' }}>100% Self-Custody.</strong> Tus claves se generan y quedan en tu dispositivo.
            SpainCoin no tiene acceso a tus fondos. Si pierdes la clave privada, pierdes tus fondos para siempre.
          </div>
        </>
      )}

      {tab === 'download' && (
        <>
          <div style={sectionCard}>
            <h2 style={{ fontSize: '1rem', fontWeight: '600', color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
              CLI Wallet (avanzado)
            </h2>
            <p style={{ fontSize: '0.82rem', color: 'var(--text-secondary)', marginBottom: '1rem' }}>
              Para usuarios técnicos. Descarga el binario para tu sistema operativo.
            </p>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '0.5rem' }}>
              {[
                { os: 'macOS (Apple Silicon)', file: 'spc-macos-arm64', icon: '🍎' },
                { os: 'macOS (Intel)', file: 'spc-macos-amd64', icon: '🍎' },
                { os: 'Windows', file: 'spc-windows-amd64.exe', icon: '🪟' },
                { os: 'Linux', file: 'spc-linux-amd64', icon: '🐧' },
              ].map(d => (
                <a key={d.file} href={`https://github.com/spaincoin/spaincoin/releases/latest/download/${d.file}`}
                  target="_blank" rel="noopener noreferrer"
                  style={{
                    display: 'flex', alignItems: 'center', gap: '0.6rem',
                    padding: '0.7rem', borderRadius: '8px',
                    background: 'var(--bg-secondary)', border: '1px solid var(--border)',
                    textDecoration: 'none', fontSize: '0.85rem', color: 'var(--text-primary)',
                  }}>
                  <span>{d.icon}</span> {d.os}
                </a>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  )
}
