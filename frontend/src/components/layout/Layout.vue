<template>
  <div class="wrapper">
<!--    <Navbar />-->
<!--    <Sidebar />-->

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
