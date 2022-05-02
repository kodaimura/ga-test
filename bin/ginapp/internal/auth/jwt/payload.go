package jwt

import (
	"time"
	"errors"

    "github.com/gin-gonic/gin"
    jwtpackage "github.com/golang-jwt/jwt/v4"

    "ginapp/internal/model/repository"
)

/*
Payloadに項目を追加する場合
・JwtPayloadに項目追加
・GeneratePayloadにロジック追加
・Extract項目名 (getter)追加
*/


const JwtExpires time.Duration = 108000

type JwtPayload struct {
	jwtpackage.StandardClaims

	UId int `json:"uid"`
    UserName string `json:"username"`
}


func GeneratePayload(uid int) (JwtPayload, error) {
	pl := JwtPayload{}

	ur := repository.NewUserRepository()
	user, err := ur.SelectByUId(uid)

	if err != nil {
		return pl, errors.New("GeneratePayload error")
	}

	pl.IssuedAt =  time.Now().Unix()
    pl.ExpiresAt = time.Now().Add(time.Second * JwtExpires).Unix()
	pl.UId = uid
	pl.UserName = user.UserName

	return pl, nil
}


func ExtractUId (c *gin.Context) (int, error) {
	pl := c.Keys["payload"]
	if pl == nil {
		return -1, errors.New("ExtractUId error")
	} else {
		return pl.(JwtPayload).UId, nil
	}
	
}


func ExtractUserName (c *gin.Context) (string, error) {
	pl := c.Keys["payload"]
	if pl == nil {
		return "", errors.New("ExtractUserName error")
	} else {
		return pl.(JwtPayload).UserName, nil
	}
}
