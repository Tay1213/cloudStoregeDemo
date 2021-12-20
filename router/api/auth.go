package api

import (
	"cloudStoregeDemo/pkg/app"
	"cloudStoregeDemo/pkg/constant"
	"cloudStoregeDemo/pkg/e"
	"cloudStoregeDemo/pkg/gredis"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		defer func() {
			if r := recover(); r != nil {
				c.Abort()
				appG.Respond(http.StatusUnauthorized, e.TOKEN_INVALID, nil)
			}
		}()

		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			panic("没有Authorization参数")
		}
		auth = strings.Fields(auth)[1]
		// 校验token
		claims, err := parseToken(auth)
		if err != nil {
			panic("token校验失败")
		}
		res := gredis.Exists(constant.TOKEN_PRIFIX + ":" + claims.Id + ":" + auth)
		if res {
			println("token 正确")
			gredis.Expire(constant.TOKEN_PRIFIX+":"+claims.Id+":"+auth, constant.EXPIRE_TIME)
		} else {
			panic("token过期")
		}
		c.Set("userId", claims.Id)
		c.Next()
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		defer func() {
			if r := recover(); r != nil {
				c.Abort()
				appG.Respond(http.StatusUnauthorized, e.TOKEN_INVALID, nil)
			}
		}()
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			panic("没有Authorization参数")
		}
		//可能有panic
		auth = strings.Fields(auth)[1]
		// 校验token
		claims, err := parseToken(auth)
		if err != nil {
			panic("token校验失败")
		}
		res := gredis.Exists(constant.TOKEN_PRIFIX + ":" + claims.Id + ":" + auth)
		if !res {
			panic("token过期")
		}
		_, err = gredis.Delete(constant.TOKEN_PRIFIX + ":" + claims.Id + ":" + auth)
		if err != nil {
			panic("redis删除失败")
		}
		_, err = gredis.Delete(constant.LOGIN_PRIFIX + ":" + claims.Id)
		if err != nil {
			panic("redis删除失败")
		}
		c.Next()
	}
}

func parseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(constant.JWT_SECRET), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}
