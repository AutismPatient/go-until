package convert

import (
	"encoding/json"
	"github.com/goinggo/mapstructure"
	"io"
	"log"
	"reflect"
)

/*
	Map转结构体
*/
func MapToStruct(mapInstance map[string]interface{}, structI interface{}) (err error) {
	return mapstructure.Decode(mapInstance, structI)
}

/*
	结构体转MAP
*/
func StructToMap(obj interface{}) (data map[string]interface{}) {
	var (
		tf = reflect.TypeOf(obj)
		tv = reflect.ValueOf(obj)
	)
	data = make(map[string]interface{})
	for i := 0; i < tf.NumField(); i++ {
		data[tf.Field(i).Name] = tv.Field(i).Interface()
	}
	return data
}

/*
	json转map
*/
func JsonToMap(r io.Reader) (mData map[string]interface{}) {
	var (
		decoder = json.NewDecoder(r)
		jsonVal interface{}
	)
	err := decoder.Decode(&jsonVal)
	if err != nil {
		panic(err)
	}
	return StructToMap(jsonVal)
}

/*
	MAP转Json
*/
func MapToJson(m map[string]interface{}) (jVal []byte) {
	var (
		structI interface{}
	)
	err := MapToStruct(m, structI)
	if err != nil {
		panic(err)
	}
	jVal, err = json.Marshal(&structI)
	if err != nil {
		log.Fatal(err.Error())
	}
	return jVal
}
