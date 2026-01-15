-- Таблица report_paid_storage из m250603_214130_create_report_paid_storage_table
CREATE TABLE IF NOT EXISTS report_paid_storage (
                                                   id SERIAL PRIMARY KEY,
                                                   id_task VARCHAR(500) NOT NULL,
    id_user INT NOT NULL,
    date_from DATE NOT NULL,
    date_to DATE NOT NULL,
    pay VARCHAR(500) NOT NULL DEFAULT '0'
    );

CREATE INDEX idx_report_paid_storage_id_task ON report_paid_storage(id_task);
CREATE INDEX idx_report_paid_storage_id_user ON report_paid_storage(id_user);
CREATE INDEX idx_report_paid_storage_id ON report_paid_storage(id);