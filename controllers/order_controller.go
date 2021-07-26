package controllers

import (
	R "KillShopping/response"
	"KillShopping/services"
	"github.com/gin-gonic/gin"
	"net/http"

)

type OrderControllerImp interface {
	Get(ctx *gin.Context)
}

type OrderController struct {
	OrderService services.OrderServiceImp `inject:""`
}

func (c *OrderController) Get(ctx *gin.Context) {
	var getOrderPageService services.GetOrderPageService
	if err := ctx.ShouldBind(&getOrderPageService); err == nil {
		if data, err := c.OrderService.GetOrder(&getOrderPageService); err == nil {
			R.Ok(ctx, "成功", data)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
	}
	return
}
