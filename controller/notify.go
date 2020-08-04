/**
 * @Author: Resynz
 * @Date: 2020/8/4 11:06
 */
package controller

import (
	"ding-talk-notify-service/config"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	// todo 获取解析内容，处理相关逻辑

	// todo 根据相关逻辑，推送至回调队列

	ctx.JSON(200, response)
}
