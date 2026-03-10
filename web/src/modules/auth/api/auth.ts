import { AUTH_API } from '@/core/constants/api'
import type { User, LoginRequest, RegisterRequest } from '../models/auth'

export async function login(payload: LoginRequest): Promise<User> {
  const response = await fetch(`${AUTH_API}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    const data = await response.json().catch(() => ({}))
    throw new Error(data.error ?? 'Login failed')
  }
  return response.json()
}

export async function logout(): Promise<void> {
  const response = await fetch(`${AUTH_API}/logout`, {
    method: 'POST',
    credentials: 'include',
  })
  if (!response.ok) {
    throw new Error('Logout failed')
  }
}

export async function register(payload: RegisterRequest): Promise<User> {
  const response = await fetch(`${AUTH_API}/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    const data = await response.json().catch(() => ({}))
    throw new Error(data.error ?? 'Registration failed')
  }
  return response.json()
}

export async function getMe(): Promise<User> {
  const response = await fetch(`${AUTH_API}/me`, {
    credentials: 'include',
  })
  if (!response.ok) {
    throw new Error('Not authenticated')
  }
  return response.json()
}
