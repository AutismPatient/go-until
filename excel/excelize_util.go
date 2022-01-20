// Package excel
// 基于 go excelize 封装的工具包
package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
)

func NewExcelByStruct(sheetName string, data []interface{}) error {

	f := excelize.NewFile()

	var (
		ColumnHeader = make([]string, 0)
		NowRow       = 1
	)

	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	for _, v := range data {
		t := reflect.TypeOf(v)
		value := reflect.ValueOf(v)
		row := make([]interface{}, 0)

		fields := GetStructFields(t)

		for i, fd := range fields {
			if tag := fd.Tag.Get("excel"); tag != "" {
				val := value.Field(i).Interface()
				row = append(row, val)
				ColumnHeader = append(ColumnHeader, tag)
			}
		}

		if NowRow == 1 {
			err := f.SetSheetRow(sheetName, "A1", &ColumnHeader)
			if err != nil {
				return err
			}
		}

		NowRow++
		err := f.SetSheetRow(sheetName, fmt.Sprintf("A%d", NowRow), &row)
		if err != nil {
			return err
		}
	}
	if err := f.SaveAs("test.xlsx"); err != nil {
		return err
	}

	return nil
}

func GetStructFields(t reflect.Type) (fields []reflect.StructField) {
	fields = make([]reflect.StructField, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			fields = append(fields, GetStructFields(field.Type)...)
		} else {
			fields = append(fields, field)
		}
	}

	return fields
}
