package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email     string `json:"email" gorm:"type:varchar(100);unique_index;not null"`
	Tel       string `json:"tel" gorm:"type:varchar(11);unique_index;not null"`
	UserName  string `json:"user_name" gorm:"unique;not null"`
	PassWord  string `json:"pass_word" gorm:"not null"`
	Authority int    `json:"authority" gorm:"DEFAULT:1;not null"`
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PassWord), []byte(password))
	return err == nil
}

// SetPassword 调用生成Password, bcrypt会生成一个 hash码, 然后将哈希码存入数据库User的Password中
// 后续拿着这个hash码和新传入进来的比较即可
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.PassWord = string(bytes)
	return nil
}
