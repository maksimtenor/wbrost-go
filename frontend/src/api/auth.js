import apiClient from './client'

export const authAPI = {
    async login(credentials) {
        try {
            const response = await apiClient.post('/auth/login', credentials)
            return {
                success: true,
                data: response.data
            }
        } catch (error) {
            return {
                success: false,
                error: error.response?.data?.error || error.message || 'Login failed',
                validationErrors: error.response?.data?.errors,
                status: error.response?.status
            }
        }
    },

    async signup(userData) {
        try {
            const response = await apiClient.post('/auth/signup', userData)
            return {
                success: true,
                data: response.data
            }
        } catch (error) {
            return {
                success: false,
                error: error.response?.data?.error || error.message || 'Registration failed',
                validationErrors: error.response?.data?.errors,
                status: error.response?.status
            }
        }
    },

    logout() {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        delete apiClient.defaults.headers.common['Authorization']
    },

    setAuthData(token, user) {
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
        apiClient.defaults.headers.common['Authorization'] = `Bearer ${token}`
    },

    isAuthenticated() {
        return !!localStorage.getItem('token')
    },

    getCurrentUser() {
        const userStr = localStorage.getItem('user')
        try {
            return userStr ? JSON.parse(userStr) : null
        } catch {
            return null
        }
    },

    // Только один метод для получения свежих данных
    async getFreshUserData() {
        try {
            const response = await apiClient.get('/auth/me')
            return {
                success: true,
                data: response.data
            }
        } catch (error) {
            return {
                success: false,
                error: error.response?.data?.error || error.message || 'Failed to get user data',
                status: error.response?.status
            }
        }
    }
}