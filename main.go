package main

func main() {

	// fmt.Println("Hello World")

	/*
		test of timer

	*/

	//tr := timer.NewTimer()
	//
	//if err := tr.AddTask(timer.SetSecond(1), func() {
	//	fmt.Println("Hello World !")
	//}); err != nil {
	//	fmt.Printf("error to add task:%s", err)
	//	os.Exit(-1)
	//}
	//
	//tr.Start()
	//
	//// 阻塞
	//select {}

	/*
		测试文件哈希 PASS
	*/
	//file, err := os.Open("C:\\Users\\Administrator\\Desktop\\QQ图片20201208155409.png")
	//if err != nil {
	//	panic(err)
	//}
	//
	//var str = string2.Helper.CreateFileHash(file)
	//
	//fmt.Println(str)

	// 测试 logrus

	//log := logger.NewLog(logger.Parameter{
	//	Level:               4,
	//	ReportCaller:        true,
	//	Fields:              nil,
	//	Hook:                nil,
	//	IO:                  os.Stdout,
	//	RegisterExitHandler: nil,
	//	DeferExitHandler:    nil,
	//})
	//
	//log.SetFormatter(&logger.LogoutsFormatter{})
	//
	//log.Infoln("这是一条来自logger日志的输出信息！")

}
