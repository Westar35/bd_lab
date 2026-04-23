package handlers

import "net/http"

type loginPageData struct {
	Title     string
	Error     string
	Flash     string
	Username  string
	LoggedOut bool
}

// LoginPage отображает форму входа.
func (a *App) LoginPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.authMW.CurrentUser(r); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := loginPageData{
		Title:     "Вход в систему",
		Flash:     "",
		LoggedOut: r.URL.Query().Get("logout") == "1",
	}

	a.render(w, r, "login.html", data)
}

// Login проверяет логин/пароль и создает сессию.
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Некорректные данные формы", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if !a.authSvc.Authenticate(username, password) {
		data := loginPageData{
			Title:    "Вход в систему",
			Error:    "Неверный логин или пароль",
			Username: username,
		}
		a.render(w, r, "login.html", data)
		return
	}

	if err := a.authMW.Login(w, r, username); err != nil {
		a.logger.Printf("[auth] ошибка создания сессии: %v", err)
		http.Error(w, "Ошибка входа в систему", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout завершает сессию.
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	if err := a.authMW.Logout(w, r); err != nil {
		a.logger.Printf("[auth] ошибка выхода: %v", err)
	}
	http.Redirect(w, r, "/login?logout=1", http.StatusSeeOther)
}
