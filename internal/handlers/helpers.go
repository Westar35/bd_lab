package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	dbx "bd_lab_3/internal/db"
	"bd_lab_3/internal/models"
	"bd_lab_3/internal/repositories"

	"github.com/go-chi/chi/v5"
)

func (a *App) getEntity(slug string) (models.EntityConfig, bool) {
	e, ok := a.entities[slug]
	return e, ok
}

func (a *App) baseData(w http.ResponseWriter, r *http.Request, title, active string) models.BasePageData {
	username, auth := a.authMW.CurrentUser(r)
	activeDB := a.activeDB(r)

	return models.BasePageData{
		Title:        title,
		ActiveMenu:   active,
		Username:     username,
		Auth:         auth,
		ActiveDB:     string(activeDB),
		ActiveDBName: activeDB.DisplayName(),
		Flash:        a.authMW.PullFlash(w, r),
		NavItems:     a.navItems(active),
	}
}

func (a *App) navItems(active string) []models.NavItem {
	items := make([]models.NavItem, 0, 16)
	items = append(items, models.NavItem{Name: "Главная", Path: "/", Active: active == "dashboard"})

	for _, slug := range a.entityOrder {
		entity := a.entities[slug]
		if entity.Category == "lookup" {
			continue
		}
		items = append(items, models.NavItem{
			Name:   entity.Title,
			Path:   "/" + slug,
			Active: active == slug,
		})
	}

	items = append(items,
		models.NavItem{Name: "Справочники", Path: "/lookups", Active: active == "lookups"},
		models.NavItem{Name: "Поиск", Path: "/search", Active: active == "search"},
		models.NavItem{Name: "Отчеты", Path: "/reports", Active: active == "reports"},
		models.NavItem{Name: "Справка", Path: "/help", Active: active == "help"},
	)

	return items
}

func (a *App) activeDB(r *http.Request) dbx.DBType {
	return a.authMW.ActiveDB(r, dbx.ParseDBType(a.cfg.DefaultDB, dbx.DBPostgres))
}

func (a *App) render(w http.ResponseWriter, r *http.Request, templateName string, data any) {
	if err := a.renderer.Render(w, templateName, data); err != nil {
		a.logger.Printf("[render] ошибка шаблона %s: %v", templateName, err)
		http.Error(w, "Внутренняя ошибка шаблона", http.StatusInternalServerError)
	}
}

func (a *App) parseID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.ParseInt(idStr, 10, 64)
}

func (a *App) parseListParams(r *http.Request, entity models.EntityConfig) repositories.ListParams {
	q := r.URL.Query()
	page := parseIntDefault(q.Get("page"), 1)
	perPage := parseIntDefault(q.Get("per_page"), 20)
	sort := q.Get("sort")
	dir := q.Get("dir")
	textQ := q.Get("q")

	filters := make(map[string]string)
	for _, f := range entity.Fields {
		if f.Filterable {
			key := "f_" + f.Name
			filters[f.Name] = q.Get(key)
		}
	}

	return repositories.ListParams{
		Page:    page,
		PerPage: perPage,
		Sort:    sort,
		Dir:     dir,
		Q:       textQ,
		Filters: filters,
	}
}

func buildQueryTail(params repositories.ListParams) template.URL {
	values := url.Values{}
	if params.PerPage > 0 {
		values.Set("per_page", strconv.Itoa(params.PerPage))
	}
	if params.Sort != "" {
		values.Set("sort", params.Sort)
	}
	if params.Dir != "" {
		values.Set("dir", params.Dir)
	}
	if params.Q != "" {
		values.Set("q", params.Q)
	}
	for key, val := range params.Filters {
		if val != "" {
			values.Set("f_"+key, val)
		}
	}

	encoded := values.Encode()
	if encoded == "" {
		return ""
	}
	return template.URL("&" + encoded)
}

func buildFilterTail(params repositories.ListParams) template.URL {
	values := url.Values{}
	if params.PerPage > 0 {
		values.Set("per_page", strconv.Itoa(params.PerPage))
	}
	if params.Q != "" {
		values.Set("q", params.Q)
	}
	for key, val := range params.Filters {
		if val != "" {
			values.Set("f_"+key, val)
		}
	}

	encoded := values.Encode()
	if encoded == "" {
		return ""
	}
	return template.URL("&" + encoded)
}

func parseIntDefault(value string, fallback int) int {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func pageNumbers(totalPages int) []int {
	if totalPages < 1 {
		return []int{1}
	}
	pages := make([]int, 0, totalPages)
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	return pages
}

func (a *App) notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (a *App) redirectWithFlash(w http.ResponseWriter, r *http.Request, target, message string) {
	_ = a.authMW.SetFlash(w, r, message)
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func parseDateOrDefault(input string, fallback string) string {
	if input == "" {
		return fallback
	}
	return input
}

func mapToValues(row map[string]string, fields []models.Field) map[string]string {
	values := make(map[string]string, len(fields))
	for _, f := range fields {
		values[f.Name] = row[f.Name]
	}
	return values
}

func draftValues(r *http.Request, fields []models.Field) map[string]string {
	values := make(map[string]string, len(fields))
	q := r.URL.Query()
	for _, f := range fields {
		values[f.Name] = q.Get("draft_" + f.Name)
		if selected := q.Get("selected_" + f.Name); selected != "" {
			values[f.Name] = selected
		}
	}
	return values
}

func returnToWithDraft(r *http.Request) string {
	returnTo := r.URL.Query().Get("return_to")
	if returnTo == "" {
		return ""
	}
	values := url.Values{}
	for key, vals := range r.URL.Query() {
		if strings.HasPrefix(key, "draft_") {
			for _, val := range vals {
				values.Add(key, val)
			}
		}
	}
	if encoded := values.Encode(); encoded != "" {
		sep := "?"
		if strings.Contains(returnTo, "?") {
			sep = "&"
		}
		returnTo += sep + encoded
	}
	return returnTo
}

func appendQuery(rawURL, key, value string) string {
	sep := "?"
	if strings.Contains(rawURL, "?") {
		sep = "&"
	}
	return rawURL + sep + url.QueryEscape(key) + "=" + url.QueryEscape(value)
}

func relatedSlugs(entity models.EntityConfig) map[string]string {
	result := make(map[string]string)
	for _, f := range entity.FormFields() {
		if f.Type != "select" || f.RefTable == "" {
			continue
		}
		if slug := slugForTable(f.RefTable); slug != "" {
			result[f.Name] = slug
		}
	}
	return result
}

func slugForTable(table string) string {
	switch table {
	case "vehicle":
		return "vehicles"
	case "driver":
		return "drivers"
	case "department":
		return "departments"
	case "vehicle_class":
		return "vehicle-classes"
	case "vehicle_status":
		return "vehicle-statuses"
	case "fuel_type":
		return "fuel-types"
	case "transmission_type":
		return "transmission-types"
	case "acquisition_type":
		return "acquisition-types"
	case "maintenance_type":
		return "maintenance-types"
	case "payment_type":
		return "payment-types"
	case "contract_type":
		return "contract-types"
	case "contract_status":
		return "contract-statuses"
	case "counterparty":
		return "counterparties"
	case "contract":
		return "contracts"
	default:
		return ""
	}
}

func makeErrorMap(err error) map[string]string {
	if err == nil {
		return map[string]string{}
	}
	return map[string]string{"_form": fmt.Sprint(err)}
}
