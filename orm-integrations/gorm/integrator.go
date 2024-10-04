package admingorm

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type Integrator struct {
	DB *gorm.DB
}

func NewIntegrator(db *gorm.DB) *Integrator {
	return &Integrator{DB: db}
}

func (i *Integrator) GetPrimaryKeyValue(model interface{}) (interface{}, error) {
	modelValue := reflect.ValueOf(model)

	if modelValue.Kind() == reflect.Ptr {
		if modelValue.IsNil() {
			return nil, fmt.Errorf("model pointer is nil")
		}
		modelValue = modelValue.Elem()
	} else if modelValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model is neither a struct nor a pointer to a struct")
	}

	modelType := modelValue.Type()

	stmt := &gorm.Statement{DB: i.DB}
	err := stmt.Parse(reflect.New(modelType).Interface())
	if err != nil {
		return nil, fmt.Errorf("failed to parse model: %v", err)
	}

	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return nil, fmt.Errorf("no primary field found for model %s", modelType.Name())
	}

	primaryKeyValue := modelValue.FieldByName(primaryField.Name)
	if !primaryKeyValue.IsValid() {
		return nil, fmt.Errorf("primary key field %s not found in model %s", primaryField.Name, modelType.Name())
	}

	return primaryKeyValue.Interface(), nil
}

func (i *Integrator) GetPrimaryKeyType(model interface{}) (reflect.Type, error) {
	modelValue := reflect.ValueOf(model)

	if modelValue.Kind() == reflect.Ptr {
		if modelValue.IsNil() {
			return nil, fmt.Errorf("model pointer is nil")
		}
		modelValue = modelValue.Elem()
	} else if modelValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model is neither a struct nor a pointer to a struct")
	}

	modelType := modelValue.Type()

	stmt := &gorm.Statement{DB: i.DB}
	err := stmt.Parse(reflect.New(modelType).Interface())
	if err != nil {
		return nil, fmt.Errorf("failed to parse model: %v", err)
	}

	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return nil, fmt.Errorf("no primary field found for model %s", modelType.Name())
	}

	primaryKeyField, found := modelType.FieldByName(primaryField.Name)

	if !found {
		return nil, fmt.Errorf("primary key field %s not found in model %s", primaryField.Name, modelType.Name())
	}

	return primaryKeyField.Type, nil
}

func (i *Integrator) FetchInstances(model interface{}) (interface{}, error) {
	modelType := reflect.TypeOf(model).Elem()
	sliceType := reflect.SliceOf(modelType)
	instances := reflect.New(sliceType).Interface()

	err := i.DB.Find(instances, model).Error
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func (i *Integrator) FetchInstancesOnlyFields(model interface{}, fields []string) (interface{}, error) {
	modelType := reflect.TypeOf(model).Elem()
	sliceType := reflect.SliceOf(modelType)
	instances := reflect.New(sliceType).Interface()

	selectFields := make([]string, len(fields))
	for idx, fieldName := range fields {
		field, found := modelType.FieldByName(fieldName)
		if found {
			selectFields[idx] = getGormColumnName(field)
		} else {
			return nil, fmt.Errorf("field %s not foun in model", fieldName)
		}
	}

	selectFieldStr := strings.Join(selectFields, ", ")

	err := i.DB.Select(selectFieldStr).Find(instances, model).Error
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func (i *Integrator) FetchInstancesOnlyFieldWithSearch(model interface{}, fields []string, query string, searchFields []string) (interface{}, error) {
	modelType := reflect.TypeOf(model).Elem()
	sliceType := reflect.SliceOf(modelType)
	instances := reflect.New(sliceType).Interface()

	selectFields := make([]string, len(fields))
	for idx, fieldName := range fields {
		field, found := modelType.FieldByName(fieldName)
		if found {
			selectFields[idx] = getGormColumnName(field)
		} else {
			return nil, fmt.Errorf("field %s not found in model", fieldName)
		}
	}
	selectFieldStr := strings.Join(selectFields, ", ")

	var searchConditions []string
	var searchArgs []interface{}
	for _, searchField := range searchFields {
		field, found := modelType.FieldByName(searchField)
		if found {
			columnName := getGormColumnName(field)
			searchConditions = append(searchConditions, fmt.Sprintf("%s LIKE ?", columnName))
			searchArgs = append(searchArgs, "%"+query+"%")
		} else {
			return nil, fmt.Errorf("field %s not found in model", searchField)
		}
	}
	searchConditionStr := strings.Join(searchConditions, " OR ")

	err := i.DB.Select(selectFieldStr).Where(searchConditionStr, searchArgs...).Find(instances, model).Error
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func (i *Integrator) DeleteInstance(model interface{}, instanceID interface{}) error {
	modelType := reflect.TypeOf(model).Elem()

	stmt := &gorm.Statement{DB: i.DB}
	err := stmt.Parse(model)
	if err != nil {
		return fmt.Errorf("failed to parse model: %v", err)
	}

	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return fmt.Errorf("no primary field found for model %s", modelType.Name())
	}

	primaryKey := primaryField.DBName

	condition := map[string]interface{}{primaryKey: instanceID}

	result := i.DB.Where(condition).Delete(model)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("Deleted %d instances of %s\n", result.RowsAffected, modelType.Name())

	if result.RowsAffected == 0 {
		return fmt.Errorf("no instance found with %s = %v", primaryKey, instanceID)
	}

	return nil
}

func (i *Integrator) FetchInstanceOnlyFields(model interface{}, id interface{}, fields []string) (interface{}, error) {
	modelType := reflect.TypeOf(model).Elem()
	instance := reflect.New(modelType).Interface()

	selectFields := make([]string, len(fields))
	for idx, fieldName := range fields {
		field, found := modelType.FieldByName(fieldName)
		if found {
			selectFields[idx] = getGormColumnName(field)
		} else {
			return nil, fmt.Errorf("field %s not found in model", fieldName)
		}
	}
	selectFieldStr := strings.Join(selectFields, ", ")

	stmt := &gorm.Statement{DB: i.DB}
	err := stmt.Parse(model)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model: %v", err)
	}

	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return nil, fmt.Errorf("no primary field found for model %s", modelType.Name())
	}

	primaryKey := primaryField.DBName
	err = i.DB.Select(selectFieldStr).Where(fmt.Sprintf("%s = ?", primaryKey), id).First(instance, model).Error
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (i *Integrator) FetchInstance(model interface{}, instanceID interface{}) (interface{}, error) {
	modelType := reflect.TypeOf(model).Elem()
	instance := reflect.New(modelType).Interface()

	stmt := &gorm.Statement{DB: i.DB}
	err := stmt.Parse(model)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model: %v", err)
	}

	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return nil, fmt.Errorf("no primary field found for model %s", modelType.Name())
	}

	primaryKey := primaryField.DBName
	err = i.DB.Where(fmt.Sprintf("%s = ?", primaryKey), instanceID).First(instance, model).Error
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (i *Integrator) CreateInstance(instance interface{}) error {
	return i.DB.Create(instance).Error
}

func (i *Integrator) UpdateInstance(instance interface{}, primaryKey interface{}) error {
	modelType := reflect.TypeOf(instance).Elem()

	stmt := &gorm.Statement{DB: i.DB}
	if err := stmt.Parse(instance); err != nil {
		return fmt.Errorf("failed to parse model: %v", err)
	}

	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return fmt.Errorf("no primary field found for model %s", modelType.Name())
	}

	primaryKeyDBName := primaryField.DBName

	return i.DB.Model(instance).Where(fmt.Sprintf("%s = ?", primaryKeyDBName), primaryKey).Save(instance).Error
}

func (i *Integrator) CreateInstanceOnlyFields(instance interface{}, fields []string) error {
	if len(fields) == 0 {
		return i.DB.Create(instance).Error
	}

	return i.DB.Select(fields).Create(instance).Error
}

func (i *Integrator) UpdateInstanceOnlyFields(instance interface{}, fields []string, primaryKey interface{}) error {
	modelType := reflect.TypeOf(instance).Elem()
	modelValue := reflect.ValueOf(instance).Elem()

	stmt := &gorm.Statement{DB: i.DB}
	if err := stmt.Parse(instance); err != nil {
		return fmt.Errorf("failed to parse model: %v", err)
	}

	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return fmt.Errorf("no primary field found for model %s", modelType.Name())
	}

	primaryKeyDBName := primaryField.DBName

	if len(fields) == 0 {
		return i.DB.Model(instance).Where(fmt.Sprintf("%s = ?", primaryKeyDBName), primaryKey).Save(instance).Error
	}

	updateData := make(map[string]interface{})
	for _, fieldName := range fields {
		fieldValue := modelValue.FieldByName(fieldName)
		if !fieldValue.IsValid() {
			return fmt.Errorf("field %s not found in model", fieldName)
		}
		updateData[fieldName] = fieldValue.Interface()
	}
	if _, exists := updateData[primaryField.Name]; !exists {
		updateData[primaryField.Name] = primaryKey
		fields = append(fields, primaryField.Name)
	}

	return i.DB.Model(instance).Where(fmt.Sprintf("%s = ?", primaryKeyDBName), primaryKey).Select(fields).Updates(updateData).Error
}
