import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { useDarkModeStore } from '@/stores/darkMode'
import { getMe, AuthRequiredError } from '@/modules/auth/api/auth'
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
    {
      path: '/expenses',
      name: 'expenses',
      component: () => import('@/modules/expense/views/ExpensesListPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/reports',
      name: 'reports',
      component: () => import('@/modules/reports/views/ReportsPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories',
      name: 'categories',
      component: () => import('@/modules/categories/views/CategoriesPage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/modules/settings/views/SettingsPage.vue'),
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
      const dark = useDarkModeStore()
      const theme = useThemeStore()
      await theme.load(dark.isDark)
    } catch (err) {
      // Only clear auth on explicit auth failure. Network/5xx → leave state alone
      // so transient outages don't kick users out.
      if (err instanceof AuthRequiredError) {
        authStore.clearUser()
      } else {
        console.warn('Auth bootstrap failed (non-auth error):', err)
      }
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
