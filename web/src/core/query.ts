import { QueryClient, QueryCache, MutationCache } from '@tanstack/vue-query'
import { isAuthRequiredError } from './auth-error'

// Lazy router import to avoid circular dep on module init.
async function bounceToLogin() {
  const [{ default: router }, { useAuthStore }] = await Promise.all([
    import('@/router'),
    import('@/stores/auth'),
  ])
  try {
    useAuthStore().clearUser()
  } catch {
    // pinia not yet ready — best-effort
  }
  if (router.currentRoute.value.name !== 'login') {
    void router.push({ name: 'login', query: { redirect: router.currentRoute.value.fullPath } })
  }
}

function handleAuthFailure(err: unknown) {
  if (isAuthRequiredError(err)) {
    void bounceToLogin()
  }
}

export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: handleAuthFailure,
  }),
  mutationCache: new MutationCache({
    onError: handleAuthFailure,
  }),
  defaultOptions: {
    queries: {
      retry: (failureCount, error) => {
        if (isAuthRequiredError(error)) return false
        return failureCount < 2
      },
    },
  },
})
