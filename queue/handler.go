/**
 * @Author: Resynz
 * @Date: 2020/8/4 10:36
 */
package queue

import (
	"ding-talk-notify-service/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func SetSignalHandler() {
	sign := make(chan os.Signal, 1)
	notifyExitChan = make(chan bool, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		defer func() {
			log.Println("关闭队列中...")
			defer close(notifyQueue)
			defer close(notifyExitChan)
			log.Println("程序退出.")
			os.Exit(0)
		}()
		s := <-sign
		log.Printf("接收到退出信号:%v\n", s)
		stopQueue = true
		notifyQueue <- nil
		log.Println("正在等待任务队列退出...")
		log.Println("正在等待回调队列退出...")
		<-notifyExitChan
		log.Println("回调队列退出完毕!")

		if len(config.NotifyRegisterMap) > 0 {
			log.Println("正在保存尚未回调的资源...")
			if err := config.StoreUnHandleRegister(); err != nil {
				log.Printf("保存失败！error:%v\n", err)
			}
		}
	}()
}
