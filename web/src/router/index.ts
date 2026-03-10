import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getMe } from '@/modules/auth/api/auth'
import { queryClient } from '@/core/query'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/modules/auth/views/LoginPage.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/modules/auth/views/RegisterPage.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      name: 'home',
      component: () => import('@/modules/expense/views/ExpensePage.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  if (!authStore.bootstrapped) {
    try {
      const user = await queryClient.fetchQuery({ queryKey: ['auth', 'me'], queryFn: getMe })
      authStore.setUser(user)
    } catch {
      authStore.clearUser()
    }
    authStore.setBootstrapped()
  }

  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (!requiresAuth && authStore.isAuthenticated) {
    return { name: 'home' }
  }
})

export default router
