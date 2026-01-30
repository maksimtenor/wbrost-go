<script setup>
import BaseLayout from "@/components/layout/BaseLayout.vue";
</script>

<template>
  <BaseLayout>
    <!-- Заголовок страницы -->
    <template #title-icon>
      <svg class="title-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
           stroke-width="2">
        <path
            d="M20.42 4.58a5.4 5.4 0 0 0-7.65 0l-.77.78-.77-.78a5.4 5.4 0 0 0-7.65 0C1.46 6.7 1.33 10.28 4 13l8 8 8-8c2.67-2.72 2.54-6.3.42-8.42z"></path>
        <line x1="12" y1="8" x2="12" y2="16"></line>
        <line x1="8" y1="12" x2="16" y2="12"></line>
      </svg>
    </template>
    <template #title>Поддержка проекта</template>

    <div class="donation-container">
      <!-- Основная карточка -->
      <div class="donation-card">
        <div class="donation-header">
          <div class="donation-icon">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                 stroke-width="2">
              <path
                  d="M20.42 4.58a5.4 5.4 0 0 0-7.65 0l-.77.78-.77-.78a5.4 5.4 0 0 0-7.65 0C1.46 6.7 1.33 10.28 4 13l8 8 8-8c2.67-2.72 2.54-6.3.42-8.42z"></path>
            </svg>
          </div>
          <h2 class="donation-title">Поддержите развитие проекта</h2>
          <p class="donation-subtitle">Ваша поддержка помогает нам становиться лучше каждый день</p>
        </div>

        <div class="donation-content">
          <div class="benefits-grid">
            <div class="benefit">
              <div class="benefit-icon">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
                  <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"></polyline>
                  <polyline points="17 6 23 6 23 12"></polyline>
                </svg>
              </div>
              <div class="benefit-content">
                <h4>Новые функции</h4>
                <p>Разработка полезных инструментов для анализа</p>
              </div>
            </div>

            <div class="benefit">
              <div class="benefit-icon">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                </svg>
              </div>
              <div class="benefit-content">
                <h4>Стабильная работа</h4>
                <p>Обеспечение бесперебойной работы серверов</p>
              </div>
            </div>

            <div class="benefit">
              <div class="benefit-icon">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
                  <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                  <polyline points="15 3 21 3 21 9"></polyline>
                  <line x1="10" y1="14" x2="21" y2="3"></line>
                </svg>
              </div>
              <div class="benefit-content">
                <h4>Быстрая поддержка</h4>
                <p>Моментальное решение ваших вопросов</p>
              </div>
            </div>

            <div class="benefit">
              <div class="benefit-icon">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="2" y1="12" x2="22" y2="12"></line>
                  <path
                      d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
                </svg>
              </div>
              <div class="benefit-content">
                <h4>Улучшение платформы</h4>
                <p>Постоянная оптимизация и обновления</p>
              </div>
            </div>
          </div>

          <div class="donation-message">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                 stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="16" x2="12" y2="12"></line>
              <line x1="12" y1="8" x2="12.01" y2="8"></line>
            </svg>
            <p>Каждое пожертвование, независимо от суммы, приближает нас к реализации новых идей!</p>
          </div>
        </div>
      </div>

      <!-- Карточка с донатом -->
      <div class="donation-actions">
        <div class="donation-info">
          <h3>Способы поддержки</h3>
          <p>Вы можете поддержать проект финансово через удобную для вас платформу</p>

          <div class="payment-methods">
            <div class="payment-method">
              <div class="method-icon">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
                  <rect x="1" y="4" width="22" height="16" rx="2" ry="2"></rect>
                  <line x1="1" y1="10" x2="23" y2="10"></line>
                </svg>
              </div>
              <div class="method-content">
                <h4>ЮMoney</h4>
                <p>Быстрый и безопасный перевод</p>
              </div>
            </div>

            <div class="payment-method">
              <div class="method-icon">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="8" y1="15" x2="16" y2="15"></line>
                  <line x1="9" y1="9" x2="9.01" y2="9"></line>
                  <line x1="15" y1="9" x2="15.01" y2="9"></line>
                </svg>
              </div>
              <div class="method-content">
                <h4>Другие способы</h4>
                <p>Напишите нам для получения реквизитов</p>
              </div>
            </div>
          </div>
        </div>

        <div class="donation-widget">
          <h3>Сделать пожертвование</h3>
          <p class="widget-subtitle">Выберите удобный способ перевода через ЮMoney</p>

          <div class="widget-content">
            <div class="donation-amounts">
              <div class="amount-option">
                <span class="amount">100 ₽</span>
                <span class="amount-label">Поддержка</span>
              </div>
              <div class="amount-option popular">
                <span class="amount">300 ₽</span>
                <span class="amount-label">Популярно</span>
              </div>
              <div class="amount-option">
                <span class="amount">500 ₽</span>
                <span class="amount-label">Спасибо!</span>
              </div>
              <div class="amount-option">
                <span class="amount">1000 ₽</span>
                <span class="amount-label">Меценат</span>
              </div>
            </div>

            <div class="yoomoney-widget">
              <iframe
                  src="https://yoomoney.ru/quickpay/fundraise/button?billNumber=16I97N3TS3F.241116&"
                  width="100%"
                  height="50"
                  frameborder="0"
                  allowtransparency="true"
                  scrolling="no">
              </iframe>
            </div>

            <div class="widget-note">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                   stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="16" x2="12" y2="12"></line>
                <line x1="12" y1="8" x2="12.01" y2="8"></line>
              </svg>
              <p>Переводы обрабатываются безопасно через ЮKassa. Данные вашей карты защищены.</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Благодарности -->
      <div class="thank-you-section">
        <div class="thank-you-content">
          <h2>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                 stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            Спасибо за вашу поддержку!
          </h2>
          <p>Благодаря вам мы можем продолжать развивать платформу и делать её лучше для всех пользователей.</p>

          <div class="thank-you-stats">
            <div class="stat">
              <span class="stat-number">100%</span>
              <span class="stat-label">Идёт на развитие</span>
            </div>
            <div class="stat">
              <span class="stat-number">24/7</span>
              <span class="stat-label">Поддержка</span>
            </div>
            <div class="stat">
              <span class="stat-number">∞</span>
              <span class="stat-label">Благодарность</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </BaseLayout>
</template>

<style>
@import '@/assets/css/views/site/donation.css';
</style>