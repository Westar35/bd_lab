package handlers

import (
	"net/http"

	"bd_lab_3/internal/models"
)

type dashboardSection struct {
	Name        string
	Path        string
	Description string
}

type dashboardData struct {
	BaseData models.BasePageData
	Sections []dashboardSection
}

type helpData struct {
	BaseData models.BasePageData
}

// Dashboard отображает главную страницу с навигацией.
func (a *App) Dashboard(w http.ResponseWriter, r *http.Request) {
	sections := []dashboardSection{
		{Name: "Автомобили", Path: "/vehicles", Description: "Карточки транспорта, статус, пробег, типы топлива и классы."},
		{Name: "Водители", Path: "/drivers", Description: "Учет водителей и их удостоверений."},
		{Name: "Подразделения", Path: "/departments", Description: "Структура подразделений автопарка."},
		{Name: "Назначения ТС", Path: "/assignments", Description: "Закрепление автомобилей за водителями и подразделениями по периодам."},
		{Name: "Путевые листы", Path: "/trip-sheets", Description: "Поездки, маршруты и пробег."},
		{Name: "Топливные операции", Path: "/fuel-txns", Description: "Заправки, суммы, тип оплаты и одометр."},
		{Name: "Обслуживание и ремонты", Path: "/maintenance-orders", Description: "Работы по ТО/ремонту и их стоимость."},
		{Name: "Контрагенты", Path: "/counterparties", Description: "Справочник контрагентов (юр/физ лица)."},
		{Name: "Договоры", Path: "/contracts", Description: "Договоры с типом, сроками и статусами."},
		{Name: "События аренды", Path: "/rental-events", Description: "Выдача/возврат автомобилей по договорам аренды."},
		{Name: "Поиск", Path: "/search", Description: "Обязательные учебные поисковые запросы."},
		{Name: "Отчеты", Path: "/reports", Description: "Сводные и статистические отчеты по автопарку."},
		{Name: "Справка", Path: "/help", Description: "Описание предметной области и инструкции по работе."},
	}

	data := dashboardData{
		BaseData: a.baseData(w, r, "Главная", "dashboard"),
		Sections: sections,
	}

	a.render(w, r, "dashboard.html", data)
}

// HelpPage отображает страницу справки/о проекте.
func (a *App) HelpPage(w http.ResponseWriter, r *http.Request) {
	data := helpData{BaseData: a.baseData(w, r, "Справка", "help")}
	a.render(w, r, "help.html", data)
}
