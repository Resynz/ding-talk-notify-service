/**
 * @Author: Resynz
 * @Date: 2020/7/30 11:31
 */
package config

import (
	"ding-talk-notify-service/lib"
	"ding-talk-notify-service/structs"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	AppPort           = 10010       // 服务端口号
	LogPath           = "./logs"    // 日志文件路径
	LogName           = "app.log"   // 日志文件名称
	ConfPath          = "./configs" // 配置文件路径
	Mode              = "debug"     // debug or release
	DingTalkHandler   *lib.DingTalkHandler
	NotifyRegisterMap structs.RegisterMap
	StorePath         = "./store" // 尚未完成的回调任务存储路径
	QueueLen          = 10
)

func getEnv(e *string, key string) {
	if c := os.Getenv(key); c != "" {
		*e = c
	}
}

func initUnHandleRegister() error {
	NotifyRegisterMap = make(structs.RegisterMap)
	f, err := os.OpenFile(fmt.Sprintf("%s/register.json", StorePath), os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	var registers []*structs.Register
	dc := json.NewDecoder(f)
	if err = dc.Decode(&registers); err != nil {
		return err
	}
	for _, v := range registers {
		NotifyRegisterMap[v.InstanceId] = v.NotifyUrl
	}
	return nil
}

func StoreUnHandleRegister() error {
	registers := make([]*structs.Register, 0)
	for k, v := range NotifyRegisterMap {
		r := &structs.Register{
			InstanceId: k,
			NotifyUrl:  v,
		}
		registers = append(registers, r)
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/register.json", StorePath), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	ec := json.NewEncoder(f)
	if err = ec.Encode(&registers); err != nil {
		return err
	}
	return nil
}

func initDingTalkHandler() error {
	f, err := os.Open(fmt.Sprintf("%s/dingtalk.json", ConfPath))
	if err != nil {
		return err
	}
	defer f.Close()
	var conf structs.DingTalkConfig
	d := json.NewDecoder(f)
	if err = d.Decode(&conf); err != nil {
		return err
	}
	DingTalkHandler, err = lib.NewDingTalkHandler(&conf)
	return err
}

func printEnv() {
	log.Printf("---Service Config ---\n")
	log.Printf("AppPort:%d\n", AppPort)
	log.Printf("Mode:%s\n", Mode)
	log.Printf("ConfPath:%s\n", ConfPath)
	log.Printf("LogPath:%s\n", LogPath)
	log.Printf("LogName:%s\n", LogName)
	log.Printf("---------------------\n")
}

// 初始化环境变量
func InitEnv() error {
	if p, err := strconv.Atoi(os.Getenv("APP_PORT")); err == nil && p > 0 {
		AppPort = p
	}
	getEnv(&LogPath, "LOG_PATH")
	getEnv(&LogName, "LOG_NAME")
	getEnv(&ConfPath, "CONF_PATH")
	getEnv(&Mode, "MODE")
	if err := initDingTalkHandler(); err != nil {
		return err
	}
	if err := initUnHandleRegister(); err != nil {
		return err
	}
	printEnv()
	return nil
}
