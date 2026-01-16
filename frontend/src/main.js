import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import { dataRefreshService } from '@/service/dataRefresh'

// Импорты стилей остаются без изменений
import '@/assets/css/v4-shims.css'
import '@fortawesome/fontawesome-free/css/all.css'
import '@/assets/css/landing/style.css'
import '@/assets/css/landing/font/inter.css'
import '@/assets/css/landing/main.min.css'
import '@/assets/css/landing/main.min.rtl.css'
import '@/assets/css/site.css'
import '@/assets/css/adminlte.css'
import '@/assets/css/adminlte.min.css'

// Создаем приложение
const app = createApp(App)

// Загружаем данные пользователя при старте
async function initApp() {
    const token = localStorage.getItem('token')

    if (token) {
        try {
            // Устанавливаем токен в store
            store.commit('SET_TOKEN', token)

            // Загружаем свежие данные с сервера
            await store.dispatch('loadUserData')

            // Запускаем автоматическое обновление
            dataRefreshService.start()

            console.log('App initialized with user data')
        } catch (error) {
            console.error('Init error:', error)
        }
    }
}

// Инициализируем и монтируем
initApp().then(() => {
    app.use(store)
    app.use(router)
    app.mount('#app')
})

// Останавливаем обновление при выходе
router.beforeEach((to) => {
    if (to.path === '/logout') {
        dataRefreshService.stop()
    }
})