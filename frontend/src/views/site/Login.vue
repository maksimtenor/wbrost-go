<template>
  <div class="hold-transition login-page">
    <div class="login-box">
      <div class="login-logo">
        <router-link to="/">
          <img src="/img/wbrost-logo_2.png" alt="logo" height="32">
        </router-link>
      </div>

      <div class="card">
        <div class="card-body login-card-body">
          <p class="login-box-msg">Авторизуйтесь для продолжения</p>

          <form @submit.prevent="handleLogin">
            <div class="form-group has-feedback" :class="{ 'has-error': errors.username }">
              <input
                  type="text"
                  class="form-control"
                  placeholder="Имя пользователя"
                  v-model="form.username"
              />
              <span class="fa fa-user form-control-feedback"></span>
              <span class="help-block-error" v-if="errors.username">{{ errors.username }}</span>
            </div>

            <div class="form-group has-feedback" :class="{ 'has-error': errors.password }">
              <input
                  type="password"
                  class="form-control"
                  placeholder="Пароль"
                  v-model="form.password"
              />
              <span class="fa fa-lock form-control-feedback"></span>
              <span class="help-block-error" v-if="errors.password">{{ errors.password }}</span>
            </div>

            <div class="row">
              <div class="col-8">
                <div class="icheck-primary">
                  <input type="checkbox" id="remember" v-model="form.rememberMe">
                  <label for="remember">Запомнить меня</label>
                </div>
              </div>
              <div class="col-4">
                <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
                  <span v-if="loading">Загрузка...</span>
                  <span v-else>Войти</span>
                </button>
              </div>
              <div class="col-12 mt-3">
                <router-link to="/signup">Регистрация</router-link>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const store = useStore()

    const form = ref({
      username: '',
      password: '',
      rememberMe: false
    })

    const errors = ref({})
    const loading = ref(false)

    const handleLogin = async () => {
      errors.value = {}
      loading.value = true

      try {
        // Временная заглушка для тестирования
        store.commit('SET_TOKEN', 'test-token')
        store.commit('SET_USER', {
          name: 'Тестовый k',
          username: 'test',
          pro_account: 'PRO',
          admin: 1
        })
        router.push('/')
      } catch (error) {
        errors.value = { general: 'Ошибка авторизации' }
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      errors,
      loading,
      handleLogin
    }
  }
}
</script>