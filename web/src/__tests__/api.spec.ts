import { describe, it, expect, vi, afterEach } from 'vitest'
import { apiFetch } from '@/core/api'
import { AuthRequiredError } from '@/core/auth-error'

function mockFetchOnce(status: number, body: unknown = {}): void {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValueOnce({
      status,
      ok: status >= 200 && status < 300,
      json: () => Promise.resolve(body),
    }),
  )
}

describe('apiFetch', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('returns the response on success', async () => {
    mockFetchOnce(200, { id: 1 })
    const res = await apiFetch('/api/test')
    expect(res.ok).toBe(true)
  })

  it('always passes credentials: include', async () => {
    const spy = vi.fn().mockResolvedValueOnce({ status: 200, ok: true, json: () => Promise.resolve({}) })
    vi.stubGlobal('fetch', spy)
    await apiFetch('/api/test')
    expect(spy).toHaveBeenCalledWith('/api/test', expect.objectContaining({ credentials: 'include' }))
  })

  it('merges caller-supplied init without losing credentials', async () => {
    const spy = vi.fn().mockResolvedValueOnce({ status: 200, ok: true })
    vi.stubGlobal('fetch', spy)
    await apiFetch('/api/test', { method: 'POST' })
    expect(spy).toHaveBeenCalledWith(
      '/api/test',
      expect.objectContaining({ credentials: 'include', method: 'POST' }),
    )
  })

  it('throws AuthRequiredError on 401', async () => {
    mockFetchOnce(401)
    await expect(apiFetch('/api/secret')).rejects.toThrow(AuthRequiredError)
  })

  it('throws AuthRequiredError on 403', async () => {
    mockFetchOnce(403)
    await expect(apiFetch('/api/secret')).rejects.toThrow(AuthRequiredError)
  })

  it('does NOT throw on non-ok statuses other than 401/403', async () => {
    mockFetchOnce(404)
    const res = await apiFetch('/api/missing')
    expect(res.ok).toBe(false)
    expect(res.status).toBe(404)
  })
})
