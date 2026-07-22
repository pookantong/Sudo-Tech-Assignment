<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { getMovies } from '@/services/movie.api'
import type { Movie } from '@/types/movie'

const router = useRouter()

const movies = ref<Movie[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    movies.value = await getMovies()
  } catch (err) {
    console.error(err)
    error.value = 'Failed to load movies'
  } finally {
    loading.value = false
  }
})

function openMovie(movieId: string) {
  router.push(`/movies/${movieId}`)
}
</script>

<template>
  <main class="movies-page">
    <!-- Header -->
    <header class="page-header">
      <div>
        <h1>
          Movies
        </h1>

        <p class="page-description">
          Choose your favorite movie and book your seats now.
        </p>
      </div>
    </header>

    <!-- Loading -->
    <section
      v-if="loading"
      class="movie-grid"
    >
      <div
        v-for="index in 6"
        :key="index"
        class="movie-card skeleton-card"
      >
        <div class="skeleton skeleton-poster" />

        <div class="skeleton skeleton-title" />
      </div>
    </section>

    <!-- Error -->
    <section
      v-else-if="error"
      class="state-container"
    >
      <div class="state-icon">
        ⚠️
      </div>

      <h2>
        Something went wrong
      </h2>

      <p>
        {{ error }}
      </p>
    </section>

    <!-- Empty -->
    <section
      v-else-if="movies.length === 0"
      class="state-container"
    >
      <div class="state-icon">
        🎬
      </div>

      <h2>
        No movies available
      </h2>

      <p>
        There are currently no movies showing.
      </p>
    </section>

    <!-- Movies -->
    <section
      v-else
      class="movie-grid"
    >
      <article
        v-for="movie in movies"
        :key="movie.id"
        class="movie-card"
        @click="openMovie(movie.id)"
      >
        <div class="poster-wrapper">
          <img
            :src="movie.image_url"
            :alt="movie.title"
            class="movie-poster"
          >

          <div class="poster-overlay">
            <span>
              View Details
            </span>
          </div>
        </div>

        <div class="movie-info">
          <h2>
            {{ movie.title }}
          </h2>

          <button
            class="book-button"
            @click.stop="openMovie(movie.id)"
          >
            View Showtimes
          </button>
        </div>
      </article>
    </section>
  </main>
</template>

<style scoped>
.movies-page {
  max-width: 1280px;
  margin: 0 auto;
  padding: 48px 32px 80px;
}

/* =========================
   Header
========================= */

.page-header {
  margin-bottom: 40px;
}

.eyebrow {
  margin: 0 0 8px;

  color: #6366f1;

  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.12em;
}

.page-header h1 {
  margin: 0;

  font-size: 42px;
  font-weight: 800;
  letter-spacing: -0.03em;

  color: #fff;
}

.page-description {
  margin: 12px 0 0;

  color: #6b7280;

  font-size: 16px;
}

/* =========================
   Movie Grid
========================= */

.movie-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 28px;
}

/* =========================
   Movie Card
========================= */

.movie-card {
  overflow: hidden;

  background: #ffffff;

  border: 1px solid #e5e7eb;
  border-radius: 16px;

  cursor: pointer;

  transition:
    transform 0.25s ease,
    box-shadow 0.25s ease,
    border-color 0.25s ease;
}

.movie-card:hover {
  transform: translateY(-6px);

  border-color: #c7d2fe;

  box-shadow:
    0 20px 40px rgb(0 0 0 / 10%);
}

/* =========================
   Poster
========================= */

.poster-wrapper {
  position: relative;

  overflow: hidden;

  aspect-ratio: 2 / 3;

  background: #f3f4f6;
}

.movie-poster {
  display: block;

  width: 100%;
  height: 100%;

  object-fit: cover;

  transition: transform 0.4s ease;
}

.movie-card:hover .movie-poster {
  transform: scale(1.05);
}

/* =========================
   Overlay
========================= */

.poster-overlay {
  position: absolute;
  inset: 0;

  display: flex;
  align-items: flex-end;
  justify-content: center;

  padding: 20px;

  background: linear-gradient(
    to top,
    rgb(0 0 0 / 70%),
    transparent 50%
  );

  opacity: 0;

  transition: opacity 0.25s ease;
}

.movie-card:hover .poster-overlay {
  opacity: 1;
}

.poster-overlay span {
  padding: 8px 16px;

  color: #ffffff;

  background: rgb(0 0 0 / 60%);

  border-radius: 999px;

  font-size: 13px;
  font-weight: 600;
}

/* =========================
   Movie Info
========================= */

.movie-info {
  padding: 18px;
}

.movie-info h2 {
  overflow: hidden;

  margin: 0 0 16px;

  color: #111827;

  font-size: 18px;
  font-weight: 700;

  white-space: nowrap;
  text-overflow: ellipsis;
}

.book-button {
  width: 100%;

  padding: 11px 16px;

  color: #ffffff;

  background: #4f46e5;

  border: 0;
  border-radius: 8px;

  font-size: 14px;
  font-weight: 600;

  cursor: pointer;

  transition:
    background 0.2s ease,
    transform 0.2s ease;
}

.book-button:hover {
  background: #4338ca;
}

.book-button:active {
  transform: scale(0.98);
}

/* =========================
   Empty / Error State
========================= */

.state-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  min-height: 360px;

  text-align: center;
}

.state-icon {
  margin-bottom: 16px;

  font-size: 48px;
}

.state-container h2 {
  margin: 0 0 8px;

  color: #111827;

  font-size: 24px;
}

.state-container p {
  margin: 0;

  color: #6b7280;
}

/* =========================
   Skeleton Loading
========================= */

.skeleton-card {
  cursor: default;
}

.skeleton-card:hover {
  transform: none;

  border-color: #e5e7eb;

  box-shadow: none;
}

.skeleton {
  position: relative;

  overflow: hidden;

  background: #e5e7eb;
}

.skeleton::after {
  position: absolute;
  inset: 0;

  content: '';

  background: linear-gradient(
    90deg,
    transparent,
    rgb(255 255 255 / 50%),
    transparent
  );

  animation: shimmer 1.5s infinite;
}

.skeleton-poster {
  aspect-ratio: 2 / 3;
}

.skeleton-title {
  width: 70%;
  height: 20px;

  margin: 20px;
}

@keyframes shimmer {
  from {
    transform: translateX(-100%);
  }

  to {
    transform: translateX(100%);
  }
}
</style>