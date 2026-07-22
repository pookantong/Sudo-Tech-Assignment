<script setup lang="ts">
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth.store'

const router = useRouter()
const authStore = useAuthStore()

const { token, user } = storeToRefs(authStore)

function logout() {
  authStore.logout()

  router.push({
    name: 'login',
  })
}
</script>

<template>
  <nav class="navbar">
    <div class="navbar-container">
      <div class="nav-left">
        <RouterLink
          to="/movies"
          class="logo"
        >
          CINEMA
        </RouterLink>

        <RouterLink
          v-if="
            token &&
            user?.role === 'ADMIN'
          "
          to="/admin"
          class="dashboard-button"
        >
          Dashboard
        </RouterLink>
      </div>

      <div class="nav-right">
        <template v-if="token">
          <span class="profile">
            {{ user?.name ?? user?.email }}
          </span>

          <button
            class="logout-button"
            @click="logout"
          >
            Logout
          </button>
        </template>

        <RouterLink
          v-else
          to="/login"
          class="login-button"
        >
          Login
        </RouterLink>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  position: sticky;

  top: 0;

  z-index: 100;

  width: 100%;

  border-bottom: 1px solid #333;

  background: #000;
}

.navbar-container {
  display: flex;

  align-items: center;

  justify-content: space-between;

  max-width: 1200px;

  height: 72px;

  margin: 0 auto;

  padding: 0 24px;
}

.nav-left {
  display: flex;

  align-items: center;

  gap: 28px;
}

.logo {
  color: #d6ad55;

  font-size: 24px;

  font-weight: 800;

  text-decoration: none;
}

.nav-right {
  display: flex;

  align-items: center;

  gap: 20px;
}

.profile {
  color: #fff;

  font-size: 16px;

  font-weight: 600;
}

.dashboard-button,
.logout-button,
.login-button {
  padding: 8px 18px;

  border-radius: 6px;

  font-size: 14px;

  font-weight: 600;

  cursor: pointer;
}

.dashboard-button {
  border: 1px solid #d6ad55;

  color: #d6ad55;

  text-decoration: none;
}

.dashboard-button:hover {
  color: #111;

  background: #d6ad55;
}

.logout-button {
  border: 1px solid #8b3a3a;

  color: #e57373;

  background: transparent;
}

.logout-button:hover {
  color: #fff;

  background: #8b3a3a;
}

.login-button {
  border: 1px solid #d6ad55;

  color: #d6ad55;

  text-decoration: none;
}

.login-button:hover {
  color: #111;

  background: #d6ad55;
}
</style>