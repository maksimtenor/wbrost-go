-- Таблица stat из m250603_214156_create_stat_table
CREATE TABLE IF NOT EXISTS stat (
                                    id SERIAL PRIMARY KEY,
                                    hash_info VARCHAR(2000),
    user_id INT NOT NULL,
    nm_id BIGINT,
    ppvz_for_pay VARCHAR(500),
    supplier_oper_name INT,
    delivery_rub DECIMAL(10, 2) DEFAULT 0.00,
    penalty DECIMAL(10, 2) DEFAULT 0.00,
    additional_payment DECIMAL(10, 2) DEFAULT 0.00,
    storage_fee DECIMAL(10, 2) DEFAULT 0.00,
    rebill_logistic_cost VARCHAR(500) DEFAULT '0.00',
    acquiring_fee DECIMAL(10, 2) DEFAULT 0.00,
    acquiring_percent DECIMAL(10, 2) DEFAULT 0.00,
    ppvz_sales_commission DECIMAL(10, 2) DEFAULT 0.00,
    deduction DECIMAL(10, 2) DEFAULT 0.00,
    ppvz_spp_prc VARCHAR(500),
    ppvz_kvw_prc_base VARCHAR(500),
    ppvz_kvw_prc VARCHAR(500),
    acceptance DECIMAL(10, 2) DEFAULT 0.00,
    dlv_prc DECIMAL(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    rr_dt TIMESTAMP,
    shk_id BIGINT,
    sticker_id VARCHAR(500),
    gi_id BIGINT,
    realizationreport_id BIGINT,
    barcode VARCHAR(1000),
    bonus_type_name VARCHAR(1000),
    last_error VARCHAR(500),
    brand_name VARCHAR(500),
    ppvz_office_id BIGINT,
    assembly_id BIGINT,
    sa_name VARCHAR(1000),
    ppvz_vw_nds VARCHAR(500),
    ppvz_vw VARCHAR(500),
    gi_box_type_name VARCHAR(500),
    subject_name VARCHAR(1000),
    ts_name VARCHAR(500),
    quantity INT,
    retail_price DECIMAL(10, 2) DEFAULT 0.00,
    retail_amount DECIMAL(10, 2) DEFAULT 0.00,
    commission_percent DECIMAL(10, 2) DEFAULT 0.00,
    office_name VARCHAR(1000),
    order_dt TIMESTAMP,
    sale_dt TIMESTAMP,
    delivery_amount INT,
    return_amount INT,
    report_type INT,
    srid VARCHAR(2000),
    rid BIGINT
    );

CREATE INDEX idx_stat_user_id ON stat(user_id);
CREATE INDEX idx_stat_nm_id ON stat(nm_id);
CREATE INDEX idx_stat_ppvz_for_pay ON stat(ppvz_for_pay);
CREATE INDEX idx_stat_quantity ON stat(quantity);
CREATE INDEX idx_stat_realizationreport_id ON stat(realizationreport_id);
CREATE INDEX idx_stat_rebill_logistic_cost ON stat(rebill_logistic_cost);
CREATE INDEX idx_stat_report_type ON stat(report_type);
CREATE INDEX idx_stat_return_amount ON stat(return_amount);
CREATE INDEX idx_stat_rid ON stat(rid);
CREATE INDEX idx_stat_supplier_oper_name ON stat(supplier_oper_name);
CREATE INDEX idx_stat_rr_dt ON stat(rr_dt);