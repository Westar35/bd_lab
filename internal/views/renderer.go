package views

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// Renderer рендерит HTML шаблоны.
type Renderer struct {
	tmpl *template.Template
}

func NewRenderer(templatesDir string) (*Renderer, error) {
	funcMap := template.FuncMap{
		"eq": func(a, b any) bool { return fmt.Sprint(a) == fmt.Sprint(b) },
		"contains": func(text, part string) bool {
			return strings.Contains(strings.ToLower(text), strings.ToLower(part))
		},
		"formatDate": func(value string) string {
			if value == "" {
				return ""
			}
			if len(value) >= 10 {
				return value[:10]
			}
			return value
		},
		"formatDateTime": func(value string) string {
			if value == "" {
				return ""
			}
			layouts := []string{time.RFC3339, "2006-01-02 15:04:05-07", "2006-01-02 15:04:05", "2006-01-02 15:04"}
			for _, layout := range layouts {
				t, err := time.Parse(layout, value)
				if err == nil {
					return t.Format("02.01.2006 15:04")
				}
			}
			if len(value) >= 16 {
				return strings.Replace(value[:16], "T", " ", 1)
			}
			return value
		},
	}

	pattern := filepath.Join(templatesDir, "*.html")
	tmpl, err := template.New("all").Funcs(funcMap).ParseGlob(pattern)
	if err != nil {
		return nil, err
	}

	return &Renderer{tmpl: tmpl}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.tmpl.ExecuteTemplate(w, name, data)
}
