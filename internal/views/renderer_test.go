package views

import "testing"

func TestNewRendererParsesTemplates(t *testing.T) {
	_, err := NewRenderer("../../templates")
	if err != nil {
		t.Fatalf("ожидалось успешное чтение шаблонов, получено: %v", err)
	}
}
