export type SeatStatus = 'AVAILABLE' | 'LOCKED' | 'BOOKED'

export interface Seat {
  id: string
  label: string
  row: number
  col: number
  status: SeatStatus
  price: number
}