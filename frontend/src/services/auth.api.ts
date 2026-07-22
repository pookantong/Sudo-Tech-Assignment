import api from './api'

export interface GoogleLoginRequest {
  id_token: string
}

export interface AuthResponse {
  token: string
  user: {
    id: string
    email: string
    name: string
    role: string
  }
}

export async function googleLogin(idToken: string): Promise<AuthResponse> {
  const response = await api.post<AuthResponse>('/auth/google', {
    id_token: idToken,
  })

  return response.data
}
