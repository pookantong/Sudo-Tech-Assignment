<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'

import { useRoute, useRouter } from 'vue-router'

import { getSeatMap } from '@/services/movie.api'

import { selectSeat } from '@/services/booking.api'

import type { Seat } from '@/types/seat'

const route = useRoute()
const router = useRouter()

const seats = ref<Seat[]>([])
const selectedSeat = ref<Seat | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

let socket: WebSocket | undefined

const showtimeId = route.params.id as string

/**
 * จัดกลุ่ม Seat ตาม Row
 *
 * เช่น
 *
 * row 0
 * A1 A2 A3 A4
 *
 * row 1
 * B1 B2 B3 B4
 *
 * และเรียงตาม col
 */
const seatRows = computed(() => {
  const rows = new Map<number, Seat[]>()

  for (const seat of seats.value) {
    if (!rows.has(seat.row)) {
      rows.set(seat.row, [])
    }

    rows.get(seat.row)!.push(seat)
  }

  return Array.from(rows.entries())
    .sort(([rowA], [rowB]) => rowA - rowB)
    .map(([, rowSeats]) => rowSeats.sort((a, b) => a.col - b.col))
})

onMounted(async () => {
  try {
    seats.value = await getSeatMap(showtimeId)

    connectWebSocket()
  } catch (err) {
    console.error(err)

    error.value = 'Failed to load seats'
  } finally {
    loading.value = false
  }
})

function connectWebSocket() {
  const wsUrl = `${import.meta.env.VITE_WS_URL}/ws`

  socket = new WebSocket(wsUrl)

  socket.onopen = () => {
    console.log('[WebSocket] Connected')
  }

  socket.onmessage = (event) => {
    try {
      const update = JSON.parse(event.data)

      console.log('[WebSocket] Seat update:', update)

      if (update.showtime_id !== showtimeId) {
        return
      }

      const seat = seats.value.find((item) => item.id === update.seat_id)

      if (!seat) {
        return
      }

      // อัปเดตสถานะที่นั่งแบบ realtime
      seat.status = update.status

      // ถ้าที่นั่งที่กำลังเลือก
      // ถูกคนอื่น Lock หรือ Book
      if (selectedSeat.value?.id === seat.id && seat.status !== 'AVAILABLE') {
        selectedSeat.value = null
      }
    } catch (err) {
      console.error('[WebSocket] Invalid message:', err)
    }
  }

  socket.onerror = (event) => {
    console.error('[WebSocket] Error:', event)
  }

  socket.onclose = () => {
    console.log('[WebSocket] Disconnected')
  }
}

function chooseSeat(seat: Seat) {
  if (seat.status !== 'AVAILABLE') {
    return
  }

  /**
   * กดที่นั่งเดิมอีกครั้ง
   * เพื่อยกเลิกการเลือก
   */
  if (selectedSeat.value?.id === seat.id) {
    selectedSeat.value = null
    return
  }

  selectedSeat.value = seat
}

async function continueBooking() {
  if (!selectedSeat.value) {
    return
  }

  try {
    const booking = await selectSeat(
      showtimeId,
      selectedSeat.value.id,
    )

    // เก็บ Booking ที่ได้จาก Backend
    sessionStorage.setItem(
      'pending-booking',
      JSON.stringify(booking),
    )

    router.push({
      name: 'booking-confirm',
      params: {
        id: booking.id,
      },
    })
  } catch (err) {
    console.error(err)

    error.value = 'Failed to select seat'

    try {
      seats.value = await getSeatMap(
        showtimeId,
      )
    } catch (refreshErr) {
      console.error(refreshErr)
    }

    selectedSeat.value = null
  }
}

function goBack() {
  router.back()
}

onUnmounted(() => {
  socket?.close()
})
</script>

<template>
  <main class="seat-page">
    <div class="container">
      <!-- ==================== -->
      <!-- Step Header -->
      <!-- ==================== -->

      <div class="steps">
        <div class="step" @click="goBack">เลือกรอบภาพยนตร์</div>

        <div class="step active">เลือกที่นั่ง</div>

        <div class="step">ซื้อตั๋ว</div>
      </div>

      <h1>Select Seat</h1>

      <!-- ==================== -->
      <!-- Loading -->
      <!-- ==================== -->

      <p v-if="loading" class="status">Loading seats...</p>

      <!-- ==================== -->
      <!-- Error -->
      <!-- ==================== -->

      <p v-else-if="error" class="status error">
        {{ error }}
      </p>

      <!-- ==================== -->
      <!-- Seat Content -->
      <!-- ==================== -->

      <template v-else>
        <!-- Screen -->

        <div class="screen">SCREEN</div>

        <!-- ==================== -->
        <!-- Legend -->
        <!-- ==================== -->

        <div class="legend">
          <div class="legend-item">
            <span class="legend-seat available" />

            <span> Available </span>
          </div>

          <div class="legend-item">
            <span class="legend-seat selected" />

            <span> Selected </span>
          </div>

          <div class="legend-item">
            <span class="legend-seat locked" />

            <span> Locked </span>
          </div>

          <div class="legend-item">
            <span class="legend-seat booked" />

            <span> Booked </span>
          </div>
        </div>

        <!-- ==================== -->
        <!-- Seat Map -->
        <!-- ==================== -->

        <div class="seat-map">
          <!-- Row -->

          <div v-for="(rowSeats, rowIndex) in seatRows" :key="rowIndex" class="seat-row">
            <!-- Seats -->

            <button
              v-for="seat in rowSeats"
              :key="seat.id"
              class="seat"
              :class="{
                available: seat.status === 'AVAILABLE',

                locked: seat.status === 'LOCKED',

                booked: seat.status === 'BOOKED',

                selected: selectedSeat?.id === seat.id,
              }"
              :disabled="seat.status !== 'AVAILABLE'"
              @click="chooseSeat(seat)"
            >
              {{ seat.label }}
            </button>
          </div>
        </div>

        <!-- ==================== -->
        <!-- Bottom Action -->
        <!-- ==================== -->

        <div class="booking-action">
          <div v-if="selectedSeat" class="selected-info">
            <div>
              Selected seat:

              <strong>
                {{ selectedSeat.label }}
              </strong>
            </div>

            <div class="price-info">
              Price:

              <strong>
                ฿{{ selectedSeat.price.toFixed(2) }}
              </strong>
            </div>
          </div>

          <div v-else class="selected-info">
            Please select a seat
          </div>

          <button
            class="continue-button"
            :disabled="!selectedSeat"
            @click="continueBooking"
          >
            Continue
          </button>
        </div>
      </template>
    </div>
  </main>
</template>

<style scoped>
/* ==================== */
/* Page */
/* ==================== */

.seat-page {
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

  margin-bottom: 32px;

  overflow: hidden;

  border: 1px solid #fff;

  border-radius: 12px;
}

.step {
  display: flex;

  align-items: center;

  justify-content: center;

  color: #fff;

  background: #000;

  font-size: 18px;

  font-weight: 600;

  cursor: default;
}

.step.active {
  color: #111;

  background: #fff;
}

.step:first-child {
  cursor: pointer;
}

/* ==================== */
/* Screen */
/* ==================== */

.screen {
  width: 80%;

  max-width: 700px;

  margin: 40px auto 48px;

  padding: 12px;

  color: #111;

  background: #fff;

  border-radius: 50% 50% 0 0;

  text-align: center;

  font-size: 14px;

  font-weight: 700;
}

/* ==================== */
/* Legend */
/* ==================== */

.legend {
  display: flex;

  justify-content: center;

  flex-wrap: wrap;

  gap: 24px;

  margin-bottom: 32px;

  color: #aaa;
}

.legend-item {
  display: flex;

  align-items: center;

  gap: 8px;

  font-size: 14px;
}

.legend-seat {
  width: 22px;

  height: 22px;

  border-radius: 5px;
}

/* ==================== */
/* Seat Map */
/* ==================== */

.seat-map {
  display: flex;

  flex-direction: column;

  align-items: center;

  gap: 12px;

  max-width: 700px;

  margin: 0 auto 40px;
}

.seat-row {
  display: flex;

  justify-content: center;

  gap: 12px;
}

/* ==================== */
/* Seat */
/* ==================== */

.seat {
  width: 48px;

  height: 42px;

  border: none;

  border-radius: 6px;

  color: #111;

  background: #fff;

  font-size: 12px;

  font-weight: 600;

  cursor: pointer;

  transition: 0.2s ease;
}

.seat.available:hover {
  background: #d6ad55;

  transform: translateY(-2px);
}

.seat.selected {
  color: #111;

  background: #d6ad55;

  transform: translateY(-2px);
}

.seat.locked {
  color: #666;
  background: #b8b8b8;

  cursor: not-allowed;
}

.seat.booked {
  color: #bbb;

  background: #333;

  cursor: not-allowed;
}

/* ==================== */
/* Legend Colors */
/* ==================== */

.legend-seat.available {
  background: #fff;
}

.legend-seat.selected {
  background: #d6ad55;
}

.legend-seat.locked {
  background: #b8b8b8;
}

.legend-seat.booked {
  background: #333;
}

/* ==================== */
/* Action */
/* ==================== */

.booking-action {
  display: flex;

  align-items: center;

  justify-content: space-between;

  max-width: 700px;

  margin: 0 auto;

  padding-top: 24px;

  border-top: 1px solid #444;
}

.selected-info {
  color: #aaa;
}

.selected-info strong {
  margin-left: 8px;

  color: #d6ad55;

  font-size: 20px;
}

.continue-button {
  padding: 12px 28px;

  border: none;

  border-radius: 8px;

  color: #111;

  background: #d6ad55;

  font-size: 16px;

  font-weight: 700;

  cursor: pointer;
}

.continue-button:disabled {
  opacity: 0.4;

  cursor: not-allowed;
}

/* ==================== */
/* Status */
/* ==================== */

.status {
  padding: 60px 0;

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
  .seat-page {
    padding: 16px 0 40px;
  }

  .steps {
    height: 48px;
  }

  .step {
    font-size: 14px;
  }

  .seat-map {
    gap: 8px;

    overflow-x: auto;

    align-items: flex-start;

    padding: 0 16px;
  }

  .seat-row {
    gap: 8px;
  }

  .seat {
    flex: 0 0 42px;

    width: 42px;

    height: 38px;
  }

  .booking-action {
    flex-direction: column;

    align-items: stretch;

    gap: 20px;
  }

  .continue-button {
    width: 100%;
  }
}
</style>
