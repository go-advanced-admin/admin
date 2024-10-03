package adminpanel

import (
	"fmt"
	"github.com/go-advanced-admin/admin/internal/form"
	"github.com/go-advanced-admin/admin/internal/form/forms"
	"github.com/go-advanced-admin/admin/internal/utils"
	"net/http"
	"reflect"
)

type Instance struct {
	InstanceID  interface{}
	Data        interface{}
	Model       *Model
	Permissions Permissions
}

func (i *Instance) GetLink() string {
	return fmt.Sprintf("%s/%v/view", i.Model.GetLink(), i.InstanceID)
}

func (i *Instance) GetFullLink() string {
	return i.Model.App.Panel.Config.GetLink(i.GetLink())
}

func (i *Instance) GetEditLink() string {
	return fmt.Sprintf("%s/%v/edit", i.Model.GetLink(), i.InstanceID)
}

func (i *Instance) GetFullEditLink() string {
	return i.Model.App.Panel.Config.GetLink(i.GetEditLink())
}

func (m *Model) GetInstanceDeleteHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		instanceIDStr := m.App.Panel.Web.GetPathParam(data, "id")
		if instanceIDStr == "" {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("instance id is required"))
		}

		primaryKeyValue := reflect.New(m.PrimaryKeyType).Elem()
		if err := utils.SetStringsAsType(primaryKeyValue, instanceIDStr); err != nil {
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

		err = m.App.Panel.ORM.DeleteInstance(m.PTR, instanceIDInterface)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		return http.StatusFound, m.GetLink()
	}
}

func (m *Model) GetInstanceViewHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		instanceIDStr := m.App.Panel.Web.GetPathParam(data, "id")
		if instanceIDStr == "" {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("instance id is required"))
		}

		primaryKeyValue := reflect.New(m.PrimaryKeyType).Elem()
		if err := utils.SetStringsAsType(primaryKeyValue, instanceIDStr); err != nil {
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

		instanceData, err := m.App.Panel.ORM.FetchInstanceOnlyFields(m.PTR, instanceIDInterface, fieldsToFetch)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		html, err := m.App.Panel.Config.Renderer.RenderTemplate("instance", map[string]interface{}{
			"model": m,
			"apps":  apps, "navBarItems": m.App.Panel.Config.GetNavBarItems(data),
			"instance": instanceData,
		})
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		return http.StatusOK, html
	}
}

type ModelAddForm struct {
	forms.BaseForm
	Model *Model
}

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
		if field.AddFormField != nil {
			fieldsToInclude = append(fieldsToInclude, field.Name)
		}
	}

	err = f.Model.App.Panel.ORM.CreateInstanceOnlyFields(instancePtr.Interface(), fieldsToInclude)
	if err != nil {
		return nil, err
	}

	return instancePtr.Interface(), nil
}

type ModelEditForm struct {
	forms.BaseForm
	Model      *Model
	InstanceID interface{}
}

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

	err = f.Model.App.Panel.ORM.UpdateInstanceOnlyFields(instancePtr.Interface(), fieldsToInclude, f.InstanceID)
	if err != nil {
		return nil, err
	}

	return instancePtr.Interface(), nil
}

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

				html, err := m.App.Panel.Config.Renderer.RenderTemplate("new_instance", map[string]interface{}{
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
			instanceID := m.PrimaryKeyGetter(instance)
			if instanceID == nil {
				return GetErrorHTML(http.StatusInternalServerError, fmt.Errorf("instance id is nil"))
			}

			instanceLink := fmt.Sprintf("%s/%v/view", m.GetFullLink(), instanceID)

			return http.StatusFound, instanceLink
		} else {
			return GetErrorHTML(http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		}
	}
}

func (m *Model) GetEditHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		instanceIDStr := m.App.Panel.Web.GetPathParam(data, "id")
		if instanceIDStr == "" {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("instance id is required"))
		}

		primaryKeyValue := reflect.New(m.PrimaryKeyType).Elem()
		if err := utils.SetStringsAsType(primaryKeyValue, instanceIDStr); err != nil {
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
			instanceID := m.PrimaryKeyGetter(instance)
			if instanceID == nil {
				return GetErrorHTML(http.StatusInternalServerError, fmt.Errorf("instance id is nil"))
			}

			instanceLink := fmt.Sprintf("%s/%v/view", m.GetFullLink(), instanceID)

			return http.StatusFound, instanceLink
		} else {
			return GetErrorHTML(http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		}
	}
}
