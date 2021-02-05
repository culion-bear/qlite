package network

import (
	"github.com/dgrijalva/jwt-go"
	jardiniere "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"math/rand"
	"time"
)

func initToken() (tokenHandle *jardiniere.Middleware){
	tokenHandle= jardiniere.New(jardiniere.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(TokenKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(ctx iris.Context, err error){
			sendError(ctx,ErrToken,err.Error())
		},
	})
	return
}

func tokenMessage(database int32) (tokenString string,err error){
	message:=jwt.MapClaims{
		"iat":time.Now().Unix(),
		"exp":time.Now().Add(24*time.Hour * time.Duration(1)).Unix(),
		"iss":"QLite",
		"number":database,
		"rand":rand.Int(),
	}
	tokenString,err=jwt.NewWithClaims(jwt.SigningMethodHS256,message).SignedString([]byte(TokenKey))
	return
}

func getTokenNumber(ctx iris.Context) int{
	return int(ctx.Values().Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)["number"].(float64))
}