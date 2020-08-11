/**
 * @Author: Resynz
 * @Date: 2020/8/4 10:36
 */
package queue

import (
	"bytes"
	"ding-talk-notify-service/config"
	"ding-talk-notify-service/structs"
	"log"
	"net/http"
)

var (
	notifyQueue    chan *structs.NotifyTask
	notifyExitChan chan bool
	stopQueue      = false
)

func doNotify(task *structs.NotifyTask) {
	// todo 失败重试

	log.Printf("准备发送回调请求,notifyUrl:%s\n", task.NotifyUrl)
	request, err := http.NewRequest("POST", task.NotifyUrl, bytes.NewBuffer(task.Body))
	if err != nil {
		log.Printf("new http request failed! error:%v\n", err)
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("client.Do failed! error:%v\n", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("bad response status code:%d\n", resp.StatusCode)
		return
	}
	log.Println("回调成功")
}

func PushToTaskQueue(task *structs.NotifyTask) {
	notifyQueue <- task
}

func StartQueue() {
	notifyQueue = make(chan *structs.NotifyTask, config.QueueLen)
	for !stopQueue {
		task := <-notifyQueue
		if task == nil {
			continue
		}
		doNotify(task)
	}
	log.Println("回调队列结束")
	notifyExitChan <- true
}
