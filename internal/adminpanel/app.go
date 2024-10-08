package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"github.com/go-advanced-admin/admin/internal/form/fields"
	"github.com/go-advanced-admin/admin/internal/logging"
	"github.com/go-advanced-admin/admin/internal/utils"
	"net/http"
	"reflect"
	"strings"
)

// App represents an application within the admin panel, grouping related models together.
type App struct {
	Name        string
	DisplayName string
	Models      map[string]*Model
	ModelsSlice []*Model
	Panel       *AdminPanel
	ORM         ORMIntegrator
}

// CreateViewLog creates a log entry when the app is viewed.
func (a *App) CreateViewLog(ctx interface{}) error {
	return a.Panel.Config.CreateLog(ctx, logging.LogStoreLevelPanelView, a.Name, nil, "", "")
}

// GetORM returns the ORM integrator for the app.
func (a *App) GetORM() ORMIntegrator {
	if a.ORM != nil {
		return a.ORM
	}
	return a.Panel.GetORM()
}

// RegisterModel registers a model with the app, making it available in the admin interface.
func (a *App) RegisterModel(model interface{}, orm ORMIntegrator) (*Model, error) {
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

	var fieldConfigs []FieldConfig
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldName := field.Name
		fieldDisplayName := utils.HumanizeName(fieldName)
		includeInList := true
		includeInFetch := true
		includeInSearch := true
		includeInInstanceView := true
		includeInAddForm := true
		includeInEditForm := true
		var formAddField form.Field
		var formEditField form.Field

		tag := field.Tag.Get("admin")
		if tag != "" {
			listFetchTagPresent := false
			parsedTags := strings.Split(tag, ";")
			for _, t := range parsedTags {
				pair := strings.SplitN(t, ":", 2)
				var key, value string
				if len(pair) >= 2 {
					key, value = pair[0], pair[1]
				} else {
					key = pair[0]
				}

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
				case "editForm":
					if value == "exclude" {
						includeInEditForm = false
					} else if value == "include" {
						includeInEditForm = true
					} else {
						return nil, fmt.Errorf("invalid value for 'editForm' tag: %s", value)
					}
				case "displayName":
					fieldDisplayName = value
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

		var formField form.Field
		if includeInAddForm || includeInEditForm {
			switch fieldType.Kind() {
			case reflect.String:
				formField = &fields.TextField{}
				if tag != "" {
					parsedTags := strings.Split(tag, ";")
					for _, t := range parsedTags {
						pair := strings.SplitN(t, ":", 2)
						var key, value string
						if len(pair) >= 2 {
							key, value = pair[0], pair[1]
						} else {
							key = pair[0]
						}

						switch key {
						case "placeholder":
							formField.(*fields.TextField).Placeholder = &value
						case "required":
							formField.(*fields.TextField).Required = true
						case "regex":
							formField.(*fields.TextField).Regex = &value
						case "maxLength":
							maxLengthInterface, err := utils.ConvertStringToType(value, reflect.TypeOf(uint(0)))
							if err != nil {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							maxLength, ok := maxLengthInterface.(uint)
							if !ok {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							formField.(*fields.TextField).MaxLength = &maxLength
						case "minLength":
							minLengthInterface, err := utils.ConvertStringToType(value, reflect.TypeOf(uint(0)))
							if err != nil {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							minLength, ok := minLengthInterface.(uint)
							if !ok {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							formField.(*fields.TextField).MinLength = &minLength
						}
					}
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				formField = &fields.IntegerField{}
				if tag != "" {
					parsedTags := strings.Split(tag, ";")
					for _, t := range parsedTags {
						pair := strings.SplitN(t, ":", 2)
						var key, value string
						if len(pair) >= 2 {
							key, value = pair[0], pair[1]
						} else {
							key = pair[0]
						}

						switch key {
						case "required":
							formField.(*fields.IntegerField).Required = true
						case "max":
							maxInterface, err := utils.ConvertStringToType(value, reflect.TypeOf(0))
							if err != nil {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							maxValue, ok := maxInterface.(int)
							if !ok {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							formField.(*fields.IntegerField).MaxValue = &maxValue
						case "min":
							minInterface, err := utils.ConvertStringToType(value, reflect.TypeOf(0))
							if err != nil {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							minValue, ok := minInterface.(int)
							if !ok {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							formField.(*fields.IntegerField).MinValue = &minValue
						}
					}
				}
			case reflect.Float32, reflect.Float64:
				formField = &fields.FloatField{}
				if tag != "" {
					parsedTags := strings.Split(tag, ";")
					for _, t := range parsedTags {
						pair := strings.SplitN(t, ":", 2)
						var key, value string
						if len(pair) >= 2 {
							key, value = pair[0], pair[1]
						} else {
							key = pair[0]
						}

						switch key {
						case "required":
							formField.(*fields.FloatField).Required = true
						case "max":
							maxInterface, err := utils.ConvertStringToType(value, reflect.TypeOf(float64(0)))
							if err != nil {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							maxValue, ok := maxInterface.(float64)
							if !ok {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							formField.(*fields.FloatField).MaxValue = &maxValue
						case "min":
							minInterface, err := utils.ConvertStringToType(value, reflect.TypeOf(float64(0)))
							if err != nil {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							minValue, ok := minInterface.(float64)
							if !ok {
								return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
							}
							formField.(*fields.FloatField).MinValue = &minValue
						}
					}
				}
			case reflect.Bool:
				formField = &fields.BooleanField{}
				if tag != "" {
					parsedTags := strings.Split(tag, ";")
					for _, t := range parsedTags {
						pair := strings.SplitN(t, ":", 2)
						key := pair[0]

						switch key {
						case "required":
							formField.(*fields.BooleanField).Required = true
						}
					}
				}
			default:
				func() {}() // Nothing happens for this
			}
			if formField != nil && tag != "" {
				parsedTags := strings.Split(tag, ";")
				for _, t := range parsedTags {
					pair := strings.SplitN(t, ":", 2)
					var key, value string
					if len(pair) >= 2 {
						key, value = pair[0], pair[1]
					} else {
						key = pair[0]
					}

					switch key {
					case "initial":
						convertedValue, err := utils.ConvertStringToType(value, fieldType)
						if err != nil {
							return nil, fmt.Errorf("error converting value '%s' to type '%s': %w", value, fieldType.Name(), err)
						}
						formField.RegisterInitialValue(convertedValue)
					}
				}
			}
		}

		if includeInAddForm {
			formAddField = formField
		}
		if includeInEditForm {
			formEditField = formField
		}

		fieldGenerator, implemented := model.(AdminFormFieldInterface)
		if implemented {
			formFieldForAdd := fieldGenerator.AdminFormField(fieldName, false)
			if formFieldForAdd != nil {
				formAddField = formFieldForAdd
			}
			formFieldForEdit := fieldGenerator.AdminFormField(fieldName, true)
			if formFieldForEdit != nil {
				formEditField = formFieldForEdit
			}
		}

		fieldConfigs = append(fieldConfigs, FieldConfig{
			Name:                  fieldName,
			DisplayName:           fieldDisplayName,
			FieldType:             fieldType,
			IncludeInListDisplay:  includeInList,
			IncludeInListFetch:    includeInFetch,
			IncludeInSearch:       includeInSearch,
			IncludeInInstanceView: includeInInstanceView,
			AddFormField:          formAddField,
			EditFormField:         formEditField,
		})
	}

	modelInstance := &Model{
		Name:        name,
		DisplayName: displayName,
		PTR:         model,
		App:         a,
		Fields:      fieldConfigs,
		ORM:         orm,
	}
	a.Panel.Web.HandleRoute("GET", a.Panel.Config.GetPrefix()+modelInstance.GetLink(), modelInstance.GetViewHandler())
	a.Panel.Web.HandleRoute("GET", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/:id/view", modelInstance.GetInstanceViewHandler())
	a.Panel.Web.HandleRoute("DELETE", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/:id/view", modelInstance.GetInstanceDeleteHandler())
	a.Panel.Web.HandleRoute("GET", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/add", modelInstance.GetAddHandler())
	a.Panel.Web.HandleRoute("POST", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/add", modelInstance.GetAddHandler())
	a.Panel.Web.HandleRoute("GET", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/:id/edit", modelInstance.GetEditHandler())
	a.Panel.Web.HandleRoute("POST", a.Panel.Config.GetPrefix()+modelInstance.GetLink()+"/:id/edit", modelInstance.GetEditHandler())
	a.ModelsSlice = append(a.ModelsSlice, modelInstance)
	a.Models[name] = modelInstance
	return modelInstance, nil
}

// GetHandler returns the HTTP handler function for the app's main page.
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

		html, err := a.Panel.Config.Renderer.RenderTemplate("app", map[string]interface{}{"app": a, "models": models, "navBarItems": a.Panel.Config.GetNavBarItems(data)})
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		err = a.CreateViewLog(data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		return http.StatusOK, html
	}
}

// GetLink returns the relative URL path to the app.
func (a *App) GetLink() string {
	return fmt.Sprintf("/a/%s", a.Name)
}

// GetFullLink returns the full URL path to the app, including the admin prefix.
func (a *App) GetFullLink() string {
	return a.Panel.Config.GetLink(a.GetLink())
}
