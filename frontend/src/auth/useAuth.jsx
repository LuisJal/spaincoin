import { useState, useEffect, createContext, useContext } from 'react'
import { getMe } from '../api/client.js'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)      // {email, address, balance_spc}
  const [token, setToken] = useState(() => localStorage.getItem('spc_token'))
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) { setLoading(false); return }
    getMe(token)
      .then(data => setUser(data))
      .catch(() => { localStorage.removeItem('spc_token'); setToken(null) })
      .finally(() => setLoading(false))
  }, [token])

  function saveAuth(data) {
    localStorage.setItem('spc_token', data.token)
    setToken(data.token)
    setUser({ email: data.email, address: data.address, balance_spc: 0 })
  }

  function logout() {
    localStorage.removeItem('spc_token')
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, token, loading, saveAuth, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}
