package main

import (
	"encoding/json"
	"os"
)

type ConfigInfo struct {
	configfile string // 配置文件
	Laddr      string // 完整地址   eg: 127.0.0.1:2220
	DB_url     string // mongoDB数据库地址 eg:192.168.95.130:27017
	Log        struct {
		Bak_mode     string
		Dir          string
		Filename     string
		Bak_num      int32
		File_size    int64
		Unit         string
		Output_Level string
	}
	FileServer_Laddr string
}

func NewConfig(configfile string) *ConfigInfo {
	return &ConfigInfo{
		configfile: configfile,
	}
}

func (self *ConfigInfo) LoadConfig() error {

	file, err := os.Open(self.configfile)
	if err != nil {
		consoleOutput(err.Error())
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(&self)
	if err != nil {
		consoleOutput(err.Error())
		return err
	}
	return nil
}
