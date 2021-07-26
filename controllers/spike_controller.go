package controllers

import (
	R "KillShopping/response"
	"KillShopping/services"
	"KillShopping/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

)

type SpikeController struct {
	SpikeService services.SpikeServiceImp
}

func (c *SpikeController) Shopping(ctx *gin.Context) {
	s, ok := ctx.Get("spikeServiceUri")
	spikeServiceUri := s.(services.SpikeServiceUri)
	if ok {
		userInfo, ok := ctx.Get("jwtUserInfo")
		if !ok {
			R.Error(ctx, "系统错误", nil)
			return
		}
		info := userInfo.(utils.JwtUserInfo)
		if err := c.SpikeService.Shopping(&info, spikeServiceUri.Id); err == nil {
			R.Ok(ctx, "抢购成功！", nil)
			return
		} else {
			fmt.Println(err.Error())
			R.Response(ctx, http.StatusCreated, err.Error(), nil, http.StatusCreated)
			return
		}
	} else {
		R.Response(ctx, http.StatusUnprocessableEntity, "参数错误", nil, http.StatusUnprocessableEntity)
		return
	}
}
