package router

import (
	"KillShopping/controllers"
	"KillShopping/middleware"
	"KillShopping/models"
	"KillShopping/repositories"
	"KillShopping/services"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"log"
)

func InitRouter() *gin.Engine {
	var userController controllers.UserController
	var commodityController controllers.CommodityController
	var orderController controllers.OrderController

	// 依赖注入
	var injector inject.Graph
	err := injector.Provide(
		&inject.Object{Value: &repositories.UserManagerRepository{Db: models.MysqlHandler}},
		&inject.Object{Value: &services.UserService{}},
		&inject.Object{Value: &userController},

		&inject.Object{Value: &repositories.CommodityRepository{Db: models.MysqlHandler}},
		&inject.Object{Value: &services.CommodityService{}},
		&inject.Object{Value: &commodityController},

		&inject.Object{Value: &repositories.OrderRepository{Db: models.MysqlHandler}},
		&inject.Object{Value: &services.OrderService{}},
		&inject.Object{Value: &orderController},
	)
	if err != nil {
		log.Fatal("inject fatal: ", err)
	}
	if err := injector.Populate(); err != nil {
		log.Fatal("inject fatal: ", err)
	}

	//gin
	app := gin.Default()
	api := app.Group("/api")
	{
		// login界面就会生成jwtToken, 来以备后续使用
		api.POST("/login", userController.Login)
		api.POST("/register", userController.Register)
		// 这个根据邮箱获取用户的个人信息， 相当于展示个人信息的页面
		api.GET("/me", middleware.Auth(), userController.Info)
		// 下面的加了一个 middleware.Admin() 是检测是否有购买的权限 authority
		api.POST("/", middleware.Auth(), middleware.Admin(), commodityController.AddCommodity)
		api.DELETE("/commodity/:id", middleware.Auth(), middleware.Admin(), commodityController.DelCommodity)
		api.GET("/commodity/:id", middleware.Auth(), middleware.Admin(), commodityController.GetCommodityById)
		api.GET("/commodity", middleware.Auth(), middleware.Admin(), commodityController.GetCommodity)
		api.PUT("/commodity/:id", middleware.Auth(), middleware.Admin(), commodityController.UpdateCommodity)

		api.GET("/order", middleware.Auth(), middleware.Admin(), orderController.Get)
	}
	return app
}
