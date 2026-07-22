<script setup lang="ts">
import {
  computed,
  onMounted,
  reactive,
  ref,
} from 'vue'

import {
  listAdminBookings,
} from '@/services/booking.api'

import type {
  Booking,
  BookingStatus,
} from '@/types/booking'

const bookings = ref<Booking[]>([])

const loading = ref(false)

const error = ref('')

const filter = reactive({
  user: '',
  showtime: '',
  status: '',
})

const statusOptions = [
  {
    label: 'ทั้งหมด',
    value: '',
  },
  {
    label: 'Pending',
    value: 'PENDING',
  },
  {
    label: 'Success',
    value: 'SUCCESS',
  },
  {
    label: 'Timeout',
    value: 'TIMEOUT',
  },
  {
    label: 'Failed',
    value: 'FAILED',
  },
]

const totalBookings = computed(() => {
  return bookings.value.length
})

const successBookings = computed(() => {
  return bookings.value.filter(
    (booking) =>
      booking.status === 'SUCCESS',
  ).length
})

const pendingBookings = computed(() => {
  return bookings.value.filter(
    (booking) =>
      booking.status === 'PENDING',
  ).length
})

const failedBookings = computed(() => {
  return bookings.value.filter(
    (booking) =>
      booking.status === 'FAILED' ||
      booking.status === 'TIMEOUT',
  ).length
})

const totalRevenue = computed(() => {
  return bookings.value
    .filter(
      (booking) =>
        booking.status === 'SUCCESS',
    )
    .reduce(
      (total, booking) =>
        total + booking.price,
      0,
    )
})

async function fetchBookings() {
  loading.value = true

  error.value = ''

  try {
    bookings.value =
      await listAdminBookings({
        user:
          filter.user ||
          undefined,

        showtime:
          filter.showtime ||
          undefined,

        status:
          filter.status ||
          undefined,
      })
  } catch (err) {
    console.error(err)

    error.value =
      'ไม่สามารถโหลดข้อมูลการจองได้'
  } finally {
    loading.value = false
  }
}

async function clearFilter() {
  filter.user = ''

  filter.showtime = ''

  filter.status = ''

  await fetchBookings()
}

function formatPrice(
  price: number,
) {
  return `${price.toLocaleString()} บาท`
}

function formatDate(
  date: string,
) {
  return new Date(
    date,
  ).toLocaleString(
    'th-TH',
  )
}

function getStatusClass(
  status: BookingStatus,
) {
  switch (status) {
    case 'SUCCESS':
      return 'success'

    case 'PENDING':
      return 'pending'

    case 'TIMEOUT':
      return 'timeout'

    case 'FAILED':
      return 'failed'

    default:
      return ''
  }
}

function getStatusLabel(
  status: BookingStatus,
) {
  switch (status) {
    case 'SUCCESS':
      return 'Success'

    case 'PENDING':
      return 'Pending'

    case 'TIMEOUT':
      return 'Timeout'

    case 'FAILED':
      return 'Failed'

    default:
      return status
  }
}

onMounted(() => {
  fetchBookings()
})
</script>

<template>
  <div class="dashboard">
    <!-- Header -->
    <header class="dashboard-header">
      <div>
        <h1>
          Admin Dashboard
        </h1>

        <p>
          จัดการและตรวจสอบรายการจองภาพยนตร์
        </p>
      </div>

      <button
        class="refresh-button"
        :disabled="loading"
        @click="fetchBookings"
      >
        <span
          v-if="loading"
        >
          กำลังโหลด...
        </span>

        <span
          v-else
        >
          รีเฟรช
        </span>
      </button>
    </header>

    <!-- Statistics -->
    <section class="stats">
      <!-- Total -->
      <div class="stat-card">
        <div class="stat-content">
          <span class="stat-label">
            รายการจองทั้งหมด
          </span>

          <strong class="stat-value">
            {{ totalBookings }}
          </strong>
        </div>

        <div class="stat-icon">
          📋
        </div>
      </div>

      <!-- Success -->
      <div
        class="
          stat-card
          success-card
        "
      >
        <div class="stat-content">
          <span class="stat-label">
            สำเร็จ
          </span>

          <strong class="stat-value">
            {{ successBookings }}
          </strong>
        </div>

        <div class="stat-icon">
          ✓
        </div>
      </div>

      <!-- Pending -->
      <div
        class="
          stat-card
          pending-card
        "
      >
        <div class="stat-content">
          <span class="stat-label">
            กำลังดำเนินการ
          </span>

          <strong class="stat-value">
            {{ pendingBookings }}
          </strong>
        </div>

        <div class="stat-icon">
          ⏳
        </div>
      </div>

      <!-- Revenue -->
      <div
        class="
          stat-card
          revenue-card
        "
      >
        <div class="stat-content">
          <span class="stat-label">
            รายได้จากรายการสำเร็จ
          </span>

          <strong class="stat-value">
            {{ formatPrice(totalRevenue) }}
          </strong>
        </div>

        <div class="stat-icon">
          ฿
        </div>
      </div>
    </section>

    <!-- Filter -->
    <section class="filter-card">
      <div class="filter-header">
        <div>
          <h2>
            ค้นหารายการจอง
          </h2>

          <p>
            Filter bookings by user, showtime, or status
          </p>
        </div>
      </div>

      <div class="filter-form">
        <!-- User -->
        <div class="filter-group">
          <label>
            User ID
          </label>

          <input
            v-model="filter.user"
            type="text"
            placeholder="กรอก User ID"
            @keyup.enter="fetchBookings"
          />
        </div>

        <!-- Showtime -->
        <div class="filter-group">
          <label>
            Showtime ID
          </label>

          <input
            v-model="filter.showtime"
            type="text"
            placeholder="กรอก Showtime ID"
            @keyup.enter="fetchBookings"
          />
        </div>

        <!-- Status -->
        <div class="filter-group">
          <label>
            Status
          </label>

          <select
            v-model="filter.status"
          >
            <option
              v-for="option in statusOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </option>
          </select>
        </div>

        <!-- Actions -->
        <div class="filter-actions">
          <button
            class="search-button"
            @click="fetchBookings"
          >
            ค้นหา
          </button>

          <button
            class="clear-button"
            @click="clearFilter"
          >
            ล้าง
          </button>
        </div>
      </div>
    </section>

    <!-- Error -->
    <div
      v-if="error"
      class="error-message"
    >
      {{ error }}
    </div>

    <!-- Bookings -->
    <section class="table-card">
      <!-- Table Header -->
      <div class="table-header">
        <div>
          <h2>
            รายการจอง
          </h2>

          <p>
            รายการทั้งหมด
            {{ bookings.length }}
            รายการ
          </p>
        </div>

        <button
          class="small-refresh-button"
          :disabled="loading"
          @click="fetchBookings"
        >
          ↻
        </button>
      </div>

      <!-- Loading -->
      <div
        v-if="loading"
        class="loading-state"
      >
        <div class="loading-spinner"></div>

        <p>
          กำลังโหลดข้อมูล...
        </p>
      </div>

      <!-- Empty -->
      <div
        v-else-if="
          bookings.length === 0
        "
        class="empty-state"
      >
        <div class="empty-icon">
          📭
        </div>

        <h3>
          ไม่พบข้อมูลการจอง
        </h3>

        <p>
          ลองเปลี่ยนเงื่อนไขการค้นหา
        </p>
      </div>

      <!-- Table -->
      <div
        v-else
        class="table-wrapper"
      >
        <table>
          <thead>
            <tr>
              <th>
                #
              </th>

              <th>
                Booking ID
              </th>

              <th>
                User ID
              </th>

              <th>
                Showtime ID
              </th>

              <th>
                Seat
              </th>

              <th>
                Price
              </th>

              <th>
                Status
              </th>

              <th>
                Created At
              </th>

              <th>
                Updated At
              </th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="(
                booking,
                index
              ) in bookings"
              :key="booking.id"
            >
              <!-- Number -->
              <td>
                <span
                  class="row-number"
                >
                  {{ index + 1 }}
                </span>
              </td>

              <!-- Booking ID -->
              <td>
                <span
                  class="id-text"
                  :title="
                    booking.id
                  "
                >
                  {{ booking.id }}
                </span>
              </td>

              <!-- User ID -->
              <td>
                <span
                  class="id-text"
                  :title="
                    booking.user_id
                  "
                >
                  {{ booking.user_id }}
                </span>
              </td>

              <!-- Showtime ID -->
              <td>
                <span
                  class="id-text"
                  :title="
                    booking.showtime_id
                  "
                >
                  {{ booking.showtime_id }}
                </span>
              </td>

              <!-- Seat -->
              <td>
                <span
                  class="seat"
                >
                  {{ booking.seat_label }}
                </span>
              </td>

              <!-- Price -->
              <td>
                <span
                  class="price"
                >
                  {{
                    formatPrice(
                      booking.price,
                    )
                  }}
                </span>
              </td>

              <!-- Status -->
              <td>
                <span
                  class="status"
                  :class="
                    getStatusClass(
                      booking.status,
                    )
                  "
                >
                  <span
                    class="status-dot"
                  ></span>

                  {{
                    getStatusLabel(
                      booking.status,
                    )
                  }}
                </span>
              </td>

              <!-- Created -->
              <td>
                <span
                  class="date"
                >
                  {{
                    formatDate(
                      booking.created_at,
                    )
                  }}
                </span>
              </td>

              <!-- Updated -->
              <td>
                <span
                  class="date"
                >
                  {{
                    formatDate(
                      booking.updated_at,
                    )
                  }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<style scoped>
* {
  box-sizing: border-box;
}

.dashboard {
  min-height: 100vh;

  padding: 32px;

  background: #f5f7fb;

  color: #1f2937;
}

/* =========================
   Header
========================= */

.dashboard-header {
  display: flex;

  justify-content: space-between;

  align-items: center;

  margin-bottom: 28px;
}

.dashboard-header h1 {
  margin: 0;

  color: #111827;

  font-size: 30px;

  font-weight: 700;
}

.dashboard-header p {
  margin: 8px 0 0;

  color: #6b7280;

  font-size: 14px;
}

.refresh-button {
  border: none;

  border-radius: 8px;

  padding: 11px 20px;

  background: #111827;

  color: white;

  font-size: 14px;

  font-weight: 600;

  cursor: pointer;

  transition: 0.2s;
}

.refresh-button:hover {
  background: #374151;
}

.refresh-button:disabled {
  opacity: 0.6;

  cursor: not-allowed;
}

/* =========================
   Statistics
========================= */

.stats {
  display: grid;

  grid-template-columns:
    repeat(4, 1fr);

  gap: 18px;

  margin-bottom: 24px;
}

.stat-card {
  display: flex;

  justify-content: space-between;

  align-items: center;

  padding: 22px;

  background: white;

  border-radius: 12px;

  border-left: 4px solid #6366f1;

  box-shadow:
    0 2px 8px
    rgba(
      0,
      0,
      0,
      0.05
    );
}

.success-card {
  border-left-color: #22c55e;
}

.pending-card {
  border-left-color: #f59e0b;
}

.revenue-card {
  border-left-color: #8b5cf6;
}

.stat-content {
  display: flex;

  flex-direction: column;

  gap: 8px;
}

.stat-label {
  color: #6b7280;

  font-size: 13px;

  font-weight: 500;
}

.stat-value {
  color: #111827;

  font-size: 28px;

  font-weight: 700;
}

.stat-icon {
  display: flex;

  justify-content: center;

  align-items: center;

  width: 42px;

  height: 42px;

  border-radius: 10px;

  background: #f3f4f6;

  color: #374151;

  font-size: 20px;

  font-weight: 700;
}

/* =========================
   Filter
========================= */

.filter-card {
  margin-bottom: 24px;

  padding: 22px;

  background: white;

  border-radius: 12px;

  box-shadow:
    0 2px 8px
    rgba(
      0,
      0,
      0,
      0.05
    );
}

.filter-header {
  margin-bottom: 20px;
}

.filter-header h2 {
  margin: 0;

  color: #111827;

  font-size: 18px;

  font-weight: 700;
}

.filter-header p {
  margin: 6px 0 0;

  color: #6b7280;

  font-size: 13px;
}

.filter-form {
  display: grid;

  grid-template-columns:
    1fr
    1fr
    220px
    auto;

  gap: 16px;

  align-items: end;
}

.filter-group {
  display: flex;

  flex-direction: column;

  gap: 8px;
}

.filter-group label {
  color: #374151;

  font-size: 13px;

  font-weight: 600;
}

.filter-group input,
.filter-group select {
  width: 100%;

  height: 42px;

  padding: 0 12px;

  border: 1px solid #d1d5db;

  border-radius: 8px;

  outline: none;

  background: white;

  color: #1f2937;

  font-size: 14px;

  transition: 0.2s;
}

.filter-group input:focus,
.filter-group select:focus {
  border-color: #6366f1;

  box-shadow:
    0 0 0 3px
    rgba(
      99,
      102,
      241,
      0.1
    );
}

.filter-actions {
  display: flex;

  gap: 8px;
}

.search-button,
.clear-button {
  height: 42px;

  padding: 0 18px;

  border: none;

  border-radius: 8px;

  font-size: 14px;

  font-weight: 600;

  cursor: pointer;
}

.search-button {
  background: #4f46e5;

  color: white;
}

.search-button:hover {
  background: #4338ca;
}

.clear-button {
  background: #e5e7eb;

  color: #374151;
}

.clear-button:hover {
  background: #d1d5db;
}

/* =========================
   Error
========================= */

.error-message {
  margin-bottom: 20px;

  padding: 14px 16px;

  border-radius: 8px;

  background: #fee2e2;

  color: #b91c1c;

  font-size: 14px;
}

/* =========================
   Table
========================= */

.table-card {
  overflow: hidden;

  background: white;

  border-radius: 12px;

  box-shadow:
    0 2px 8px
    rgba(
      0,
      0,
      0,
      0.05
    );
}

.table-header {
  display: flex;

  justify-content: space-between;

  align-items: center;

  padding: 20px 22px;

  border-bottom: 1px solid #e5e7eb;
}

.table-header h2 {
  margin: 0;

  color: #111827;

  font-size: 18px;
}

.table-header p {
  margin: 6px 0 0;

  color: #6b7280;

  font-size: 13px;
}

.small-refresh-button {
  width: 36px;

  height: 36px;

  border: none;

  border-radius: 8px;

  background: #f3f4f6;

  color: #374151;

  font-size: 20px;

  cursor: pointer;
}

.small-refresh-button:hover {
  background: #e5e7eb;
}

.small-refresh-button:disabled {
  opacity: 0.5;

  cursor: not-allowed;
}

.table-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;

  min-width: 1250px;

  border-collapse: collapse;
}

th,
td {
  padding: 16px 18px;

  text-align: left;

  border-bottom: 1px solid #f1f5f9;

  white-space: nowrap;
}

th {
  background: #f8fafc;

  color: #475569;

  font-size: 12px;

  font-weight: 700;

  text-transform: uppercase;
}

td {
  color: #374151;

  font-size: 13px;
}

tbody tr:hover {
  background: #fafafa;
}

.row-number {
  color: #9ca3af;
}

.id-text {
  display: block;

  max-width: 180px;

  overflow: hidden;

  text-overflow: ellipsis;

  color: #4b5563;

  font-family: monospace;

  font-size: 12px;
}

.seat {
  display: inline-block;

  padding: 5px 10px;

  border-radius: 6px;

  background: #eef2ff;

  color: #4338ca;

  font-size: 12px;

  font-weight: 700;
}

.price {
  color: #374151;

  font-weight: 600;
}

.date {
  color: #6b7280;

  font-size: 12px;
}

.status {
  display: inline-flex;

  align-items: center;

  gap: 6px;

  padding: 5px 10px;

  border-radius: 999px;

  font-size: 11px;

  font-weight: 700;
}

.status-dot {
  width: 6px;

  height: 6px;

  border-radius: 50%;

  background: currentColor;
}

.status.success {
  background: #dcfce7;

  color: #166534;
}

.status.pending {
  background: #fef3c7;

  color: #92400e;
}

.status.timeout {
  background: #f3f4f6;

  color: #4b5563;
}

.status.failed {
  background: #fee2e2;

  color: #991b1b;
}

/* =========================
   Loading
========================= */

.loading-state {
  display: flex;

  flex-direction: column;

  justify-content: center;

  align-items: center;

  min-height: 300px;

  color: #6b7280;
}

.loading-spinner {
  width: 32px;

  height: 32px;

  margin-bottom: 16px;

  border: 3px solid #e5e7eb;

  border-top-color: #4f46e5;

  border-radius: 50%;

  animation:
    spin
    0.8s
    linear
    infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* =========================
   Empty
========================= */

.empty-state {
  display: flex;

  flex-direction: column;

  justify-content: center;

  align-items: center;

  min-height: 300px;

  color: #6b7280;
}

.empty-icon {
  margin-bottom: 12px;

  font-size: 42px;
}

.empty-state h3 {
  margin: 0;

  color: #374151;

  font-size: 16px;
}

.empty-state p {
  margin-top: 8px;

  font-size: 13px;
}
</style>