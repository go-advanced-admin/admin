package adminpanel

type Model struct {
	Name        string
	DisplayName string
	PTR         interface{}
	App         *App
}

type AdminModelNameInterface interface {
	AdminName() string
}

type AdminModelDisplayNameInterface interface {
	AdminDisplayName() string
}
