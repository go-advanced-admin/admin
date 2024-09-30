package adminpanel

type Model struct {
	Name string
	PTR  interface{}
}

type AdminModelNameInterface interface {
	AdminName() string
}
