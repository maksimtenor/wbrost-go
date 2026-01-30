import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
const BASE_API_URL = process.env.VITE_BASE_API_URL // базовый урл для vite конфига (с 8080 портом)
const BASE_FRONTEND_PORT = process.env.VITE_BASE_FRONTEND_PORT // базовый урл для vite конфига (с 8080 портом)

export default defineConfig({
    plugins: [vue()],
    resolve: {
        alias: {
            '@': resolve(__dirname, 'src')
        }
    },
    server: {
        port: parseInt(BASE_FRONTEND_PORT),
        proxy: {
            '/api': {
                target: BASE_API_URL,
                changeOrigin: true
            }
        }
    },
    build: {
        outDir: 'dist',
        emptyOutDir: true
    }
})