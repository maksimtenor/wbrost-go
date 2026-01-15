wbros-go/                    # Корень проекта
├── backend/                 # Go бэкенд (REST API)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go     # Точка входа
│   ├── internal/
│   │   ├── config/         # Конфигурация
│   │   ├── entity/         # Сущности (модели)
│   │   │   ├── user.go
│   │   │   ├── product.go
│   │   │   └── order.go
│   │   ├── handler/        # HTTP обработчики (контроллеры)
│   │   │   ├── user_handler.go
│   │   │   ├── auth_handler.go
│   │   │   └── ...
│   │   ├── service/        # Бизнес-логика (сервисы)
│   │   │   ├── user_service.go
│   │   │   └── ...
│   │   ├── repository/     # Работа с БД (репозитории)
│   │   │   ├── user_repository.go
│   │   │   └── ...
│   │   └── middleware/     # Middleware (аутентификация, логирование)
│   ├── pkg/                # Общедоступные пакеты
│   │   ├── database/
│   │   ├── auth/
│   │   └── ...
│   ├── migrations/         # SQL миграции
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/               # Vue.js фронтенд
│   ├── src/
│   │   ├── api/           # API клиенты для Go бэкенда
│   │   │   ├── auth.js
│   │   │   ├── users.js
│   │   │   └── ...
│   │   ├── assets/        # Статика (картинки, шрифты)
│   │   ├── components/    # Vue компоненты
│   │   │   ├── common/
│   │   │   ├── layout/
│   │   │   └── ...
│   │   ├── router/        # Vue Router
│   │   │   └── index.js
│   │   ├── store/         # Vuex/Pinia хранилище
│   │   │   └── index.js
│   │   ├── views/         # Страницы (Yii2 views → Vue views)
│   │   │   ├── Auth/
│   │   │   │   ├── Login.vue
│   │   │   │   └── Register.vue
│   │   │   ├── Users/
│   │   │   │   ├── Index.vue    # users/index
│   │   │   │   ├── Create.vue   # users/create
│   │   │   │   └── Edit.vue     # users/update
│   │   │   └── ...
│   │   ├── App.vue
│   │   └── main.js
│   ├── public/
│   │   ├── index.html
│   │   └── ...
│   ├── package.json
│   ├── vite.config.js     # или vue.config.js
│   └── Dockerfile
│
├── docker-compose.yml     # Общий docker-compose
├── nginx/                 # Nginx для проксирования
│   └── nginx.conf
├── .env.example           # Переменные окружения
└── README.md