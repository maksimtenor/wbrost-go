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
            <span class="account-type-badge" :class="userTypeClass">
              {{ userTypeAccount }}
            </span>
          </span>
        </div>
      </div>

      <!-- Telegram Links -->
      <div class="user-panel mt-3 pb-3 mb-3 telegram-links">
        <p>Полезное в телеграм:</p>
        <p>
          <span class="telegram-icon fa fa-send"></span>
          <a href="https://t.me/marevichh" target="_blank">Техподдержка</a>
        </p>
        <p>
          <span class="telegram-icon fa fa-send"></span>
          <a href="https://t.me/wb_supply_helper_bot" target="_blank">Бот: Проверка поставок</a>
        </p>
        <p>
          <span class="telegram-icon fa fa-send"></span>
          <a href="https://t.me/wb_halyavka" target="_blank">Канал: скидки WB</a>
        </p>
      </div>

      <!-- Sidebar Menu -->
      <nav class="mt-2">
        <ul class="nav nav-pills nav-sidebar flex-column" data-widget="treeview" role="menu" data-accordion="false">

          <!-- Меню для авторизованных -->
          <template v-if="user">
            <!-- Статистика -->
            <li class="nav-item has-treeview" :class="{ 'menu-open': activeMenu === 'statistics' }">
              <a href="#" class="nav-link" @click.prevent="toggleMenu('statistics')">
                <i class="nav-icon fas fa-chart-bar"></i>
                <p>
                  Статистика
                  <i class="right fas fa-angle-left"></i>
                </p>
              </a>
              <ul class="nav nav-treeview">
                <li class="nav-item">
                  <router-link to="/stat/detail" class="nav-link" @click="closeAllMenus">
                    <i class="fas fa-chart-pie nav-icon"></i>
                    <p>Детальная статистика</p>
                  </router-link>
                </li>
                <li class="nav-item">
                  <router-link to="/stat/get-reports" class="nav-link" @click="closeAllMenus">
                    <i class="fas fa-reply nav-icon"></i>
                    <p>Запросить отчет</p>
                  </router-link>
                </li>
              </ul>
            </li>

            <!-- Товары -->
            <li class="nav-item">
              <router-link to="/profile/carts" class="nav-link">
                <i class="nav-icon fas fa-list"></i>
                <p>Товары</p>
              </router-link>
            </li>
            <!-- Пользователи -->
            <li class="nav-item" v-if="user.admin">
              <router-link to="/users" class="nav-link">
                <i class="nav-icon fas fa-atom"></i>
                <p>Пользователи</p>
              </router-link>
            </li>
            <!-- Профиль -->
            <li class="nav-item has-treeview" :class="{ 'menu-open': activeMenu === 'profile' }">
              <a href="#" class="nav-link" @click.prevent="toggleMenu('profile')">
                <i class="nav-icon fas fa-user-cog"></i>
                <p>
                  Профиль
                  <i class="right fas fa-angle-left"></i>
                </p>
              </a>
              <ul class="nav nav-treeview">
                <li class="nav-item">
                  <router-link to="/profile" class="nav-link" @click="closeAllMenus">
                    <i class="fas fa-edit nav-icon"></i>
                    <p>Личный кабинет</p>
                  </router-link>
                </li>
                <li class="nav-item">
                  <router-link to="/profile/ApiKeys" class="nav-link" @click="closeAllMenus">
                    <i class="fas fa-key nav-icon"></i>
                    <p>Мои API ключи</p>
                  </router-link>
                </li>
              </ul>
            </li>
            <li class="nav-header">Прочее</li>
            <!-- Тариф -->
            <li class="nav-item">
              <a href="#" class="nav-link tariff-item">
                <i class="nav-icon fas fa-crown"></i>
                <p>
                  Тариф
                  <span class="tariff-badge" :class="proBadgeClass">
                    <i class="fas" :class="proIconClass"></i>
                    {{ proAccountText }}
                  </span>
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
                <i class="nav-icon fas fa-user-plus"></i>
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
import { mapState } from 'vuex'

export default {
  name: 'Sidebar',

  data() {
    return {
      activeMenu: null
    }
  },

  computed: {
    // Используем mapState для реактивной связи с store
    ...mapState({
      user: state => state.user,
      isAuthenticated: state => state.isAuthenticated
    }),

    userName() {
      if (!this.user) return ''
      return this.user.name || this.user.username || ''
    },

    proAccountText() {
      if (!this.user || this.user.pro === undefined) return 'trial'
      return this.user.pro === 1 ? 'PRO' : 'FREE'
    },

    proBadgeClass() {
      if (!this.user || this.user.pro === undefined) return 'badge-secondary'
      return this.user.pro === 1 ? 'badge-success' : 'badge-secondary'
    },

    proIconClass() {
      if (!this.user || this.user.pro === undefined) return 'fa-user'
      return this.user.pro === 1 ? 'fa-crown' : 'fa-user'
    },

    userTypeAccount() {
      if (!this.user) return 'гость'

      if (this.user.admin === 2) {
        return 'владелец'
      } else if (this.user.admin === 1) {
        return 'админ'
      }
      return 'пользователь'
    },

    userTypeClass() {
      if (!this.user) return 'account-type-guest'
      if (this.user.admin === 2) {
        return 'account-type-owner'
      } else if (this.user.admin === 1) {
        return 'account-type-admin'
      }
      return 'account-type-user'
    }
  },

  methods: {
    toggleMenu(menuName) {
      if (this.activeMenu === menuName) {
        this.activeMenu = null
      } else {
        this.activeMenu = menuName
      }

      // Инициализируем treeview после изменения состояния
      this.$nextTick(() => {
        this.initTreeview()
      })
    },

    closeAllMenus() {
      this.activeMenu = null
    },

    initTreeview() {
      // Инициализируем treeview плагин AdminLTE
      this.$nextTick(() => {
        if (typeof window !== 'undefined' && window.$ && window.$.fn && window.$.fn.treeview) {
          try {
            // Инициализируем все treeview элементы
            $('[data-widget="treeview"]').each(function() {
              const treeview = $(this)
              // Проверяем, инициализирован ли уже
              if (!treeview.data('lte.treeview')) {
                treeview.Treeview('init')
              }
            })
          } catch (error) {
            console.log('Treeview инициализация:', error)
          }
        }
      })
    }
  },

  mounted() {
    // Инициализируем treeview при монтировании
    this.initTreeview()

    // Если пользователь авторизован, обновляем данные при монтировании
    if (this.isAuthenticated) {
      this.$store.dispatch('loadUserData')
    }
  },

  updated() {
    // Инициализируем treeview при обновлении компонента
    this.initTreeview()
  },

  // Следим за изменениями маршрута
  watch: {
    $route() {
      // Закрываем меню при переходе на другую страницу
      this.closeAllMenus()

      // При каждом переходе по страницам обновляем данные
      if (this.isAuthenticated) {
        // Делаем с небольшой задержкой чтобы не нагружать
        setTimeout(() => {
          this.$store.dispatch('loadUserData').catch(console.error)
        }, 1000)
      }

      // Инициализируем treeview после перехода
      this.$nextTick(() => {
        this.initTreeview()
      })
    },

    user() {
      // При изменении пользователя инициализируем treeview
      this.$nextTick(() => {
        this.initTreeview()
      })
    }
  }
}
</script>

<style>
@import '@/assets/css/components/sidebar.css';
</style>