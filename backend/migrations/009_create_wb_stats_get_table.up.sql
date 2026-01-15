-- Таблица wb_stats_get из m250603_214429_create_wb_stats_get_table
CREATE TABLE IF NOT EXISTS wb_stats_get (
                                            id SERIAL PRIMARY KEY,
                                            id_user INT NOT NULL,
                                            status INT,
                                            date_from VARCHAR(255) NOT NULL,
    date_to VARCHAR(255) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_error VARCHAR(1000)
    );