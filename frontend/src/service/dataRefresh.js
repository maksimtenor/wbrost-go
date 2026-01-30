import { authAPI } from '@/api/auth'
import store from '@/store'
const REFRESH_INTERVAL = import.meta.env.REFRESH_INTERVAL

class DataRefreshService {
    constructor() {
        this.intervalId = null
        this.refreshInterval = REFRESH_INTERVAL // 30 секунд
    }

    start() {
        if (this.intervalId) return

        // Первое обновление сразу
        this.refreshUserData()

        // Затем по интервалу
        this.intervalId = setInterval(() => {
            this.refreshUserData()
        }, this.refreshInterval)
    }

    stop() {
        if (this.intervalId) {
            clearInterval(this.intervalId)
            this.intervalId = null
        }
    }

    async refreshUserData() {
        if (!store.state.isAuthenticated) return

        try {
            const result = await authAPI.getFreshUserData()
            if (result.success) {
                store.commit('SET_USER', result.data)
                console.log('User data refreshed automatically')
            } else if (result.status === 401) {
                // Токен истек
                this.stop()
                store.dispatch('logout')
            }
        } catch (error) {
            console.error('Auto-refresh error:', error)
        }
    }
}

// Создаем глобальный экземпляр
export const dataRefreshService = new DataRefreshService()