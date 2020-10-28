package main

import (
	"fmt"
	"go-until/sync/timer"
	"os"
)

func main() {

	// fmt.Println("Hello World")

	/*
		test of timer

	*/

	tr := timer.NewTimer()

	if err := tr.AddTask("@every 1s", func() {
		fmt.Println("Hello World !")
	}); err != nil {
		fmt.Printf("error to add task:%s", err)
		os.Exit(-1)
	}

	tr.Start()

	// 阻塞
	select {}
}
