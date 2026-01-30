<template>
  <div class="wrapper">
    <router-view />
  </div>
</template>

<script>
import Navbar from './Navbar.vue'
import Sidebar from './Sidebar.vue'

export default {
  components: { Navbar, Sidebar }
}
</script>
