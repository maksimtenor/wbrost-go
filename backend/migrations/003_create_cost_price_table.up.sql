-- Таблица cost_price из m250603_214114_create_cost_price_table
CREATE TABLE IF NOT EXISTS cost_price (
                                          id SERIAL PRIMARY KEY,
                                          id_user INT NOT NULL,
                                          articul VARCHAR(255) NOT NULL,
    price_cost VARCHAR(255) NOT NULL,
    date_create DATE NOT NULL
    );

CREATE INDEX idx_cost_price_id_user ON cost_price(id_user);