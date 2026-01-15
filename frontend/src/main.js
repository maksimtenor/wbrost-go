import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

// Импортируем стили
import '@/assets/css/bootstrap-icons.css'
import '@/assets/css/v4-shims.css'
// import '@/assets/css/site.css'
import '@/assets/css/landing/style.css'
import '@/assets/css/landing/font/inter.css'
import '@/assets/css/landing/main.min.css'
import '@/assets/css/landing/main.min.rtl.css'
import '@/assets/css/site.css'

// Подключаем AdminLTE CSS (если есть отдельный файл)
import '@/assets/css/adminlte.min.css'
const app = createApp(App)

app.use(store)
app.use(router)

app.mount('#app')