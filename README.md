### Запуск Vue.js или Фронденд

1. Скопируйте .env.example и переименуйте в .env
2. Установите докер на линукс или если винда(https://www.docker.com/products/docker-desktop/):
```bash 
   docker-compose logs backend (логи бэкент контейнера)
   docker-compose logs frontend (логи фронтенд контейнера)
   docker-compose ps (посмотреть статус контейнеров, собрались ли и запустились ли)
   docker-compose down -v (Остановить все контейнеры с удалением)
   docker-compose down (просто остановить все контейнеры)
   docker-compose down frontend(остановить только контейнер фронтенд)
   docker-compose restart frontend(перезагрузить только контейнер фронтенд)
   docker-compose restart(перезагрузить всё контейнеры)
   docker-compose up --build -d (запустить контейнеры со сборкой и оставить свободной консоль)
   docker-compose up --build (запустить контейнеры со сборкой и быть зависимой от консоли но видеть все события)
   docker-compose up -d (просто запустить контейнеры)
```
3. Введите команды для установки:
```bash
docker-compose up --build -d
Далее проверяем что все запустилось docker-compose ps и так же логи фронта и бэка на предмет ошибок
Заходим на сайт по адресу: http://localhost:3001/ и кайфуем :)
```
4. Инфо об основных скриптах:
```bash 
   go mod tidy (устанавливает зависимости)
   go run migrate.go (выполнит миграции в базу данных)
   go run cmd/app/main.go (запуск сервера)
   go run ./cmd/worker/main.go -once (временная команда для подтягивания статистики от ВБ)
   go run ./cmd/articles_worker/main.go -once (временая команда для подтягивания артикулов/карточек из ВБ)
```

### Git - ведение версионности Semantic Versioning (SemVer)
```
MAJOR - несовместимые изменения API
MINOR - новая функциональность с обратной совместимостью
PATCH - исправление багов
------------------------------------------------------------
v1.0.0     - Первый релиз
v1.0.1     - Исправление бага
v1.1.0     - Добавлена новая функциональность
v2.0.0     - Ломающие изменения
```

### На будущее - (пример) крон на запуск воркера по запросу статы
```bash
# Запускать каждые 5 минут
*/5 * * * * root cd /var/www/wbrost-go/backend && /usr/local/go/bin/go run cmd/worker/main.go --once >> /var/log/wbrost-worker.log 2>&1
```