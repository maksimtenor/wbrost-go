<template>
  <!-- Всегда показываем Layout для простоты -->
  <Layout>
    <router-view />
  </Layout>
</template>

<script>
import Layout from './components/layout/Layout.vue'

export default {
  name: 'App',
  components: {
    Layout
  }
}
</script>

<style>
/* Глобальные стили */
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}
</style>