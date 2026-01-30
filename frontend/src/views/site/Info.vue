<script setup>
import BaseLayout from "@/components/layout/BaseLayout.vue";
</script>

<template>
  <BaseLayout>
    <template #title-icon>
      <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
           stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
        <line x1="12" y1="17" x2="12.01" y2="17"></line>
      </svg>
    </template>
    <template #title>Как пользоваться платформой</template>

    <div class="info-container">
      <!-- Карточка с подключением -->
      <div class="instruction-card">
        <div class="card-header">
          <div class="card-icon">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                 stroke-width="2">
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
              <p>Зайдите в ваш аккаунт на сайте <a href="https://seller.wildberries.ru" target="_blank" class="link">seller.wildberries.ru</a>
              </p>
            </div>
          </div>

          <div class="step">
            <div class="step-number">2</div>
            <div class="step-content">
              <p>В личном кабинете нажмите на имя профиля и выберите:</p>
              <p class="step-detail">Настройки → Доступ к API</p>
              <a href="https://seller.wildberries.ru/supplier-settings/access-to-api" target="_blank" class="step-link">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
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
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
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
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
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
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
               stroke-width="2">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <span>Готово! Теперь вам доступна возможность загрузки статистики из Wildberries</span>
        </div>

        <div class="step-note">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
               stroke-width="2">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
          </svg>
          <span>Проверить активацию токена: Профиль → Мои API ключи</span>
        </div>
      </div>

      <!-- Карточка с загрузкой данных -->
      <div class="instruction-card">
        <div class="card-header">
          <div class="card-icon">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                 stroke-width="2">
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
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
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
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
               stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
          </svg>
          <span>Возможна ошибка при повторной загрузке одинаковых данных</span>
        </div>

        <div class="step-detail">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
               stroke-width="2">
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
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                 stroke-width="2">
              <line x1="12" y1="1" x2="12" y2="23"></line>
              <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
            </svg>
          </div>
          <h2 class="card-title">Как добавить себестоимость товара</h2>
        </div>

        <div class="warning-message">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
               stroke-width="2">
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
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                       stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                  </svg>
                  Редактировать
                </button>
                <button class="demo-btn save">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                       stroke-width="2">
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
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
               stroke-width="2">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <span>Готово! Эти данные будут учитываться при расчёте статистики</span>
        </div>
      </div>

      <!-- Важное примечание -->
      <div class="important-note">
        <div class="note-header">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
               stroke-width="2">
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
  </BaseLayout>
</template>

<style scoped>
@import '@/assets/css/views/site/info.css';
</style>