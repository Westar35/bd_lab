package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const sessionName = "fleet_session"

// AuthMiddleware управляет cookie-сессиями и проверкой входа.
type AuthMiddleware struct {
	store *sessions.CookieStore
}

func NewAuthMiddleware(sessionKey string, secure bool) *AuthMiddleware {
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 12,
	}
	return &AuthMiddleware{store: store}
}

func (m *AuthMiddleware) session(r *http.Request) (*sessions.Session, error) {
	return m.store.Get(r, sessionName)
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := m.CurrentUser(r)
		if !ok {
			_ = m.SetFlash(w, r, "Сначала выполните вход в систему")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) Login(w http.ResponseWriter, r *http.Request, username string) error {
	sess, err := m.session(r)
	if err != nil {
		return err
	}
	sess.Values["auth"] = true
	sess.Values["username"] = username
	return sess.Save(r, w)
}

func (m *AuthMiddleware) Logout(w http.ResponseWriter, r *http.Request) error {
	sess, err := m.session(r)
	if err != nil {
		return err
	}
	sess.Options.MaxAge = -1
	return sess.Save(r, w)
}

func (m *AuthMiddleware) CurrentUser(r *http.Request) (string, bool) {
	sess, err := m.session(r)
	if err != nil {
		return "", false
	}

	authVal, ok := sess.Values["auth"].(bool)
	if !ok || !authVal {
		return "", false
	}

	username, ok := sess.Values["username"].(string)
	if !ok || username == "" {
		return "", false
	}

	return username, true
}

func (m *AuthMiddleware) SetFlash(w http.ResponseWriter, r *http.Request, message string) error {
	sess, err := m.session(r)
	if err != nil {
		return err
	}
	sess.AddFlash(message)
	return sess.Save(r, w)
}

func (m *AuthMiddleware) PullFlash(w http.ResponseWriter, r *http.Request) string {
	sess, err := m.session(r)
	if err != nil {
		return ""
	}
	flashes := sess.Flashes()
	_ = sess.Save(r, w)
	if len(flashes) == 0 {
		return ""
	}
	if msg, ok := flashes[0].(string); ok {
		return msg
	}
	return ""
}
