-- Таблица для Telegram бота
CREATE TABLE IF NOT EXISTS bot_user (
                                        id SERIAL PRIMARY KEY,
                                        telegram_chat_id BIGINT NOT NULL UNIQUE,
                                        wb_api_key VARCHAR(1000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_bot_user_telegram_chat_id ON bot_user(telegram_chat_id);