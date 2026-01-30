<!-- components/layout/BaseLayout.vue -->
<template>
  <div class="wrapper hold-transition sidebar-mini">
    <Navbar />
    <Sidebar />

    <div class="content-wrapper">
      <div class="content-header">
        <slot name="header">
          <!-- Заголовок по умолчанию -->
          <div class="container-fluid">
            <div class="row mb-2">
              <div class="col-sm-6">
                <h1 class="page-title">
                  <slot name="title-icon">
                    <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                      <path d="M12 2L2 7l10 5 10-5-10-5z"></path>
                    </svg>
                  </slot>
                  <slot name="title">Заголовок</slot>
                </h1>
              </div>
              <div class="col-sm-6 text-right">
                <slot name="header-right"></slot>
              </div>
            </div>
          </div>
        </slot>
      </div>

      <div class="content">
        <slot>
          <!-- Сюда будет подставляться контент страницы -->
          <p>Контент страницы</p>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import Navbar from "./Navbar.vue";
import Sidebar from "./Sidebar.vue";
</script>