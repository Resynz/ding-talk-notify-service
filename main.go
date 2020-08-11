/**
 * @Author: Resynz
 * @Date: 2020/7/30 11:54
 */
package main

import (
	"ding-talk-notify-service/config"
	"ding-talk-notify-service/queue"
	"ding-talk-notify-service/server"
	"log"
)

func main() {
	if err := config.InitEnv(); err != nil {
		log.Fatalf("init env failed! error:%v\n", err)
	}
	log.Println("starting service exit listener...")
	queue.SetSignalHandler()
	go queue.StartQueue()
	log.Println("\033[42;30m DONE \033[0m[DingTalkNotifyService] Start Success!")
	server.StartServer()
}
