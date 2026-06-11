import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import type { User } from '@/modules/auth/models/auth'

const alice: User = { id: 1, username: 'alice' }

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts unauthenticated and unbootstrapped', () => {
    const store = useAuthStore()
    expect(store.isAuthenticated).toBe(false)
    expect(store.currentUser).toBeNull()
    expect(store.bootstrapped).toBe(false)
  })

  it('setUser populates currentUser and flips isAuthenticated', () => {
    const store = useAuthStore()
    store.setUser(alice)
    expect(store.currentUser).toEqual(alice)
    expect(store.isAuthenticated).toBe(true)
  })

  it('clearUser resets state', () => {
    const store = useAuthStore()
    store.setUser(alice)
    store.setBootstrapped()
    store.clearUser()
    expect(store.currentUser).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(store.bootstrapped).toBe(false)
  })

  it('setBootstrapped marks the store as bootstrapped', () => {
    const store = useAuthStore()
    expect(store.bootstrapped).toBe(false)
    store.setBootstrapped()
    expect(store.bootstrapped).toBe(true)
  })

  it('setUser(null) leaves isAuthenticated false', () => {
    const store = useAuthStore()
    store.setUser(null)
    expect(store.isAuthenticated).toBe(false)
  })
})
