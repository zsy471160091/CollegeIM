package main

import (
	"flag"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"net"
	"os"
	"process"
	"session"
)

var InputConfigFile = flag.String("conf_file", "config.json", "input config name")

var Conf *ConfigInfo

func main() {
	// 读取配置
	err := getConfig()
	CheckError(err)

	// 配置日志
	consoleOutput("配置日志模块...")
	setLogger()

	// 给处理模块传入数据库服务器地址
	consoleOutput("配置数据库服务器...")
	err = process.InitDB(Conf.DB_url)
	CheckError(err)

	process.InitFileAddress(Conf.FileServer_Laddr, Conf.File_Save_Dir)
	go FileServer()
	// logger.Debug(*Conf)

	// 开启监听
	netListen, err := net.Listen("tcp", Conf.Laddr)

	CheckError(err)

	defer netListen.Close()

	consoleOutput("等待客户端连接...")
	logger.Info("等待客户端连接...")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		// consoleOutput(conn.RemoteAddr().String(), " tcp connect success")
		logger.Info(conn.RemoteAddr().String(), " tcp 成功连接")
		s := session.NewSession(conn)
		s.Run()
		// s := session.NewSession(conn)
		// s.Run()
	}

}

// 读取并解析配置文件
func getConfig() error {
	// 解析命令行参数
	consoleOutput("解析命令行")
	flag.Parse()

	// 读取配置文件
	consoleOutput("读取配置文件...")
	Conf = NewConfig(*InputConfigFile)
	err := Conf.LoadConfig()
	if err != nil {
		consoleOutput("解析配置文件出错")
		return err
	}
	return nil
}

// 配置日志模块
func setLogger() {

	// 关闭日志模块的控制台输出
	logger.SetConsole(false)

	// 选择日志备份方式
	if Conf.Log.Bak_mode == "Daily" {
		consoleOutput("日志备份方式为日期")
		logger.SetRollingDaily(Conf.Log.Dir, Conf.Log.Filename)
	} else if Conf.Log.Bak_mode == "Size" {
		unit := logger.KB
		switch Conf.Log.Unit {
		case "KB":
			unit = logger.KB
		case "MB":
			unit = logger.MB
		case "GB":
			unit = logger.GB
		case "TB":
			unit = logger.TB
		default:
			consoleOutput("配置日志文件的单位出错，使用默认配置KB")
		}

		consoleOutput("日志备份方式为文件大小")
		logger.SetRollingFile(Conf.Log.Dir, Conf.Log.Filename, Conf.Log.Bak_num, Conf.Log.File_size, unit)
	}

	// 指定日志级别（高于设定级别的才输出）
	level := logger.ALL
	switch Conf.Log.Output_Level {
	case "ALL":
		level = logger.ALL
	case "DEBUG":
		level = logger.DEBUG
	case "INFO":
		level = logger.INFO
	case "WARN":
		level = logger.WARN
	case "ERROR":
		level = logger.ERROR
	case "FATAL":
		level = logger.FATAL
	case "OFF":
		level = logger.OFF
	default:
		consoleOutput("指定日志级别出错，使用默认配置ALL")
	}
	logger.SetLevel(level)
}

// 输出到控制台
func consoleOutput(v ...interface{}) {
	fmt.Println(v...)
}

// 检查关键点的错误
func CheckError(err error) {
	if err != nil {
		logger.Fatal(err)
		fmt.Fprintln(os.Stderr, "Fatal error:", err.Error())
		os.Exit(1)
	}
}
