package string2

import (
	"fmt"
	"testing"
)

/*
	API server listening at: [::]:59485
	=== RUN   TestStrConv
	4713F404BA283A49F648DC61F3B69305
	--- PASS: TestStrConv (7.53s)
	PASS
*/
func TestStrConv(t *testing.T) {
	var helper IStringHelper
	helper = new(stringHelper)
	fmt.Println(helper.RandToken(9))
}
