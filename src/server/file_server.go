package main

import (
	"bytes"
	"encoding/binary"
	"github.com/donnie4w/go-logger/logger"
	"io"
	"net"
	"os"
)

func FileServer() {

	listener, err := net.Listen("tcp", Conf.FileServer_Laddr)

	CheckError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// consoleOutput(conn.RemoteAddr().String(), " tcp connect success")
		logger.Info(conn.RemoteAddr().String(), " tcp 连接文件服务器")
		go handle(conn)
		// go recvFile(conn)
	}
}

func handle(conn net.Conn) {

	data := make([]byte, 1024)

	c, err := conn.Read(data)
	if err != nil {
		logger.Warn(conn.RemoteAddr().String(), "命令出错，err =", err.Error())
		// os.Remove(filePath)
		conn.Close()
	} else {
		if string(data[0:c]) == "upload" {
			go recvFile(conn)
		} else if string(data[0:c]) == "download" {
			go sendFile(conn)
		}
	}

}

func recvFile(conn net.Conn) {
	defer conn.Close()

	data := make([]byte, 10*1024)

	//获取文件名（路径）
	wc, err := conn.Read(data)
	filePath := string(data[0:wc])
	fo, err := os.Create(filePath)
	if err != nil {
		logger.Warn(conn.RemoteAddr().String(), "传输文件失败，err =", err.Error())
		return
	}
	defer fo.Close()

	logger.Debug(filePath)
	for {
		c, err := conn.Read(data)
		if err != nil {
			logger.Warn(conn.RemoteAddr().String(), "传输文件失败，err =", err.Error())
			os.Remove(filePath)
			return
		}

		if string(data[0:c]) == "Finished" {
			logger.Info("文件接收完成")
			return
		}

		//write to the file
		_, err = fo.Write(data[0:c])
		if err != nil {
			logger.Warn(conn.RemoteAddr().String(), "传输文件失败，err = ", err.Error())
		}
	}
}

func sendFile(conn net.Conn) {
	defer conn.Close()

	data := make([]byte, 10*1024)
	//获取文件名（路径）
	wc, err := conn.Read(data)
	filePath := string(data[0:wc])

	// fo, err := os.Create(filePath)
	if err != nil {
		logger.Warn(conn.RemoteAddr().String(), "传输文件失败，err =", err.Error())
		return
	}
	logger.Warn("filePath", filePath)

	//打开文件
	fi, err := os.Open(filePath)
	if err != nil {
		// panic(err)
		return
	}

	defer fi.Close()
	fiinfo, err := fi.Stat()
	logger.Debug("the size of file is ", fiinfo.Size(), "bytes") //fiinfo.Size() return int64 type
	defer conn.Close()

	//发送文件大小
	// n, err := conn.Write([]byte(string(fiinfo.Size())))
	n, err := conn.Write(intToBytes(fiinfo.Size()))
	if err != nil {
		logger.Debug("conn.Write", err.Error())
		return
	}
	logger.Debug("send size ", n)

	// 发送文件
	buff := make([]byte, 1024)
	for {
		_, err := fi.Read(buff)
		// logger.Debug("conn.Write", buff)
		if err != nil && err != io.EOF {
			// panic(err)
			return
		}
		if n == 0 {
			conn.Write([]byte("filerecvend"))
			logger.Debug("filerecvend")
			return
		}
		n, err = conn.Write(buff)
		if err != nil {
			logger.Debug(err.Error())
			return
		}
	}
}

//整形转换成字节
func intToBytes(n int64) []byte {
	x := int64(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
