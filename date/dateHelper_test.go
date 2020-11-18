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
