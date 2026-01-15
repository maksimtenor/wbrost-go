<template>
  <div class="wrapper">
    <Navbar />
    <Sidebar />

    <div class="content-wrapper">
      <router-view />
    </div>
  </div>
</template>

<script>
import Navbar from './Navbar.vue'
import Sidebar from './Sidebar.vue'

export default {
  components: { Navbar, Sidebar }
}
</script>
