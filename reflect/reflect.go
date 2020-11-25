package reflect

import (
	"fmt"
	"reflect"
)

// 设置结构体某列字段的值
func SetValueByTag(result interface{}, tagName string, tagMap map[string]interface{}) error {
	t := reflect.TypeOf(result)
	if t.Name() != "" {
		return fmt.Errorf("result have to be a point")
	}
	v := reflect.ValueOf(result).Elem()
	t = v.Type()
	fieldNum := v.NumField()
	for i := 0; i < fieldNum; i++ {
		fieldInfo := t.Field(i)
		tag := fieldInfo.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		if value, ok := tagMap[tag]; ok {
			if reflect.ValueOf(value).Type() == v.FieldByName(fieldInfo.Name).Type() {
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}

// 获取结构体
func GetReflectFields(res interface{}) (fields []string) {
	var (
		ref = reflect.TypeOf(res)
	)
	for i := 0; i < ref.NumField(); i++ {
		fields = append(fields, ref.Field(i).Name)
	}
	return fields
}
