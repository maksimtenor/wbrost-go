<script setup>
</script>

<template>
  <div class="basic-form-container">
    <!-- Сообщения об ошибках/успехе -->
    <div v-if="$slots.messages" class="messages-section">
      <slot name="messages"></slot>
    </div>

    <!-- Основное содержимое формы -->
    <div class="form-content">
      <slot></slot>
    </div>

    <!-- Футер формы (опционально) -->
<!--    <div v-if="$slots.footer" class="form-footer">-->
<!--      <slot name="footer"></slot>-->
<!--    </div>-->
  </div>
</template>

<style scoped>
.basic-form-container {
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.messages-section {
  margin-bottom: 20px;
}

.form-content {
  padding: 10px 0;
}

.form-footer {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}
</style>