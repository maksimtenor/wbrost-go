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

<style scoped>
/* Основные стили */
.main-sidebar {
  min-height: 100vh;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 1038;
}

.sidebar {
  height: 100%;
  overflow-y: auto;
  padding-bottom: 20px;
}

.user-panel {
  padding: 10px;
  border-bottom: 1px solid #4f5962;
}

.user-panel .image {
  padding-right: 10px;
}

.user-panel .image img {
  width: 40px;
  height: 40px;
  object-fit: cover;
}

.user-panel .info {
  color: #fff;
  flex: 1;
}

.user-panel .info .d-block {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  font-weight: 500;
}

/* Стили для бейджей типов аккаунтов */
.account-type-badge {
  font-size: 8px !important;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-left: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 40px;
  height: 16px;
  transition: all 0.2s ease;
}

.account-type-guest {
  background: linear-gradient(135deg, #6c757d, #495057) !important;
  color: #fff !important;
  border: 1px solid #5a6268;
}

.account-type-owner {
  background: linear-gradient(135deg, #0dcaf0, #0aa2c0) !important; /* Tiffany color */
  color: #fff !important;
  border: 1px solid #0ba5cc;
  box-shadow: 0 2px 4px rgba(13, 202, 240, 0.3);
}

.account-type-admin {
  background: linear-gradient(135deg, #dc3545, #c82333) !important;
  color: #fff !important;
  border: 1px solid #bd2130;
  box-shadow: 0 2px 4px rgba(220, 53, 69, 0.3);
}

.account-type-user {
  background: linear-gradient(135deg, #20c997, #17a589) !important;
  color: #fff !important;
  border: 1px solid #14967c;
}

/* Стили для иконок телеграма */
.telegram-links {
  display: flex;
  flex-direction: column;
  padding-left: 10%;
  border-bottom: 1px solid #4f5962;
}

.telegram-links p {
  margin-bottom: 8px;
  color: #adb5bd;
  font-size: 13px;
}

.telegram-links p:first-child {
  color: #fff;
  font-weight: 500;
  margin-bottom: 12px;
}

.telegram-links a {
  color: #299fdb;
  text-decoration: none;
  transition: color 0.2s;
}

.telegram-links a:hover {
  color: #1c8bc8;
  text-decoration: underline;
}

.telegram-icon {
  color: #299fdb;
  margin-right: 10px;
  width: 16px;
  text-align: center;
}

/* Стили для бейджа тарифа */
.tariff-badge {
  font-size: 10px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  min-width: 60px;
  height: 20px;
  transition: all 0.3s ease;
  border: 1px solid transparent;
}

.tariff-badge i {
  font-size: 8px;
}

.tariff-badge.badge-success {
  background: linear-gradient(135deg, #28a745, #1e7e34) !important;
  color: #fff !important;
  border-color: #1c7430;
  box-shadow: 0 2px 4px rgba(40, 167, 69, 0.3);
}

.tariff-badge.badge-secondary {
  background: linear-gradient(135deg, #6c757d, #545b62) !important;
  color: #fff !important;
  border-color: #4e555b;
}

.tariff-item .nav-icon {
  color: #ffc107; /* Золотой цвет для короны/тарифа */
}

/* Анимации для бейджей */
.account-type-badge:hover,
.tariff-badge:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

/* Стили для кастомного treeview если AdminLTE не работает */
.nav-sidebar .nav-treeview {
  display: none;
}

.nav-sidebar .has-treeview.menu-open > .nav-treeview {
  display: block;
  animation: slideDown 0.3s ease-out;
}

/* Стили для стрелочек */
.nav-sidebar .has-treeview > a > .fa-angle-left {
  transition: transform 0.3s;
}

.nav-sidebar .has-treeview.menu-open > a > .fa-angle-left {
  transform: rotate(-90deg);
}

/* Стили для активного состояния */
.nav-sidebar .nav-link.active {
  background-color: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.nav-sidebar .nav-treeview .nav-link {
  padding-left: 2rem;
}

/* Улучшение иконок в меню */
.nav-icon {
  width: 20px;
  text-align: center;
  margin-right: 10px;
  transition: color 0.2s;
}

.nav-link:hover .nav-icon {
  //color: #3c8dbc;
}

.nav-link.active .nav-icon {
  color: #fff;
}

/* Анимация выезжания */
@keyframes slideDown {
  from {
    opacity: 0;
    max-height: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    max-height: 500px;
    transform: translateY(0);
  }
}

/* Улучшаем внешний вид для мобильных устройств */
@media (max-width: 767.98px) {
  .nav-sidebar .nav-treeview .nav-link {
    padding-left: 1.5rem;
  }

  .account-type-badge {
    font-size: 7px !important;
    padding: 1px 4px;
    min-width: 35px;
    height: 14px;
  }

  .tariff-badge {
    font-size: 9px;
    padding: 3px 6px;
    min-width: 50px;
    height: 18px;
  }
}

/* Эффект свечения для PRO аккаунта */
@keyframes proGlow {
  0%, 100% {
    box-shadow: 0 0 5px rgba(40, 167, 69, 0.5);
  }
  50% {
    box-shadow: 0 0 10px rgba(40, 167, 69, 0.8);
  }
}

.tariff-badge.badge-success {
  animation: proGlow 2s infinite;
}

/* Эффект свечения для владельца */
@keyframes ownerGlow {
  0%, 100% {
    box-shadow: 0 0 5px rgba(13, 202, 240, 0.5);
  }
  50% {
    box-shadow: 0 0 10px rgba(13, 202, 240, 0.8);
  }
}

.account-type-owner {
  animation: ownerGlow 2s infinite;
}

/* Эффект пульсации для админа */
@keyframes adminPulse {
  0%, 100% {
    box-shadow: 0 0 5px rgba(220, 53, 69, 0.5);
  }
  50% {
    box-shadow: 0 0 12px rgba(220, 53, 69, 0.9);
  }
}

.account-type-admin {
  animation: adminPulse 1.5s infinite;
}
</style>