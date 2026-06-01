export class AuthRequiredError extends Error {
  constructor() {
    super('Authentication required')
    this.name = 'AuthRequiredError'
  }
}

export function isAuthRequiredError(e: unknown): e is AuthRequiredError {
  if (e instanceof AuthRequiredError) return true
  // Cross-module instanceof can fail with HMR / bundle splits; fall back to name.
  return typeof e === 'object' && e !== null && (e as { name?: string }).name === 'AuthRequiredError'
}
