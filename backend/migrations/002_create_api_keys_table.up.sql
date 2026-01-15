-- Таблица api_keys из m250603_214056_create_api_keys_table
CREATE TABLE IF NOT EXISTS api_keys (
                                        id SERIAL PRIMARY KEY,
                                        id_user INT NOT NULL,
                                        api_key VARCHAR(1000) NOT NULL,
    title VARCHAR(255),
    wb INT NOT NULL DEFAULT 0,
    ozon INT NOT NULL DEFAULT 0
    );

CREATE INDEX idx_api_keys_id_user ON api_keys(id_user);