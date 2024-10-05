package adminpanel

import (
	"encoding/json"
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"github.com/go-advanced-admin/admin/internal/form/forms"
	"github.com/go-advanced-admin/admin/internal/logging"
	"github.com/go-advanced-admin/admin/internal/utils"
	"net/http"
	"reflect"
)

// Instance represents a single instance of a model in the admin panel.
type Instance struct {
	InstanceID  interface{}
	Data        interface{}
	Model       *Model
	Permissions Permissions
}

// AdminInstanceReprInterface allows customizing the string representation of an instance.
type AdminInstanceReprInterface interface {
	// AdminInstanceRepr returns a string representation of the instance.
	AdminInstanceRepr() string
}

// GetRepr returns a string representation of the instance.
func (i *Instance) GetRepr() string {
	if repr, ok := i.Data.(AdminInstanceReprInterface); ok {
		return repr.AdminInstanceRepr()
	}
	return fmt.Sprint(i.Data)
}

// CreateViewLog creates a log entry when the instance is viewed.
func (i *Instance) CreateViewLog(ctx interface{}) error {
	return i.Model.App.Panel.Config.CreateLog(ctx, logging.LogStoreLevelInstanceView, fmt.Sprintf("%s | %s", i.Model.App.Name, i.Model.DisplayName), i.InstanceID, i.GetRepr(), "")
}

// CreateUpdateLog creates a log entry when the instance is updated.
func (i *Instance) CreateUpdateLog(ctx interface{}, updates map[string]interface{}) error {
	message, err := json.Marshal(updates)
	if err != nil {
		return err
	}
	return i.Model.App.Panel.Config.CreateLog(ctx, logging.LogStoreLevelUpdate, fmt.Sprintf("%s | %s", i.Model.App.Name, i.Model.DisplayName), i.InstanceID, i.GetRepr(), string(message))
}

// CreateCreateLog creates a log entry when the instance is created.
func (i *Instance) CreateCreateLog(ctx interface{}) error {
	message, err := json.Marshal(i.Data)
	if err != nil {
		return err
	}
	return i.Model.App.Panel.Config.CreateLog(ctx, logging.LogStoreLevelCreate, fmt.Sprintf("%s | %s", i.Model.App.Name, i.Model.DisplayName), i.InstanceID, i.GetRepr(), string(message))
}

// CreateDeleteLog creates a log entry when the instance is deleted.
func (i *Instance) CreateDeleteLog(ctx interface{}) error {
	return i.Model.App.Panel.Config.CreateLog(ctx, logging.LogStoreLevelDelete, fmt.Sprintf("%s | %s", i.Model.App.Name, i.Model.DisplayName), i.InstanceID, i.GetRepr(), "")
}

// GetLink returns the relative URL to view the instance.
func (i *Instance) GetLink() string {
	return fmt.Sprintf("%s/%v/view", i.Model.GetLink(), i.InstanceID)
}

// GetFullLink returns the full URL to view the instance.
func (i *Instance) GetFullLink() string {
	return i.Model.App.Panel.Config.GetLink(i.GetLink())
}

// GetEditLink returns the relative URL to edit the instance.
func (i *Instance) GetEditLink() string {
	return fmt.Sprintf("%s/%v/edit", i.Model.GetLink(), i.InstanceID)
}

// GetFullEditLink returns the full URL to edit the instance.
func (i *Instance) GetFullEditLink() string {
	return i.Model.App.Panel.Config.GetLink(i.GetEditLink())
}

func (m *Model) GetInstanceDeleteHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		instanceIDStr := m.App.Panel.Web.GetPathParam(data, "id")
		if instanceIDStr == "" {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("instance id is required"))
		}

		primaryKeyType, err := m.GetPrimaryKeyType()
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		primaryKeyValuePtr := reflect.New(primaryKeyType)
		primaryKeyValue := primaryKeyValuePtr.Elem()

		if err = utils.SetStringsAsType(primaryKeyValue, instanceIDStr); err != nil {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("invalid instance id: %v", err))
		}

		instanceIDInterface := primaryKeyValue.Interface()

		allowed, err := m.App.Panel.PermissionChecker.HasInstanceDeletePermission(m.App.Name, m.Name, instanceIDInterface, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("you are not allowed to delete this instance"))
		}

		err = m.GetORM().DeleteInstance(m.PTR, instanceIDInterface)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		instance := &Instance{
			InstanceID: instanceIDInterface,
			Model:      m,
		}

		err = instance.CreateDeleteLog(data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		return http.StatusSeeOther, m.GetLink()
	}
}

func (m *Model) GetInstanceViewHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		instanceIDStr := m.App.Panel.Web.GetPathParam(data, "id")
		if instanceIDStr == "" {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("instance id is required"))
		}

		primaryKeyType, err := m.GetPrimaryKeyType()
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		primaryKeyValuePtr := reflect.New(primaryKeyType)
		primaryKeyValue := primaryKeyValuePtr.Elem()

		if err = utils.SetStringsAsType(primaryKeyValue, instanceIDStr); err != nil {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("invalid instance id: %v", err))
		}

		instanceIDInterface := primaryKeyValue.Interface()

		allowed, err := m.App.Panel.PermissionChecker.HasInstanceReadPermission(m.App.Name, m.Name, instanceIDInterface, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("you are not allowed to view this instance"))
		}

		apps, err := GetAppsWithReadPermissions(m.App.Panel, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		var fieldsToFetch []string
		for _, fieldConfig := range m.Fields {
			if fieldConfig.IncludeInInstanceView {
				fieldsToFetch = append(fieldsToFetch, fieldConfig.Name)
			}
		}

		instanceData, err := m.GetORM().FetchInstanceOnlyFields(m.PTR, instanceIDInterface, fieldsToFetch)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		html, err := m.App.Panel.Config.Renderer.RenderTemplate("instance", map[string]interface{}{
			"model":       m,
			"apps":        apps,
			"navBarItems": m.App.Panel.Config.GetNavBarItems(data),
			"instance":    instanceData,
		})
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		instance := &Instance{
			InstanceID: instanceIDInterface,
			Model:      m,
		}
		err = instance.CreateViewLog(data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		return http.StatusOK, html
	}
}

// ModelAddForm represents the form used to add a new instance of a model.
type ModelAddForm struct {
	forms.BaseForm
	Model *Model
}

// Save processes the form data and creates a new instance of the model.
func (f *ModelAddForm) Save(values map[string]form.HTMLType) (interface{}, error) {
	cleanValues, err := form.GetCleanData(f, values)
	if err != nil {
		return nil, err
	}

	modelType := reflect.TypeOf(f.Model.PTR).Elem()
	instancePtr := reflect.New(modelType)
	instanceVal := instancePtr.Elem()

	for fieldName, value := range cleanValues {
		fieldVal := instanceVal.FieldByName(fieldName)

		if !fieldVal.IsValid() {
			return nil, fmt.Errorf("field %s not found in model", fieldName)
		}

		if !fieldVal.CanSet() {
			return nil, fmt.Errorf("field %s is not settable", fieldName)
		}

		if value == nil {
			if fieldVal.Kind() == reflect.Ptr {
				fieldVal.Set(reflect.Zero(fieldVal.Type()))
			}
			continue
		}

		val := reflect.ValueOf(value)

		if val.Type().AssignableTo(fieldVal.Type()) {
			fieldVal.Set(val)
		} else if val.Type().ConvertibleTo(fieldVal.Type()) {
			fieldVal.Set(val.Convert(fieldVal.Type()))
		} else {
			return nil, fmt.Errorf("field %s has invalid type", fieldName)
		}
	}

	fieldsToInclude := make([]string, 0)
	for _, field := range f.Model.Fields {
		if field.AddFormField != nil {
			fieldsToInclude = append(fieldsToInclude, field.Name)
		}
	}

	err = f.Model.GetORM().CreateInstanceOnlyFields(instancePtr.Interface(), fieldsToInclude)
	if err != nil {
		return nil, err
	}

	return instancePtr.Interface(), nil
}

// ModelEditForm represents the form used to edit an existing instance of a model.
type ModelEditForm struct {
	forms.BaseForm
	Model      *Model
	InstanceID interface{}
}

// Save processes the form data and updates the existing instance of the model.
func (f *ModelEditForm) Save(values map[string]form.HTMLType) (interface{}, error) {
	cleanValues, err := form.GetCleanData(f, values)
	if err != nil {
		return nil, err
	}

	modelType := reflect.TypeOf(f.Model.PTR).Elem()
	instancePtr := reflect.New(modelType)
	instanceVal := instancePtr.Elem()

	for fieldName, value := range cleanValues {
		fieldVal := instanceVal.FieldByName(fieldName)
		if !fieldVal.IsValid() {
			continue
		}
		if !fieldVal.CanSet() {
			return nil, fmt.Errorf("field %s is not settable", fieldName)
		}
		val := reflect.ValueOf(value)
		if val.Type().AssignableTo(fieldVal.Type()) {
			fieldVal.Set(val)
		} else if val.Type().ConvertibleTo(fieldVal.Type()) {
			fieldVal.Set(val.Convert(fieldVal.Type()))
		} else {
			return nil, fmt.Errorf("field %s has invalid type", fieldName)
		}
	}

	fieldsToInclude := make([]string, 0)
	for _, field := range f.Model.Fields {
		if field.EditFormField != nil {
			fieldsToInclude = append(fieldsToInclude, field.Name)
		}
	}

	err = f.Model.GetORM().UpdateInstanceOnlyFields(instancePtr.Interface(), fieldsToInclude, f.InstanceID)
	if err != nil {
		return nil, err
	}

	return instancePtr.Interface(), nil
}

// NewAddForm creates a new form for adding an instance of the model.
func (m *Model) NewAddForm() (form.Form, error) {
	f := &ModelAddForm{
		Model: m,
	}

	for _, fieldConfig := range m.Fields {
		if fieldConfig.AddFormField == nil {
			continue
		}

		err := f.AddField(fieldConfig.Name, fieldConfig.AddFormField)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

// NewEditForm creates a new form for editing an existing instance of the model.
func (m *Model) NewEditForm(instanceID interface{}) (form.Form, error) {
	f := &ModelEditForm{
		Model:      m,
		InstanceID: instanceID,
	}

	for _, fieldConfig := range m.Fields {
		if fieldConfig.EditFormField == nil {
			continue
		}

		err := f.AddField(fieldConfig.Name, fieldConfig.EditFormField)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

// GetAddHandler returns the HTTP handler function for adding a new instance.
func (m *Model) GetAddHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		allowed, err := m.App.Panel.PermissionChecker.HasModelCreatePermission(m.App.Name, m.Name, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("forbidden"))
		}

		formInstance, err := m.NewAddForm()
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		method := m.App.Panel.Web.GetRequestMethod(data)
		if method == "GET" {
			apps, err := GetAppsWithReadPermissions(m.App.Panel, data)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}

			html, err := m.App.Panel.Config.Renderer.RenderTemplate("new_instance", map[string]interface{}{
				"apps":        apps,
				"navBarItems": m.App.Panel.Config.GetNavBarItems(data),
				"form":        formInstance,
				"model":       m,
				"formErrs":    make([]error, 0),
				"fieldErrs":   make(map[string][]error),
			})
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			return http.StatusOK, html
		} else if method == "POST" {
			formData := m.App.Panel.Web.GetFormData(data)
			if formData == nil {
				return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("form data is required"))
			}
			convertedFormData, err := form.ConvertFormDataToHTMLTypeMap(formData)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			cleanFormData, err := form.GetCleanData(formInstance, convertedFormData)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			formErrs, fieldErrs, err := form.ValuesAreValid(formInstance, cleanFormData)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			containsError := false
			if len(formErrs) > 0 {
				containsError = true
			}
			for _, errs := range fieldErrs {
				if len(errs) > 0 {
					containsError = true
					break
				}
			}

			if containsError {
				err = formInstance.RegisterInitialValues(cleanFormData)

				apps, err := GetAppsWithReadPermissions(m.App.Panel, data)
				if err != nil {
					return GetErrorHTML(http.StatusInternalServerError, err)
				}

				html, err := m.App.Panel.Config.Renderer.RenderTemplate("new_instance", map[string]interface{}{
					"apps":        apps,
					"navBarItems": m.App.Panel.Config.GetNavBarItems(data),
					"form":        formInstance,
					"model":       m,
					"formErrs":    formErrs,
					"fieldErrs":   fieldErrs,
				})
				if err != nil {
					return GetErrorHTML(http.StatusInternalServerError, err)
				}
				return http.StatusOK, html
			}

			instanceInterface, err := formInstance.Save(convertedFormData)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}

			instance := instanceInterface
			instanceID, err := m.GetPrimaryKeyValue(instance)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			if instanceID == nil {
				return GetErrorHTML(http.StatusInternalServerError, fmt.Errorf("instance id is nil"))
			}

			instanceLink := fmt.Sprintf("%s/%v/view", m.GetFullLink(), instanceID)

			instanceInstance := &Instance{
				InstanceID: instanceID,
				Data:       instance,
				Model:      m,
			}

			err = instanceInstance.CreateCreateLog(data)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}

			return http.StatusSeeOther, instanceLink
		} else {
			return GetErrorHTML(http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		}
	}
}

// GetEditHandler returns the HTTP handler function for editing an existing instance.
func (m *Model) GetEditHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		instanceIDStr := m.App.Panel.Web.GetPathParam(data, "id")
		if instanceIDStr == "" {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("instance id is required"))
		}

		primaryKeyType, err := m.GetPrimaryKeyType()
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		primaryKeyValuePtr := reflect.New(primaryKeyType)
		primaryKeyValue := primaryKeyValuePtr.Elem()

		if err = utils.SetStringsAsType(primaryKeyValue, instanceIDStr); err != nil {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("invalid instance id: %v", err))
		}

		instanceIDInterface := primaryKeyValue.Interface()

		allowed, err := m.App.Panel.PermissionChecker.HasInstanceUpdatePermission(m.App.Name, m.Name, instanceIDInterface, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("you are not allowed to view this instance"))
		}

		var fieldsToFetch []string
		for _, fieldConfig := range m.Fields {
			if fieldConfig.IncludeInInstanceView {
				fieldsToFetch = append(fieldsToFetch, fieldConfig.Name)
			}
		}

		instanceData, err := m.App.Panel.ORM.FetchInstanceOnlyFields(m.PTR, instanceIDInterface, fieldsToFetch)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		formInstance, err := m.NewEditForm(instanceIDInterface)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		initialValuesMap := make(map[string]interface{})
		for _, field := range m.Fields {
			if field.EditFormField == nil {
				continue
			}
			initialValuesMap[field.Name] = reflect.ValueOf(instanceData).Elem().FieldByName(field.Name).Interface()
		}

		err = formInstance.RegisterInitialValues(initialValuesMap)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		method := m.App.Panel.Web.GetRequestMethod(data)
		if method == "GET" {
			apps, err := GetAppsWithReadPermissions(m.App.Panel, data)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}

			html, err := m.App.Panel.Config.Renderer.RenderTemplate("edit_instance", map[string]interface{}{
				"apps": apps, "navBarItems": m.App.Panel.Config.GetNavBarItems(data),
				"form":      formInstance,
				"model":     m,
				"formErrs":  make([]error, 0),
				"fieldErrs": make(map[string][]error),
			})
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			return http.StatusOK, html
		} else if method == "POST" {
			formData := m.App.Panel.Web.GetFormData(data)
			if formData == nil {
				return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("form data is required"))
			}
			convertedFormData, err := form.ConvertFormDataToHTMLTypeMap(formData)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			cleanFormData, err := form.GetCleanData(formInstance, convertedFormData)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			formErrs, fieldErrs, err := form.ValuesAreValid(formInstance, cleanFormData)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			containsError := false
			if len(formErrs) > 0 {
				containsError = true
			}
			for _, errs := range fieldErrs {
				if len(errs) > 0 {
					containsError = true
					break
				}
			}

			if containsError {
				err = formInstance.RegisterInitialValues(cleanFormData)

				apps, err := GetAppsWithReadPermissions(m.App.Panel, data)
				if err != nil {
					return GetErrorHTML(http.StatusInternalServerError, err)
				}

				html, err := m.App.Panel.Config.Renderer.RenderTemplate("edit_instance", map[string]interface{}{
					"apps": apps, "navBarItems": m.App.Panel.Config.GetNavBarItems(data),
					"form":      formInstance,
					"model":     m,
					"formErrs":  formErrs,
					"fieldErrs": fieldErrs,
				})
				if err != nil {
					return GetErrorHTML(http.StatusInternalServerError, err)
				}
				return http.StatusOK, html
			}

			instanceInterface, err := formInstance.Save(convertedFormData)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}

			instance := instanceInterface
			instanceID, err := m.GetPrimaryKeyValue(instance)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}
			if instanceID == nil {
				return GetErrorHTML(http.StatusInternalServerError, fmt.Errorf("instance id is nil"))
			}

			instanceLink := fmt.Sprintf("%s/%v/view", m.GetFullLink(), instanceID)

			instanceInstance := &Instance{
				InstanceID: instanceID,
				Data:       instance,
				Model:      m,
			}

			err = instanceInstance.CreateUpdateLog(data, cleanFormData)
			if err != nil {
				return GetErrorHTML(http.StatusInternalServerError, err)
			}

			return http.StatusSeeOther, instanceLink
		} else {
			return GetErrorHTML(http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		}
	}
}
