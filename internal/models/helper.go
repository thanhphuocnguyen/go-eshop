package models

import (
	"errors"
	"fmt"
	"reflect"
)

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func (s *BrandRule) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *CategoryRule) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *CustomerSegmentRule) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FirstTimeBuyerRule) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *PurchaseQuantityRule) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
