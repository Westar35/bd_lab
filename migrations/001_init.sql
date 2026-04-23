CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS vehicle_class (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS vehicle_status (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS fuel_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS transmission_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS acquisition_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS maintenance_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS payment_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS contract_type (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS contract_status (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS vehicle (
    id BIGSERIAL PRIMARY KEY,
    vin CHAR(17) NOT NULL UNIQUE,
    reg_number VARCHAR(20) NOT NULL UNIQUE,
    make VARCHAR(64) NOT NULL,
    model VARCHAR(64) NOT NULL,
    year INT NOT NULL CHECK (year BETWEEN 1980 AND 2100),
    vehicle_class_id BIGINT NOT NULL REFERENCES vehicle_class(id),
    status_id BIGINT NOT NULL REFERENCES vehicle_status(id),
    fuel_type_id BIGINT NOT NULL REFERENCES fuel_type(id),
    transmission_type_id BIGINT NOT NULL REFERENCES transmission_type(id),
    acquisition_type_id BIGINT NOT NULL REFERENCES acquisition_type(id),
    acquisition_date DATE NOT NULL,
    acquisition_cost NUMERIC(14,2) NOT NULL CHECK (acquisition_cost >= 0),
    current_odometer_km INT NOT NULL CHECK (current_odometer_km >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (char_length(vin) = 17)
);

CREATE TABLE IF NOT EXISTS driver (
    id BIGSERIAL PRIMARY KEY,
    fio VARCHAR(150) NOT NULL,
    license_number VARCHAR(50) NOT NULL UNIQUE,
    phone VARCHAR(40),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS department (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vehicle_assignment (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGINT NOT NULL REFERENCES vehicle(id),
    driver_id BIGINT REFERENCES driver(id),
    department_id BIGINT NOT NULL REFERENCES department(id),
    date_from DATE NOT NULL,
    date_to DATE,
    is_primary BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (date_to IS NULL OR date_to >= date_from)
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'vehicle_assignment_no_overlap'
    ) THEN
        ALTER TABLE vehicle_assignment
            ADD CONSTRAINT vehicle_assignment_no_overlap
            EXCLUDE USING gist (
                vehicle_id WITH =,
                daterange(date_from, COALESCE(date_to + 1, 'infinity'::date), '[]') WITH &&
            );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS trip_sheet (
    id BIGSERIAL PRIMARY KEY,
    trip_date DATE NOT NULL,
    vehicle_id BIGINT NOT NULL REFERENCES vehicle(id),
    driver_id BIGINT NOT NULL REFERENCES driver(id),
    department_id BIGINT NOT NULL REFERENCES department(id),
    odometer_start INT NOT NULL CHECK (odometer_start >= 0),
    odometer_end INT NOT NULL CHECK (odometer_end >= 0),
    route VARCHAR(255),
    purpose VARCHAR(255),
    distance_km INT NOT NULL CHECK (distance_km >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (odometer_end >= odometer_start),
    CHECK (distance_km = odometer_end - odometer_start)
);

CREATE TABLE IF NOT EXISTS fuel_txn (
    id BIGSERIAL PRIMARY KEY,
    txn_ts TIMESTAMPTZ NOT NULL,
    vehicle_id BIGINT NOT NULL REFERENCES vehicle(id),
    liters NUMERIC(10,2) NOT NULL CHECK (liters > 0),
    amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
    station VARCHAR(120) NOT NULL,
    odometer_km INT NOT NULL CHECK (odometer_km >= 0),
    payment_type_id BIGINT NOT NULL REFERENCES payment_type(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS maintenance_order (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGINT NOT NULL REFERENCES vehicle(id),
    maintenance_type_id BIGINT NOT NULL REFERENCES maintenance_type(id),
    open_date DATE NOT NULL,
    close_date DATE,
    service_name VARCHAR(180) NOT NULL,
    cost NUMERIC(12,2) NOT NULL CHECK (cost >= 0),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (close_date IS NULL OR close_date >= open_date)
);

CREATE TABLE IF NOT EXISTS counterparty (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(32) NOT NULL CHECK (type IN ('company', 'individual')),
    name VARCHAR(180) NOT NULL,
    inn VARCHAR(20),
    phone VARCHAR(40),
    email VARCHAR(120),
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS contract (
    id BIGSERIAL PRIMARY KEY,
    contract_type_id BIGINT NOT NULL REFERENCES contract_type(id),
    counterparty_id BIGINT NOT NULL REFERENCES counterparty(id),
    number VARCHAR(64) NOT NULL UNIQUE,
    date_from DATE NOT NULL,
    date_to DATE NOT NULL,
    status_id BIGINT NOT NULL REFERENCES contract_status(id),
    total_amount NUMERIC(14,2) CHECK (total_amount >= 0),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (date_to >= date_from)
);

CREATE TABLE IF NOT EXISTS rental_event (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id),
    vehicle_id BIGINT NOT NULL REFERENCES vehicle(id),
    pickup_ts TIMESTAMPTZ NOT NULL,
    return_ts TIMESTAMPTZ,
    price_per_day NUMERIC(12,2) NOT NULL CHECK (price_per_day >= 0),
    deposit NUMERIC(12,2) CHECK (deposit >= 0),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (return_ts IS NULL OR return_ts >= pickup_ts)
);

CREATE TABLE IF NOT EXISTS audit_log (
    id BIGSERIAL PRIMARY KEY,
    ts TIMESTAMPTZ NOT NULL DEFAULT now(),
    actor VARCHAR(100) NOT NULL,
    action VARCHAR(32) NOT NULL,
    entity VARCHAR(100) NOT NULL,
    entity_id BIGINT NOT NULL,
    details_before JSONB,
    details_after JSONB
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_vehicle_updated_at ON vehicle;
CREATE TRIGGER trg_vehicle_updated_at
BEFORE UPDATE ON vehicle
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_driver_updated_at ON driver;
CREATE TRIGGER trg_driver_updated_at
BEFORE UPDATE ON driver
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_department_updated_at ON department;
CREATE TRIGGER trg_department_updated_at
BEFORE UPDATE ON department
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_counterparty_updated_at ON counterparty;
CREATE TRIGGER trg_counterparty_updated_at
BEFORE UPDATE ON counterparty
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_contract_updated_at ON contract;
CREATE TRIGGER trg_contract_updated_at
BEFORE UPDATE ON contract
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_rental_event_updated_at ON rental_event;
CREATE TRIGGER trg_rental_event_updated_at
BEFORE UPDATE ON rental_event
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX IF NOT EXISTS idx_vehicle_make ON vehicle(make);
CREATE INDEX IF NOT EXISTS idx_vehicle_reg_number ON vehicle(reg_number);
CREATE INDEX IF NOT EXISTS idx_vehicle_status_id ON vehicle(status_id);
CREATE INDEX IF NOT EXISTS idx_vehicle_class_id ON vehicle(vehicle_class_id);

CREATE INDEX IF NOT EXISTS idx_assignment_vehicle ON vehicle_assignment(vehicle_id);
CREATE INDEX IF NOT EXISTS idx_assignment_driver ON vehicle_assignment(driver_id);
CREATE INDEX IF NOT EXISTS idx_assignment_department ON vehicle_assignment(department_id);

CREATE INDEX IF NOT EXISTS idx_trip_sheet_trip_date ON trip_sheet(trip_date);
CREATE INDEX IF NOT EXISTS idx_trip_sheet_driver_vehicle ON trip_sheet(driver_id, vehicle_id, trip_date);
CREATE INDEX IF NOT EXISTS idx_trip_sheet_vehicle ON trip_sheet(vehicle_id);

CREATE INDEX IF NOT EXISTS idx_fuel_txn_ts ON fuel_txn(txn_ts);
CREATE INDEX IF NOT EXISTS idx_fuel_txn_vehicle ON fuel_txn(vehicle_id);

CREATE INDEX IF NOT EXISTS idx_maintenance_open_date ON maintenance_order(open_date);
CREATE INDEX IF NOT EXISTS idx_maintenance_vehicle ON maintenance_order(vehicle_id);

CREATE INDEX IF NOT EXISTS idx_contract_dates ON contract(date_from, date_to);
CREATE INDEX IF NOT EXISTS idx_contract_status ON contract(status_id);

CREATE INDEX IF NOT EXISTS idx_rental_pickup_return ON rental_event(pickup_ts, return_ts);
CREATE INDEX IF NOT EXISTS idx_rental_vehicle ON rental_event(vehicle_id);
