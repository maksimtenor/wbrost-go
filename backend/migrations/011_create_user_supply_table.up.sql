-- Таблица для поставок пользователя
CREATE TABLE IF NOT EXISTS user_supply (
                                           id SERIAL PRIMARY KEY,
                                           user_id INT NOT NULL,
                                           supply_data JSONB,
                                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_supply_user_id ON user_supply(user_id);

ALTER TABLE user_supply
    ADD CONSTRAINT fk_user_supply_user
        FOREIGN KEY (user_id)
            REFERENCES users(id_user)
            ON DELETE CASCADE;