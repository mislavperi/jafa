import { describe, it, expect, vi, afterEach } from 'vitest'
import { login, logout, register, deleteAccount } from '@/modules/auth/api/auth'

// auth.ts calls fetch directly (not apiFetch) for login/logout/register/deleteAccount
function stubFetch(status: number, body: unknown = {}): ReturnType<typeof vi.fn> {
  const mock = vi.fn().mockResolvedValue({
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  })
  vi.stubGlobal('fetch', mock)
  return mock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('login', () => {
  it('returns user on success', async () => {
    const user = { id: 1, username: 'alice' }
    stubFetch(200, user)
    const result = await login({ username: 'alice', password: 'secret' })
    expect(result).toEqual(user)
  })

  it('throws with server error message on failure', async () => {
    stubFetch(401, { error: 'invalid credentials' })
    await expect(login({ username: 'x', password: 'y' })).rejects.toThrow('invalid credentials')
  })

  it('falls back to generic message when response has no error field', async () => {
    stubFetch(500, {})
    await expect(login({ username: 'x', password: 'y' })).rejects.toThrow('Login failed')
  })

  it('sends POST to /api/auth/login with credentials: include', async () => {
    const spy = stubFetch(200, { id: 1, username: 'alice' })
    await login({ username: 'alice', password: 'secret' })
    expect(spy).toHaveBeenCalledWith(
      expect.stringContaining('/login'),
      expect.objectContaining({ method: 'POST', credentials: 'include' }),
    )
  })
})

describe('logout', () => {
  it('resolves without value on success', async () => {
    stubFetch(200)
    await expect(logout()).resolves.toBeUndefined()
  })

  it('throws on failure', async () => {
    stubFetch(500)
    await expect(logout()).rejects.toThrow('Logout failed')
  })
})

describe('register', () => {
  it('returns created user on success', async () => {
    const user = { id: 5, username: 'bob' }
    stubFetch(200, user)
    const result = await register({ username: 'bob', password: 'pw' })
    expect(result).toEqual(user)
  })

  it('throws with server error message on failure', async () => {
    stubFetch(409, { error: 'username taken' })
    await expect(register({ username: 'bob', password: 'pw' })).rejects.toThrow('username taken')
  })

  it('falls back to generic message', async () => {
    stubFetch(400, {})
    await expect(register({ username: 'x', password: 'y' })).rejects.toThrow('Registration failed')
  })
})

describe('deleteAccount', () => {
  it('resolves without value on success', async () => {
    stubFetch(200)
    await expect(deleteAccount()).resolves.toBeUndefined()
  })

  it('throws with server error message on failure', async () => {
    stubFetch(403, { error: 'forbidden' })
    await expect(deleteAccount()).rejects.toThrow('forbidden')
  })

  it('falls back to generic message', async () => {
    stubFetch(500, {})
    await expect(deleteAccount()).rejects.toThrow('Failed to delete account')
  })
})
