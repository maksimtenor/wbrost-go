<script setup>
import Sidebar from "../../components/layout/Sidebar.vue";
import Navbar from "../../components/layout/Navbar.vue";
</script>

<template>
  <Navbar />
  <Sidebar />

  <div class="content-wrapper">
    <div class="content-header">
      <div class="container-fluid">
        <div class="row mb-2">
          <div class="col-sm-6">
            <h1 class="page-title">
              <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
              Как пользоваться платформой
            </h1>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="info-container">
        <!-- Карточка с подключением -->
        <div class="instruction-card">
          <div class="card-header">
            <div class="card-icon">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                <circle cx="8.5" cy="7" r="4"></circle>
                <line x1="20" y1="8" x2="20" y2="14"></line>
                <line x1="23" y1="11" x2="17" y2="11"></line>
              </svg>
            </div>
            <h2 class="card-title">Как подключить Wildberries</h2>
          </div>

          <div class="steps-container">
            <div class="step">
              <div class="step-number">1</div>
              <div class="step-content">
                <p>Зайдите в ваш аккаунт на сайте <a href="https://seller.wildberries.ru" target="_blank" class="link">seller.wildberries.ru</a></p>
              </div>
            </div>

            <div class="step">
              <div class="step-number">2</div>
              <div class="step-content">
                <p>В личном кабинете нажмите на имя профиля и выберите:</p>
                <p class="step-detail">Настройки → Доступ к API</p>
                <a href="https://seller.wildberries.ru/supplier-settings/access-to-api" target="_blank" class="step-link">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                    <polyline points="15 3 21 3 21 9"></polyline>
                    <line x1="10" y1="14" x2="21" y2="3"></line>
                  </svg>
                  Или перейдите по ссылке
                </a>
              </div>
            </div>

            <div class="step">
              <div class="step-number">3</div>
              <div class="step-content">
                <p>Выберите, с какими категориями API вы будете работать</p>
                <div class="step-note">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="12"></line>
                    <line x1="12" y1="16" x2="12.01" y2="16"></line>
                  </svg>
                  <span>Рекомендуем выбрать все категории для корректной работы</span>
                </div>
              </div>
            </div>

            <div class="step">
              <div class="step-number">4</div>
              <div class="step-content">
                <p>Нажмите <span class="badge">Создать токен</span></p>
              </div>
            </div>

            <div class="step">
              <div class="step-number">5</div>
              <div class="step-content">
                <p>Скопируйте и сохраните токен в безопасном месте</p>
                <div class="step-warning">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="12"></line>
                    <line x1="12" y1="16" x2="12.01" y2="16"></line>
                  </svg>
                  <span>Потом его нельзя будет посмотреть в личном кабинете WB Селлера</span>
                </div>
              </div>
            </div>

            <div class="step">
              <div class="step-number">6</div>
              <div class="step-content">
                <p>Добавьте токен на нашей платформе:</p>
                <p class="step-detail">Профиль → Личный кабинет → WB API KEY</p>
              </div>
            </div>

            <div class="step">
              <div class="step-number">7</div>
              <div class="step-content">
                <p>Нажмите <span class="badge-success">Сохранить</span></p>
              </div>
            </div>
          </div>

          <div class="success-message">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>Готово! Теперь вам доступна возможность загрузки статистики из Wildberries</span>
          </div>

          <div class="step-note">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
            </svg>
            <span>Проверить активацию токена: Профиль → Мои API ключи</span>
          </div>
        </div>

        <!-- Карточка с загрузкой данных -->
        <div class="instruction-card">
          <div class="card-header">
            <div class="card-icon">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="7 10 12 15 17 10"></polyline>
                <line x1="12" y1="15" x2="12" y2="3"></line>
              </svg>
            </div>
            <h2 class="card-title">Как загрузить статистику</h2>
          </div>

          <div class="steps-container">
            <div class="step">
              <div class="step-number">1</div>
              <div class="step-content">
                <p>Перейдите в раздел <span class="badge">Статистика → Запросить отчёт</span></p>
              </div>
            </div>

            <div class="step">
              <div class="step-number">2</div>
              <div class="step-content">
                <p>Выберите период дат за который хотите загрузить отчет</p>
                <div class="step-note">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="12"></line>
                    <line x1="12" y1="16" x2="12.01" y2="16"></line>
                  </svg>
                  <span>Советуем выбирать максимальный период для всей информации</span>
                </div>
              </div>
            </div>

            <div class="step">
              <div class="step-number">3</div>
              <div class="step-content">
                <p>Нажмите на кнопку <span class="badge-primary">Запросить</span></p>
              </div>
            </div>

            <div class="step">
              <div class="step-number">4</div>
              <div class="step-content">
                <p>Дождитесь завершения обработки запроса</p>
              </div>
            </div>
          </div>

          <div class="step-note warning">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>Возможна ошибка при повторной загрузке одинаковых данных</span>
          </div>

          <div class="step-detail">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>Платформа уведомит, что строки которые уже есть в базе данных - не сохранились</span>
          </div>
        </div>

        <!-- Карточка с себестоимостью -->
        <div class="instruction-card">
          <div class="card-header">
            <div class="card-icon">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="1" x2="12" y2="23"></line>
                <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
              </svg>
            </div>
            <h2 class="card-title">Как добавить себестоимость товара</h2>
          </div>

          <div class="warning-message">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>Для заполнения себестоимости и отображения фотографий и названий в детальной статистике, зайдите в Товары и подтяните их из WB</span>
          </div>

          <div class="steps-container">
            <div class="step">
              <div class="step-number">1</div>
              <div class="step-content">
                <p>Перейдите в раздел <span class="badge">Товары</span></p>
              </div>
            </div>

            <div class="step">
              <div class="step-number">2</div>
              <div class="step-content">
                <p>Нажмите кнопку <span class="btn-icon">↻</span> - это подгрузит информацию о ваших карточках</p>
              </div>
            </div>

            <div class="step">
              <div class="step-number">3</div>
              <div class="step-content">
                <p>Для каждой карточки доступны кнопки:</p>
                <div class="buttons-demo">
                  <button class="demo-btn edit">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                    </svg>
                    Редактировать
                  </button>
                  <button class="demo-btn save">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
                      <polyline points="17 21 17 13 7 13 7 21"></polyline>
                      <polyline points="7 3 7 8 15 8"></polyline>
                    </svg>
                    Сохранить
                  </button>
                </div>
              </div>
            </div>

            <div class="step">
              <div class="step-number">4</div>
              <div class="step-content">
                <p>Нажмите "Редактировать", введите себестоимость 1 единицы товара и сохраните</p>
              </div>
            </div>
          </div>

          <div class="success-message">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>Готово! Эти данные будут учитываться при расчёте статистики</span>
          </div>
        </div>

        <!-- Важное примечание -->
        <div class="important-note">
          <div class="note-header">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="16" x2="12" y2="12"></line>
              <line x1="12" y1="8" x2="12.01" y2="8"></line>
            </svg>
            <h3>Важное примечание</h3>
          </div>
          <div class="note-content">
            <p>Не забывайте добавлять налоговый % в вашем личном кабинете!</p>
            <p>Правильная настройка налогов позволяет получать точную статистику по чистой прибыли.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 28px;
  font-weight: 700;
  color: #1a202c;
  margin: 0;
}

.title-icon {
  width: 32px;
  height: 32px;
  color: #4f46e5;
}

.info-container {
  //max-width: 1000px;
  margin: 0 auto;
  padding: 0 20px;
}

.instruction-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  padding: 30px;
  margin-bottom: 30px;
  border: 1px solid #e5e7eb;
  transition: transform 0.3s ease;
}

.instruction-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.12);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 2px solid #e5e7eb;
}

.card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%);
  border-radius: 12px;
  color: white;
}

.card-icon svg {
  width: 24px;
  height: 24px;
}

.card-title {
  font-size: 22px;
  font-weight: 700;
  color: #1a202c;
  margin: 0;
}

.steps-container {
  margin-bottom: 25px;
}

.step {
  display: flex;
  gap: 20px;
  margin-bottom: 25px;
  position: relative;
}

.step:not(:last-child)::after {
  content: '';
  position: absolute;
  left: 24px;
  top: 50px;
  width: 2px;
  height: calc(100% - 20px);
  background: #e5e7eb;
}

.step-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  color: white;
  border-radius: 12px;
  font-size: 20px;
  font-weight: 700;
  flex-shrink: 0;
  z-index: 1;
}

.step-content {
  flex: 1;
  padding-top: 8px;
}

.step-content p {
  font-size: 16px;
  line-height: 1.6;
  color: #374151;
  margin: 0 0 10px 0;
}

.step-detail {
  font-size: 15px;
  color: #6b7280;
  background: #f9fafb;
  padding: 10px 15px;
  border-radius: 8px;
  margin: 10px 0;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.step-detail svg {
  width: 16px;
  height: 16px;
  color: #9ca3af;
}

.step-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #4f46e5;
  text-decoration: none;
  font-weight: 500;
  font-size: 15px;
  margin-top: 10px;
  transition: color 0.2s ease;
}

.step-link:hover {
  color: #4338ca;
}

.step-link svg {
  width: 16px;
  height: 16px;
}

.step-note {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 16px;
  background: #f0f9ff;
  border-left: 4px solid #0ea5e9;
  border-radius: 8px;
  margin: 10px 0;
}

.step-note.warning {
  background: #fffbeb;
  border-left-color: #f59e0b;
}

.step-note svg {
  width: 18px;
  height: 18px;
  color: #0ea5e9;
  flex-shrink: 0;
  margin-top: 2px;
}

.step-note.warning svg {
  color: #f59e0b;
}

.step-note span {
  font-size: 15px;
  color: #374151;
  line-height: 1.5;
}

.step-warning {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 16px;
  background: #fef2f2;
  border-left: 4px solid #ef4444;
  border-radius: 8px;
  margin: 10px 0;
}

.step-warning svg {
  width: 18px;
  height: 18px;
  color: #ef4444;
  flex-shrink: 0;
  margin-top: 2px;
}

.success-message {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
  border: 1px solid #10b981;
  border-radius: 12px;
  margin: 20px 0;
}

.success-message svg {
  width: 24px;
  height: 24px;
  color: #065f46;
}

.success-message span {
  font-size: 16px;
  font-weight: 600;
  color: #065f46;
}

.warning-message {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%);
  border: 1px solid #f59e0b;
  border-radius: 12px;
  margin: 20px 0 30px;
}

.warning-message svg {
  width: 24px;
  height: 24px;
  color: #92400e;
}

.warning-message span {
  font-size: 16px;
  font-weight: 600;
  color: #92400e;
}

.badge {
  display: inline-block;
  background: #e5e7eb;
  color: #374151;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  margin: 0 2px;
}

.badge-success {
  background: #10b981;
  color: white;
}

.badge-primary {
  background: #4f46e5;
  color: white;
}

.btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: #4f46e5;
  color: white;
  border-radius: 6px;
  font-weight: 700;
  margin: 0 4px;
}

.buttons-demo {
  display: flex;
  gap: 15px;
  margin-top: 15px;
}

.demo-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  background: white;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  cursor: default;
  transition: all 0.2s ease;
}

.demo-btn.edit {
  border-color: #f59e0b;
  color: #92400e;
}

.demo-btn.edit svg {
  color: #f59e0b;
  width: 16px;
  height: 16px;
}

.demo-btn.save {
  border-color: #10b981;
  color: #065f46;
}

.demo-btn.save svg {
  color: #10b981;
  width: 16px;
  height: 16px;
}

.important-note {
  background: linear-gradient(135deg, #fef2f2 0%, #fee2e2 100%);
  border: 1px solid #ef4444;
  border-radius: 16px;
  padding: 25px;
  margin-top: 30px;
}

.note-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
}

.note-header svg {
  width: 28px;
  height: 28px;
  color: #dc2626;
}

.note-header h3 {
  font-size: 20px;
  font-weight: 700;
  color: #dc2626;
  margin: 0;
}

.note-content p {
  font-size: 16px;
  line-height: 1.6;
  color: #374151;
  margin: 0 0 10px 0;
}

.note-content p:first-child {
  font-weight: 600;
  color: #991b1b;
}

.link {
  color: #4f46e5;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.2s ease;
}

.link:hover {
  color: #4338ca;
  text-decoration: underline;
}

/* Адаптивность */
@media (max-width: 768px) {
  .info-container {
    padding: 0 15px;
  }

  .instruction-card {
    padding: 20px;
  }

  .card-header {
    flex-direction: column;
    text-align: center;
    gap: 15px;
  }

  .step {
    flex-direction: column;
    gap: 15px;
  }

  .step:not(:last-child)::after {
    left: 50%;
    top: 48px;
    height: calc(100% - 48px);
    transform: translateX(-50%);
  }

  .step-number {
    align-self: center;
  }

  .buttons-demo {
    flex-direction: column;
  }
}
</style>