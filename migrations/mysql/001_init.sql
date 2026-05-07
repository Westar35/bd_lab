CREATE TABLE IF NOT EXISTS vehicle_class (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS vehicle_status (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS fuel_type (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS transmission_type (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS acquisition_type (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS maintenance_type (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS payment_type (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS contract_type (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS contract_status (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS vehicle (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    vin CHAR(17) NOT NULL UNIQUE,
    reg_number VARCHAR(20) NOT NULL UNIQUE,
    make VARCHAR(64) NOT NULL,
    model VARCHAR(64) NOT NULL,
    year INT NOT NULL CHECK (year BETWEEN 1980 AND 2100),
    vehicle_class_id BIGINT NOT NULL,
    status_id BIGINT NOT NULL,
    fuel_type_id BIGINT NOT NULL,
    transmission_type_id BIGINT NOT NULL,
    acquisition_type_id BIGINT NOT NULL,
    acquisition_date DATE NOT NULL,
    acquisition_cost DECIMAL(14,2) NOT NULL CHECK (acquisition_cost >= 0),
    current_odometer_km INT NOT NULL CHECK (current_odometer_km >= 0),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CHECK (CHAR_LENGTH(vin) = 17),
    CONSTRAINT fk_vehicle_class FOREIGN KEY (vehicle_class_id) REFERENCES vehicle_class(id),
    CONSTRAINT fk_vehicle_status FOREIGN KEY (status_id) REFERENCES vehicle_status(id),
    CONSTRAINT fk_vehicle_fuel_type FOREIGN KEY (fuel_type_id) REFERENCES fuel_type(id),
    CONSTRAINT fk_vehicle_transmission FOREIGN KEY (transmission_type_id) REFERENCES transmission_type(id),
    CONSTRAINT fk_vehicle_acquisition FOREIGN KEY (acquisition_type_id) REFERENCES acquisition_type(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS driver (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    fio VARCHAR(150) NOT NULL,
    license_number VARCHAR(50) NOT NULL UNIQUE,
    phone VARCHAR(40),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS department (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS vehicle_assignment (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    vehicle_id BIGINT NOT NULL,
    driver_id BIGINT NULL,
    department_id BIGINT NOT NULL,
    date_from DATE NOT NULL,
    date_to DATE,
    is_primary BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (date_to IS NULL OR date_to >= date_from),
    CONSTRAINT fk_assignment_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicle(id),
    CONSTRAINT fk_assignment_driver FOREIGN KEY (driver_id) REFERENCES driver(id),
    CONSTRAINT fk_assignment_department FOREIGN KEY (department_id) REFERENCES department(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS trip_sheet (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    trip_date DATE NOT NULL,
    vehicle_id BIGINT NOT NULL,
    driver_id BIGINT NOT NULL,
    department_id BIGINT NOT NULL,
    odometer_start INT NOT NULL CHECK (odometer_start >= 0),
    odometer_end INT NOT NULL CHECK (odometer_end >= 0),
    route VARCHAR(255),
    purpose VARCHAR(255),
    distance_km INT NOT NULL CHECK (distance_km >= 0),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (odometer_end >= odometer_start),
    CHECK (distance_km = odometer_end - odometer_start),
    CONSTRAINT fk_trip_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicle(id),
    CONSTRAINT fk_trip_driver FOREIGN KEY (driver_id) REFERENCES driver(id),
    CONSTRAINT fk_trip_department FOREIGN KEY (department_id) REFERENCES department(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS fuel_txn (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    txn_ts DATETIME NOT NULL,
    vehicle_id BIGINT NOT NULL,
    liters DECIMAL(10,2) NOT NULL CHECK (liters > 0),
    amount DECIMAL(12,2) NOT NULL CHECK (amount >= 0),
    station VARCHAR(120) NOT NULL,
    odometer_km INT NOT NULL CHECK (odometer_km >= 0),
    payment_type_id BIGINT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_fuel_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicle(id),
    CONSTRAINT fk_fuel_payment FOREIGN KEY (payment_type_id) REFERENCES payment_type(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS maintenance_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    vehicle_id BIGINT NOT NULL,
    maintenance_type_id BIGINT NOT NULL,
    open_date DATE NOT NULL,
    close_date DATE,
    service_name VARCHAR(180) NOT NULL,
    cost DECIMAL(12,2) NOT NULL CHECK (cost >= 0),
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (close_date IS NULL OR close_date >= open_date),
    CONSTRAINT fk_maintenance_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicle(id),
    CONSTRAINT fk_maintenance_type FOREIGN KEY (maintenance_type_id) REFERENCES maintenance_type(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS counterparty (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(32) NOT NULL CHECK (type IN ('company', 'individual')),
    name VARCHAR(180) NOT NULL,
    inn VARCHAR(20),
    phone VARCHAR(40),
    email VARCHAR(120),
    address TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS contract (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    contract_type_id BIGINT NOT NULL,
    counterparty_id BIGINT NOT NULL,
    number VARCHAR(64) NOT NULL UNIQUE,
    date_from DATE NOT NULL,
    date_to DATE NOT NULL,
    status_id BIGINT NOT NULL,
    total_amount DECIMAL(14,2) CHECK (total_amount >= 0),
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CHECK (date_to >= date_from),
    CONSTRAINT fk_contract_type FOREIGN KEY (contract_type_id) REFERENCES contract_type(id),
    CONSTRAINT fk_contract_counterparty FOREIGN KEY (counterparty_id) REFERENCES counterparty(id),
    CONSTRAINT fk_contract_status FOREIGN KEY (status_id) REFERENCES contract_status(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS rental_event (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    contract_id BIGINT NOT NULL,
    vehicle_id BIGINT NOT NULL,
    pickup_ts DATETIME NOT NULL,
    return_ts DATETIME,
    price_per_day DECIMAL(12,2) NOT NULL CHECK (price_per_day >= 0),
    deposit DECIMAL(12,2) CHECK (deposit >= 0),
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CHECK (return_ts IS NULL OR return_ts >= pickup_ts),
    CONSTRAINT fk_rental_contract FOREIGN KEY (contract_id) REFERENCES contract(id),
    CONSTRAINT fk_rental_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicle(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS audit_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actor VARCHAR(100) NOT NULL,
    action VARCHAR(32) NOT NULL,
    entity VARCHAR(100) NOT NULL,
    entity_id BIGINT NOT NULL,
    details_before JSON,
    details_after JSON
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_vehicle_make ON vehicle(make);
CREATE INDEX idx_vehicle_reg_number ON vehicle(reg_number);
CREATE INDEX idx_vehicle_status_id ON vehicle(status_id);
CREATE INDEX idx_vehicle_class_id ON vehicle(vehicle_class_id);

CREATE INDEX idx_assignment_vehicle ON vehicle_assignment(vehicle_id);
CREATE INDEX idx_assignment_driver ON vehicle_assignment(driver_id);
CREATE INDEX idx_assignment_department ON vehicle_assignment(department_id);

CREATE INDEX idx_trip_sheet_trip_date ON trip_sheet(trip_date);
CREATE INDEX idx_trip_sheet_driver_vehicle ON trip_sheet(driver_id, vehicle_id, trip_date);
CREATE INDEX idx_trip_sheet_vehicle ON trip_sheet(vehicle_id);

CREATE INDEX idx_fuel_txn_ts ON fuel_txn(txn_ts);
CREATE INDEX idx_fuel_txn_vehicle ON fuel_txn(vehicle_id);

CREATE INDEX idx_maintenance_open_date ON maintenance_order(open_date);
CREATE INDEX idx_maintenance_vehicle ON maintenance_order(vehicle_id);

CREATE INDEX idx_contract_dates ON contract(date_from, date_to);
CREATE INDEX idx_contract_status ON contract(status_id);
CREATE INDEX idx_contract_counterparty ON contract(counterparty_id);

CREATE INDEX idx_rental_pickup_return ON rental_event(pickup_ts, return_ts);
CREATE INDEX idx_rental_vehicle ON rental_event(vehicle_id);
CREATE INDEX idx_rental_contract ON rental_event(contract_id);
