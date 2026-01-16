<!--<template>-->
<!--  <div class="wrapper hold-transition sidebar-mini">-->
<!--    &lt;!&ndash; Navbar &ndash;&gt;-->
<!--&lt;!&ndash;    <Navbar />&ndash;&gt;-->

<!--&lt;!&ndash;    &lt;!&ndash; Main Sidebar Container &ndash;&gt;&ndash;&gt;-->
<!--&lt;!&ndash;    <Sidebar />&ndash;&gt;-->

<!--    &lt;!&ndash; Content Wrapper &ndash;&gt;-->
<!--    <div class="content-wrapper">-->
<!--      &lt;!&ndash; Content Header (Page header) &ndash;&gt;-->
<!--      <ContentHeader :title="pageTitle" :breadcrumbs="breadcrumbs" />-->

<!--      &lt;!&ndash; Main content &ndash;&gt;-->
<!--      <div class="content">-->
<!--        <div class="container-fluid">-->
<!--          <slot />-->
<!--        </div>-->
<!--      </div>-->
<!--    </div>-->

<!--    &lt;!&ndash; Main Footer &ndash;&gt;-->
<!--&lt;!&ndash;    <Footer />&ndash;&gt;-->

<!--    &lt;!&ndash; Control Sidebar &ndash;&gt;-->
<!--    <ControlSidebar />-->
<!--  </div>-->
<!--</template>-->

<!--<script>-->
<!--import Navbar from './Navbar.vue'-->
<!--import Sidebar from './Sidebar.vue'-->
<!--// import Footer from './Footer.vue'-->
<!--// import ControlSidebar from './ControlSidebar.vue'-->
<!--// import ContentHeader from './ContentHeader.vue'-->

<!--export default {-->
<!--  name: 'AppLayout',-->
<!--  components: {-->
<!--    Navbar,-->
<!--    Sidebar,-->
<!--    // Footer,-->
<!--    ControlSidebar,-->
<!--    ContentHeader-->
<!--  },-->
<!--  props: {-->
<!--    pageTitle: {-->
<!--      type: String,-->
<!--      default: ''-->
<!--    },-->
<!--    breadcrumbs: {-->
<!--      type: Array,-->
<!--      default: () => []-->
<!--    }-->
<!--  }-->
<!--}-->
<!--</script>-->

<!--<style>-->
<!--/* Глобальные стили AdminLTE */-->
<!--.wrapper {-->
<!--  min-height: 100vh;-->
<!--}-->

<!--.hold-transition .content-wrapper,-->
<!--.hold-transition .main-sidebar,-->
<!--.hold-transition .main-header {-->
<!--  transition: none !important;-->
<!--}-->
<!--</style>-->