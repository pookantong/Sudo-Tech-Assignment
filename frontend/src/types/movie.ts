export interface Movie {
  id: string
  title: string
  image_url: string
}

export interface Showtime {
  id: string
  movie_id: string
  hall_id: string
  hall_name: string
  cinema_id: string
  cinema_name: string
  starts_at: string
  price: number
}
