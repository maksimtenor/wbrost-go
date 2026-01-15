<template>
  <aside class="main-sidebar sidebar-dark-primary elevation-4 bg-black"
         style="background-image: url(/img/1.png); background-position: inherit">

    <div class="sidebar">
      <!-- Логотип -->
      <div style="display: block; padding-left: .8rem; text-align: center; margin-top: 15px">
        <router-link to="/">
          <img src="/img/wbrost-logo_2.png" alt="logo" height="32">
        </router-link>
      </div>

      <!-- User Panel -->
      <div class="user-panel mt-3 pb-3 mb-3 d-flex">
        <div class="image">
          <img src="/img/emptyAvatarSmall.jpg" class="img-circle elevation-2" alt="User Image">
        </div>
        <div class="info">
          <span class="d-block">
            {{ userName }}
            <span class="right badge badge-secondary">{{ user.typeAccount || 'гость' }}</span>
          </span>
        </div>
      </div>

      <!-- Telegram Links -->
      <div class="user-panel mt-3 pb-3 mb-3"
           style="display: flex; text-align: left; border-bottom: 1px solid #4f5962; flex-direction: column; flex-wrap: nowrap; padding-left: 10%">
        <p>Полезное в телеграм:</p>
        <p>
          <span style="color: #299fdb; margin-right: 10px;" class="fa fa-send"></span>
          <a href="https://t.me/marevichh" target="_blank" style="text-decoration: none">Техподдержка</a>
        </p>
        <p>
          <span style="color: #299fdb; margin-right: 10px;" class="fa fa-send"></span>
          <a href="https://t.me/wb_supply_helper_bot" target="_blank" style="text-decoration: none">Бот: Проверка поставок</a>
        </p>
        <p>
          <span style="color: #299fdb; margin-right: 10px;" class="fa fa-send"></span>
          <a href="https://t.me/wb_halyavka" target="_blank" style="text-decoration: none">Канал: скидки WB</a>
        </p>
      </div>

      <!-- Sidebar Menu -->
      <nav class="mt-2">
        <ul class="nav nav-pills nav-sidebar flex-column" data-widget="treeview" role="menu" data-accordion="false">

          <!-- Меню для авторизованных -->
          <template v-if="user">
            <!-- Статистика -->
            <li class="nav-item has-treeview">
              <a href="#" class="nav-link">
                <i class="nav-icon fas fa-chart-bar"></i>
                <p>
                  Статистика
                  <i class="right fas fa-angle-left"></i>
                </p>
              </a>
              <ul class="nav nav-treeview">
                <li class="nav-item">
                  <router-link to="/reports/stat-detail" class="nav-link">
                    <i class="far fa-chart-pie nav-icon"></i>
                    <p>Детальная статистика</p>
                  </router-link>
                </li>
                <li class="nav-item">
                  <router-link to="/reports" class="nav-link">
                    <i class="fas fa-reply nav-icon"></i>
                    <p>Запросить отчет</p>
                  </router-link>
                </li>
              </ul>
            </li>

            <!-- Товары -->
            <li class="nav-item">
              <router-link to="/reports/artlist" class="nav-link">
                <i class="nav-icon fas fa-list"></i>
                <p>Товары</p>
              </router-link>
            </li>

            <!-- Профиль -->
            <li class="nav-item has-treeview">
              <a href="#" class="nav-link">
                <i class="nav-icon fas fa-user-cog"></i>
                <p>
                  Профиль
                  <i class="right fas fa-angle-left"></i>
                </p>
              </a>
              <ul class="nav nav-treeview">
                <li class="nav-item">
                  <router-link to="/profile" class="nav-link">
                    <i class="fas fa-edit nav-icon"></i>
                    <p>Личный кабинет</p>
                  </router-link>
                </li>
                <li class="nav-item">
                  <router-link to="/profile/mykeys" class="nav-link">
                    <i class="fas fa-key nav-icon"></i>
                    <p>Мои API ключи</p>
                  </router-link>
                </li>
              </ul>
            </li>

            <!-- Тариф -->
            <li class="nav-item">
              <a href="#" class="nav-link">
                <i class="nav-icon fas fa-dashboard"></i>
                <p>
                  Тариф
                  <span class="right badge badge-danger">{{ user.proAccount || '' }}</span>
                </p>
              </a>
            </li>
          </template>

          <!-- Меню для гостей -->
          <template v-else>
            <li class="nav-item">
              <router-link to="/login" class="nav-link">
                <i class="nav-icon fas fa-sign-in-alt"></i>
                <p>Вход</p>
              </router-link>
            </li>
            <li class="nav-item">
              <router-link to="/signup" class="nav-link">
                <i class="nav-icon fas fa-"></i>
                <p>Регистрация</p>
              </router-link>
            </li>
          </template>

        </ul>
      </nav>
    </div>
  </aside>
</template>

<script>
export default {
  name: 'Sidebar',
  computed: {
    user() {
      return this.$store.state.user ?? '' // Или откуда берётся пользователь
    },
    userName() {
      if (!this.user) return ''
      return this.user.name || this.user.username || ''
    }
  }
}
</script>