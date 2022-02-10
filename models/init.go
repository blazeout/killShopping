package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

// MysqlHandler 全局的 *gorm.DB
var (
	MysqlHandler *gorm.DB
)

func mysqlBuild() *gorm.DB {
	var err error
	// 格式 账户名 : 密码@链接方式(地址和端口号)/表名
	DB, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			"root",
			"root",
			"127.0.0.1:3306",
			"shopping"))

	if err != nil {
		log.Panicf("models.Setup err: %v", err)
		return nil
	}

	//	设置表前缀
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return setting.DatabaseSetting.TablePrefix + defaultTableName
	//}

	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	return DB
}

func Init() {
	MysqlHandler = mysqlBuild()
}
