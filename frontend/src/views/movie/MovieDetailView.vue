<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { useRoute, useRouter } from 'vue-router'

import { getShowtimes } from '@/services/movie.api'

import type { Showtime } from '@/types/movie'

const route = useRoute()
const router = useRouter()

const showtimes = ref<Showtime[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const today = new Date()

const selectedDate = ref(formatDateValue(today))

// ====================
// Date
// ====================

const dates = computed(() => {
  const result: Date[] = []

  for (let i = 0; i < 30; i++) {
    const date = new Date(today)

    date.setDate(today.getDate() + i)

    result.push(date)
  }

  return result
})

function formatDateValue(date: Date): string {
  const year = date.getFullYear()

  const month = String(date.getMonth() + 1).padStart(2, '0')

  const day = String(date.getDate()).padStart(2, '0')

  return `${year}-${month}-${day}`
}

function formatWeekday(date: Date): string {
  return date.toLocaleDateString('en-US', {
    weekday: 'short',
  })
}

function formatMonth(date: Date): string {
  return date.toLocaleDateString('en-US', {
    month: 'short',
    year: 'numeric',
  })
}

function selectDate(date: Date) {
  selectedDate.value = formatDateValue(date)
}

// ====================
// Showtime
// ====================

const filteredShowtimes = computed(() => {
  return showtimes.value.filter((showtime) => {
    const date = new Date(showtime.starts_at)

    return formatDateValue(date) === selectedDate.value
  })
})

// ====================
// Group by Cinema
// ====================

const groupedCinemas = computed(() => {
  const cinemaMap = new Map<string, Map<string, Showtime[]>>()

  for (const showtime of filteredShowtimes.value) {
    const cinemaName = showtime.cinema_name || 'Unknown Cinema'

    const hallName = showtime.hall_name || 'Unknown Hall'

    if (!cinemaMap.has(cinemaName)) {
      cinemaMap.set(cinemaName, new Map())
    }

    const halls = cinemaMap.get(cinemaName)!

    if (!halls.has(hallName)) {
      halls.set(hallName, [])
    }

    halls.get(hallName)!.push(showtime)
  }

  return Array.from(cinemaMap.entries()).map(([cinemaName, halls]) => ({
    cinemaName,
    halls: Array.from(halls.entries()).map(([hallName, showtimes]) => ({
      hallName,
      showtimes,
    })),
  }))
})

// ====================
// Format Time
// ====================

function formatTime(date: string): string {
  return new Date(date).toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

// ====================
// Select Showtime
// ====================

function selectShowtime(showtimeID: string) {
  router.push(`/showtimes/${showtimeID}/seats`)
}

// ====================
// Load Data
// ====================

onMounted(async () => {
  try {
    showtimes.value = await getShowtimes(route.params.id as string)
  } catch (err) {
    console.error(err)

    error.value = 'Failed to load showtimes'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <main class="showtime-page">
    <div class="container">
      <!-- Step Header -->
      <div class="steps">
        <div class="step active">เลือกรอบภาพยนตร์</div>

        <div class="step">เลือกที่นั่ง</div>

        <div class="step">ซื้อตั๋ว</div>
      </div>

      <!-- Date Selector -->
      <section class="date-section">
        <div class="date-header">
          <h3>
            {{ formatMonth(dates[0]!) }}
          </h3>
        </div>

        <div class="date-list">
          <button
            v-for="date in dates"
            :key="formatDateValue(date)"
            class="date-card"
            :class="{
              selected: formatDateValue(date) === selectedDate,
            }"
            @click="selectDate(date)"
          >
            <span class="weekday">
              {{ formatWeekday(date) }}
            </span>

            <strong>
              {{ date.getDate() }}
            </strong>
          </button>
        </div>
      </section>

      <!-- Loading -->
      <p v-if="loading" class="status">Loading showtimes...</p>

      <!-- Error -->
      <p v-else-if="error" class="status error">
        {{ error }}
      </p>

      <!-- Empty -->
      <p v-else-if="groupedCinemas.length === 0" class="status">
        No showtimes available for this date
      </p>

      <!-- Cinemas -->
      <section v-else class="cinema-list">
        <div v-for="cinema in groupedCinemas" :key="cinema.cinemaName" class="cinema-section">
          <!-- Cinema Name -->
          <div class="cinema-header">
            <h2>
              {{ cinema.cinemaName }}
            </h2>
          </div>

          <!-- Halls -->
          <div v-for="hall in cinema.halls" :key="hall.hallName" class="hall-section">
            <div class="hall-header">
              <h3>
                {{ hall.hallName }}
              </h3>
            </div>

            <!-- Showtime Buttons -->
            <div class="time-list">
              <button
                v-for="showtime in hall.showtimes"
                :key="showtime.id"
                class="time-button"
                @click="selectShowtime(showtime.id)"
              >
                {{ formatTime(showtime.starts_at) }}
              </button>
            </div>
          </div>
        </div>
      </section>
    </div>
  </main>
</template>

<style scoped>
.showtime-page {
  min-height: 100vh;
  padding: 24px 0 60px;
  background: #000;
  color: #fff;
}

/* ==================== */
/* Steps */
/* ==================== */

.steps {
  display: grid;
  grid-template-columns: repeat(3, 1fr);

  height: 56px;
  margin-bottom: 28px;

  overflow: hidden;

  border: 1px solid #fff;
  border-radius: 12px;
}

.step {
  display: flex;
  align-items: center;
  justify-content: center;

  font-size: 18px;
  font-weight: 600;

  background: #000;
}

.step.active {
  color: #111;
  background: #fff;
}

/* ==================== */
/* Date */
/* ==================== */

.date-section {
  margin-bottom: 36px;
}

.date-header {
  text-align: center;
  margin-bottom: 12px;
}

.date-header h3 {
  margin: 0;

  font-size: 14px;
  font-weight: 600;
}

.date-list {
  display: flex;
  gap: 14px;

  overflow-x: auto;
  padding: 4px 0 10px;
}

.date-list::-webkit-scrollbar {
  height: 4px;
}

.date-card {
  flex: 0 0 72px;

  height: 66px;

  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  gap: 4px;

  border: 1px solid #c99d36;
  border-radius: 12px;

  color: #c99d36;
  background: #000;

  cursor: pointer;

  transition: 0.2s ease;
}

.date-card:hover {
  background: #2b210c;
}

.date-card.selected {
  color: #111;
  background: #d6ad55;
}

.weekday {
  font-size: 12px;
}

.date-card strong {
  font-size: 20px;
}

/* ==================== */
/* Cinema */
/* ==================== */

.cinema-list {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.cinema-section {
  border-bottom: 1px solid #555;
  padding-bottom: 28px;
}

.cinema-header {
  display: flex;
  align-items: center;
  justify-content: space-between;

  margin-bottom: 24px;
}

.cinema-header h2 {
  margin: 0;

  font-size: 20px;
  font-weight: 700;
}

.arrow {
  font-size: 32px;
  font-weight: 300;
}

/* ==================== */
/* Hall */
/* ==================== */

.hall-section {
  margin-bottom: 24px;
}

.hall-header {
  margin-bottom: 14px;
}

.hall-header h3 {
  margin: 0;

  font-size: 18px;
  font-weight: 600;
}

/* ==================== */
/* Times */
/* ==================== */

.time-list {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
}

.time-button {
  min-width: 100px;

  padding: 12px 20px;

  border: 1px solid #ddd;
  border-radius: 6px;

  color: #111;
  background: #fff;

  font-size: 18px;
  font-weight: 600;

  cursor: pointer;

  transition: 0.2s ease;
}

.time-button:hover {
  color: #111;
  background: #d6ad55;
  border-color: #d6ad55;
}

/* ==================== */
/* Status */
/* ==================== */

.status {
  padding: 40px 0;

  color: #aaa;
  text-align: center;
}

.error {
  color: #e57373;
}

/* ==================== */
/* Mobile */
/* ==================== */

@media (max-width: 768px) {
  .showtime-page {
    padding: 16px 0 40px;
  }

  .steps {
    height: 48px;
  }

  .step {
    font-size: 14px;
  }

  .date-card {
    flex-basis: 64px;
  }

  .cinema-header h2 {
    font-size: 18px;
  }

  .hall-header h3 {
    font-size: 16px;
  }

  .time-button {
    min-width: 88px;
    padding: 10px 14px;
    font-size: 16px;
  }
}
</style>
