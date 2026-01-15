-- Таблица users из m250603_214206_create_users_table
CREATE TABLE IF NOT EXISTS users (
                                     id_user SERIAL PRIMARY KEY,
                                     taxes INT DEFAULT 0,
                                     username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    admin SMALLINT NOT NULL DEFAULT 0,
    block INT DEFAULT 0,
    pro INT NOT NULL DEFAULT 0,
    name VARCHAR(35),
    phone VARCHAR(50),
    wb_key TEXT,
    ozon_key TEXT,
    u2782212_wbrosus INT NOT NULL DEFAULT 0,
    ozon_status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    del INT DEFAULT 0,
    last_login TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_users_admin ON users(admin);
CREATE INDEX idx_users_block ON users(block);
CREATE INDEX idx_users_del ON users(del);
CREATE INDEX idx_users_pro ON users(pro);
CREATE INDEX idx_users_taxes ON users(taxes);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_id_user ON users(id_user);

COMMENT ON COLUMN users.taxes IS 'Налоговый %';
COMMENT ON COLUMN users.wb_key IS 'API Ключ WB';
COMMENT ON COLUMN users.ozon_key IS 'API Ключ Ozon';
COMMENT ON COLUMN users.pro IS 'Активен ли про акк';
COMMENT ON COLUMN users.last_login IS 'Дата последнего входа в аккаунт';