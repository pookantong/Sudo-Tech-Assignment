<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const success = route.query.status === 'success'

function backToMovies() {
  router.push('/movies')
}
</script>

<template>
  <main class="result-page">
    <div class="container">
      <!-- Step Header -->
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

      <!-- Result Card -->
      <div
        class="result-card"
        :class="{
          success,
          failed: !success,
        }"
      >
        <div
          class="result-icon"
          :class="{
            success,
            failed: !success,
          }"
        >
          {{ success ? '✓' : '✕' }}
        </div>

        <h1 v-if="success">
          Booking Successful
        </h1>

        <h1 v-else>
          Payment Failed
        </h1>

        <p v-if="success">
          Your movie ticket has been booked successfully.
        </p>

        <p v-else>
          We couldn't complete your payment.
          Please try again.
        </p>

        <button
          class="back-button"
          @click="backToMovies"
        >
          Back to Movies
        </button>
      </div>
    </div>
  </main>
</template>

<style scoped>
.result-page {
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
  margin-bottom: 80px;

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
/* Result Card */
/* ==================== */

.result-card {
  display: flex;
  flex-direction: column;
  align-items: center;

  max-width: 520px;
  margin: 0 auto;
  padding: 48px 32px;

  text-align: center;

  border: 1px solid #333;
  border-radius: 16px;

  background: #111;
}

.result-card h1 {
  margin: 20px 0 12px;

  font-size: 32px;
}

.result-card p {
  margin-bottom: 32px;

  color: #aaa;
  font-size: 16px;
}

/* ==================== */
/* Icon */
/* ==================== */

.result-icon {
  display: flex;
  align-items: center;
  justify-content: center;

  width: 80px;
  height: 80px;

  border-radius: 50%;

  font-size: 42px;
  font-weight: 700;
}

.result-icon.success {
  color: #111;
  background: #d6ad55;
}

.result-icon.failed {
  color: #fff;
  background: #8b3a3a;
}

/* ==================== */
/* Button */
/* ==================== */

.back-button {
  padding: 14px 32px;

  border: none;
  border-radius: 8px;

  color: #111;
  background: #d6ad55;

  font-size: 16px;
  font-weight: 700;

  cursor: pointer;

  transition: 0.2s ease;
}

.back-button:hover {
  background: #e2bd6c;
}

/* ==================== */
/* Mobile */
/* ==================== */

@media (max-width: 768px) {
  .result-page {
    padding: 16px 0 40px;
  }

  .steps {
    height: 48px;
    margin-bottom: 48px;
  }

  .step {
    font-size: 14px;
  }

  .result-card {
    padding: 40px 24px;
  }

  .result-card h1 {
    font-size: 26px;
  }
}
</style>