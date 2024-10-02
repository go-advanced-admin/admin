package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/utils"
	"net/http"
	"reflect"
	"strings"
)

type App struct {
	Name        string
	DisplayName string
	Models      map[string]*Model
	ModelsSlice []*Model
	Panel       *AdminPanel
}

func (a *App) RegisterModel(model interface{}) (*Model, error) {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("admin model '%s' must be a pointer to a struct", modelType.Name())
	}

	modelType = modelType.Elem()
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("admin model '%s' must be a pointer to a struct", modelType.Name())
	}

	var name string
	namer, ok := model.(AdminModelNameInterface)
	if ok {
		name = namer.AdminName()
	} else {
		name = modelType.Name()
	}

	if !utils.IsURLSafe(name) {
		return nil, fmt.Errorf("admin model '%s' name is not URL safe", name)
	}

	var displayName string
	displayNamer, ok := model.(AdminModelDisplayNameInterface)
	if ok {
		displayName = displayNamer.AdminDisplayName()
	} else {
		displayName = utils.HumanizeName(name)
	}

	if _, exists := a.Models[name]; exists {
		return nil, fmt.Errorf("admin model '%s' already exists in app '%s'. Models cannot be registered more than once", name, a.Name)
	}

	var fields []FieldConfig
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldName := field.Name
		fieldDisplayName := utils.HumanizeName(fieldName)
		includeInList := true
		includeInFetch := true
		includeInSearch := true
		includeInInstanceView := true
		includeInAddForm := true

		tag := field.Tag.Get("admin")
		if tag != "" {
			listFetchTagPresent := false
			parsedTags := strings.Split(tag, ";")
			for _, t := range parsedTags {
				pair := strings.SplitN(t, ":", 2)
				key, value := pair[0], pair[1]

				switch key {
				case "listDisplay":
					if value == "exclude" {
						includeInList = false
					} else if value == "include" {
						includeInList = true
					} else {
						return nil, fmt.Errorf("invalid value for 'listDisplay' tag: %s", value)
					}
				case "listFetch":
					listFetchTagPresent = true
					if value == "exclude" {
						includeInFetch = false
					} else if value == "include" {
						includeInFetch = true
					} else {
						return nil, fmt.Errorf("invalid value for 'listFetch' tag: %s", value)
					}
				case "search":
					if value == "exclude" {
						includeInSearch = false
					} else if value == "include" {
						includeInSearch = true
					} else {
						return nil, fmt.Errorf("invalid value for 'search' tag: %s", value)
					}
				case "view":
					if value == "exclude" {
						includeInInstanceView = false
					} else if value == "include" {
						includeInInstanceView = true
					} else {
						return nil, fmt.Errorf("invalid value for 'view' tag: %s", value)
					}
				case "addForm":
					if value == "exclude" {
						includeInAddForm = false
					} else if value == "include" {
						includeInAddForm = true
					} else {
						return nil, fmt.Errorf("invalid value for 'addForm' tag: %s", value)
					}
				case "displayName":
					fieldDisplayName = value
				default:
					return nil, fmt.Errorf("unknown tag key: %s", key)
				}
			}
			if !listFetchTagPresent {
				if fieldName == "ID" {
					includeInFetch = true
				} else {
					includeInFetch = includeInList
				}
			}
		}

		fieldType := field.Type

		fields = append(fields, FieldConfig{
			Name:                  fieldName,
			DisplayName:           fieldDisplayName,
			FieldType:             fieldType,
			IncludeInListDisplay:  includeInList,
			IncludeInListFetch:    includeInFetch,
			IncludeInSearch:       includeInSearch,
			IncludeInInstanceView: includeInInstanceView,
			IncludeInAddForm:      includeInAddForm,
		})
	}

	var primaryKeyType reflect.Type
	primaryKeyGetter, err := GetPrimaryKeyGetter(model)
	if err != nil {
		return nil, fmt.Errorf("error determining primary key for model '%s': %w", name, err)
	}

	if idField, found := reflect.TypeOf(model).Elem().FieldByName("ID"); found {
		primaryKeyType = idField.Type
	} else if _, ok = model.(AdminModelGetIDInterface); ok {
		tempInstance := reflect.New(reflect.TypeOf(model).Elem()).Interface()
		idInterface := tempInstance.(AdminModelGetIDInterface).AdminGetID()
		primaryKeyType = reflect.TypeOf(idInterface)
	} else {
		return nil, fmt.Errorf("could not determine primary key type for model '%s'", name)
	}

	modelInstance := &Model{
		Name:             name,
		DisplayName:      displayName,
		PTR:              model,
		App:              a,
		Fields:           fields,
		PrimaryKeyGetter: primaryKeyGetter,
		PrimaryKeyType:   primaryKeyType,
	}
	a.Panel.Web.HandleRoute("GET", a.Panel.Config.GetPrefix()+modelInstance.GetLink(), modelInstance.GetViewHandler())
	a.Panel.Web.HandleRoute("GET", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/:id", modelInstance.GetInstanceViewHandler())
	a.Panel.Web.HandleRoute("DELETE", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/:id", modelInstance.GetInstanceDeleteHandler())
	a.Panel.Web.HandleRoute("GET", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/add", modelInstance.GetAddHandler())
	a.Panel.Web.HandleRoute("POST", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/add", modelInstance.GetAddHandler())
	a.ModelsSlice = append(a.ModelsSlice, modelInstance)
	a.Models[name] = modelInstance
	return modelInstance, nil
}

func (a *App) GetHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		allowed, err := a.Panel.PermissionChecker.HasAppReadPermission(a.Name, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("forbidden"))
		}

		models, err := GetModelsWithReadPermissions(a, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		html, err := a.Panel.Config.Renderer.RenderTemplate("app.html", map[string]interface{}{"app": a, "models": models})
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		return http.StatusOK, html
	}
}

func (a *App) GetLink() string {
	return fmt.Sprintf("/%s", a.Name)
}

func (a *App) GetFullLink() string {
	return a.Panel.Config.GetLink(a.GetLink())
}
