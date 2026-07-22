import { defineStore } from 'pinia'
import { ref } from 'vue'
import { googleLogin } from '@/services/auth.api'

interface AuthUser {
  id: string
  email: string
  name: string
  role: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))

  const user = ref<AuthUser | null>(JSON.parse(localStorage.getItem('user') ?? 'null'))

  const loading = ref(false)
  const error = ref<string | null>(null)

  async function loginWithGoogle(idToken: string) {
    loading.value = true
    error.value = null

    try {
      const response = await googleLogin(idToken)

      token.value = response.token
      user.value = response.user

      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
    } catch (err) {
      error.value = 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = null
    user.value = null

    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token,
    user,
    loading,
    error,
    loginWithGoogle,
    logout,
  }
})
