export type BookingStatus = 'PENDING' | 'SUCCESS' | 'TIMEOUT' | 'FAILED'

export interface Booking {
  id: string
  user_id: string
  showtime_id: string
  seat_id: string
  seat_label: string
  price: number
  status: BookingStatus
  created_at: string
  updated_at: string
}