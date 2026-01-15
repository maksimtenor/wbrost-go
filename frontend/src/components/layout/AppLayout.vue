<template>
  <div class="wrapper hold-transition sidebar-mini">
    <!-- Navbar -->
<!--    <Navbar />-->

<!--    &lt;!&ndash; Main Sidebar Container &ndash;&gt;-->
<!--    <Sidebar />-->

    <!-- Content Wrapper -->
    <div class="content-wrapper">
      <!-- Content Header (Page header) -->
      <ContentHeader :title="pageTitle" :breadcrumbs="breadcrumbs" />

      <!-- Main content -->
      <div class="content">
        <div class="container-fluid">
          <slot />
        </div>
      </div>
    </div>

    <!-- Main Footer -->
<!--    <Footer />-->

    <!-- Control Sidebar -->
    <ControlSidebar />
  </div>
</template>

<script>
import Navbar from './Navbar.vue'
import Sidebar from './Sidebar.vue'
// import Footer from './Footer.vue'
// import ControlSidebar from './ControlSidebar.vue'
// import ContentHeader from './ContentHeader.vue'

export default {
  name: 'AppLayout',
  components: {
    Navbar,
    Sidebar,
    // Footer,
    ControlSidebar,
    ContentHeader
  },
  props: {
    pageTitle: {
      type: String,
      default: ''
    },
    breadcrumbs: {
      type: Array,
      default: () => []
    }
  }
}
</script>

<style>
/* Глобальные стили AdminLTE */
.wrapper {
  min-height: 100vh;
}

.hold-transition .content-wrapper,
.hold-transition .main-sidebar,
.hold-transition .main-header {
  transition: none !important;
}
</style>