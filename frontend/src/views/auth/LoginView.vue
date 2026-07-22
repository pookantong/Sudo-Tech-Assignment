<script setup lang="ts">
import { GoogleLogin } from 'vue3-google-login'
import { useAuthStore } from '@/stores/auth.store'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

async function handleLoginSuccess(response: { credential: string }) {
  try {
    await authStore.loginWithGoogle(response.credential)
    router.push('/movies')
  } catch (error) {
    console.error(error)
  }
}

function handleLoginError() {
  authStore.error = 'Login failed'
}
</script>

<template>
  <main class="auth-page">
    <div class="login-card">
      <h1>Cinema Booking</h1>
      <p class="subtitle">Sign in to book your movie tickets</p>

      <GoogleLogin :callback="handleLoginSuccess" @error="handleLoginError" />

      <p v-if="authStore.error" class="error">
        {{ authStore.error }}
      </p>
    </div>
  </main>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.login-card {
  width: 100%;
  max-width: 420px;
  padding: 40px;
  background: white;
  border-radius: 16px;
  text-align: center;
  box-shadow: 0 10px 40px rgb(0 0 0 / 10%);
}

.login-card h1 {
  margin-bottom: 8px;
}

.subtitle {
  color: #666;
  margin-bottom: 32px;
}

.google-button {
  width: 100%;
  padding: 14px 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  font-weight: 600;
}

.google-button:hover {
  background: #f5f5f5;
}

.google-button:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.error {
  margin-top: 16px;
  color: #d32f2f;
}
</style>
