import api from './api'
import type { Booking } from '@/types/booking'

export async function selectSeat(showtimeId: string, seatId: string) {
  const { data } = await api.post<Booking>('/seats/select', {
    showtime_id: showtimeId,
    seat_id: seatId,
  })

  return data
}

export async function confirmPayment(bookingId: string, outcome: 'success' | 'fail') {
  const { data } = await api.post('/bookings/confirm', {
    booking_id: bookingId,
    outcome,
  })

  return data
}

export interface AdminBookingFilter {
  user?: string
  showtime?: string
  status?: string
}

export async function listAdminBookings(filter: AdminBookingFilter = {}) {
  const { data } = await api.get<Booking[]>('/admin/bookings', {
    params: filter,
  })

  return data
}
