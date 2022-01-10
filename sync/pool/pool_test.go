package pool

import (
	"fmt"
	"testing"
	"time"
)

/*
 testing
*/

func TestRunPool(t *testing.T) {

	p := NewPool(10)

	f := func() error {
		for i := 1; i <= 100; i++ {
			fmt.Print(i)
		}
		return nil
	}

	task := NewTask(f)

	p.goRoutines(task)

	p.Run()

	time.Sleep(time.Second * 2)
}
