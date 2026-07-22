<script setup lang="ts">
import { ref } from 'vue'

import { useRoute, useRouter } from 'vue-router'

import { confirmPayment } from '@/services/booking.api'

import type { Booking } from '@/types/booking'

const route = useRoute()
const router = useRouter()

const bookingId = route.params.id as string

const loading = ref(false)
const error = ref<string | null>(null)

// ====================
// Booking
// ====================

const storedBooking = sessionStorage.getItem(
  'pending-booking',
)

const booking = ref<Booking | null>(
  storedBooking
    ? JSON.parse(storedBooking)
    : null,
)

// ====================
// Payment
// ====================

async function pay(
  outcome: 'success' | 'fail',
) {
  loading.value = true
  error.value = null

  try {
    await confirmPayment(
      bookingId,
      outcome,
    )

    sessionStorage.removeItem(
      'pending-booking',
    )

    router.push({
      name: 'booking-result',
      query: {
        status:
          outcome === 'success'
            ? 'success'
            : 'failed',
      },
    })
  } catch (err) {
    console.error(err)

    error.value =
      'Payment failed. Please try again.'
  } finally {
    loading.value = false
  }
}

// ====================
// Format Price
// ====================

function formatPrice(
  price: number,
): string {
  return new Intl.NumberFormat(
    'th-TH',
    {
      style: 'currency',
      currency: 'THB',
    },
  ).format(price)
}

// ====================
// Back
// ====================

function goBack() {
  router.back()
}
</script>

<template>
  <main class="confirm-page">
    <div class="container">
      <!-- ==================== -->
      <!-- Step Header -->
      <!-- ==================== -->

      <div class="steps">
        <div class="step">
          เลือกรอบภาพยนตร์
        </div>

        <div class="step">
          เลือกที่นั่ง
        </div>

        <div class="step active">
          ซื้อตั๋ว
        </div>
      </div>

      <!-- ==================== -->
      <!-- Confirm Card -->
      <!-- ==================== -->

      <div class="confirm-card">
        <!-- Payment Icon -->

        <div class="payment-icon">
          ฿
        </div>

        <!-- Heading -->

        <h1>
          Confirm Booking
        </h1>

        <p class="subtitle">
          Complete your payment to confirm your booking.
        </p>

        <!-- ==================== -->
        <!-- Booking Information -->
        <!-- ==================== -->

        <div
          v-if="booking"
          class="booking-info"
        >
          <!-- Seat -->

          <div class="info-row">
            <span>
              Seat
            </span>

            <strong>
              {{ booking.seat_label }}
            </strong>
          </div>

          <!-- Price -->

          <div class="info-row total-row">
            <span>
              Total
            </span>

            <strong class="price">
              {{ formatPrice(booking.price) }}
            </strong>
          </div>
        </div>

        <!-- Booking Not Found -->

        <p
          v-else
          class="error"
        >
          Booking information not found.
        </p>

        <!-- Error -->

        <p
          v-if="error"
          class="error"
        >
          {{ error }}
        </p>

        <!-- ==================== -->
        <!-- Payment Actions -->
        <!-- ==================== -->

        <div class="payment-actions">
          <button
            class="pay-button"
            :disabled="loading || !booking"
            @click="pay('success')"
          >
            <span v-if="loading">
              Processing...
            </span>

            <span v-else>
              Pay Now
            </span>
          </button>

          <button
            class="fail-button"
            :disabled="loading || !booking"
            @click="pay('fail')"
          >
            Simulate Payment Failed
          </button>
        </div>

        <!-- ==================== -->
        <!-- Back -->
        <!-- ==================== -->

        <button
          class="back-button"
          :disabled="loading"
          @click="goBack"
        >
          Back
        </button>
      </div>
    </div>
  </main>
</template>

<style scoped>
/* ==================== */
/* Page */
/* ==================== */

.confirm-page {
  min-height: 100vh;

  padding: 24px 0 60px;

  color: #fff;
  background: #000;
}

/* ==================== */
/* Steps */
/* ==================== */

.steps {
  display: grid;

  grid-template-columns: repeat(3, 1fr);

  height: 56px;

  margin-bottom: 64px;

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
}

.step.active {
  color: #111;

  background: #fff;
}

/* ==================== */
/* Confirm Card */
/* ==================== */

.confirm-card {
  max-width: 520px;

  margin: 0 auto;

  padding: 40px 32px;

  text-align: center;

  border: 1px solid #333;

  border-radius: 16px;

  background: #111;
}

/* ==================== */
/* Payment Icon */
/* ==================== */

.payment-icon {
  display: flex;

  align-items: center;

  justify-content: center;

  width: 72px;

  height: 72px;

  margin: 0 auto 24px;

  color: #111;

  background: #d6ad55;

  border-radius: 50%;

  font-size: 36px;

  font-weight: 700;
}

/* ==================== */
/* Heading */
/* ==================== */

.confirm-card h1 {
  margin-bottom: 12px;

  font-size: 32px;
}

.subtitle {
  margin-bottom: 32px;

  color: #aaa;
}

/* ==================== */
/* Booking Info */
/* ==================== */

.booking-info {
  margin-bottom: 28px;

  padding: 16px;

  border: 1px solid #333;

  border-radius: 8px;

  background: #181818;
}

.info-row {
  display: flex;

  align-items: center;

  justify-content: space-between;

  padding: 12px 0;

  border-bottom: 1px solid #333;
}

.info-row:last-child {
  border-bottom: none;
}

.info-row span {
  color: #999;
}

.info-row strong {
  color: #fff;

  font-size: 18px;
}

.info-row .price {
  color: #d6ad55;

  font-size: 22px;
}

/* ==================== */
/* Payment Actions */
/* ==================== */

.payment-actions {
  display: flex;

  flex-direction: column;

  gap: 12px;
}

.pay-button,
.fail-button,
.back-button {
  width: 100%;

  padding: 14px 24px;

  border-radius: 8px;

  font-size: 16px;

  font-weight: 700;

  cursor: pointer;
}

/* ==================== */
/* Pay */
/* ==================== */

.pay-button {
  border: none;

  color: #111;

  background: #d6ad55;
}

.pay-button:hover {
  background: #e2bd6c;
}

/* ==================== */
/* Failed */
/* ==================== */

.fail-button {
  border: 1px solid #8b3a3a;

  color: #e57373;

  background: transparent;
}

.fail-button:hover {
  color: #fff;

  background: #8b3a3a;
}

/* ==================== */
/* Back */
/* ==================== */

.back-button {
  margin-top: 20px;

  border: 1px solid #444;

  color: #aaa;

  background: transparent;
}

.back-button:hover {
  color: #fff;

  background: #222;
}

/* ==================== */
/* Disabled */
/* ==================== */

.pay-button:disabled,
.fail-button:disabled,
.back-button:disabled {
  opacity: 0.5;

  cursor: not-allowed;
}

/* ==================== */
/* Error */
/* ==================== */

.error {
  margin-bottom: 20px;

  color: #e57373;
}

/* ==================== */
/* Mobile */
/* ==================== */

@media (max-width: 768px) {
  .confirm-page {
    padding: 16px 0 40px;
  }

  .steps {
    height: 48px;

    margin-bottom: 40px;
  }

  .step {
    font-size: 14px;
  }

  .confirm-card {
    padding: 32px 20px;
  }

  .confirm-card h1 {
    font-size: 26px;
  }
}
</style>