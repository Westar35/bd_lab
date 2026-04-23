package models

// NavItem элемент меню.
type NavItem struct {
	Name   string
	Path   string
	Active bool
}

// BasePageData общие данные для шаблонов.
type BasePageData struct {
	Title      string
	ActiveMenu string
	Username   string
	Auth       bool
	Flash      string
	Error      string
	NavItems   []NavItem
}

// ListQuery параметры списка.
type ListQuery struct {
	Page    int
	PerPage int
	Sort    string
	Dir     string
	Q       string
	Filters map[string]string
}

// ListPageData данные списка сущностей.
type ListPageData struct {
	BaseData    BasePageData
	Entity      EntityConfig
	Rows        []map[string]string
	Query       ListQuery
	Total       int
	TotalPages  int
	PageNumbers []int
	ListFields  []Field
	FilterSets  map[string][]Option
	QueryTail   string
	FilterTail  string
}

// FormPageData данные формы создания/редактирования.
type FormPageData struct {
	BaseData    BasePageData
	Entity      EntityConfig
	Fields      []Field
	Values      map[string]string
	Errors      map[string]string
	Selects     map[string][]Option
	Action      string
	SubmitLabel string
	IsEdit      bool
	ID          string
}

// DetailPageData данные карточки сущности.
type DetailPageData struct {
	BaseData BasePageData
	Entity   EntityConfig
	Row      map[string]string
	ID       string
}

// DeletePageData данные подтверждения удаления.
type DeletePageData struct {
	BaseData BasePageData
	Entity   EntityConfig
	Row      map[string]string
	ID       string
}
