package services

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
)

// HumanizeDBError преобразует техническую ошибку БД в понятный текст.
func HumanizeDBError(err error) string {
	if err == nil {
		return ""
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			switch pgErr.ConstraintName {
			case "vehicle_vin_key":
				return "Автомобиль с таким VIN уже существует"
			case "vehicle_reg_number_key":
				return "Автомобиль с таким гос. номером уже существует"
			case "driver_license_number_key":
				return "Водитель с таким номером удостоверения уже существует"
			case "department_code_key":
				return "Подразделение с таким кодом уже существует"
			case "contract_number_key":
				return "Договор с таким номером уже существует"
			default:
				return "Нарушено ограничение уникальности"
			}
		case "23503":
			return "Невозможно выполнить операцию из-за связанных записей"
		case "23514":
			return "Данные не прошли проверку ограничений БД"
		case "23P01":
			return "Периоды назначения пересекаются для выбранного автомобиля"
		case "22001":
			return "Один из текстовых параметров слишком длинный"
		}
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062:
			return "Нарушено ограничение уникальности"
		case 1451, 1452:
			return "Невозможно выполнить операцию из-за связанных записей"
		case 3819:
			return "Данные не прошли проверку ограничений БД"
		case 1406:
			return "Один из текстовых параметров слишком длинный"
		}
	}

	msg := err.Error()
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "duplicate") {
		return "Нарушено ограничение уникальности"
	}
	if strings.Contains(lower, "foreign key") || strings.Contains(lower, "constraint fails") {
		return "Невозможно выполнить операцию из-за связанных записей"
	}
	if strings.Contains(lower, "периоды назначения пересекаются") {
		return "Периоды назначения пересекаются для выбранного автомобиля"
	}

	return "Ошибка базы данных: " + msg
}
