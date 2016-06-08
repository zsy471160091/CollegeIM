package session

import (
	// "encoding/json"
	"github.com/donnie4w/go-logger/logger"
	"github.com/donnie4w/json4g"
	"io"
	"net"
	"process"
	"protocol"
	"time"
)

// 长连接超时时间
var timeOut = 120

type Session struct {
	id    string
	laddr string
	conn  net.Conn

	// 读写
	sendChan chan []byte
	recvChan chan []byte

	// 关闭session
	closeChan chan string
}

func NewSession(conn net.Conn) *Session {
	session := &Session{
		laddr:     conn.RemoteAddr().String(),
		conn:      conn,
		sendChan:  make(chan []byte, 8),
		recvChan:  make(chan []byte, 8),
		closeChan: make(chan string),
	}

	return session
}

// 开启客户端会话服务（接收、反馈、处理）
func (session *Session) Run() {
	logger.Info(session.laddr, "服务协程启动")
	go session.reader()
	go session.writer()
	go session.process()
}

// 客户端失去连接（主动离线或其它原因）后，关闭开启的协程
func (session *Session) Close() {
	logger.Info(session.laddr, "结束连接")
	session.closeChan <- "close"
	session.closeChan <- "close"
	// session.closeChan <- "close"

	defer session.conn.Close()
	defer close(session.recvChan)
	defer close(session.sendChan)
	defer close(session.closeChan)

	process.OffLine(session.id)
}

// 用于接收客户端发送的请求
func (session *Session) reader() {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//默认以1K大小从socket缓冲区中取数据（取出的数据可能刚好为一个请求/从个请求/部分请求）
	buffer := make([]byte, 1024)
	for {
		select {
		case <-session.closeChan:
			logger.Debug("session.reader 收到关闭信号，关闭读协程")
			return

		default:
			n, err := session.conn.Read(buffer)
			// logger.Info("收到消息", buffer, "大小为", n)
			if err != nil {
				if err == io.EOF {
					logger.Info("客户端", session.laddr, "主动关闭连接")

				} else {
					logger.Warn(session.laddr, " 连接出错，错误信息: ", err)
					// session.timer.Reset(time.Second * 2)
					// session.closeChan <- "close"

				}
				session.Close()
				return

			}
			// 取出一条完整请求后通过rescvChan发送给处理协程，并返回buffer缓冲区剩余的数据
			tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), session.recvChan)
		}
	}
}

// 用于给客户端反馈请求结果
func (session *Session) writer() {
	for {
		select {
		case msg := <-session.sendChan:
			logger.Debug("发送给客户端", session.laddr, ":", string(msg))
			session.conn.Write(protocol.Enpack(msg))
		case <-session.closeChan:
			logger.Debug("session.writer 收到关闭信号，关闭写协程")
			return
		}
	}
}

// 初步处理请求（判断请求类型，调用相应处理函数）
func (session *Session) process() {
	for {
		select {
		case msg := <-session.recvChan:
			{
				if len(msg) > 0 {
					// session.timer.Reset(time.Second * time.Duration(timeOut))
					session.conn.SetDeadline(time.Now().Add(time.Duration(timeOut) * time.Second))

					logger.Info("客户机", session.laddr, "的请求:", string(msg))
					cmd := getCMD(msg)

					logger.Debug("解析出请求类型:", cmd)
					opt, ok := process.CMD_PROCESS[cmd]
					if ok {
						go opt(msg, session.sendChan)
					}
					if cmd == protocol.IM_LOGIN {
						jroot, err := json4g.LoadByString(string(msg))
						if err != nil {
							logger.Debug("json解析失败:", err)
							return
						}
						jnode := jroot.GetNodeByName("id")
						session.id = jnode.ValueString
					}
				}
			}
		case <-session.closeChan:
			logger.Debug("session.process 收到关闭信号，关闭处理协程")
			return
		}
	}
}

func getCMD(msg []byte) int {
	jroot, err := json4g.LoadByString(string(msg))
	if err != nil {
		logger.Warn("json解析失败:", err)
		return -1
	}
	jnode := jroot.GetNodeByName("cmd")
	return int(jnode.ValueNumber)

}
