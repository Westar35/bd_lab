package services

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bd_lab_3/internal/models"
)

const (
	dateLayout          = "2006-01-02"
	dateTimeLocalLayout = "2006-01-02T15:04"
)

// ParseAndValidateForm парсит и валидирует поля формы сущности.
func ParseAndValidateForm(entity models.EntityConfig, form url.Values) (map[string]any, map[string]string, map[string]string) {
	parsed := make(map[string]any)
	raw := make(map[string]string)
	errs := make(map[string]string)

	for _, field := range entity.FormFields() {
		input := strings.TrimSpace(form.Get(field.Name))

		if field.Type == "checkbox" {
			if form.Get(field.Name) != "" {
				input = "true"
			} else {
				input = ""
			}
		}

		raw[field.Name] = input

		if field.Required && input == "" {
			errs[field.Name] = "Поле обязательно для заполнения"
			continue
		}

		if input == "" {
			if field.Nullable || field.Type == "checkbox" || !field.Required {
				if field.Type == "checkbox" {
					parsed[field.Name] = false
				} else if field.Type == "text" || field.Type == "textarea" || field.Type == "email" {
					parsed[field.Name] = nil
				} else if field.Type == "select" || field.Type == "date" || field.Type == "datetime-local" || field.Type == "decimal" || field.Type == "int" {
					parsed[field.Name] = nil
				} else {
					parsed[field.Name] = nil
				}
				continue
			}
		}

		switch field.Type {
		case "text", "textarea", "email":
			if field.MaxLength > 0 && len([]rune(input)) > field.MaxLength {
				errs[field.Name] = fmt.Sprintf("Максимальная длина поля: %d символов", field.MaxLength)
				continue
			}
			if field.Type == "email" && !strings.Contains(input, "@") {
				errs[field.Name] = "Введите корректный email"
				continue
			}
			parsed[field.Name] = input
		case "int":
			n, err := strconv.Atoi(input)
			if err != nil {
				errs[field.Name] = "Введите целое число"
				continue
			}
			if field.Min != "" {
				min, _ := strconv.Atoi(field.Min)
				if n < min {
					errs[field.Name] = fmt.Sprintf("Значение не может быть меньше %d", min)
					continue
				}
			}
			if field.Max != "" {
				max, _ := strconv.Atoi(field.Max)
				if n > max {
					errs[field.Name] = fmt.Sprintf("Значение не может быть больше %d", max)
					continue
				}
			}
			parsed[field.Name] = n
		case "decimal":
			f, err := strconv.ParseFloat(strings.ReplaceAll(input, ",", "."), 64)
			if err != nil {
				errs[field.Name] = "Введите корректное число"
				continue
			}
			if field.Min != "" {
				min, _ := strconv.ParseFloat(field.Min, 64)
				if f < min {
					errs[field.Name] = fmt.Sprintf("Значение не может быть меньше %s", field.Min)
					continue
				}
			}
			if field.Max != "" {
				max, _ := strconv.ParseFloat(field.Max, 64)
				if f > max {
					errs[field.Name] = fmt.Sprintf("Значение не может быть больше %s", field.Max)
					continue
				}
			}
			parsed[field.Name] = f
		case "select":
			if len(field.StaticOptions) > 0 {
				parsed[field.Name] = input
				continue
			}
			n, err := strconv.Atoi(input)
			if err != nil {
				errs[field.Name] = "Выберите значение из списка"
				continue
			}
			parsed[field.Name] = n
		case "date":
			t, err := time.Parse(dateLayout, input)
			if err != nil {
				errs[field.Name] = "Некорректная дата"
				continue
			}
			parsed[field.Name] = t
		case "datetime-local":
			t, err := time.Parse(dateTimeLocalLayout, input)
			if err != nil {
				errs[field.Name] = "Некорректная дата и время"
				continue
			}
			parsed[field.Name] = t
		case "checkbox":
			parsed[field.Name] = input == "true"
		default:
			parsed[field.Name] = input
		}
	}

	runEntityBusinessValidation(entity, parsed, errs)

	return parsed, raw, errs
}

func runEntityBusinessValidation(entity models.EntityConfig, values map[string]any, errs map[string]string) {
	switch entity.Table {
	case "vehicle":
		vin, _ := values["vin"].(string)
		if vin != "" && len(vin) != 17 {
			errs["vin"] = "VIN должен содержать 17 символов"
		}
		year, ok := values["year"].(int)
		maxYear := time.Now().Year() + 1
		if ok && (year < 1980 || year > maxYear) {
			errs["year"] = fmt.Sprintf("Год выпуска должен быть в диапазоне 1980-%d", maxYear)
		}
		if km, ok := values["current_odometer_km"].(int); ok && km < 0 {
			errs["current_odometer_km"] = "Пробег не может быть отрицательным"
		}
		if cost, ok := values["acquisition_cost"].(float64); ok && cost < 0 {
			errs["acquisition_cost"] = "Стоимость не может быть отрицательной"
		}
	case "vehicle_assignment":
		from, fromOK := values["date_from"].(time.Time)
		toAny := values["date_to"]
		if fromOK && toAny != nil {
			to, ok := toAny.(time.Time)
			if ok && to.Before(from) {
				errs["date_to"] = "Дата окончания не может быть раньше даты начала"
			}
		}
	case "trip_sheet":
		start, startOK := values["odometer_start"].(int)
		end, endOK := values["odometer_end"].(int)
		dist, distOK := values["distance_km"].(int)
		if startOK && endOK {
			if end < start {
				errs["odometer_end"] = "Конечный одометр не может быть меньше начального"
			}
			if distOK && dist != end-start {
				errs["distance_km"] = fmt.Sprintf("Пробег должен быть равен %d", end-start)
			}
		}
	case "fuel_txn":
		if liters, ok := values["liters"].(float64); ok && liters <= 0 {
			errs["liters"] = "Количество литров должно быть больше 0"
		}
		if amount, ok := values["amount"].(float64); ok && amount < 0 {
			errs["amount"] = "Сумма не может быть отрицательной"
		}
		if km, ok := values["odometer_km"].(int); ok && km < 0 {
			errs["odometer_km"] = "Одометр не может быть отрицательным"
		}
	case "maintenance_order":
		from, fromOK := values["open_date"].(time.Time)
		toAny := values["close_date"]
		if fromOK && toAny != nil {
			to, ok := toAny.(time.Time)
			if ok && to.Before(from) {
				errs["close_date"] = "Дата закрытия не может быть раньше даты открытия"
			}
		}
		if cost, ok := values["cost"].(float64); ok && cost < 0 {
			errs["cost"] = "Стоимость не может быть отрицательной"
		}
	case "contract":
		from, fromOK := values["date_from"].(time.Time)
		to, toOK := values["date_to"].(time.Time)
		if fromOK && toOK && to.Before(from) {
			errs["date_to"] = "Дата окончания не может быть раньше даты начала"
		}
		if total, ok := values["total_amount"].(float64); ok && total < 0 {
			errs["total_amount"] = "Сумма не может быть отрицательной"
		}
	case "rental_event":
		from, fromOK := values["pickup_ts"].(time.Time)
		toAny := values["return_ts"]
		if fromOK && toAny != nil {
			to, ok := toAny.(time.Time)
			if ok && to.Before(from) {
				errs["return_ts"] = "Время возврата не может быть раньше времени выдачи"
			}
		}
		if price, ok := values["price_per_day"].(float64); ok && price < 0 {
			errs["price_per_day"] = "Цена не может быть отрицательной"
		}
		if dep, ok := values["deposit"].(float64); ok && dep < 0 {
			errs["deposit"] = "Депозит не может быть отрицательным"
		}
	}
}

// FormatDBValueForInput преобразует значение из БД для подстановки в форму.
func FormatDBValueForInput(field models.Field, value string) string {
	if value == "" || value == "<nil>" {
		return ""
	}

	switch field.Type {
	case "date":
		if len(value) >= 10 {
			return value[:10]
		}
		return value
	case "datetime-local":
		if len(value) >= 16 {
			return strings.Replace(value[:16], " ", "T", 1)
		}
		return value
	case "checkbox":
		lower := strings.ToLower(value)
		if lower == "true" || value == "1" || lower == "t" {
			return "true"
		}
		return ""
	default:
		return value
	}
}
