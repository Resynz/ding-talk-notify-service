/**
 * @Author: Resynz
 * @Date: 2020/8/4 10:55
 */
package controller

import "github.com/gin-gonic/gin"

func Ping(ctx *gin.Context) {
	ctx.String(200, "success")
}
