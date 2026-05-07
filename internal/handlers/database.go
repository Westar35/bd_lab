package handlers

import (
	"net/http"

	dbx "bd_lab_3/internal/db"
)

func (a *App) SwitchDB(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Некорректные данные формы", http.StatusBadRequest)
		return
	}

	raw := r.FormValue("db")
	if raw != string(dbx.DBPostgres) && raw != string(dbx.DBMySQL) {
		a.redirectWithFlash(w, r, "/", "Неизвестная база данных")
		return
	}
	target := dbx.DBType(raw)
	if _, _, err := a.dbManager.Get(target); err != nil {
		a.redirectWithFlash(w, r, "/", "Не удалось переключить базу данных: "+err.Error())
		return
	}
	if err := a.authMW.SetActiveDB(w, r, target); err != nil {
		http.Error(w, "Не удалось сохранить активную БД", http.StatusInternalServerError)
		return
	}

	back := r.Header.Get("Referer")
	if back == "" {
		back = "/"
	}
	a.redirectWithFlash(w, r, back, "Активная база данных изменена на "+target.DisplayName()+".")
}
