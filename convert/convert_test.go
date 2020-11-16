package convert

import (
	"fmt"
	"testing"
)

type User struct {
	Name string
	Age  int
}

//=== RUN   TestStructToMap
//{Name  string  0 [0] false}
//{Age  int  16 [1] false}
//map[Age:19 Name:猪大肠]
//--- PASS: TestStructToMap (0.00s)
//PASS
func TestStructToMap(t *testing.T) {
	var (
		user User
	)
	user.Name = "猪大肠"
	user.Age = 19
	data := StructToMap(user)
	fmt.Println(data)
}

//=== RUN   TestMapToStruct
//{猪大肠 0}
//--- PASS: TestMapToStruct (0.00s)
//PASS
func TestMapToStruct(t *testing.T) {
	var (
		user User
		maps = make(map[string]interface{})
	)
	maps["Name"] = "猪大肠"
	maps["ages"] = 26
	err := MapToStruct(maps, &user)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(user)
}
