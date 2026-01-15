import { createStore } from 'vuex'
import axios from 'axios'

export default createStore({
    state: {
        token: localStorage.getItem('token') || null,
        user: JSON.parse(localStorage.getItem('user')) || null,
        isAuthenticated: !!localStorage.getItem('token')
    },
    getters: {
        isAuthenticated: state => state.isAuthenticated,
        currentUser: state => state.user,
        isAdmin: state => state.user ? state.user.admin === 1 : false,
        proAccount: state => state.user ? state.user.pro_account : '',
        usersCount: state => 0, // Заглушка
        usersRealCount: state => 0 // Заглушка
    },
    mutations: {
        SET_TOKEN(state, token) {
            state.token = token
            state.isAuthenticated = !!token
            if (token) {
                localStorage.setItem('token', token)
                axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
            } else {
                localStorage.removeItem('token')
                delete axios.defaults.headers.common['Authorization']
            }
        },
        SET_USER(state, user) {
            state.user = user
            if (user) {
                localStorage.setItem('user', JSON.stringify(user))
            } else {
                localStorage.removeItem('user')
            }
        }
    },
    actions: {
        async login({ commit }, credentials) {
            const response = await axios.post('/api/login', credentials)
            commit('SET_TOKEN', response.data.token)
            commit('SET_USER', response.data.user)
        },
        async logout({ commit }) {
            await axios.post('/api/logout')
            commit('SET_TOKEN', null)
            commit('SET_USER', null)
        },
        async checkAuth({ commit }) {
            try {
                const response = await axios.get('/api/user')
                commit('SET_USER', response.data)
            } catch (error) {
                commit('SET_TOKEN', null)
                commit('SET_USER', null)
            }
        }
    }
})