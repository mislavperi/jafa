import { AuthRequiredError } from './auth-error'

export async function apiFetch(input: RequestInfo | URL, init: RequestInit = {}): Promise<Response> {
  const response = await fetch(input, {
    credentials: 'include',
    ...init,
  })
  if (response.status === 401 || response.status === 403) {
    throw new AuthRequiredError()
  }
  return response
}
