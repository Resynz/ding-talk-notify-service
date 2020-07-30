/**
 * @Author: Resynz
 * @Date: 2020/7/30 11:31
 */
package config

import (
	"log"
	"os"
	"strconv"
)

var (
	AppPort  = 10010       // 服务端口号
	LogPath  = "./logs"    // 日志文件路径
	LogName  = "app.log"   // 日志文件名称
	ConfPath = "./configs" // 配置文件路径
	Mode     = "debug"     // debug or product
)

func getEnv(e *string, key string) {
	if c := os.Getenv(key); c != "" {
		*e = c
	}
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
	printEnv()
	return nil
}
