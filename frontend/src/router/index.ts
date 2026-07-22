import { createRouter, createWebHistory } from 'vue-router'

import MovieListView from '@/views/movie/MovieListView.vue'
import MovieDetailView from '@/views/movie/MovieDetailView.vue'
import SeatSelectionView from '@/views/booking/SeatSelectionView.vue'
import BookingConfirmView from '@/views/booking/BookingConfirmView.vue'
import BookingResultView from '@/views/booking/BookingResultView.vue'
import AdminDashboardView from '@/views/booking/AdminDashboardView.vue'
import LoginView from '@/views/auth/LoginView.vue'
import { isAdmin, getToken } from '@/utils/auth'

const router = createRouter({
  history: createWebHistory(),

  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },

    {
      path: '/',
      redirect: '/movies',
    },

    {
      path: '/movies',
      component: MovieListView,
    },

    {
      path: '/movies/:id',
      component: MovieDetailView,
    },

    {
      path: '/showtimes/:id/seats',
      component: SeatSelectionView,

      meta: {
        requiresAuth: true,
      },
    },

    {
      path: '/bookings/:id/confirm',
      name: 'booking-confirm',
      component: BookingConfirmView,

      meta: {
        requiresAuth: true,
      },
    },

    {
      path: '/booking-result',
      name: 'booking-result',
      component: BookingResultView,

      meta: {
        requiresAuth: true,
      },
    },

    {
      path: '/admin',
      name: 'admin-dashboard',
      component: AdminDashboardView,
      meta: {
        requiresAuth: true,
        requiresAdmin: true,
      },
    },
  ],
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    return {
      name: 'login',
      query: {
        redirect: to.fullPath,
      },
    }
  }

  if (to.name === 'login' && token) {
    return { path: '/movies' }
  }

  if (to.meta.requiresAdmin && !isAdmin()) {
    return {
      name: 'movies',
    }
  }

  return true
})

export default router
