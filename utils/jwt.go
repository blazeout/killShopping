package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const JwtSecret = "kih**&hgyshq##js"

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

type JwtUserInfo struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Authority int    `json:"authority"`
}

func (user *JwtUserInfo) GenerateToken() (string, error) {
	claim := jwt.MapClaims{
		"email":     user.Email,
		"id":        user.Id,
		"name":      user.Username,
		"authority": user.Authority,
		"nbf":       time.Now().Unix(),           // 生效时间
		"iat":       time.Now().Unix(),           // 签发时间
		"exp":       time.Now().Unix() + 3*60*60, // 过期时间
		"iss":       "coderth.cn",                // issuer 签发链接
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokens, err := token.SignedString([]byte(JwtSecret))
	return tokens, err
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	}
}

func (user *JwtUserInfo) ParseToken(tokens string) (err error) {
	token, err := jwt.Parse(tokens, secret())
	if err != nil {
		return
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to map-claim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}
	user.Email = claim["email"].(string)
	user.Username = claim["name"].(string)
	user.Authority = int(claim["authority"].(float64))
	user.Id = int(claim["id"].(float64))
	return err
}
