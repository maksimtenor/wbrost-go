FROM node:18-alpine

WORKDIR /app

# Копируем package.json - ТЕПЕРЬ ПРАВИЛЬНО!
# Dockerfile находится в docker/frontend/, контекст сборки - папка frontend/
# Поэтому копируем из текущей директории
COPY frontend/package*.json ./

# Устанавливаем зависимости
RUN npm install --production=false

# Копируем исходный код
COPY frontend/ .
# Копируем ВСЁ из контекста (папка frontend/)

# Открываем порт
EXPOSE 3000

# Команда запуска в режиме разработки
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]