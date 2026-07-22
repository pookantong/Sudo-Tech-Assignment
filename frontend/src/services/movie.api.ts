import api from './api'
import type { Movie, Showtime } from '@/types/movie'
import type { Seat } from '@/types/seat'

export async function getMovies() {
  const { data } = await api.get<Movie[]>('/movies')
  return data
}

export async function getShowtimes(movieId: string) {
  const { data } = await api.get<Showtime[]>(`/movies/${movieId}/showtimes`)

  return data
}

export async function getSeatMap(showtimeId: string) {
  const { data } = await api.get<Seat[]>(`/showtimes/${showtimeId}/seats`)

  return data
}
