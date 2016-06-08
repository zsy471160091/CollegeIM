package protocol

import (
	"bytes"
	"encoding/binary"
)

// 包头
const (
	ConstHeader       = "CollegeIM"
	ConstHeaderLength = 9
	ConstMsgLength    = 4
)

// 封包
func Enpack(message []byte) []byte {
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)

}

// 解包
func Depack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i++ {
		if length < i+ConstHeaderLength+ConstMsgLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstMsgLength])
			if length < i+ConstHeaderLength+ConstMsgLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstMsgLength : i+ConstHeaderLength+ConstMsgLength+messageLength]
			readerChannel <- data
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
