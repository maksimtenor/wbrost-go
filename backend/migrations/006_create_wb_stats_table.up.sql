-- Таблица wb_stats из m250603_214156_create_wb_stats_table
CREATE TABLE IF NOT EXISTS wb_stats (
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

CREATE INDEX idx_wb_stats_user_id ON wb_stats(user_id);
CREATE INDEX idx_wb_stats_nm_id ON wb_stats(nm_id);
CREATE INDEX idx_wb_stats_ppvz_for_pay ON wb_stats(ppvz_for_pay);
CREATE INDEX idx_wb_stats_quantity ON wb_stats(quantity);
CREATE INDEX idx_wb_stats_realizationreport_id ON wb_stats(realizationreport_id);
CREATE INDEX idx_wb_stats_rebill_logistic_cost ON wb_stats(rebill_logistic_cost);
CREATE INDEX idx_wb_stats_report_type ON wb_stats(report_type);
CREATE INDEX idx_wb_stats_return_amount ON wb_stats(return_amount);
CREATE INDEX idx_wb_stats_rid ON wb_stats(rid);
CREATE INDEX idx_wb_stats_supplier_oper_name ON wb_stats(supplier_oper_name);
CREATE INDEX idx_wb_stats_rr_dt ON wb_stats(rr_dt);