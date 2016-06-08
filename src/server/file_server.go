package main

// func FileServer() {
// 	exit := make(chan bool)
// 	go func() {
// 		listener, err := net.Listen("tcp", Conf.FileServer_Laddr) //TCPListener listen
// 		if err != nil {
// 			logger.Warn("Initialize error", err.Error())
// 			exit <- true
// 			return
// 		}

// 		logger.Info("FileServer listening...")
// 		tcpcon, err := listener.AcceptTCP() //TCPConn client
// 		if err != nil {
// 			logger.Warn(err.Error())
// 			//continue
// 		}

// 		logger.Info("Client connect")
// 		data := make([]byte, 1024)
// 		if err != nil {
// 			fmt.Println("tcpcon.Read(data)" + err.Error())
// 		}

// 		//recv file name
// 		wc, err := tcpcon.Read(data)
// 		fo, err := os.Create("./" + string(data[0:wc]))
// 		if err != nil {
// 			fmt.Println("os.Create" + err.Error())
// 		}
// 		fmt.Println("the file's name is:", string(data[0:wc]))
// 		//recb file size
// 		wc, err = tcpcon.Read(data)
// 		fmt.Println("the file's size is:", string(data[0:wc]))
// 		defer fo.Close()

// 		for {
// 			c, err := tcpcon.Read(data) //???为何调用conn类的Read
// 			if err != nil {
// 				fmt.Println("tcpcon.Read(data)" + err.Error())
// 			}
// 			if string(data[0:c]) == "filerecvend" {
// 				fmt.Println("string(data[0:c]) == filerecvend is true")
// 				tcpcon.Write([]byte("file recv finished!\r\n"))
// 				tcpcon.Close()
// 				break
// 			}
// 			//write to the file
// 			_, err = fo.Write(data[0:c])
// 			if err != nil {
// 				fmt.Println("write err" + err.Error())
// 			}
// 		}
// 	}()
// 	<-exit
// 	fmt.Println(Show("server close!"))
// }
