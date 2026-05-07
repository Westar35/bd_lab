package handlers

import (
	"log"
	"net/http"

	"bd_lab_3/internal/config"
	dbx "bd_lab_3/internal/db"
	"bd_lab_3/internal/middleware"
	"bd_lab_3/internal/models"
	"bd_lab_3/internal/repositories"
	"bd_lab_3/internal/services"
	"bd_lab_3/internal/views"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// App объединяет зависимости HTTP слоя.
type App struct {
	cfg         *config.Config
	logger      *log.Logger
	renderer    *views.Renderer
	authMW      *middleware.AuthMiddleware
	authSvc     *services.AuthService
	dbManager   *dbx.Manager
	entities    map[string]models.EntityConfig
	entityOrder []string
	lookupOrder []string
}

func NewApp(cfg *config.Config, logger *log.Logger, dbManager *dbx.Manager) (*App, error) {
	renderer, err := views.NewRenderer("templates")
	if err != nil {
		return nil, err
	}

	app := &App{
		cfg:       cfg,
		logger:    logger,
		renderer:  renderer,
		authMW:    middleware.NewAuthMiddleware(cfg.SessionKey, cfg.AppEnv == "production"),
		authSvc:   services.NewAuthService(cfg.AdminUsername, cfg.AdminPassword),
		dbManager: dbManager,
		entities:  models.EntityMap(),
		entityOrder: []string{
			"vehicles",
			"drivers",
			"departments",
			"assignments",
			"trip-sheets",
			"fuel-txns",
			"maintenance-orders",
			"counterparties",
			"contracts",
			"rental-events",
			"audit-log",
			"vehicle-classes",
			"vehicle-statuses",
			"fuel-types",
			"transmission-types",
			"acquisition-types",
			"maintenance-types",
			"payment-types",
			"contract-types",
			"contract-statuses",
		},
		lookupOrder: []string{
			"vehicle-classes",
			"vehicle-statuses",
			"fuel-types",
			"transmission-types",
			"acquisition-types",
			"maintenance-types",
			"payment-types",
			"contract-types",
			"contract-statuses",
		},
	}

	return app, nil
}

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Logger)

	fileServer(r, "/static", http.Dir("static"))

	r.Get("/login", a.LoginPage)
	r.Post("/login", a.Login)

	r.Group(func(pr chi.Router) {
		pr.Use(a.authMW.RequireAuth)

		pr.Get("/", a.Dashboard)
		pr.Post("/logout", a.Logout)
		pr.Post("/switch-db", a.SwitchDB)

		pr.Get("/help", a.HelpPage)
		pr.Get("/lookups", a.LookupsPage)

		pr.Get("/api/options/vehicles", a.APIVehicles)
		pr.Get("/api/options/drivers", a.APIDrivers)
		pr.Get("/api/options/models", a.APIModels)
		pr.Get("/api/options/reg-numbers", a.APIRegNumbers)
		pr.Get("/api/options/vins", a.APIVINs)
		pr.Get("/api/options/contracts", a.APIContracts)
		pr.Get("/api/options/rental-events", a.APIRentalEvents)
		pr.Get("/api/options/trip-sheets", a.APITripSheets)
		pr.Get("/api/options/fuel-transactions", a.APIFuelTransactions)
		pr.Get("/api/options/maintenance-orders", a.APIMaintenanceOrders)

		pr.Get("/search", a.SearchIndex)
		pr.Get("/search/vehicles-by-make", a.SearchVehiclesByMake)
		pr.Get("/search/trips", a.SearchTrips)
		pr.Get("/search/mileage", a.SearchMileage)

		pr.Get("/reports", a.ReportsIndex)
		pr.Get("/reports/mileage", a.ReportMileage)
		pr.Get("/reports/fuel", a.ReportFuel)
		pr.Get("/reports/maintenance", a.ReportMaintenance)
		pr.Get("/reports/by-status", a.ReportByStatus)
		pr.Get("/reports/by-class", a.ReportByClass)
		pr.Get("/reports/current-rentals", a.ReportCurrentRentals)
		pr.Get("/reports/summary", a.ReportSummary)

		for _, slug := range a.entityOrder {
			entitySlug := slug
			basePath := "/" + entitySlug
			pr.Get(basePath, a.EntityList(entitySlug))
			pr.Get(basePath+"/new", a.EntityCreatePage(entitySlug))
			pr.Post(basePath+"/new", a.EntityCreate(entitySlug))
			pr.Get(basePath+"/{id:[0-9]+}", a.EntityDetail(entitySlug))
			pr.Get(basePath+"/{id:[0-9]+}/edit", a.EntityEditPage(entitySlug))
			pr.Post(basePath+"/{id:[0-9]+}/edit", a.EntityEdit(entitySlug))
			pr.Get(basePath+"/{id:[0-9]+}/delete", a.EntityDeletePage(entitySlug))
			pr.Post(basePath+"/{id:[0-9]+}/delete", a.EntityDelete(entitySlug))
		}
	})

	return r
}

func (a *App) repositoriesForRequest(r *http.Request) (*repositories.EntityRepository, *repositories.AnalyticsRepository, dbx.DBType, error) {
	active := a.activeDB(r)
	conn, actual, err := a.dbManager.Get(active)
	if err != nil {
		return nil, nil, actual, err
	}
	return repositories.NewEntityRepository(conn, actual), repositories.NewAnalyticsRepository(conn, actual), actual, nil
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] == '/' {
		panic("path must not end with slash")
	}

	fs := http.StripPrefix(path, http.FileServer(root))
	r.Get(path+"/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
