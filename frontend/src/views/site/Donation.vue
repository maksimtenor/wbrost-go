<script>
import Navbar from "../../components/layout/Navbar.vue";
import Sidebar from "../../components/layout/Sidebar.vue";

export default {
  name: "Donation",
  components: {Sidebar, Navbar}
}
</script>

<template>
      <Navbar />
      <Sidebar />
  <div class="content-wrapper">
  <div class="content-header">
    <div class="container-fluid">
      <div class="row mb-2">
        <div class="col-sm-6">
          <h1 class="m-0">
            Помощь проекту                    </h1>
        </div><!-- /.col -->
        <div class="col-sm-6">
        </div><!-- /.col -->
      </div><!-- /.row -->
    </div><!-- /.container-fluid -->
  </div>
  <div class="content">
    <h5>Ваши пожертвования позволяют сайту развиваться, спасибо!</h5>
    <iframe src="https://yoomoney.ru/quickpay/fundraise/button?billNumber=16I97N3TS3F.241116&" width="330" height="50" frameborder="0" allowtransparency="true" scrolling="no"></iframe>
  </div>
  </div>
</template>

<style scoped>
/* Пустой - все стили в site.css */
</style>