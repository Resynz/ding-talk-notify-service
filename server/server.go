/**
 * @Author: Resynz
 * @Date: 2020/8/4 10:53
 */
package server

import (
	"ding-talk-notify-service/config"
	"ding-talk-notify-service/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func StartServer() {
	gin.SetMode(config.Mode)
	router := gin.New()
	router.MaxMultipartMemory = 8 << 20 //8mb`
	router.GET("/ping", controller.Ping)
	router.POST("/register", controller.Register)
	router.GET("/check/:instanceId", controller.Check)
	router.POST("/eventReceive", controller.EventReceive)
	if err := router.Run(fmt.Sprintf(":%d", config.AppPort)); err != nil {
		log.Fatalf("start server failed! error:%v\n", err)
	}
}
