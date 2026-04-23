package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"bd_lab_3/internal/models"
)

// AnalyticsRepository хранит поисковые и отчетные SQL-запросы.
type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) VehiclesByMake(ctx context.Context, makeValue string) ([]map[string]string, error) {
	query := `
SELECT
    v.make,
    v.model,
    v.year::text AS year,
    v.reg_number,
    vs.name AS status_name
FROM vehicle v
JOIN vehicle_status vs ON vs.id = v.status_id
WHERE v.make ILIKE $1
ORDER BY v.make, v.model, v.year DESC`

	rows, err := r.db.QueryContext(ctx, query, "%"+strings.TrimSpace(makeValue)+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) TripSheetsByDriverVehiclePeriod(ctx context.Context, driverID *int, regNumber string, startDate, endDate time.Time) ([]map[string]string, error) {
	args := []any{startDate, endDate}
	whereParts := []string{"ts.trip_date BETWEEN $1 AND $2"}

	if driverID != nil {
		args = append(args, *driverID)
		whereParts = append(whereParts, fmt.Sprintf("ts.driver_id = $%d", len(args)))
	}

	if strings.TrimSpace(regNumber) != "" {
		args = append(args, "%"+strings.TrimSpace(regNumber)+"%")
		whereParts = append(whereParts, fmt.Sprintf("v.reg_number ILIKE $%d", len(args)))
	}

	query := fmt.Sprintf(`
SELECT
    ts.trip_date::text AS trip_date,
    d.fio AS driver_fio,
    v.reg_number,
    ts.odometer_start::text AS odometer_start,
    ts.odometer_end::text AS odometer_end,
    ts.distance_km::text AS distance_km
FROM trip_sheet ts
JOIN driver d ON d.id = ts.driver_id
JOIN vehicle v ON v.id = ts.vehicle_id
WHERE %s
ORDER BY ts.trip_date DESC, d.fio, v.reg_number`, strings.Join(whereParts, " AND "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) TotalMileageByDriverVehicle(ctx context.Context, startDate, endDate time.Time) ([]map[string]string, error) {
	query := `
SELECT
    d.fio AS driver_fio,
    v.make,
    v.model,
    SUM(ts.distance_km)::text AS total_km
FROM trip_sheet ts
JOIN driver d ON d.id = ts.driver_id
JOIN vehicle v ON v.id = ts.vehicle_id
WHERE ts.trip_date BETWEEN $1 AND $2
GROUP BY d.fio, v.make, v.model
ORDER BY d.fio, total_km::numeric DESC`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) MileageReport(ctx context.Context, startDate, endDate time.Time) ([]map[string]string, error) {
	query := `
SELECT
    v.reg_number,
    v.make,
    v.model,
    SUM(ts.distance_km)::text AS total_km,
    COUNT(ts.id)::text AS trips_count
FROM trip_sheet ts
JOIN vehicle v ON v.id = ts.vehicle_id
WHERE ts.trip_date BETWEEN $1 AND $2
GROUP BY v.id, v.reg_number, v.make, v.model
ORDER BY total_km::numeric DESC`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) FuelExpensesReport(ctx context.Context, startDate, endDate time.Time) ([]map[string]string, error) {
	query := `
SELECT
    v.reg_number,
    v.make,
    v.model,
    ROUND(SUM(f.liters)::numeric, 2)::text AS total_liters,
    ROUND(SUM(f.amount)::numeric, 2)::text AS total_amount
FROM fuel_txn f
JOIN vehicle v ON v.id = f.vehicle_id
WHERE f.txn_ts::date BETWEEN $1 AND $2
GROUP BY v.id, v.reg_number, v.make, v.model
ORDER BY total_amount::numeric DESC`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) MaintenanceExpensesReport(ctx context.Context, startDate, endDate time.Time) ([]map[string]string, error) {
	query := `
SELECT
    v.reg_number,
    v.make,
    v.model,
    ROUND(SUM(m.cost)::numeric, 2)::text AS total_cost,
    COUNT(m.id)::text AS orders_count
FROM maintenance_order m
JOIN vehicle v ON v.id = m.vehicle_id
WHERE m.open_date BETWEEN $1 AND $2
GROUP BY v.id, v.reg_number, v.make, v.model
ORDER BY total_cost::numeric DESC`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) VehiclesByStatus(ctx context.Context) ([]map[string]string, error) {
	query := `
SELECT
    vs.name AS status_name,
    COUNT(v.id)::text AS vehicles_count
FROM vehicle_status vs
LEFT JOIN vehicle v ON v.status_id = vs.id
GROUP BY vs.id, vs.name
ORDER BY vehicles_count::int DESC, vs.name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) VehiclesByClass(ctx context.Context) ([]map[string]string, error) {
	query := `
SELECT
    vc.name AS class_name,
    COUNT(v.id)::text AS vehicles_count
FROM vehicle_class vc
LEFT JOIN vehicle v ON v.vehicle_class_id = vc.id
GROUP BY vc.id, vc.name
ORDER BY vehicles_count::int DESC, vc.name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) CurrentRentals(ctx context.Context) ([]map[string]string, error) {
	query := `
SELECT
    re.id::text AS rental_id,
    c.number AS contract_number,
    cp.name AS counterparty,
    v.reg_number,
    v.make,
    v.model,
    re.pickup_ts::text AS pickup_ts,
    re.price_per_day::text AS price_per_day,
    COALESCE(re.deposit::text, '') AS deposit
FROM rental_event re
JOIN contract c ON c.id = re.contract_id
JOIN counterparty cp ON cp.id = c.counterparty_id
JOIN vehicle v ON v.id = re.vehicle_id
WHERE re.return_ts IS NULL
ORDER BY re.pickup_ts DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) DatabaseSummary(ctx context.Context) ([]map[string]string, error) {
	query := `
SELECT 'Автомобили' AS entity_name, COUNT(*)::text AS records_count FROM vehicle
UNION ALL SELECT 'Водители', COUNT(*)::text FROM driver
UNION ALL SELECT 'Подразделения', COUNT(*)::text FROM department
UNION ALL SELECT 'Назначения ТС', COUNT(*)::text FROM vehicle_assignment
UNION ALL SELECT 'Путевые листы', COUNT(*)::text FROM trip_sheet
UNION ALL SELECT 'Топливные операции', COUNT(*)::text FROM fuel_txn
UNION ALL SELECT 'Обслуживание/ремонт', COUNT(*)::text FROM maintenance_order
UNION ALL SELECT 'Контрагенты', COUNT(*)::text FROM counterparty
UNION ALL SELECT 'Договоры', COUNT(*)::text FROM contract
UNION ALL SELECT 'События аренды', COUNT(*)::text FROM rental_event
UNION ALL SELECT 'Аудит-лог', COUNT(*)::text FROM audit_log`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) DriverOptions(ctx context.Context) ([]models.Option, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id::text, fio FROM driver ORDER BY fio`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := make([]models.Option, 0)
	for rows.Next() {
		var o models.Option
		if err := rows.Scan(&o.ID, &o.Label); err != nil {
			return nil, err
		}
		options = append(options, o)
	}
	return options, rows.Err()
}

func (r *AnalyticsRepository) VehicleOptions(ctx context.Context) ([]models.Option, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id::text, make || ' ' || model || ' (' || reg_number || ')' AS label FROM vehicle ORDER BY make, model, reg_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := make([]models.Option, 0)
	for rows.Next() {
		var o models.Option
		if err := rows.Scan(&o.ID, &o.Label); err != nil {
			return nil, err
		}
		options = append(options, o)
	}
	return options, rows.Err()
}
