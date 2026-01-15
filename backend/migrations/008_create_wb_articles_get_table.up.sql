-- Таблица wb_articles_get из m250603_214410_create_wb_articles_get_table
CREATE TABLE IF NOT EXISTS wb_articles_get (
                                               id SERIAL PRIMARY KEY,
                                               id_user INT NOT NULL,
                                               status INT,
                                               created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                               updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                               last_error VARCHAR(1000)
    );