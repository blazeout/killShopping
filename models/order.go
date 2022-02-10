package models

import "github.com/jinzhu/gorm"

// Order 一个订单的具体属性

type Order struct {
	gorm.Model      // 里面包含了Order ID
	UserId     uint `json:"user_id" gorm:"type:;not null"`
	User       User `gorm:"foreignkey:UserId"`

	CommodityId uint      `json:"commodity_id" gorm:"type:;not null"`
	Commodity   Commodity `gorm:"foreignkey:CommodityId"`

	OrderId string `json:"order_id" gorm:"not null"`
}
