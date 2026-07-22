import { jwtDecode } from 'jwt-decode'

interface TokenPayload {
  user_id: string
  role: 'USER' | 'ADMIN'
  exp?: number
}

export function getToken(): string | null {
  return localStorage.getItem('token')
}

export function getUserRole(): 'USER' | 'ADMIN' | null {
  const token = getToken()

  if (!token) {
    return null
  }

  try {
    const payload = jwtDecode<TokenPayload>(token)

    return payload.role
  } catch {
    return null
  }
}

export function isAdmin(): boolean {
  return getUserRole() === 'ADMIN'
}
