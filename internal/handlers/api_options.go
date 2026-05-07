package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bd_lab_3/internal/models"
	"bd_lab_3/internal/repositories"
)

func (a *App) APIVehicles(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.VehicleOptionsFiltered(r.Context(), intPtrFromQuery(r, "driver_id"), intPtrFromQuery(r, "department_id"))
	})
}

func (a *App) APIDrivers(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.DriverOptionsFiltered(r.Context(), intPtrFromQuery(r, "vehicle_id"), intPtrFromQuery(r, "department_id"))
	})
}

func (a *App) APIModels(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.ModelOptions(r.Context(), r.URL.Query().Get("make"))
	})
}

func (a *App) APIRegNumbers(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.RegNumberOptions(r.Context(), r.URL.Query().Get("make"), r.URL.Query().Get("model"))
	})
}

func (a *App) APIVINs(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.VINOptions(r.Context(), r.URL.Query().Get("reg_number"))
	})
}

func (a *App) APIContracts(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.ContractOptionsFiltered(r.Context(), intPtrFromQuery(r, "counterparty_id"))
	})
}

func (a *App) APIRentalEvents(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.RentalEventOptionsFiltered(r.Context(), intPtrFromQuery(r, "contract_id"))
	})
}

func (a *App) APITripSheets(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.TripSheetOptionsFiltered(r.Context(), intPtrFromQuery(r, "vehicle_id"), intPtrFromQuery(r, "driver_id"))
	})
}

func (a *App) APIFuelTransactions(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.FuelTransactionOptionsFiltered(r.Context(), intPtrFromQuery(r, "vehicle_id"))
	})
}

func (a *App) APIMaintenanceOrders(w http.ResponseWriter, r *http.Request) {
	a.writeOptions(w, r, func(repo *repositories.AnalyticsRepository) ([]models.Option, error) {
		return repo.MaintenanceOrderOptionsFiltered(r.Context(), intPtrFromQuery(r, "vehicle_id"))
	})
}

func (a *App) writeOptions(w http.ResponseWriter, r *http.Request, load func(*repositories.AnalyticsRepository) ([]models.Option, error)) {
	_, analyticsRepo, _, err := a.repositoriesForRequest(r)
	if err != nil {
		http.Error(w, "Ошибка выбора активной базы данных", http.StatusInternalServerError)
		return
	}

	options, err := load(analyticsRepo)
	if err != nil {
		a.logger.Printf("[api options] %v", err)
		http.Error(w, "Ошибка загрузки вариантов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(options); err != nil {
		a.logger.Printf("[api options encode] %v", err)
	}
}

func intPtrFromQuery(r *http.Request, key string) *int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return nil
	}
	return &value
}
