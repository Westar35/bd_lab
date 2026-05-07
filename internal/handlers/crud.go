package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"bd_lab_3/internal/models"
	"bd_lab_3/internal/services"
)

// EntityList показывает список записей сущности.
func (a *App) EntityList(slug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity, ok := a.getEntity(slug)
		if !ok {
			a.notFound(w, r)
			return
		}
		entityRepo, _, _, err := a.repositoriesForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := a.parseListParams(r, entity)
		result, err := entityRepo.List(r.Context(), entity, params)
		if err != nil {
			a.logger.Printf("[list] %s: %v", entity.Table, err)
			http.Error(w, "Ошибка загрузки списка", http.StatusInternalServerError)
			return
		}

		listFields := make([]models.Field, 0, len(entity.ListColumns)+1)
		if idField, ok := entity.FieldByName("id"); ok {
			listFields = append(listFields, idField)
		}
		for _, col := range entity.ListColumns {
			if f, exists := entity.FieldByName(col); exists {
				listFields = append(listFields, f)
			}
		}

		filterSets := make(map[string][]models.Option)
		for _, f := range entity.Fields {
			if !f.Filterable {
				continue
			}
			switch f.Type {
			case "checkbox":
				filterSets[f.Name] = []models.Option{
					{ID: "true", Label: "Да"},
					{ID: "false", Label: "Нет"},
				}
			case "select":
				opts, err := entityRepo.SelectOptions(r.Context(), f)
				if err != nil {
					a.logger.Printf("[filters] %s/%s: %v", entity.Table, f.Name, err)
					continue
				}
				filterSets[f.Name] = opts
			default:
				continue
			}
		}

		data := models.ListPageData{
			BaseData: a.baseData(w, r, entity.Title, slug),
			Entity:   entity,
			Rows:     result.Rows,
			Query: models.ListQuery{
				Page:    params.Page,
				PerPage: params.PerPage,
				Sort:    params.Sort,
				Dir:     params.Dir,
				Q:       params.Q,
				Filters: params.Filters,
			},
			Total:       result.Total,
			TotalPages:  result.TotalPages,
			PageNumbers: pageNumbers(result.TotalPages),
			ListFields:  listFields,
			FilterSets:  filterSets,
			QueryTail:   buildQueryTail(params),
			FilterTail:  buildFilterTail(params),
		}

		a.render(w, r, "entity_list.html", data)
	}
}

// EntityCreatePage отображает форму создания записи.
func (a *App) EntityCreatePage(slug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity, ok := a.getEntity(slug)
		if !ok {
			a.notFound(w, r)
			return
		}
		if entity.ReadOnly {
			http.Error(w, "Эта таблица доступна только для просмотра", http.StatusForbidden)
			return
		}

		selects, err := a.selectSets(r, entity)
		if err != nil {
			a.logger.Printf("[form-selects] %s: %v", entity.Table, err)
			http.Error(w, "Ошибка загрузки справочников", http.StatusInternalServerError)
			return
		}

		data := models.FormPageData{
			BaseData:     a.baseData(w, r, "Создание: "+entity.TitleSingle, slug),
			Entity:       entity,
			Fields:       entity.FormFields(),
			Values:       draftValues(r, entity.FormFields()),
			Errors:       map[string]string{},
			Selects:      selects,
			RelatedSlugs: relatedSlugs(entity),
			Action:       "/" + slug + "/new",
			CurrentPath:  "/" + slug + "/new",
			ReturnTo:     returnToWithDraft(r),
			ReturnField:  r.URL.Query().Get("field"),
			SubmitLabel:  "Создать",
			IsEdit:       false,
		}

		a.render(w, r, "entity_form.html", data)
	}
}

// EntityCreate сохраняет новую запись.
func (a *App) EntityCreate(slug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity, ok := a.getEntity(slug)
		if !ok {
			a.notFound(w, r)
			return
		}
		if entity.ReadOnly {
			http.Error(w, "Эта таблица доступна только для просмотра", http.StatusForbidden)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Некорректные данные формы", http.StatusBadRequest)
			return
		}

		parsed, raw, valErrs := services.ParseAndValidateForm(entity, r.Form)
		if len(valErrs) > 0 {
			a.renderFormWithErrors(w, r, entity, raw, valErrs, false, "", "/"+slug+"/new")
			return
		}

		entityRepo, _, _, err := a.repositoriesForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		actor, _ := a.authMW.CurrentUser(r)
		createdID, err := entityRepo.Create(r.Context(), entity, parsed, actor)
		if err != nil {
			valErrs["_form"] = services.HumanizeDBError(err)
			a.renderFormWithErrors(w, r, entity, raw, valErrs, false, "", "/"+slug+"/new")
			return
		}

		if returnTo := r.FormValue("_return_to"); returnTo != "" {
			a.redirectWithFlash(w, r, appendQuery(returnTo, "selected_"+r.FormValue("_field"), fmt.Sprint(createdID)), "Связанная запись создана")
			return
		}
		a.redirectWithFlash(w, r, "/"+slug, "Запись успешно создана")
	}
}

// EntityDetail показывает карточку записи.
func (a *App) EntityDetail(slug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity, ok := a.getEntity(slug)
		if !ok {
			a.notFound(w, r)
			return
		}
		id, err := a.parseID(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		entityRepo, _, _, err := a.repositoriesForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		row, err := entityRepo.GetByID(r.Context(), entity, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			a.logger.Printf("[detail] %s/%d: %v", entity.Table, id, err)
			http.Error(w, "Ошибка загрузки карточки", http.StatusInternalServerError)
			return
		}

		data := models.DetailPageData{
			BaseData: a.baseData(w, r, "Карточка: "+entity.TitleSingle, slug),
			Entity:   entity,
			Row:      row,
			ID:       fmt.Sprint(id),
		}

		a.render(w, r, "entity_detail.html", data)
	}
}

// EntityEditPage отображает форму редактирования.
func (a *App) EntityEditPage(slug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity, ok := a.getEntity(slug)
		if !ok {
			a.notFound(w, r)
			return
		}
		if entity.ReadOnly {
			http.Error(w, "Эта таблица доступна только для просмотра", http.StatusForbidden)
			return
		}

		id, err := a.parseID(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		entityRepo, _, _, err := a.repositoriesForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		row, err := entityRepo.GetRawByID(r.Context(), entity, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			a.logger.Printf("[edit-page] %s/%d: %v", entity.Table, id, err)
			http.Error(w, "Ошибка загрузки записи", http.StatusInternalServerError)
			return
		}

		values := make(map[string]string)
		for _, f := range entity.FormFields() {
			values[f.Name] = services.FormatDBValueForInput(f, row[f.Name])
		}
		for key, value := range draftValues(r, entity.FormFields()) {
			if value != "" {
				values[key] = value
			}
		}

		selects, err := a.selectSets(r, entity)
		if err != nil {
			a.logger.Printf("[form-selects-edit] %s: %v", entity.Table, err)
			http.Error(w, "Ошибка загрузки справочников", http.StatusInternalServerError)
			return
		}

		data := models.FormPageData{
			BaseData:     a.baseData(w, r, "Редактирование: "+entity.TitleSingle, slug),
			Entity:       entity,
			Fields:       entity.FormFields(),
			Values:       values,
			Errors:       map[string]string{},
			Selects:      selects,
			RelatedSlugs: relatedSlugs(entity),
			Action:       fmt.Sprintf("/%s/%d/edit", slug, id),
			CurrentPath:  fmt.Sprintf("/%s/%d/edit", slug, id),
			SubmitLabel:  "Сохранить",
			IsEdit:       true,
			ID:           fmt.Sprint(id),
		}

		a.render(w, r, "entity_form.html", data)
	}
}

// EntityEdit сохраняет изменения.
func (a *App) EntityEdit(slug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity, ok := a.getEntity(slug)
		if !ok {
			a.notFound(w, r)
			return
		}
		if entity.ReadOnly {
			http.Error(w, "Эта таблица доступна только для просмотра", http.StatusForbidden)
			return
		}

		id, err := a.parseID(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if err = r.ParseForm(); err != nil {
			http.Error(w, "Некорректные данные формы", http.StatusBadRequest)
			return
		}

		parsed, raw, valErrs := services.ParseAndValidateForm(entity, r.Form)
		action := fmt.Sprintf("/%s/%d/edit", slug, id)
		if len(valErrs) > 0 {
			a.renderFormWithErrors(w, r, entity, raw, valErrs, true, fmt.Sprint(id), action)
			return
		}

		entityRepo, _, _, err := a.repositoriesForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		actor, _ := a.authMW.CurrentUser(r)
		err = entityRepo.Update(r.Context(), entity, id, parsed, actor)
		if err != nil {
			valErrs["_form"] = services.HumanizeDBError(err)
			a.renderFormWithErrors(w, r, entity, raw, valErrs, true, fmt.Sprint(id), action)
			return
		}

		a.redirectWithFlash(w, r, "/"+slug, "Изменения сохранены")
	}
}

// EntityDeletePage показывает подтверждение удаления.
func (a *App) EntityDeletePage(slug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity, ok := a.getEntity(slug)
		if !ok {
			a.notFound(w, r)
			return
		}
		if entity.ReadOnly {
			http.Error(w, "Эта таблица доступна только для просмотра", http.StatusForbidden)
			return
		}

		id, err := a.parseID(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		entityRepo, _, _, err := a.repositoriesForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		row, err := entityRepo.GetByID(r.Context(), entity, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			a.logger.Printf("[delete-page] %s/%d: %v", entity.Table, id, err)
			http.Error(w, "Ошибка загрузки записи", http.StatusInternalServerError)
			return
		}

		data := models.DeletePageData{
			BaseData: a.baseData(w, r, "Удаление: "+entity.TitleSingle, slug),
			Entity:   entity,
			Row:      row,
			ID:       fmt.Sprint(id),
		}

		a.render(w, r, "entity_delete.html", data)
	}
}

// EntityDelete выполняет удаление.
func (a *App) EntityDelete(slug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity, ok := a.getEntity(slug)
		if !ok {
			a.notFound(w, r)
			return
		}
		if entity.ReadOnly {
			http.Error(w, "Эта таблица доступна только для просмотра", http.StatusForbidden)
			return
		}

		id, err := a.parseID(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		entityRepo, _, _, err := a.repositoriesForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		actor, _ := a.authMW.CurrentUser(r)
		err = entityRepo.Delete(r.Context(), entity, id, actor)
		if err != nil {
			a.logger.Printf("[delete] %s/%d: %v", entity.Table, id, err)
			a.redirectWithFlash(w, r, "/"+slug, services.HumanizeDBError(err))
			return
		}

		a.redirectWithFlash(w, r, "/"+slug, "Запись удалена")
	}
}

func (a *App) renderFormWithErrors(
	w http.ResponseWriter,
	r *http.Request,
	entity models.EntityConfig,
	values map[string]string,
	errs map[string]string,
	isEdit bool,
	id string,
	action string,
) {
	selects, err := a.selectSets(r, entity)
	if err != nil {
		a.logger.Printf("[render-form-errors] %s: %v", entity.Table, err)
		http.Error(w, "Ошибка загрузки справочников", http.StatusInternalServerError)
		return
	}

	title := "Создание: " + entity.TitleSingle
	submit := "Создать"
	if isEdit {
		title = "Редактирование: " + entity.TitleSingle
		submit = "Сохранить"
	}

	data := models.FormPageData{
		BaseData:     a.baseData(w, r, title, entity.Slug),
		Entity:       entity,
		Fields:       entity.FormFields(),
		Values:       values,
		Errors:       errs,
		Selects:      selects,
		RelatedSlugs: relatedSlugs(entity),
		Action:       action,
		CurrentPath:  action,
		ReturnTo:     r.FormValue("_return_to"),
		ReturnField:  r.FormValue("_field"),
		SubmitLabel:  submit,
		IsEdit:       isEdit,
		ID:           id,
	}

	a.render(w, r, "entity_form.html", data)
}

func (a *App) selectSets(r *http.Request, entity models.EntityConfig) (map[string][]models.Option, error) {
	entityRepo, _, _, err := a.repositoriesForRequest(r)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]models.Option)
	for _, f := range entity.FormFields() {
		if f.Type != "select" {
			continue
		}

		opts, err := entityRepo.SelectOptions(r.Context(), f)
		if err != nil {
			return nil, err
		}
		result[f.Name] = opts
	}
	return result, nil
}
