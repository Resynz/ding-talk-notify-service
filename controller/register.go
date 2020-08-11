/**
 * @Author: Resynz
 * @Date: 2020/8/11 09:33
 */
package controller

import (
	"ding-talk-notify-service/config"
	"ding-talk-notify-service/structs"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(ctx *gin.Context) {
	body, err := ctx.GetRawData()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, nil)
		return
	}
	var reg structs.Register
	if err = json.Unmarshal(body, &reg); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, nil)
		return
	}
	config.NotifyRegisterMap[reg.InstanceId] = reg.NotifyUrl
	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "OK",
	})
}

func Check(ctx *gin.Context) {
	instanceId := ctx.Param("instanceId")
	notifyUrl, ok := config.NotifyRegisterMap[instanceId]
	if !ok {
		ctx.JSON(200, map[string]interface{}{
			"code": 101,
			"msg":  "instanceId not found",
		})
		return
	}
	data := &structs.Register{
		InstanceId: instanceId,
		NotifyUrl:  notifyUrl,
	}
	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "OK",
		"data": data,
	})
}
