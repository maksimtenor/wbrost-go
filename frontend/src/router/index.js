import { createRouter, createWebHistory } from 'vue-router'

// Views
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import Signup from '@/views/Signup.vue'
import Info from '@/views/Info.vue'
import Donation from '@/views/Donation.vue'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
        meta: { requiresAuth: true }
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: { guestOnly: true }
    },
    {
        path: '/signup',
        name: 'Signup',
        component: Signup,
        meta: { guestOnly: true }
    },
    {
        path: '/info',
        name: 'Info',
        component: Info,
        meta: { guestOnly: true }
    },
    {
        path: '/donation',
        name: 'Donation',
        component: Donation,
        meta: { guestOnly: true }
    },
    // Можно добавить catch-all маршрут в конце
    {
        path: '/:pathMatch(.*)*',
        redirect: '/'
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Навигационные guards
router.beforeEach((to, from, next) => {
    // Временно отключаем проверку auth для тестирования
    // TODO: Добавить реальную проверку когда будет store
    next()
})

export default router