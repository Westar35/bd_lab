package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	dbx "bd_lab_3/internal/db"
	"bd_lab_3/internal/models"
)

// AnalyticsRepository хранит поисковые и отчетные SQL-запросы.
type AnalyticsRepository struct {
	db     *sql.DB
	dbType dbx.DBType
}

func NewAnalyticsRepository(db *sql.DB, dbType dbx.DBType) *AnalyticsRepository {
	return &AnalyticsRepository{db: db, dbType: dbType}
}

func (r *AnalyticsRepository) VehiclesByMake(ctx context.Context, makeValue string) ([]map[string]string, error) {
	query := fmt.Sprintf(`
SELECT
    v.make,
    v.model,
    v.year AS year,
    v.reg_number,
    vs.name AS status_name
FROM vehicle v
JOIN vehicle_status vs ON vs.id = v.status_id
WHERE %s
ORDER BY v.make, v.model, v.year DESC`, r.likeExpr("v.make", 1))

	rows, err := r.db.QueryContext(ctx, query, "%"+strings.ToLower(strings.TrimSpace(makeValue))+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) VehiclesByAttributes(ctx context.Context, makeValue, modelValue, regNumber, vin string) ([]map[string]string, error) {
	args := make([]any, 0, 4)
	where := make([]string, 0, 4)
	if strings.TrimSpace(makeValue) != "" {
		args = append(args, "%"+strings.ToLower(strings.TrimSpace(makeValue))+"%")
		where = append(where, r.likeExpr("v.make", len(args)))
	}
	if strings.TrimSpace(modelValue) != "" {
		args = append(args, strings.ToLower(strings.TrimSpace(modelValue)))
		where = append(where, fmt.Sprintf("LOWER(v.model) = %s", r.placeholder(len(args))))
	}
	if strings.TrimSpace(regNumber) != "" {
		args = append(args, strings.TrimSpace(regNumber))
		where = append(where, fmt.Sprintf("v.reg_number = %s", r.placeholder(len(args))))
	}
	if strings.TrimSpace(vin) != "" {
		args = append(args, strings.TrimSpace(vin))
		where = append(where, fmt.Sprintf("v.vin = %s", r.placeholder(len(args))))
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	query := fmt.Sprintf(`
SELECT
    v.make,
    v.model,
    v.year AS year,
    v.reg_number,
    v.vin,
    vc.name AS class_name,
    ft.name AS fuel_type,
    vs.name AS status_name
FROM vehicle v
JOIN vehicle_status vs ON vs.id = v.status_id
JOIN vehicle_class vc ON vc.id = v.vehicle_class_id
JOIN fuel_type ft ON ft.id = v.fuel_type_id
%s
ORDER BY v.make, v.model, v.year DESC, v.reg_number`, whereSQL)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) TripSheetsByDriverVehiclePeriod(ctx context.Context, driverID, vehicleID *int, startDate, endDate time.Time) ([]map[string]string, error) {
	args := []any{startDate, endDate}
	whereParts := []string{fmt.Sprintf("ts.trip_date BETWEEN %s AND %s", r.placeholder(1), r.placeholder(2))}

	if driverID != nil {
		args = append(args, *driverID)
		whereParts = append(whereParts, fmt.Sprintf("ts.driver_id = %s", r.placeholder(len(args))))
	}

	if vehicleID != nil {
		args = append(args, *vehicleID)
		whereParts = append(whereParts, fmt.Sprintf("ts.vehicle_id = %s", r.placeholder(len(args))))
	}

	query := fmt.Sprintf(`
SELECT
    ts.trip_date AS trip_date,
    d.fio AS driver_fio,
    v.reg_number,
    ts.odometer_start AS odometer_start,
    ts.odometer_end AS odometer_end,
    ts.distance_km AS distance_km
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
    SUM(ts.distance_km) AS total_km
FROM trip_sheet ts
JOIN driver d ON d.id = ts.driver_id
JOIN vehicle v ON v.id = ts.vehicle_id
WHERE ts.trip_date BETWEEN %s AND %s
GROUP BY d.fio, v.make, v.model
ORDER BY d.fio, SUM(ts.distance_km) DESC, v.make, v.model`

	query = fmt.Sprintf(query, r.placeholder(1), r.placeholder(2))

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
    SUM(ts.distance_km) AS total_km,
    COUNT(ts.id) AS trips_count
FROM trip_sheet ts
JOIN vehicle v ON v.id = ts.vehicle_id
WHERE ts.trip_date BETWEEN %s AND %s
GROUP BY v.id, v.reg_number, v.make, v.model
ORDER BY SUM(ts.distance_km) DESC, v.reg_number`

	query = fmt.Sprintf(query, r.placeholder(1), r.placeholder(2))

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
    ROUND(SUM(f.liters), 2) AS total_liters,
    ROUND(SUM(f.amount), 2) AS total_amount
FROM fuel_txn f
JOIN vehicle v ON v.id = f.vehicle_id
WHERE %s BETWEEN %s AND %s
GROUP BY v.id, v.reg_number, v.make, v.model
ORDER BY SUM(f.amount) DESC, v.reg_number`

	query = fmt.Sprintf(query, r.dateExpr("f.txn_ts"), r.placeholder(1), r.placeholder(2))

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
    ROUND(SUM(m.cost), 2) AS total_cost,
    COUNT(m.id) AS orders_count
FROM maintenance_order m
JOIN vehicle v ON v.id = m.vehicle_id
WHERE m.open_date BETWEEN %s AND %s
GROUP BY v.id, v.reg_number, v.make, v.model
ORDER BY SUM(m.cost) DESC, v.reg_number`

	query = fmt.Sprintf(query, r.placeholder(1), r.placeholder(2))

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
    COUNT(v.id) AS vehicles_count
FROM vehicle_status vs
LEFT JOIN vehicle v ON v.status_id = vs.id
GROUP BY vs.id, vs.name
ORDER BY COUNT(v.id) DESC, vs.name`

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
    COUNT(v.id) AS vehicles_count
FROM vehicle_class vc
LEFT JOIN vehicle v ON v.vehicle_class_id = vc.id
GROUP BY vc.id, vc.name
ORDER BY COUNT(v.id) DESC, vc.name`

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
    re.id AS rental_id,
    c.number AS contract_number,
    cp.name AS counterparty,
    v.reg_number,
    v.make,
    v.model,
    re.pickup_ts AS pickup_ts,
    re.price_per_day AS price_per_day,
    COALESCE(%s, '') AS deposit
FROM rental_event re
JOIN contract c ON c.id = re.contract_id
JOIN counterparty cp ON cp.id = c.counterparty_id
JOIN vehicle v ON v.id = re.vehicle_id
WHERE re.return_ts IS NULL
ORDER BY re.pickup_ts DESC`

	query = fmt.Sprintf(query, r.castText("re.deposit"))

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) DatabaseSummary(ctx context.Context) ([]map[string]string, error) {
	query := `
SELECT 'Автомобили' AS entity_name, COUNT(*) AS records_count FROM vehicle
UNION ALL SELECT 'Водители', COUNT(*) FROM driver
UNION ALL SELECT 'Подразделения', COUNT(*) FROM department
UNION ALL SELECT 'Назначения ТС', COUNT(*) FROM vehicle_assignment
UNION ALL SELECT 'Путевые листы', COUNT(*) FROM trip_sheet
UNION ALL SELECT 'Топливные операции', COUNT(*) FROM fuel_txn
UNION ALL SELECT 'Обслуживание/ремонт', COUNT(*) FROM maintenance_order
UNION ALL SELECT 'Контрагенты', COUNT(*) FROM counterparty
UNION ALL SELECT 'Договоры', COUNT(*) FROM contract
UNION ALL SELECT 'События аренды', COUNT(*) FROM rental_event
UNION ALL SELECT 'Аудит-лог', COUNT(*) FROM audit_log`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func (r *AnalyticsRepository) DriverOptions(ctx context.Context) ([]models.Option, error) {
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT %s, fio FROM driver ORDER BY fio`, r.castText("id")))
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
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT %s, %s AS label FROM vehicle ORDER BY make, model, reg_number`, r.castText("id"), r.vehicleLabel("")))
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

func (r *AnalyticsRepository) VehicleOptionsFiltered(ctx context.Context, driverID, departmentID *int) ([]models.Option, error) {
	args := make([]any, 0, 2)
	where := make([]string, 0, 2)
	if driverID != nil {
		args = append(args, *driverID, *driverID)
		where = append(where, fmt.Sprintf(`(
EXISTS (SELECT 1 FROM vehicle_assignment va WHERE va.vehicle_id = v.id AND va.driver_id = %s)
OR EXISTS (SELECT 1 FROM trip_sheet ts WHERE ts.vehicle_id = v.id AND ts.driver_id = %s)
)`, r.placeholder(len(args)-1), r.placeholder(len(args))))
	}
	if departmentID != nil {
		args = append(args, *departmentID, *departmentID)
		where = append(where, fmt.Sprintf(`(
EXISTS (SELECT 1 FROM vehicle_assignment va WHERE va.vehicle_id = v.id AND va.department_id = %s)
OR EXISTS (SELECT 1 FROM trip_sheet ts WHERE ts.vehicle_id = v.id AND ts.department_id = %s)
)`, r.placeholder(len(args)-1), r.placeholder(len(args))))
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	query := fmt.Sprintf(`
SELECT DISTINCT %s AS id, %s AS label
FROM vehicle v
%s
ORDER BY label`, r.castText("v.id"), r.vehicleLabel("v"), whereSQL)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) DriverOptionsFiltered(ctx context.Context, vehicleID, departmentID *int) ([]models.Option, error) {
	args := make([]any, 0, 2)
	where := make([]string, 0, 2)
	if vehicleID != nil {
		args = append(args, *vehicleID, *vehicleID)
		where = append(where, fmt.Sprintf(`(
EXISTS (SELECT 1 FROM vehicle_assignment va WHERE va.driver_id = d.id AND va.vehicle_id = %s)
OR EXISTS (SELECT 1 FROM trip_sheet ts WHERE ts.driver_id = d.id AND ts.vehicle_id = %s)
)`, r.placeholder(len(args)-1), r.placeholder(len(args))))
	}
	if departmentID != nil {
		args = append(args, *departmentID, *departmentID)
		where = append(where, fmt.Sprintf(`(
EXISTS (SELECT 1 FROM vehicle_assignment va WHERE va.driver_id = d.id AND va.department_id = %s)
OR EXISTS (SELECT 1 FROM trip_sheet ts WHERE ts.driver_id = d.id AND ts.department_id = %s)
)`, r.placeholder(len(args)-1), r.placeholder(len(args))))
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	query := fmt.Sprintf(`
SELECT DISTINCT %s AS id, d.fio AS label
FROM driver d
%s
ORDER BY d.fio`, r.castText("d.id"), whereSQL)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) ModelOptions(ctx context.Context, makeValue string) ([]models.Option, error) {
	args := make([]any, 0, 1)
	whereSQL := ""
	if strings.TrimSpace(makeValue) != "" {
		args = append(args, strings.ToLower(strings.TrimSpace(makeValue)))
		whereSQL = fmt.Sprintf("WHERE LOWER(make) = %s", r.placeholder(1))
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT DISTINCT model AS id, model AS label
FROM vehicle
%s
ORDER BY model`, whereSQL), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) RegNumberOptions(ctx context.Context, makeValue, modelValue string) ([]models.Option, error) {
	args := make([]any, 0, 2)
	where := make([]string, 0, 2)
	if strings.TrimSpace(makeValue) != "" {
		args = append(args, strings.ToLower(strings.TrimSpace(makeValue)))
		where = append(where, fmt.Sprintf("LOWER(make) = %s", r.placeholder(len(args))))
	}
	if strings.TrimSpace(modelValue) != "" {
		args = append(args, strings.ToLower(strings.TrimSpace(modelValue)))
		where = append(where, fmt.Sprintf("LOWER(model) = %s", r.placeholder(len(args))))
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT reg_number AS id, %s AS label
FROM vehicle
%s
ORDER BY reg_number`, r.concat("reg_number", "' - '", "vin"), whereSQL), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) VINOptions(ctx context.Context, regNumber string) ([]models.Option, error) {
	args := make([]any, 0, 1)
	whereSQL := ""
	if strings.TrimSpace(regNumber) != "" {
		args = append(args, strings.TrimSpace(regNumber))
		whereSQL = fmt.Sprintf("WHERE reg_number = %s", r.placeholder(1))
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT vin AS id, vin AS label
FROM vehicle
%s
ORDER BY vin`, whereSQL), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) ContractOptionsFiltered(ctx context.Context, counterpartyID *int) ([]models.Option, error) {
	args := make([]any, 0, 1)
	whereSQL := ""
	if counterpartyID != nil {
		args = append(args, *counterpartyID)
		whereSQL = fmt.Sprintf("WHERE c.counterparty_id = %s", r.placeholder(1))
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT %s AS id, %s AS label
FROM contract c
JOIN counterparty cp ON cp.id = c.counterparty_id
%s
ORDER BY c.number`, r.castText("c.id"), r.concat("c.number", "' - '", "cp.name"), whereSQL), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) RentalEventOptionsFiltered(ctx context.Context, contractID *int) ([]models.Option, error) {
	args := make([]any, 0, 1)
	whereSQL := ""
	if contractID != nil {
		args = append(args, *contractID)
		whereSQL = fmt.Sprintf("WHERE re.contract_id = %s", r.placeholder(1))
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT %s AS id, %s AS label
FROM rental_event re
JOIN contract c ON c.id = re.contract_id
JOIN vehicle v ON v.id = re.vehicle_id
%s
ORDER BY re.pickup_ts DESC`, r.castText("re.id"), r.concat("'Аренда #'", r.castText("re.id"), "' - '", "c.number", "' - '", r.vehicleLabel("v")), whereSQL), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) TripSheetOptionsFiltered(ctx context.Context, vehicleID, driverID *int) ([]models.Option, error) {
	args := make([]any, 0, 2)
	where := make([]string, 0, 2)
	if vehicleID != nil {
		args = append(args, *vehicleID)
		where = append(where, fmt.Sprintf("ts.vehicle_id = %s", r.placeholder(len(args))))
	}
	if driverID != nil {
		args = append(args, *driverID)
		where = append(where, fmt.Sprintf("ts.driver_id = %s", r.placeholder(len(args))))
	}
	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT %s AS id, %s AS label
FROM trip_sheet ts
JOIN driver d ON d.id = ts.driver_id
JOIN vehicle v ON v.id = ts.vehicle_id
%s
ORDER BY ts.trip_date DESC, ts.id DESC
LIMIT 200`, r.castText("ts.id"), r.concat(r.castText("ts.trip_date"), "' - '", "d.fio", "' - '", "v.reg_number"), whereSQL), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) FuelTransactionOptionsFiltered(ctx context.Context, vehicleID *int) ([]models.Option, error) {
	args := make([]any, 0, 1)
	whereSQL := ""
	if vehicleID != nil {
		args = append(args, *vehicleID)
		whereSQL = fmt.Sprintf("WHERE f.vehicle_id = %s", r.placeholder(1))
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT %s AS id, %s AS label
FROM fuel_txn f
JOIN vehicle v ON v.id = f.vehicle_id
%s
ORDER BY f.txn_ts DESC, f.id DESC
LIMIT 200`, r.castText("f.id"), r.concat(r.castText("f.txn_ts"), "' - '", "v.reg_number", "' - '", "f.station"), whereSQL), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) MaintenanceOrderOptionsFiltered(ctx context.Context, vehicleID *int) ([]models.Option, error) {
	args := make([]any, 0, 1)
	whereSQL := ""
	if vehicleID != nil {
		args = append(args, *vehicleID)
		whereSQL = fmt.Sprintf("WHERE m.vehicle_id = %s", r.placeholder(1))
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT %s AS id, %s AS label
FROM maintenance_order m
JOIN vehicle v ON v.id = m.vehicle_id
JOIN maintenance_type mt ON mt.id = m.maintenance_type_id
%s
ORDER BY m.open_date DESC, m.id DESC
LIMIT 200`, r.castText("m.id"), r.concat(r.castText("m.open_date"), "' - '", "v.reg_number", "' - '", "mt.name"), whereSQL), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOptions(rows)
}

func (r *AnalyticsRepository) placeholder(n int) string {
	return r.dbType.Placeholder(n)
}

func (r *AnalyticsRepository) likeExpr(expr string, arg int) string {
	if r.dbType == dbx.DBMySQL {
		return fmt.Sprintf("LOWER(CAST(%s AS CHAR)) LIKE %s", expr, r.placeholder(arg))
	}
	return fmt.Sprintf("LOWER(%s::text) LIKE %s", expr, r.placeholder(arg))
}

func (r *AnalyticsRepository) dateExpr(expr string) string {
	if r.dbType == dbx.DBMySQL {
		return fmt.Sprintf("DATE(%s)", expr)
	}
	return fmt.Sprintf("%s::date", expr)
}

func (r *AnalyticsRepository) castText(expr string) string {
	if r.dbType == dbx.DBMySQL {
		return fmt.Sprintf("CAST(%s AS CHAR)", expr)
	}
	return fmt.Sprintf("%s::text", expr)
}

func (r *AnalyticsRepository) vehicleLabel(alias string) string {
	prefix := ""
	if alias != "" {
		prefix = alias + "."
	}
	if r.dbType == dbx.DBMySQL {
		return fmt.Sprintf("CONCAT(%smake, ' ', %smodel, ' (', %sreg_number, ')')", prefix, prefix, prefix)
	}
	return fmt.Sprintf("%smake || ' ' || %smodel || ' (' || %sreg_number || ')'", prefix, prefix, prefix)
}

func (r *AnalyticsRepository) concat(parts ...string) string {
	if r.dbType == dbx.DBMySQL {
		return "CONCAT(" + strings.Join(parts, ", ") + ")"
	}
	return strings.Join(parts, " || ")
}

func scanOptions(rows *sql.Rows) ([]models.Option, error) {
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
