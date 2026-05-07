package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	dbx "bd_lab_3/internal/db"
	"bd_lab_3/internal/models"
)

var ErrAssignmentOverlap = errors.New("периоды назначения пересекаются для выбранного автомобиля")

// ListParams параметры выборки списка.
type ListParams struct {
	Page    int
	PerPage int
	Sort    string
	Dir     string
	Q       string
	Filters map[string]string
}

// ListResult результат выдачи списка.
type ListResult struct {
	Rows       []map[string]string
	Total      int
	TotalPages int
}

// EntityRepository репозиторий универсального CRUD.
type EntityRepository struct {
	db     *sql.DB
	dbType dbx.DBType
}

func NewEntityRepository(db *sql.DB, dbType dbx.DBType) *EntityRepository {
	return &EntityRepository{db: db, dbType: dbType}
}

// List возвращает страницу списка сущности.
func (r *EntityRepository) List(ctx context.Context, entity models.EntityConfig, params ListParams) (ListResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage <= 0 || params.PerPage > 100 {
		params.PerPage = 20
	}

	sortField := r.resolveSortField(entity, params.Sort)
	dir := strings.ToUpper(params.Dir)
	if dir != "ASC" && dir != "DESC" {
		dir = strings.ToUpper(entity.DefaultDir)
		if dir != "ASC" && dir != "DESC" {
			dir = "DESC"
		}
	}

	selectParts := []string{r.textExpr("t", "id", "id")}
	listFieldSet := make(map[string]struct{}, len(entity.ListColumns))
	for _, col := range entity.ListColumns {
		listFieldSet[col] = struct{}{}
	}

	for _, f := range entity.Fields {
		if _, ok := listFieldSet[f.Name]; !ok {
			continue
		}
		selectParts = append(selectParts, r.displayExpr(f, "t"))
	}

	whereSQL, args := r.buildWhereClause(entity, params)

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM %s t %s", entity.Table, whereSQL)
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return ListResult{}, err
	}

	offset := (params.Page - 1) * params.PerPage
	argsWithPaging := append([]any{}, args...)
	argsWithPaging = append(argsWithPaging, params.PerPage, offset)

	query := fmt.Sprintf(
		"SELECT %s FROM %s t %s ORDER BY t.%s %s LIMIT %s OFFSET %s",
		strings.Join(selectParts, ", "),
		entity.Table,
		whereSQL,
		sortField,
		dir,
		r.placeholder(len(args)+1),
		r.placeholder(len(args)+2),
	)

	rows, err := r.db.QueryContext(ctx, query, argsWithPaging...)
	if err != nil {
		return ListResult{}, err
	}
	defer rows.Close()

	resultRows, err := scanRowsToMaps(rows)
	if err != nil {
		return ListResult{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.PerPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	return ListResult{Rows: resultRows, Total: total, TotalPages: totalPages}, nil
}

// GetByID возвращает карточку сущности в отображаемом виде.
func (r *EntityRepository) GetByID(ctx context.Context, entity models.EntityConfig, id int64) (map[string]string, error) {
	selectParts := make([]string, 0, len(entity.Fields)+1)
	selectParts = append(selectParts, r.textExpr("t", "id", "id"))

	for _, f := range entity.DetailFields() {
		if f.Name == "id" {
			continue
		}
		selectParts = append(selectParts, r.displayExpr(f, "t"))
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s t WHERE t.id = %s",
		strings.Join(selectParts, ", "),
		entity.Table,
		r.placeholder(1),
	)

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapped, err := scanRowsToMaps(rows)
	if err != nil {
		return nil, err
	}
	if len(mapped) == 0 {
		return nil, sql.ErrNoRows
	}

	return mapped[0], nil
}

// GetRawByID возвращает запись в сыром виде для формы редактирования.
func (r *EntityRepository) GetRawByID(ctx context.Context, entity models.EntityConfig, id int64) (map[string]string, error) {
	formFields := entity.FormFields()
	selectParts := make([]string, 0, len(formFields))
	for _, f := range formFields {
		selectParts = append(selectParts, r.textExpr("t", f.Name, f.Name))
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s t WHERE t.id = %s",
		strings.Join(selectParts, ", "),
		entity.Table,
		r.placeholder(1),
	)

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapped, err := scanRowsToMaps(rows)
	if err != nil {
		return nil, err
	}
	if len(mapped) == 0 {
		return nil, sql.ErrNoRows
	}

	return mapped[0], nil
}

// Create создает запись и пишет аудит.
func (r *EntityRepository) Create(ctx context.Context, entity models.EntityConfig, values map[string]any, actor string) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	if err = r.validateAssignmentOverlap(ctx, tx, entity, values, 0); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	cols := make([]string, 0)
	placeholders := make([]string, 0)
	args := make([]any, 0)

	for _, f := range entity.FormFields() {
		cols = append(cols, f.Name)
		placeholders = append(placeholders, r.placeholder(len(args)+1))
		args = append(args, values[f.Name])
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		entity.Table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	var id int64
	if r.dbType == dbx.DBPostgres {
		if err = tx.QueryRowContext(ctx, query+" RETURNING id", args...).Scan(&id); err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	} else {
		result, execErr := tx.ExecContext(ctx, query, args...)
		if execErr != nil {
			_ = tx.Rollback()
			return 0, execErr
		}
		id, err = result.LastInsertId()
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	after, err := fetchRowJSON(ctx, tx, r.dbType, entity.Table, id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	if err = insertAudit(ctx, tx, r.dbType, actor, "INSERT", entity.Table, id, nil, after); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

// Update обновляет запись и пишет аудит.
func (r *EntityRepository) Update(ctx context.Context, entity models.EntityConfig, id int64, values map[string]any, actor string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = r.validateAssignmentOverlap(ctx, tx, entity, values, id); err != nil {
		_ = tx.Rollback()
		return err
	}

	before, err := fetchRowJSON(ctx, tx, r.dbType, entity.Table, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	sets := make([]string, 0)
	args := make([]any, 0)

	for _, f := range entity.FormFields() {
		sets = append(sets, fmt.Sprintf("%s = %s", f.Name, r.placeholder(len(args)+1)))
		args = append(args, values[f.Name])
	}
	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = %s",
		entity.Table,
		strings.Join(sets, ", "),
		r.placeholder(len(args)),
	)

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if affected == 0 {
		_ = tx.Rollback()
		return sql.ErrNoRows
	}

	after, err := fetchRowJSON(ctx, tx, r.dbType, entity.Table, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = insertAudit(ctx, tx, r.dbType, actor, "UPDATE", entity.Table, id, before, after); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Delete удаляет запись и пишет аудит.
func (r *EntityRepository) Delete(ctx context.Context, entity models.EntityConfig, id int64, actor string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	before, err := fetchRowJSON(ctx, tx, r.dbType, entity.Table, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	result, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = %s", entity.Table, r.placeholder(1)), id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if affected == 0 {
		_ = tx.Rollback()
		return sql.ErrNoRows
	}

	if err = insertAudit(ctx, tx, r.dbType, actor, "DELETE", entity.Table, id, before, nil); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// SelectOptions возвращает значения для select-поля.
func (r *EntityRepository) SelectOptions(ctx context.Context, field models.Field) ([]models.Option, error) {
	if len(field.StaticOptions) > 0 {
		return field.StaticOptions, nil
	}

	if field.RefTable == "" || field.RefLabelExpr == "" {
		return nil, nil
	}

	query := fmt.Sprintf(
		"SELECT id::text AS id, %s AS label FROM %s ORDER BY label",
		r.refLabelExpr(field.RefTable, ""),
		field.RefTable,
	)
	if r.dbType == dbx.DBMySQL {
		query = fmt.Sprintf(
			"SELECT CAST(id AS CHAR) AS id, %s AS label FROM %s ORDER BY label",
			r.refLabelExpr(field.RefTable, ""),
			field.RefTable,
		)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := make([]models.Option, 0)
	for rows.Next() {
		var opt models.Option
		if err := rows.Scan(&opt.ID, &opt.Label); err != nil {
			return nil, err
		}
		options = append(options, opt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return options, nil
}

func (r *EntityRepository) validateAssignmentOverlap(ctx context.Context, tx *sql.Tx, entity models.EntityConfig, values map[string]any, currentID int64) error {
	if entity.Table != "vehicle_assignment" {
		return nil
	}

	vehicleID, ok := values["vehicle_id"].(int)
	if !ok || vehicleID <= 0 {
		return nil
	}
	dateFrom, ok := values["date_from"].(time.Time)
	if !ok {
		return nil
	}
	dateTo := values["date_to"]

	query := fmt.Sprintf(`
SELECT COUNT(*)
FROM vehicle_assignment
WHERE vehicle_id = %s
  AND id <> %s
  AND date_from <= COALESCE(%s::date, DATE '9999-12-31')
  AND COALESCE(date_to, DATE '9999-12-31') >= %s`,
		r.placeholder(1),
		r.placeholder(2),
		r.placeholder(3),
		r.placeholder(4),
	)
	if r.dbType == dbx.DBMySQL {
		query = fmt.Sprintf(`
SELECT COUNT(*)
FROM vehicle_assignment
WHERE vehicle_id = %s
  AND id <> %s
  AND date_from <= COALESCE(%s, DATE('9999-12-31'))
  AND COALESCE(date_to, DATE('9999-12-31')) >= %s`,
			r.placeholder(1),
			r.placeholder(2),
			r.placeholder(3),
			r.placeholder(4),
		)
	}

	var count int
	if err := tx.QueryRowContext(ctx, query, vehicleID, currentID, dateTo, dateFrom).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return ErrAssignmentOverlap
	}
	return nil
}

func (r *EntityRepository) resolveSortField(entity models.EntityConfig, sort string) string {
	if sort == "" {
		return entity.DefaultSort
	}
	if sort == "id" {
		return "id"
	}
	for _, f := range entity.Fields {
		if f.Name == sort && f.Sortable {
			return sort
		}
	}
	return entity.DefaultSort
}

func (r *EntityRepository) displayExpr(field models.Field, alias string) string {
	safeAlias := field.Name

	switch field.Type {
	case "select":
		if len(field.StaticOptions) > 0 {
			cases := make([]string, 0, len(field.StaticOptions))
			for _, opt := range field.StaticOptions {
				cases = append(cases, fmt.Sprintf("WHEN %s.%s = '%s' THEN '%s'", alias, field.Name, opt.ID, opt.Label))
			}
			return fmt.Sprintf("COALESCE((CASE %s ELSE %s END), '') AS %s", strings.Join(cases, " "), r.castExpr(alias, field.Name), safeAlias)
		}
		if field.RefTable != "" && field.RefLabelExpr != "" {
			if r.dbType == dbx.DBMySQL {
				return fmt.Sprintf(
					"COALESCE((SELECT %s FROM %s ref WHERE ref.id = %s.%s), '') AS %s",
					r.refLabelExpr(field.RefTable, "ref"),
					field.RefTable,
					alias,
					field.Name,
					safeAlias,
				)
			}
			return fmt.Sprintf(
				"COALESCE((SELECT %s FROM %s ref WHERE ref.id = %s.%s), '')::text AS %s",
				r.refLabelExpr(field.RefTable, "ref"),
				field.RefTable,
				alias,
				field.Name,
				safeAlias,
			)
		}
		return r.textExpr(alias, field.Name, safeAlias)
	case "checkbox":
		return fmt.Sprintf("CASE WHEN %s.%s THEN 'Да' ELSE 'Нет' END AS %s", alias, field.Name, safeAlias)
	default:
		return r.textExpr(alias, field.Name, safeAlias)
	}
}

func (r *EntityRepository) buildWhereClause(entity models.EntityConfig, params ListParams) (string, []any) {
	clauses := make([]string, 0)
	args := make([]any, 0)

	if strings.TrimSpace(params.Q) != "" {
		sub := make([]string, 0, len(entity.SearchColumns)+4)
		queryText := strings.TrimSpace(params.Q)
		for _, col := range entity.SearchColumns {
			args = append(args, "%"+strings.ToLower(queryText)+"%")
			sub = append(sub, r.likeExpr(fmt.Sprintf("t.%s", col), len(args)))
		}
		for _, field := range entity.Fields {
			if field.Type != "select" || field.RefTable == "" || field.RefLabelExpr == "" {
				continue
			}
			args = append(args, "%"+strings.ToLower(queryText)+"%")
			sub = append(sub, fmt.Sprintf(
				"EXISTS (SELECT 1 FROM %s ref WHERE ref.id = t.%s AND %s)",
				field.RefTable,
				field.Name,
				r.likeRawExpr(r.refLabelExpr(field.RefTable, "ref"), len(args)),
			))
		}
		if len(sub) == 0 {
			clauses = append(clauses, "1=0")
		} else {
			clauses = append(clauses, "("+strings.Join(sub, " OR ")+")")
		}
	}

	for key, val := range params.Filters {
		if strings.TrimSpace(val) == "" {
			continue
		}

		field, ok := entity.FieldByName(key)
		if !ok || !field.Filterable {
			continue
		}

		switch field.Type {
		case "checkbox":
			boolVal := strings.EqualFold(val, "true") || val == "1"
			args = append(args, boolVal)
			clauses = append(clauses, fmt.Sprintf("t.%s = %s", key, r.placeholder(len(args))))
		default:
			if (field.Type == "select" && len(field.StaticOptions) == 0) || field.Type == "int" {
				intVal, err := strconv.Atoi(val)
				if err != nil {
					continue
				}
				args = append(args, intVal)
				clauses = append(clauses, fmt.Sprintf("t.%s = %s", key, r.placeholder(len(args))))
				continue
			}
			args = append(args, "%"+strings.ToLower(val)+"%")
			clauses = append(clauses, r.likeExpr(fmt.Sprintf("t.%s", key), len(args)))
		}
	}

	if len(clauses) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(clauses, " AND "), args
}

func scanRowsToMaps(rows *sql.Rows) ([]map[string]string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]string, 0)
	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]string, len(cols))
		for i, c := range cols {
			rowMap[c] = anyToString(values[i])
		}
		result = append(result, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func anyToString(v any) string {
	switch val := v.(type) {
	case nil:
		return ""
	case []byte:
		return string(val)
	case string:
		return val
	case time.Time:
		return val.Format(time.RFC3339)
	case bool:
		if val {
			return "Да"
		}
		return "Нет"
	default:
		return fmt.Sprint(val)
	}
}

func fetchRowJSON(ctx context.Context, tx *sql.Tx, dbType dbx.DBType, table string, id int64) (json.RawMessage, error) {
	placeholder := dbType.Placeholder(1)
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id = %s", table, placeholder), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapped, err := scanRowsToMaps(rows)
	if err != nil {
		return nil, err
	}
	if len(mapped) == 0 {
		return nil, sql.ErrNoRows
	}
	return json.Marshal(mapped[0])
}

func insertAudit(ctx context.Context, tx *sql.Tx, dbType dbx.DBType, actor, action, entity string, entityID int64, before, after json.RawMessage) error {
	query := `INSERT INTO audit_log (ts, actor, action, entity, entity_id, details_before, details_after)
         VALUES (now(), $1, $2, $3, $4, $5::jsonb, $6::jsonb)`
	if dbType == dbx.DBMySQL {
		query = `INSERT INTO audit_log (ts, actor, action, entity, entity_id, details_before, details_after)
         VALUES (CURRENT_TIMESTAMP, ?, ?, ?, ?, ?, ?)`
	}
	_, err := tx.ExecContext(ctx,
		query,
		actor,
		action,
		entity,
		entityID,
		nullableJSON(before),
		nullableJSON(after),
	)
	return err
}

func nullableJSON(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	return string(raw)
}

func (r *EntityRepository) placeholder(n int) string {
	return r.dbType.Placeholder(n)
}

func (r *EntityRepository) textExpr(alias, column, asName string) string {
	return fmt.Sprintf("COALESCE(%s, '') AS %s", r.castExpr(alias, column), asName)
}

func (r *EntityRepository) castExpr(alias, column string) string {
	if r.dbType == dbx.DBMySQL {
		return fmt.Sprintf("CAST(%s.%s AS CHAR)", alias, column)
	}
	return fmt.Sprintf("%s.%s::text", alias, column)
}

func (r *EntityRepository) likeExpr(expr string, arg int) string {
	return r.likeRawExpr(expr, arg)
}

func (r *EntityRepository) likeRawExpr(expr string, arg int) string {
	if r.dbType == dbx.DBMySQL {
		return fmt.Sprintf("LOWER(CAST(%s AS CHAR)) LIKE %s", expr, r.placeholder(arg))
	}
	return fmt.Sprintf("LOWER(%s::text) LIKE %s", expr, r.placeholder(arg))
}

func (r *EntityRepository) refLabelExpr(table, alias string) string {
	prefix := ""
	if alias != "" {
		prefix = alias + "."
	}
	if r.dbType == dbx.DBMySQL {
		switch table {
		case "vehicle":
			return fmt.Sprintf("CONCAT(%smake, ' ', %smodel, ' (', %sreg_number, ')')", prefix, prefix, prefix)
		case "driver":
			return prefix + "fio"
		case "department":
			return prefix + "name"
		case "contract":
			return prefix + "number"
		case "counterparty":
			return prefix + "name"
		default:
			return prefix + "name"
		}
	}
	switch table {
	case "vehicle":
		return fmt.Sprintf("%smake || ' ' || %smodel || ' (' || %sreg_number || ')'", prefix, prefix, prefix)
	case "driver":
		return prefix + "fio"
	case "department":
		return prefix + "name"
	case "contract":
		return prefix + "number"
	case "counterparty":
		return prefix + "name"
	default:
		return prefix + "name"
	}
}
