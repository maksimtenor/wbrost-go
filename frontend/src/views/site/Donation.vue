<script setup>
import Navbar from "../../components/layout/Navbar.vue";
import Sidebar from "../../components/layout/Sidebar.vue";
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
                <path d="M20.42 4.58a5.4 5.4 0 0 0-7.65 0l-.77.78-.77-.78a5.4 5.4 0 0 0-7.65 0C1.46 6.7 1.33 10.28 4 13l8 8 8-8c2.67-2.72 2.54-6.3.42-8.42z"></path>
                <line x1="12" y1="8" x2="12" y2="16"></line>
                <line x1="8" y1="12" x2="16" y2="12"></line>
              </svg>
              Поддержка проекта
            </h1>
          </div>
        </div>
      </div>
    </div>

    <div class="content">
      <div class="donation-container">
        <!-- Основная карточка -->
        <div class="donation-card">
          <div class="donation-header">
            <div class="donation-icon">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20.42 4.58a5.4 5.4 0 0 0-7.65 0l-.77.78-.77-.78a5.4 5.4 0 0 0-7.65 0C1.46 6.7 1.33 10.28 4 13l8 8 8-8c2.67-2.72 2.54-6.3.42-8.42z"></path>
              </svg>
            </div>
            <h2 class="donation-title">Поддержите развитие проекта</h2>
            <p class="donation-subtitle">Ваша поддержка помогает нам становиться лучше каждый день</p>
          </div>

          <div class="donation-content">
            <div class="benefits-grid">
              <div class="benefit">
                <div class="benefit-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="2" y1="12" x2="22" y2="12"></line>
                    <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
                  </svg>
                </div>
                <div class="benefit-content">
                  <h4>Улучшение платформы</h4>
                  <p>Постоянная оптимизация и обновления</p>
                </div>
              </div>
            </div>

            <div class="donation-message">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
  color: #ef4444;
}

.donation-container {
  //max-width: 1000px;
  margin: 0 auto;
  padding: 0 20px;
}

/* Основная карточка */
.donation-card {
  background: linear-gradient(135deg, #ff6b6b 0%, #ff8e8e 100%);
  border-radius: 20px;
  padding: 40px;
  margin-bottom: 30px;
  color: white;
  box-shadow: 0 8px 30px rgba(255, 107, 107, 0.3);
  position: relative;
  overflow: hidden;
}

.donation-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 1%, transparent 1%);
  background-size: 20px 20px;
  opacity: 0.3;
  animation: float 20s linear infinite;
}

@keyframes float {
  0% { transform: translate(0, 0) rotate(0deg); }
  100% { transform: translate(-50px, -50px) rotate(360deg); }
}

.donation-header {
  text-align: center;
  margin-bottom: 40px;
  position: relative;
  z-index: 1;
}

.donation-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  margin-bottom: 20px;
  backdrop-filter: blur(10px);
}

.donation-icon svg {
  width: 40px;
  height: 40px;
}

.donation-title {
  font-size: 32px;
  font-weight: 800;
  margin: 0 0 10px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.donation-subtitle {
  font-size: 18px;
  opacity: 0.9;
  margin: 0;
}

/* Преимущества */
.benefits-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 25px;
  margin-bottom: 30px;
  position: relative;
  z-index: 1;
}

.benefit {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 25px;
  display: flex;
  align-items: flex-start;
  gap: 20px;
  transition: transform 0.3s ease;
}

.benefit:hover {
  transform: translateY(-5px);
  background: rgba(255, 255, 255, 0.25);
}

.benefit-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  flex-shrink: 0;
}

.benefit-icon svg {
  width: 24px;
  height: 24px;
}

.benefit-content h4 {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 8px 0;
}

.benefit-content p {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
  line-height: 1.5;
}

.donation-message {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 15px;
  position: relative;
  z-index: 1;
}

.donation-message svg {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}

.donation-message p {
  font-size: 16px;
  font-weight: 500;
  margin: 0;
}

/* Карточка с действиями */
.donation-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
  margin-bottom: 30px;
}

@media (max-width: 768px) {
  .donation-actions {
    grid-template-columns: 1fr;
  }
}

.donation-info {
  background: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.donation-info h3 {
  font-size: 24px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 15px 0;
}

.donation-info p {
  font-size: 15px;
  color: #6b7280;
  margin: 0 0 25px 0;
  line-height: 1.6;
}

.payment-methods {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.payment-method {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 20px;
  background: #f9fafb;
  border-radius: 12px;
  transition: transform 0.2s ease;
}

.payment-method:hover {
  transform: translateX(5px);
  background: #f3f4f6;
}

.method-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.method-icon svg {
  width: 24px;
  height: 24px;
  color: #4f46e5;
}

.method-content h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 5px 0;
}

.method-content p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

/* Виджет доната */
.donation-widget {
  background: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.donation-widget h3 {
  font-size: 24px;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 10px 0;
}

.widget-subtitle {
  font-size: 15px;
  color: #6b7280;
  margin: 0 0 25px 0;
}

.donation-amounts {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  margin-bottom: 25px;
}

@media (max-width: 480px) {
  .donation-amounts {
    grid-template-columns: 1fr;
  }
}

.amount-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.amount-option:hover {
  border-color: #4f46e5;
  transform: translateY(-2px);
}

.amount-option.popular {
  border-color: #ff6b6b;
  background: #fff5f5;
}

.amount-option.popular .amount {
  color: #ff6b6b;
}

.amount {
  font-size: 20px;
  font-weight: 700;
  color: #1a202c;
  margin-bottom: 5px;
}

.amount-label {
  font-size: 12px;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.yoomoney-widget {
  margin-bottom: 20px;
}

.widget-note {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 15px;
  background: #f0f9ff;
  border-radius: 12px;
}

.widget-note svg {
  width: 18px;
  height: 18px;
  color: #0ea5e9;
  flex-shrink: 0;
  margin-top: 2px;
}

.widget-note p {
  font-size: 14px;
  color: #374151;
  margin: 0;
  line-height: 1.5;
}

/* Благодарности */
.thank-you-section {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  border-radius: 20px;
  padding: 40px;
  color: white;
  text-align: center;
}

.thank-you-content h2 {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 15px 0;
}

.thank-you-content h2 svg {
  width: 28px;
  height: 28px;
}

.thank-you-content p {
  font-size: 16px;
  opacity: 0.9;
  max-width: 600px;
  margin: 0 auto 30px;
  line-height: 1.6;
}

.thank-you-stats {
  display: flex;
  justify-content: center;
  gap: 40px;
  flex-wrap: wrap;
}

.stat {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-number {
  font-size: 32px;
  font-weight: 800;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* Адаптивность */
@media (max-width: 768px) {
  .donation-container {
    padding: 0 15px;
  }

  .donation-card {
    padding: 25px;
  }

  .donation-title {
    font-size: 26px;
  }

  .benefits-grid {
    grid-template-columns: 1fr;
  }

  .thank-you-section {
    padding: 30px 20px;
  }

  .thank-you-content h2 {
    font-size: 24px;
    flex-direction: column;
    gap: 10px;
  }

  .thank-you-stats {
    gap: 30px;
  }
}
</style>