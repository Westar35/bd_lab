package handlers

import (
	"net/http"
	"strconv"
	"time"

	"bd_lab_3/internal/models"
)

type searchIndexData struct {
	BaseData models.BasePageData
}

type searchVehiclesData struct {
	BaseData   models.BasePageData
	Make       string
	Model      string
	RegNumber  string
	VIN        string
	Models     []models.Option
	RegNumbers []models.Option
	VINs       []models.Option
	Rows       []map[string]string
	Searched   bool
}

type searchTripsData struct {
	BaseData   models.BasePageData
	DriverID   string
	VehicleID  string
	DateFrom   string
	DateTo     string
	Drivers    []models.Option
	Vehicles   []models.Option
	Rows       []map[string]string
	HasResults bool
	Searched   bool
}

type searchMileageData struct {
	BaseData models.BasePageData
	DateFrom string
	DateTo   string
	Rows     []map[string]string
	Searched bool
}

type reportsIndexData struct {
	BaseData models.BasePageData
}

type periodReportData struct {
	BaseData models.BasePageData
	DateFrom string
	DateTo   string
	Rows     []map[string]string
}

type simpleReportData struct {
	BaseData models.BasePageData
	Rows     []map[string]string
}

// SearchIndex страница навигации по поискам.
func (a *App) SearchIndex(w http.ResponseWriter, r *http.Request) {
	data := searchIndexData{BaseData: a.baseData(w, r, "Поиск", "search")}
	a.render(w, r, "search_index.html", data)
}

// SearchVehiclesByMake поиск 1: автомобили по марке.
func (a *App) SearchVehiclesByMake(w http.ResponseWriter, r *http.Request) {
	makeValue := r.URL.Query().Get("make")
	modelValue := r.URL.Query().Get("model")
	regNumber := r.URL.Query().Get("reg_number")
	vin := r.URL.Query().Get("vin")
	rows := make([]map[string]string, 0)
	searched := makeValue != "" || modelValue != "" || regNumber != "" || vin != ""

	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	modelsList, err := analyticsRepo.ModelOptions(r.Context(), makeValue)
	if err != nil {
		a.logger.Printf("[search vehicle models] %v", err)
		http.Error(w, "Ошибка загрузки моделей", http.StatusInternalServerError)
		return
	}
	regNumbers, err := analyticsRepo.RegNumberOptions(r.Context(), makeValue, modelValue)
	if err != nil {
		a.logger.Printf("[search vehicle reg numbers] %v", err)
		http.Error(w, "Ошибка загрузки госномеров", http.StatusInternalServerError)
		return
	}
	vins, err := analyticsRepo.VINOptions(r.Context(), regNumber)
	if err != nil {
		a.logger.Printf("[search vehicle vins] %v", err)
		http.Error(w, "Ошибка загрузки VIN", http.StatusInternalServerError)
		return
	}

	if searched {
		result, err := analyticsRepo.VehiclesByAttributes(r.Context(), makeValue, modelValue, regNumber, vin)
		if err != nil {
			a.logger.Printf("[search make] %v", err)
			http.Error(w, "Ошибка выполнения поиска", http.StatusInternalServerError)
			return
		}
		rows = result
	}

	data := searchVehiclesData{
		BaseData:   a.baseData(w, r, "Поиск автомобилей", "search"),
		Make:       makeValue,
		Model:      modelValue,
		RegNumber:  regNumber,
		VIN:        vin,
		Models:     modelsList,
		RegNumbers: regNumbers,
		VINs:       vins,
		Rows:       rows,
		Searched:   searched,
	}

	a.render(w, r, "search_vehicles_by_make.html", data)
}

// SearchTrips поиск 2: путевые листы по водителю, авто и периоду.
func (a *App) SearchTrips(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	drivers, err := analyticsRepo.DriverOptions(r.Context())
	if err != nil {
		a.logger.Printf("[search trips drivers] %v", err)
		http.Error(w, "Ошибка загрузки водителей", http.StatusInternalServerError)
		return
	}

	driverID := q.Get("driver_id")
	vehicleID := q.Get("vehicle_id")
	var driverPtr *int
	if driverID != "" {
		id, convErr := strconv.Atoi(driverID)
		if convErr == nil {
			driverPtr = &id
		}
	}
	var vehiclePtr *int
	if vehicleID != "" {
		id, convErr := strconv.Atoi(vehicleID)
		if convErr == nil {
			vehiclePtr = &id
		}
	}

	vehicles, err := analyticsRepo.VehicleOptionsFiltered(r.Context(), driverPtr, nil)
	if err != nil {
		a.logger.Printf("[search trips vehicles] %v", err)
		http.Error(w, "Ошибка загрузки автомобилей", http.StatusInternalServerError)
		return
	}

	dateFrom := parseDateOrDefault(q.Get("date_from"), "")
	dateTo := parseDateOrDefault(q.Get("date_to"), "")

	rows := make([]map[string]string, 0)
	searched := dateFrom != "" && dateTo != ""

	if searched {
		start, err := time.Parse("2006-01-02", dateFrom)
		if err != nil {
			http.Error(w, "Некорректная дата начала", http.StatusBadRequest)
			return
		}
		end, err := time.Parse("2006-01-02", dateTo)
		if err != nil {
			http.Error(w, "Некорректная дата окончания", http.StatusBadRequest)
			return
		}

		rows, err = analyticsRepo.TripSheetsByDriverVehiclePeriod(r.Context(), driverPtr, vehiclePtr, start, end)
		if err != nil {
			a.logger.Printf("[search trips] %v", err)
			http.Error(w, "Ошибка выполнения поиска", http.StatusInternalServerError)
			return
		}
	}

	data := searchTripsData{
		BaseData:   a.baseData(w, r, "Поиск путевых листов", "search"),
		DriverID:   driverID,
		VehicleID:  vehicleID,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		Drivers:    drivers,
		Vehicles:   vehicles,
		Rows:       rows,
		HasResults: len(rows) > 0,
		Searched:   searched,
	}

	a.render(w, r, "search_trips.html", data)
}

// SearchMileage поиск 3: суммарный пробег по водителям и авто за период.
func (a *App) SearchMileage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	dateFrom := parseDateOrDefault(q.Get("date_from"), "")
	dateTo := parseDateOrDefault(q.Get("date_to"), "")

	rows := make([]map[string]string, 0)
	searched := dateFrom != "" && dateTo != ""

	if searched {
		start, err := time.Parse("2006-01-02", dateFrom)
		if err != nil {
			http.Error(w, "Некорректная дата начала", http.StatusBadRequest)
			return
		}
		end, err := time.Parse("2006-01-02", dateTo)
		if err != nil {
			http.Error(w, "Некорректная дата окончания", http.StatusBadRequest)
			return
		}

		_, analyticsRepo, _, err := a.repositoriesForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rows, err = analyticsRepo.TotalMileageByDriverVehicle(r.Context(), start, end)
		if err != nil {
			a.logger.Printf("[search mileage] %v", err)
			http.Error(w, "Ошибка выполнения поиска", http.StatusInternalServerError)
			return
		}
	}

	data := searchMileageData{
		BaseData: a.baseData(w, r, "Суммарный пробег", "search"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Rows:     rows,
		Searched: searched,
	}

	a.render(w, r, "search_mileage.html", data)
}

// ReportsIndex навигация по отчетам.
func (a *App) ReportsIndex(w http.ResponseWriter, r *http.Request) {
	data := reportsIndexData{BaseData: a.baseData(w, r, "Отчеты", "reports")}
	a.render(w, r, "reports_index.html", data)
}

// ReportMileage отчет по пробегу за период.
func (a *App) ReportMileage(w http.ResponseWriter, r *http.Request) {
	dateFrom, dateTo, start, end, ok := parseReportPeriod(r)
	if !ok {
		http.Error(w, "Некорректный период отчета", http.StatusBadRequest)
		return
	}

	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := analyticsRepo.MileageReport(r.Context(), start, end)
	if err != nil {
		a.logger.Printf("[report mileage] %v", err)
		http.Error(w, "Ошибка формирования отчета", http.StatusInternalServerError)
		return
	}

	data := periodReportData{
		BaseData: a.baseData(w, r, "Отчет: пробег за период", "reports"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Rows:     rows,
	}
	a.render(w, r, "report_mileage.html", data)
}

// ReportFuel отчет по расходам на топливо.
func (a *App) ReportFuel(w http.ResponseWriter, r *http.Request) {
	dateFrom, dateTo, start, end, ok := parseReportPeriod(r)
	if !ok {
		http.Error(w, "Некорректный период отчета", http.StatusBadRequest)
		return
	}

	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := analyticsRepo.FuelExpensesReport(r.Context(), start, end)
	if err != nil {
		a.logger.Printf("[report fuel] %v", err)
		http.Error(w, "Ошибка формирования отчета", http.StatusInternalServerError)
		return
	}

	data := periodReportData{
		BaseData: a.baseData(w, r, "Отчет: расходы на топливо", "reports"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Rows:     rows,
	}
	a.render(w, r, "report_fuel.html", data)
}

// ReportMaintenance отчет по расходам на обслуживание и ремонт.
func (a *App) ReportMaintenance(w http.ResponseWriter, r *http.Request) {
	dateFrom, dateTo, start, end, ok := parseReportPeriod(r)
	if !ok {
		http.Error(w, "Некорректный период отчета", http.StatusBadRequest)
		return
	}

	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := analyticsRepo.MaintenanceExpensesReport(r.Context(), start, end)
	if err != nil {
		a.logger.Printf("[report maintenance] %v", err)
		http.Error(w, "Ошибка формирования отчета", http.StatusInternalServerError)
		return
	}

	data := periodReportData{
		BaseData: a.baseData(w, r, "Отчет: расходы на обслуживание", "reports"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Rows:     rows,
	}
	a.render(w, r, "report_maintenance.html", data)
}

// ReportByStatus отчет по количеству авто по статусам.
func (a *App) ReportByStatus(w http.ResponseWriter, r *http.Request) {
	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := analyticsRepo.VehiclesByStatus(r.Context())
	if err != nil {
		a.logger.Printf("[report status] %v", err)
		http.Error(w, "Ошибка формирования отчета", http.StatusInternalServerError)
		return
	}

	data := simpleReportData{
		BaseData: a.baseData(w, r, "Отчет: автомобили по статусам", "reports"),
		Rows:     rows,
	}
	a.render(w, r, "report_by_status.html", data)
}

// ReportByClass отчет по количеству авто по классам.
func (a *App) ReportByClass(w http.ResponseWriter, r *http.Request) {
	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := analyticsRepo.VehiclesByClass(r.Context())
	if err != nil {
		a.logger.Printf("[report class] %v", err)
		http.Error(w, "Ошибка формирования отчета", http.StatusInternalServerError)
		return
	}

	data := simpleReportData{
		BaseData: a.baseData(w, r, "Отчет: автомобили по классам", "reports"),
		Rows:     rows,
	}
	a.render(w, r, "report_by_class.html", data)
}

// ReportCurrentRentals отчет по текущим арендам.
func (a *App) ReportCurrentRentals(w http.ResponseWriter, r *http.Request) {
	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := analyticsRepo.CurrentRentals(r.Context())
	if err != nil {
		a.logger.Printf("[report rentals] %v", err)
		http.Error(w, "Ошибка формирования отчета", http.StatusInternalServerError)
		return
	}

	data := simpleReportData{
		BaseData: a.baseData(w, r, "Отчет: текущие аренды", "reports"),
		Rows:     rows,
	}
	a.render(w, r, "report_current_rentals.html", data)
}

// ReportSummary сводный отчет по количеству записей.
func (a *App) ReportSummary(w http.ResponseWriter, r *http.Request) {
	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := analyticsRepo.DatabaseSummary(r.Context())
	if err != nil {
		a.logger.Printf("[report summary] %v", err)
		http.Error(w, "Ошибка формирования отчета", http.StatusInternalServerError)
		return
	}

	data := simpleReportData{
		BaseData: a.baseData(w, r, "Отчет: сводка базы", "reports"),
		Rows:     rows,
	}
	a.render(w, r, "report_summary.html", data)
}

func parseReportPeriod(r *http.Request) (dateFrom string, dateTo string, start time.Time, end time.Time, ok bool) {
	q := r.URL.Query()

	today := time.Now()
	defaultFrom := today.AddDate(0, -1, 0).Format("2006-01-02")
	defaultTo := today.Format("2006-01-02")

	dateFrom = parseDateOrDefault(q.Get("date_from"), defaultFrom)
	dateTo = parseDateOrDefault(q.Get("date_to"), defaultTo)

	var err error
	start, err = time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, false
	}
	end, err = time.Parse("2006-01-02", dateTo)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, false
	}

	if end.Before(start) {
		return "", "", time.Time{}, time.Time{}, false
	}

	return dateFrom, dateTo, start, end, true
}
