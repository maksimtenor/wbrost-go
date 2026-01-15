import { createRouter, createWebHistory } from 'vue-router'

// Views
import Home from '@/views/site/Home.vue'
import Login from '@/views/site/Login.vue'
import Signup from '@/views/site/Signup.vue'
import Info from '@/views/site/Info.vue'
import Donation from '@/views/site/Donation.vue'
import Profile from '@/views/profile/Index.vue'
import ApiKeys from '@/views/profile/ApiKeys.vue'
import Carts from '@/views/profile/Carts.vue'
import User from '@/views/user/Index.vue'
import GetReports from '@/views/stat/GetReports.vue'
import StatDetail from '@/views/stat/StatDetail.vue'

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
    {
        path: '/profile',
        name: 'Profile',
        component: Profile,
        meta: { guestOnly: true }
    },
    {
        path: '/profile/apikeys',
        name: 'ApiKeys',
        component: ApiKeys,
        meta: { guestOnly: true }
    },
    {
        path: '/profile/carts',
        name: 'Carts',
        component: Carts,
        meta: { guestOnly: true }
    },
    {
        path: '/user',
        name: 'User',
        component: User,
        meta: { guestOnly: true }
    },
    {
        path: '/stat/detail',
        name: 'StatDetail',
        component: StatDetail,
        meta: { guestOnly: true }
    },
    {
        path: '/stat/get-reports',
        name: 'GetReports',
        component: GetReports,
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