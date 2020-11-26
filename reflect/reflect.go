package reflect

import (
	"fmt"
	"reflect"
	"strconv"
)

type Info struct {
	Fields   map[string]interface{}         `json:"fields"`
	ElemInfo map[string]reflect.StructField `json:"elem_info"`
	Size     int64                          `json:"size"`
}

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

// 获取结构体列名
func GetReflectFields(res interface{}) (fields []string) {
	var (
		ref = reflect.TypeOf(res)
	)
	for i := 0; i < ref.NumField(); i++ {
		fields = append(fields, ref.Field(i).Name)
	}
	return fields
}

/*
	获取结构体基本信息
*/
func GetReflectInfo(res interface{}) (info Info) {
	var (
		ref  = reflect.TypeOf(res)
		elem = ref.Elem()
		v    = reflect.ValueOf(elem)
	)
	info.Size, _ = strconv.ParseInt(strconv.Itoa(ref.NumField()), 0, 64)
	for i := 0; i < ref.NumField(); i++ {
		name := ref.Field(i).Name
		info.ElemInfo[name] = elem.Field(i)
		info.Fields[name] = v.Field(i).Interface()
	}
	return info
}
