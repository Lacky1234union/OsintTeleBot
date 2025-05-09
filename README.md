OSINT Telegram Bot 🔍

Telegram-бот для работы с OSINT-разведкой, хранения и анализа данных из утечек (слитые базы, дампы, соцсети).
⚙️ Технологии

    Язык: Go (Golang)

    Базы данных: PostgreSQL (основная), MongoDB (для "грязных" дампов)

    Архитектура: Clean Architecture + Modular Design

    Telegram API: go-telegram-bot-api

    Миграции: Goose

    Логирование: Zerolog/Slog

    Конфиги: ENV-переменные

🚀 Запуск проекта
1. Клонирование репозитория
bash

git clone https://github.com/yourusername/osint-bot.git
cd osint-bot

2. Настройка окружения

Создайте файл .env в корне проекта:
ini

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=osintbot

# Telegram
TELEGRAM_TOKEN=your_bot_token

3. Запуск PostgreSQL (Docker)
bash

docker-compose up -d

Файл docker-compose.yml:
yaml

version: '3'
services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: yourpassword
      POSTGRES_DB: osintbot
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
volumes:
  postgres_data:

4. Применение миграций
bash

goose -dir ./migrations postgres "user=postgres password=yourpassword dbname=osintbot sslmode=disable" up

5. Запуск бота
bash

go run ./cmd/bot/main.go

📂 Структура проекта
bash

osint-bot/
├── cmd/                  # Точки входа
├── internal/             # Основная бизнес-логика
│   ├── app/              # Ядро приложения
│   │   ├── errors/       # Кастомные ошибки
│   │   ├── models/       # Сущности данных
│   │   ├── repositories/ # Слой работы с БД
│   │   └── services/     # Бизнес-логика
│   ├── config/           # Конфигурация
│   └── transport/        # Внешние интерфейсы
│       ├── telegram/     # Telegram handlers
│       └── http/         # REST API (optional)
├── pkg/                  # Общие утилиты
│   ├── database/         # Инициализация БД
│   ├── errorhandler/     # Обработка ошибок
│   └── logger/           # Логирование
├── migrations/           # SQL миграции
├── scripts/              # Вспомогательные скрипты
└── .env.example          # Шаблон конфигурации

🛠 Команды бота

    /start — Приветствие.

    /find_email [email] — Поиск по почте.

    /find_phone [phone] — Поиск по номеру.

    /add_person — Добавить человека в базу.

    /import_dump — Загрузить дамп (файлом).

🔒 Безопасность

    Пароли шифруются (AES-256).

    Доступ к БД только по IP/VPN.

    Логирование всех операций.

📈 Планы по развитию

    Поддержка MongoDB для сырых дампов.

    Интеграция с Elasticsearch для полнотекстового поиска.

    Веб-интерфейс для администрирования.

    Аналитика частотности паролей/почт.

📜 Лицензия

MIT

Автор: [Ваше имя]
GitHub: github.com/yourusername
