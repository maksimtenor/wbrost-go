<template>
  <!-- Navbar -->
  <nav class="main-header navbar navbar-expand navbar-white navbar-light">
    <!-- Left navbar links -->
    <ul class="navbar-nav">
<!--      <li class="nav-item">-->
<!--        <a class="nav-link" data-widget="pushmenu" href="#" role="button">-->
<!--          <i class="fas fa-bars"></i>-->
<!--        </a>-->
<!--      </li>-->
      <li class="nav-item d-sm-inline-block">
        <router-link to="/" class="nav-link">Dashboard</router-link>
      </li>
      <li class="nav-item d-sm-inline-block">
        <router-link to="/info" class="nav-link">Инструкция</router-link>
      </li>
      <li class="nav-item d-sm-inline-block">
        <router-link to="/donation" class="nav-link">Помощь</router-link>
      </li>
    </ul>

    <!-- Right navbar links -->
    <ul class="navbar-nav ml-auto">
      <li class="nav-item">
        <a href="#" @click.prevent="handleAuthAction" class="nav-link">
          <i :class="authIcon"></i>
        </a>
      </li>
    </ul>
  </nav>
  <!-- /.navbar -->
</template>

<script>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

export default {
  name: 'Navbar',
  setup() {
    const router = useRouter()
    const store = useStore()

    const isAuthenticated = computed(() => store.getters.isAuthenticated)

    const authIcon = computed(() => {
      return isAuthenticated.value ? 'fas fa-sign-out-alt' : 'fas fa-sign-in-alt'
    })

    const authActionLink = computed(() => {
      return isAuthenticated.value ? '/logout' : '/login'
    })

    const handleAuthAction = async () => {
      if (isAuthenticated.value) {
        await store.dispatch('logout')
        router.push('/login')
      } else {
        router.push('/login')
      }
    }

    const openInfo = () => {
      // Логика для инструкции
      console.log('Open info')
    }

    const openDonation = () => {
      // Логика для помощи
      console.log('Open donation')
    }

    return {
      isAuthenticated,
      authIcon,
      handleAuthAction,
      openInfo,
      openDonation
    }
  }
}
</script>