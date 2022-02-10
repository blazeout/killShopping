package main

import (
	"KillShopping/controllers"
	"KillShopping/middleware"
	"KillShopping/models"
	"KillShopping/repositories"
	R "KillShopping/response"
	"KillShopping/services"
	"KillShopping/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"strconv"
	"sync"
)

//基于hash环的分布式秒杀
var (
	//分布式集群地址
	hostList = []string{"127.0.0.1", "127.0.0.1", "127.0.0.1"}
	//端口
	port = "8081"
	//记录现在的秒杀商品的数量
	commodityCache map[int]*models.Commodity
	//锁
	mutex sync.Mutex

	//hash环
	consistent utils.ConsistentHashImp
)

func main() {
	consistent = utils.NewConsistent(20)
	for _, v := range hostList {
		consistent.Add(v)
	}
	//缓存所有需要秒杀的商品的库存
	models.Init()
	models.MysqlHandler.AutoMigrate(models.Order{})
	repository := &repositories.CommodityRepository{Db: models.MysqlHandler}
	service := &services.CommodityService{CommodityRepository: repository}
	commodityList, err := service.GetCommodityAll()

	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Panic("缓存所有需要秒杀的商品的库存，获取库存失败")
		os.Exit(1)
		return
	}

	commodityCache = make(map[int]*models.Commodity)
	for _, value := range *commodityList {
		commodityCache[int(value.ID)] = &value
	}

	app := gin.Default()
	// 这个是获取本地的ip地址
	ip, err := utils.GetIp()
	if err != nil {
		utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Panic("ip获取失败")
		os.Exit(1)
		return
	}
	// ip = "127.0.0.3"
	// 生产者

	simple := services.NewRabbitMQSimple("product")
	spikeService := &services.SpikeService{
		CommodityCache:   &commodityCache,
		RabbitMqValidate: simple,
	}

	spikeController := &controllers.SpikeController{SpikeService: spikeService} //, middleware.Auth()
	// 一秒最多允许1个请求, 最大并发量无限制
	limiter := tollbooth.NewLimiter(1, nil)
	// uid是userId, id为商品id
	app.GET("/:uid/spike/:id", tollbooth_gin.LimitHandler(limiter), Ip(consistent, ip), middleware.Auth(), spikeController.Shopping)

	app.GET("/", tollbooth_gin.LimitHandler(limiter), func(context *gin.Context) {
		context.JSON(200, gin.H{"data": 1})
	})

	_ = app.Run(fmt.Sprintf(":%v", port))
}
func Ip(Consistent utils.ConsistentHashImp, LocalHost string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var spikeServiceUri services.SpikeServiceUri
		if err := c.ShouldBindUri(&spikeServiceUri); err == nil {
			c.Set("spikeServiceUri", spikeServiceUri)
			id := strconv.Itoa(spikeServiceUri.UId)
			// 这个是根据商品id去哈希环里面找, 找到这个商品顺时针距离最近的服务器节点ip地址返回
			ip, err := Consistent.Get(id)
			if err != nil {
				utils.Log.WithFields(log.Fields{"errMsg": err.Error()}).Warningln("hash环获取数据错误")
				R.Response(c, http.StatusInternalServerError, "服务器错误", err.Error(), http.StatusInternalServerError)
				c.Abort()
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			if commodityCache[spikeServiceUri.Id].Stock <= 0 {
				R.Response(c, http.StatusCreated, "商品已经卖完", nil, http.StatusCreated)
				c.Abort()
				return
			}
			// 如果相等就说明访问的是当前节点，就直接去数据层找数据，如果不是那么就走转发流程，去找目标ip节点
			if ip == LocalHost {
				c.Next()
				return
			} else {
				//代理处理
				res, _, _ := utils.GetCurl(fmt.Sprintf("http://%v:%v/%v/spike/%v", ip, port, c.Param("uid"), c.Param("id")), c.GetHeader("Authorization"))
				if res.StatusCode == 200 {
					R.Response(c, http.StatusOK, "成功抢到", nil, http.StatusOK)
					c.Abort()
					return
				} else {
					R.Response(c, http.StatusCreated, "未抢到", nil, http.StatusCreated)
					c.Abort()
					return
				}
			}
		} else {
			R.Response(c, http.StatusUnprocessableEntity, "参数错误", err.Error(), http.StatusUnprocessableEntity)
			c.Abort()
			return
		}
	}
}
