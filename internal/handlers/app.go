package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"bd_lab_3/internal/config"
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
	cfg           *config.Config
	logger        *log.Logger
	renderer      *views.Renderer
	authMW        *middleware.AuthMiddleware
	authSvc       *services.AuthService
	entityRepo    *repositories.EntityRepository
	analyticsRepo *repositories.AnalyticsRepository
	entities      map[string]models.EntityConfig
	entityOrder   []string
}

func NewApp(cfg *config.Config, logger *log.Logger, db *sql.DB) (*App, error) {
	renderer, err := views.NewRenderer("templates")
	if err != nil {
		return nil, err
	}

	app := &App{
		cfg:           cfg,
		logger:        logger,
		renderer:      renderer,
		authMW:        middleware.NewAuthMiddleware(cfg.SessionKey, cfg.AppEnv == "production"),
		authSvc:       services.NewAuthService(cfg.AdminUsername, cfg.AdminPassword),
		entityRepo:    repositories.NewEntityRepository(db),
		analyticsRepo: repositories.NewAnalyticsRepository(db),
		entities:      models.EntityMap(),
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

		pr.Get("/help", a.HelpPage)

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

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] == '/' {
		panic("path must not end with slash")
	}

	fs := http.StripPrefix(path, http.FileServer(root))
	r.Get(path+"/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
