package adminpanel

type Model struct {
	Name        string
	DisplayName string
	PTR         interface{}
}

type AdminModelNameInterface interface {
	AdminName() string
}

type AdminModelDisplayNameInterface interface {
	AdminDisplayName() string
}
