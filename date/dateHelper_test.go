package date

import (
	"fmt"
	"testing"
	"time"
)

//=== RUN   TestGetDayTimeSortText
//11-18 21:53
//--- PASS: TestGetDayTimeSortText (0.05s)
//PASS
//
//Process finished with exit code 0
func TestGetDayTimeSortText(t *testing.T) {
	fmt.Println(GetDayTimeSortText(time.Now().Unix()))
}

//=== RUN   TestGetDayTimeText
//11-19 09:21:01
//--- PASS: TestGetDayTimeText (0.01s)
//PASS
//
//Process finished with exit code 0
func TestGetDayTimeText(t *testing.T) {
	fmt.Println(GetDayTimeText(time.Now().Unix()))
}
