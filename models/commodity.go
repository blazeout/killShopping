package models

import (
	"github.com/jinzhu/gorm"
)

// Commodity 商品
// Name 属性 Link 链接属性 Price 价格属性 Stock 剩余库存属性 StartTime 开始售卖时间属性
type Commodity struct {
	gorm.Model
	Name      string `json:"name" gorm:"type:varchar(100);not null"`
	Link      string `json:"link" gorm:"type:varchar(100);not null"`
	Price     string `json:"price" gorm:"not null"`
	Stock     int    `json:"stock" gorm:"not null"`
	StartTime int64  `json:"startTime" gorm:"not null"`
}
