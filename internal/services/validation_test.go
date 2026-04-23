package services

import (
	"net/url"
	"testing"

	"bd_lab_3/internal/models"
)

func TestParseAndValidateForm_TripSheetDistanceValidation(t *testing.T) {
	entity := models.EntityMap()["trip-sheets"]

	form := url.Values{}
	form.Set("trip_date", "2026-04-01")
	form.Set("vehicle_id", "1")
	form.Set("driver_id", "1")
	form.Set("department_id", "1")
	form.Set("odometer_start", "1000")
	form.Set("odometer_end", "1100")
	form.Set("route", "Маршрут")
	form.Set("purpose", "Цель")
	form.Set("distance_km", "90")

	_, _, errs := ParseAndValidateForm(entity, form)
	if errs["distance_km"] == "" {
		t.Fatalf("ожидалась ошибка валидации поля distance_km")
	}
}

func TestParseAndValidateForm_VinLengthValidation(t *testing.T) {
	entity := models.EntityMap()["vehicles"]

	form := url.Values{}
	form.Set("vin", "SHORTVIN")
	form.Set("reg_number", "A123BC77")
	form.Set("make", "Toyota")
	form.Set("model", "Camry")
	form.Set("year", "2020")
	form.Set("vehicle_class_id", "1")
	form.Set("status_id", "1")
	form.Set("fuel_type_id", "1")
	form.Set("transmission_type_id", "1")
	form.Set("acquisition_type_id", "1")
	form.Set("acquisition_date", "2020-01-01")
	form.Set("acquisition_cost", "1000000")
	form.Set("current_odometer_km", "10000")

	_, _, errs := ParseAndValidateForm(entity, form)
	if errs["vin"] == "" {
		t.Fatalf("ожидалась ошибка длины VIN")
	}
}
