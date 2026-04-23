package models

// Option используется для select-полей.
type Option struct {
	ID    string
	Label string
}

// Field описывает поле сущности для списка/карточки/формы.
type Field struct {
	Name          string
	Label         string
	Type          string
	Required      bool
	Nullable      bool
	InList        bool
	InDetail      bool
	InForm        bool
	Sortable      bool
	Filterable    bool
	RefTable      string
	RefLabelExpr  string
	StaticOptions []Option
	Placeholder   string
}

// EntityConfig описывает метаданные сущности для универсального CRUD.
type EntityConfig struct {
	Slug          string
	Table         string
	Title         string
	TitleSingle   string
	Fields        []Field
	ListColumns   []string
	SearchColumns []string
	DefaultSort   string
	DefaultDir    string
}

// FieldByName возвращает поле по имени.
func (e EntityConfig) FieldByName(name string) (Field, bool) {
	for _, f := range e.Fields {
		if f.Name == name {
			return f, true
		}
	}
	return Field{}, false
}

// FormFields возвращает поля, доступные на форме.
func (e EntityConfig) FormFields() []Field {
	out := make([]Field, 0, len(e.Fields))
	for _, f := range e.Fields {
		if f.InForm {
			out = append(out, f)
		}
	}
	return out
}

// DetailFields возвращает поля карточки.
func (e EntityConfig) DetailFields() []Field {
	out := make([]Field, 0, len(e.Fields))
	for _, f := range e.Fields {
		if f.InDetail {
			out = append(out, f)
		}
	}
	return out
}

// EntityMap возвращает конфигурации всех обязательных бизнес-сущностей.
func EntityMap() map[string]EntityConfig {
	return map[string]EntityConfig{
		"vehicles": {
			Slug:          "vehicles",
			Table:         "vehicle",
			Title:         "Автомобили",
			TitleSingle:   "автомобиль",
			ListColumns:   []string{"vin", "reg_number", "make", "model", "year", "vehicle_class_id", "status_id", "current_odometer_km"},
			SearchColumns: []string{"vin", "reg_number", "make", "model"},
			DefaultSort:   "id",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "vin", Label: "VIN", Type: "text", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true, Placeholder: "17 символов"},
				{Name: "reg_number", Label: "Гос. номер", Type: "text", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "make", Label: "Марка", Type: "text", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "model", Label: "Модель", Type: "text", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "year", Label: "Год выпуска", Type: "int", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "vehicle_class_id", Label: "Класс ТС", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true, Filterable: true, RefTable: "vehicle_class", RefLabelExpr: "name"},
				{Name: "status_id", Label: "Статус", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true, Filterable: true, RefTable: "vehicle_status", RefLabelExpr: "name"},
				{Name: "fuel_type_id", Label: "Тип топлива", Type: "select", Required: true, InList: false, InDetail: true, InForm: true, RefTable: "fuel_type", RefLabelExpr: "name"},
				{Name: "transmission_type_id", Label: "Трансмиссия", Type: "select", Required: true, InList: false, InDetail: true, InForm: true, RefTable: "transmission_type", RefLabelExpr: "name"},
				{Name: "acquisition_type_id", Label: "Тип приобретения", Type: "select", Required: true, InList: false, InDetail: true, InForm: true, RefTable: "acquisition_type", RefLabelExpr: "name"},
				{Name: "acquisition_date", Label: "Дата приобретения", Type: "date", Required: true, InList: false, InDetail: true, InForm: true, Sortable: true},
				{Name: "acquisition_cost", Label: "Стоимость приобретения", Type: "decimal", Required: true, InList: false, InDetail: true, InForm: true},
				{Name: "current_odometer_km", Label: "Текущий одометр, км", Type: "int", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "created_at", Label: "Создан", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
				{Name: "updated_at", Label: "Обновлен", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
		"drivers": {
			Slug:          "drivers",
			Table:         "driver",
			Title:         "Водители",
			TitleSingle:   "водитель",
			ListColumns:   []string{"fio", "license_number", "phone", "created_at"},
			SearchColumns: []string{"fio", "license_number", "phone"},
			DefaultSort:   "id",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "fio", Label: "ФИО", Type: "text", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "license_number", Label: "№ водительского удостоверения", Type: "text", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "phone", Label: "Телефон", Type: "text", InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "created_at", Label: "Создан", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
				{Name: "updated_at", Label: "Обновлен", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
		"departments": {
			Slug:          "departments",
			Table:         "department",
			Title:         "Подразделения",
			TitleSingle:   "подразделение",
			ListColumns:   []string{"code", "name", "created_at"},
			SearchColumns: []string{"code", "name"},
			DefaultSort:   "id",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "code", Label: "Код", Type: "text", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "name", Label: "Название", Type: "text", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "created_at", Label: "Создан", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
				{Name: "updated_at", Label: "Обновлен", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
		"assignments": {
			Slug:          "assignments",
			Table:         "vehicle_assignment",
			Title:         "Назначения ТС",
			TitleSingle:   "назначение",
			ListColumns:   []string{"vehicle_id", "driver_id", "department_id", "date_from", "date_to", "is_primary"},
			SearchColumns: []string{"id"},
			DefaultSort:   "id",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "vehicle_id", Label: "Автомобиль", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true, Filterable: true, RefTable: "vehicle", RefLabelExpr: "make || ' ' || model || ' (' || reg_number || ')'"},
				{Name: "driver_id", Label: "Водитель", Type: "select", Nullable: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "driver", RefLabelExpr: "fio"},
				{Name: "department_id", Label: "Подразделение", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "department", RefLabelExpr: "name"},
				{Name: "date_from", Label: "Дата начала", Type: "date", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "date_to", Label: "Дата окончания", Type: "date", Nullable: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "is_primary", Label: "Основное назначение", Type: "checkbox", InList: true, InDetail: true, InForm: true, Filterable: true},
				{Name: "created_at", Label: "Создано", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
		"trip-sheets": {
			Slug:          "trip-sheets",
			Table:         "trip_sheet",
			Title:         "Путевые листы",
			TitleSingle:   "путевой лист",
			ListColumns:   []string{"trip_date", "vehicle_id", "driver_id", "department_id", "odometer_start", "odometer_end", "distance_km"},
			SearchColumns: []string{"route", "purpose"},
			DefaultSort:   "trip_date",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "trip_date", Label: "Дата поездки", Type: "date", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "vehicle_id", Label: "Автомобиль", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true, Filterable: true, RefTable: "vehicle", RefLabelExpr: "make || ' ' || model || ' (' || reg_number || ')'"},
				{Name: "driver_id", Label: "Водитель", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "driver", RefLabelExpr: "fio"},
				{Name: "department_id", Label: "Подразделение", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "department", RefLabelExpr: "name"},
				{Name: "odometer_start", Label: "Одометр начало", Type: "int", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "odometer_end", Label: "Одометр конец", Type: "int", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "route", Label: "Маршрут", Type: "text", InList: false, InDetail: true, InForm: true},
				{Name: "purpose", Label: "Цель", Type: "text", InList: false, InDetail: true, InForm: true},
				{Name: "distance_km", Label: "Пробег, км", Type: "int", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "created_at", Label: "Создан", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
		"fuel-txns": {
			Slug:          "fuel-txns",
			Table:         "fuel_txn",
			Title:         "Топливные операции",
			TitleSingle:   "топливная операция",
			ListColumns:   []string{"txn_ts", "vehicle_id", "liters", "amount", "station", "odometer_km", "payment_type_id"},
			SearchColumns: []string{"station"},
			DefaultSort:   "txn_ts",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "txn_ts", Label: "Дата/время операции", Type: "datetime-local", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "vehicle_id", Label: "Автомобиль", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "vehicle", RefLabelExpr: "make || ' ' || model || ' (' || reg_number || ')'"},
				{Name: "liters", Label: "Литры", Type: "decimal", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "amount", Label: "Сумма", Type: "decimal", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "station", Label: "АЗС", Type: "text", InList: true, InDetail: true, InForm: true},
				{Name: "odometer_km", Label: "Одометр, км", Type: "int", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "payment_type_id", Label: "Тип оплаты", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "payment_type", RefLabelExpr: "name"},
				{Name: "created_at", Label: "Создана", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
		"maintenance-orders": {
			Slug:          "maintenance-orders",
			Table:         "maintenance_order",
			Title:         "Обслуживание и ремонты",
			TitleSingle:   "заказ обслуживания",
			ListColumns:   []string{"vehicle_id", "maintenance_type_id", "open_date", "close_date", "service_name", "cost"},
			SearchColumns: []string{"service_name", "description"},
			DefaultSort:   "open_date",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "vehicle_id", Label: "Автомобиль", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "vehicle", RefLabelExpr: "make || ' ' || model || ' (' || reg_number || ')'"},
				{Name: "maintenance_type_id", Label: "Тип работ", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "maintenance_type", RefLabelExpr: "name"},
				{Name: "open_date", Label: "Дата открытия", Type: "date", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "close_date", Label: "Дата закрытия", Type: "date", Nullable: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "service_name", Label: "Сервис", Type: "text", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "cost", Label: "Стоимость", Type: "decimal", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "description", Label: "Описание", Type: "textarea", InList: false, InDetail: true, InForm: true},
				{Name: "created_at", Label: "Создан", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
		"counterparties": {
			Slug:          "counterparties",
			Table:         "counterparty",
			Title:         "Контрагенты",
			TitleSingle:   "контрагент",
			ListColumns:   []string{"type", "name", "inn", "phone", "email"},
			SearchColumns: []string{"type", "name", "inn", "phone", "email", "address"},
			DefaultSort:   "id",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "type", Label: "Тип", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, StaticOptions: []Option{{ID: "company", Label: "Юр. лицо"}, {ID: "individual", Label: "Физ. лицо"}}},
				{Name: "name", Label: "Наименование", Type: "text", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "inn", Label: "ИНН", Type: "text", InList: true, InDetail: true, InForm: true},
				{Name: "phone", Label: "Телефон", Type: "text", InList: true, InDetail: true, InForm: true},
				{Name: "email", Label: "Email", Type: "text", InList: true, InDetail: true, InForm: true},
				{Name: "address", Label: "Адрес", Type: "textarea", InList: false, InDetail: true, InForm: true},
				{Name: "created_at", Label: "Создан", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
				{Name: "updated_at", Label: "Обновлен", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
		"contracts": {
			Slug:          "contracts",
			Table:         "contract",
			Title:         "Договоры",
			TitleSingle:   "договор",
			ListColumns:   []string{"number", "contract_type_id", "counterparty_id", "date_from", "date_to", "status_id", "total_amount"},
			SearchColumns: []string{"number", "notes"},
			DefaultSort:   "id",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "contract_type_id", Label: "Тип договора", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "contract_type", RefLabelExpr: "name"},
				{Name: "counterparty_id", Label: "Контрагент", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "counterparty", RefLabelExpr: "name"},
				{Name: "number", Label: "Номер", Type: "text", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "date_from", Label: "Дата начала", Type: "date", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "date_to", Label: "Дата окончания", Type: "date", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "status_id", Label: "Статус", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "contract_status", RefLabelExpr: "name"},
				{Name: "total_amount", Label: "Сумма", Type: "decimal", Nullable: true, InList: true, InDetail: true, InForm: true},
				{Name: "notes", Label: "Примечание", Type: "textarea", Nullable: true, InList: false, InDetail: true, InForm: true},
				{Name: "created_at", Label: "Создан", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
				{Name: "updated_at", Label: "Обновлен", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
		"rental-events": {
			Slug:          "rental-events",
			Table:         "rental_event",
			Title:         "События аренды",
			TitleSingle:   "событие аренды",
			ListColumns:   []string{"contract_id", "vehicle_id", "pickup_ts", "return_ts", "price_per_day", "deposit"},
			SearchColumns: []string{"notes"},
			DefaultSort:   "pickup_ts",
			DefaultDir:    "desc",
			Fields: []Field{
				{Name: "id", Label: "ID", Type: "int", InList: true, InDetail: true, InForm: false, Sortable: true},
				{Name: "contract_id", Label: "Договор", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "contract", RefLabelExpr: "number"},
				{Name: "vehicle_id", Label: "Автомобиль", Type: "select", Required: true, InList: true, InDetail: true, InForm: true, Filterable: true, RefTable: "vehicle", RefLabelExpr: "make || ' ' || model || ' (' || reg_number || ')'"},
				{Name: "pickup_ts", Label: "Выдача", Type: "datetime-local", Required: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "return_ts", Label: "Возврат", Type: "datetime-local", Nullable: true, InList: true, InDetail: true, InForm: true, Sortable: true},
				{Name: "price_per_day", Label: "Цена/день", Type: "decimal", Required: true, InList: true, InDetail: true, InForm: true},
				{Name: "deposit", Label: "Депозит", Type: "decimal", Nullable: true, InList: true, InDetail: true, InForm: true},
				{Name: "notes", Label: "Примечание", Type: "textarea", Nullable: true, InList: false, InDetail: true, InForm: true},
				{Name: "created_at", Label: "Создано", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
				{Name: "updated_at", Label: "Обновлено", Type: "datetime", InList: false, InDetail: true, InForm: false, Sortable: true},
			},
		},
	}
}
