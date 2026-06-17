CREATE TABLE IF NOT EXISTS compounds (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    developer VARCHAR(100) NOT NULL,
    region VARCHAR(50) NOT NULL, -- 'New Cairo', 'North Coast', '6th of October'
    launch_year INT NOT NULL,
    total_units INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS sales_ledger (
    id SERIAL PRIMARY KEY,
    compound_id INT REFERENCES compounds(id),
    quarter VARCHAR(10) NOT NULL, -- e.g., 'Q1-2026'
    units_sold INT NOT NULL,
    revenue_egp NUMERIC(15, 2) NOT NULL,
    cancelled_orders INT DEFAULT 0
);

-- Seed initial records
INSERT INTO compounds (name, developer, region, launch_year, total_units) VALUES
('Mivida', 'Emaar', 'New Cairo', 2010, 5000),
('Marassi', 'Emaar', 'North Coast', 2006, 8500),
('Zed East', 'Ora Developers', 'New Cairo', 2020, 3200),
('Badya', 'Palm Hills', '6th of October', 2018, 12000);

INSERT INTO sales_ledger (compound_id, quarter, units_sold, revenue_egp, cancelled_orders) VALUES
(1, 'Q1-2026', 42, 310000000.00, 2),
(1, 'Q2-2026', 55, 415000000.00, 1),
(2, 'Q1-2026', 12, 180000000.00, 4),
(3, 'Q1-2026', 68, 520000000.00, 0),
(4, 'Q1-2026', 95, 680000000.00, 8),
(4, 'Q2-2026', 34, 250000000.00, 12);