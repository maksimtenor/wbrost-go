-- Таблица wb_articles из m250603_214347_create_wb_articles_table
CREATE TABLE IF NOT EXISTS wb_articles (
                                           id SERIAL PRIMARY KEY,
                                           articule INT NOT NULL,
                                           name VARCHAR(1000),
    photo VARCHAR(1000),
    created DATE,
    updated DATE,
    updated_at DATE,
    id_user INT NOT NULL,
    cost_price VARCHAR(255) DEFAULT '0',
    self_ransom INT NOT NULL DEFAULT 0,
    self_ransom_price VARCHAR(255) DEFAULT '0',
    rus_size VARCHAR(5) DEFAULT '',
    eu_size VARCHAR(5) DEFAULT '',
    chrt_id INT DEFAULT 0,
    barcode VARCHAR(500) DEFAULT '',
    internal_id VARCHAR(500)
    );

CREATE INDEX idx_wb_articles_articule ON wb_articles(articule);
CREATE INDEX idx_wb_articles_id_user ON wb_articles(id_user);

COMMENT ON COLUMN wb_articles.rus_size IS 'wbSize - Rus size';
COMMENT ON COLUMN wb_articles.eu_size IS 'techSize - Eu size';
COMMENT ON COLUMN wb_articles.chrt_id IS 'chrtID';
COMMENT ON COLUMN wb_articles.barcode IS 'skus';
COMMENT ON COLUMN wb_articles.internal_id IS 'nmUUID';