export interface User {
  id: number
  username: string
  first_name?: string
  last_name?: string
  email?: string
  avatar_url?: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  first_name?: string
  last_name?: string
  email?: string
}
