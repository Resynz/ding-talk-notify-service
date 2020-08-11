/**
 * @Author: Resynz
 * @Date: 2020/8/4 11:06
 */
package controller

import (
	"ding-talk-notify-service/config"
	"ding-talk-notify-service/enums"
	"ding-talk-notify-service/queue"
	"ding-talk-notify-service/structs"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type encryptBody struct {
	Encrypt string `json:"encrypt"`
}

func EventReceive(ctx *gin.Context) {
	body, err := ctx.GetRawData()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var eb encryptBody
	if err = json.Unmarshal(body, &eb); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	signature, exist := ctx.GetQuery("signature")
	if !exist {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("missing param signature"))
		return
	}
	timestamp, exist := ctx.GetQuery("timestamp")
	if !exist {
		_ = ctx.AbortWithError(http.StatusBadRequest, nil)
		return
	}
	nonce, exist := ctx.GetQuery("nonce")
	if !exist {
		_ = ctx.AbortWithError(http.StatusBadRequest, nil)
		return
	}
	log.Println("----接收到回调信息----")
	log.Printf("encrypt:%s\n", eb.Encrypt)
	log.Printf("signature:%s\n", signature)
	log.Printf("timestamp:%s\n", timestamp)
	log.Printf("nonce:%s\n", nonce)
	log.Println("---------------------")
	// 这里先校验签名是否正确
	genSign := config.DingTalkHandler.GenSignature(timestamp, nonce, eb.Encrypt)
	log.Printf("genSign:%s\n", genSign)
	if genSign != signature {
		log.Println("invalid signature")
		_ = ctx.AbortWithError(http.StatusBadRequest, nil)
		return
	}
	// 去解密
	content, result := config.DingTalkHandler.DecryptMsg(eb.Encrypt)

	if !result {
		log.Println("decrypt msg failed!")
		_ = ctx.AbortWithError(http.StatusBadRequest, nil)
		return
	}

	log.Printf("content:%s\n", string(content))

	// 这里准备返回的数据
	response, err := config.DingTalkHandler.EncryptMsg([]byte("success"))
	if err != nil {
		log.Printf("encrypt msg failed! error:%v\n", err)
		_ = ctx.AbortWithError(http.StatusBadRequest, nil)
		return
	}
	defer ctx.JSON(200, response)

	// 获取解析内容，处理相关逻辑

	// 注册回调后check
	if string(content) == "check_url" {
		return
	}

	// 获取EventType
	type eventType struct {
		EventType structs.EventType `json:"EventType"`
	}
	var et eventType
	if err = json.Unmarshal(content, &et); err != nil {
		log.Printf("解析eventType failed! error:%v\n", err)
		return
	}
	notifyUrl := ""
	switch et.EventType {
	case enums.BPMS_TASK_CHANGE, enums.BPMS_INSTANCE_CHANGE:
		type bpms struct {
			ProcessInstanceId string `json:"processInstanceId"`
			Type              string `json:"type"`
		}
		var b bpms
		if err = json.Unmarshal(content, &b); err != nil {
			log.Printf("解析 process_instance_id failed! error:%v\n", err)
			return
		}
		notify, ok := config.NotifyRegisterMap[b.ProcessInstanceId]
		if !ok {
			log.Printf("没有找到注册回调的审批实例[%s]\n", b.ProcessInstanceId)
			return
		}
		notifyUrl = notify
		// 如果是审批流结束，从map中删除
		if et.EventType == enums.BPMS_INSTANCE_CHANGE && (b.Type == "finish" || b.Type == "terminate") {
			go func() {
				time.Sleep(time.Second * 5)
				delete(config.NotifyRegisterMap, b.ProcessInstanceId)
			}()
		}
		break

		// todo more event types
	}

	if notifyUrl == "" {
		return
	}
	// 构造转发的url
	notifyUrl = fmt.Sprintf("%s?signature=%s&timestamp=%s&nonce=%s", notifyUrl, signature, timestamp, nonce)
	task := &structs.NotifyTask{
		NotifyUrl: notifyUrl,
		Body:      body,
	}
	go queue.PushToTaskQueue(task)
}
