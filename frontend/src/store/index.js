import { createStore } from 'vuex'
import { authAPI } from '@/api/auth'

export default createStore({
    state: {
        token: localStorage.getItem('token') || null,
        user: null, // Начинаем с null, загрузим при инициализации
        isAuthenticated: false // Начинаем с false
    },
    mutations: {
        SET_TOKEN(state, token) {
            state.token = token
            state.isAuthenticated = !!token
            if (token) {
                localStorage.setItem('token', token)
            } else {
                localStorage.removeItem('token')
            }
        },
        SET_USER(state, user) {
            state.user = user
            if (user) {
                localStorage.setItem('user', JSON.stringify(user))
            } else {
                localStorage.removeItem('user')
            }
        },
        LOGOUT(state) {
            state.token = null
            state.user = null
            state.isAuthenticated = false
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            localStorage.removeItem('rememberMe')
            localStorage.removeItem('savedUsername')
        }
    },
    actions: {
        login({ commit }, { token, user }) {
            commit('SET_TOKEN', token)
            commit('SET_USER', user)
        },
        logout({ commit }) {
            commit('LOGOUT')
            authAPI.logout()
        },
        // Загружает свежие данные пользователя с сервера
        async loadUserData({ commit, state }) {
            if (!state.token) return null

            try {
                const result = await authAPI.getFreshUserData()
                if (result.success) {
                    commit('SET_USER', result.data)
                    return result.data
                }
                // Если токен невалидный
                if (result.status === 401) {
                    commit('LOGOUT')
                }
                return null
            } catch (error) {
                console.error('Failed to load user data:', error)
                return null
            }
        }
    },
    getters: {
        currentUser: state => state.user,
        isAuthenticated: state => state.isAuthenticated,
        isAdmin: state => state.user?.admin === 1
    }
})