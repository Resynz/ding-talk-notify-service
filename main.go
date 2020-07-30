/**
 * @Author: Resynz
 * @Date: 2020/7/30 11:54
 */
package main

import (
	"ding-talk-notify-service/config"
	"log"
)

func main() {
	if err := config.InitEnv(); err != nil {
		log.Fatalf("init env failed! error:%v\n", err)
	}
}
