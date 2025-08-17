import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/axios'
import menu from '@/router/menu.json'
import MainLayout from '@/components/layouts/MainLayout.vue'

function loadView(view) {
  return () => import(`@/views/${view}.vue`)
}

const routes = menu.map((item) => {
  const meta = item.meta || {}

  if (meta.requiresAuth) {
    return {
      path: item.path,
      name: item.name,
      component: MainLayout,
      meta,
      children: [
        {
          path: '',
          name: `${item.name}-page`,
          component: loadView(item.component),
          meta,
        },
      ],
    }
  }

  return {
    path: item.path,
    name: item.name,
    component: loadView(item.component),
    meta,
  }
})

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore()

  if (!to.meta.requiresAuth) return next()

  if (auth.user) return next()

  try {
    const res = await api.get('/api/auth/me')
    auth.setUser(res.data)
    next()
  } catch (err) {
    auth.clearUser()
    const redirect = to.meta.redirectIfFail || 'Login'
    next({ name: redirect })
  }
})

export default router