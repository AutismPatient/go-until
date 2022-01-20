package excel

import "testing"

type SumStruct struct {
	Name  string `json:"name" excel:"名称"`
	Sport string `json:"sport" excel:"测试"`
}

type TestStruct struct {
	Title string `json:"title" excel:"标题"`
	SumStruct
}

func TestNewExcelByStruct(t *testing.T) {

	data := []TestStruct{
		{
			Title:     "A9版本测试工嗯呢该想的是请你吃几块钱氨茶碱你承接上文",
			SumStruct: SumStruct{Name: "猪大肠", Sport: "666"},
		},
		{
			Title:     "acsn我采访时看出库记那份就NFC我能去那几款财务会计是带你吃鸡",
			SumStruct: SumStruct{Name: "魏梦媛", Sport: "ppl"},
		},
	}

	var dataI = make([]interface{}, len(data))
	for i, v := range data {
		dataI[i] = v
	}

	err := NewExcelByStruct("数据表格", dataI)
	if err != nil {
		t.Fatal(err)
	}
}
