package app

import (
	"cloudStoregeDemo/pkg/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type ResultData struct {
	Code int
	Data interface{}
	Msg  string
}

func (g *Gin) Respond(httpCode, code int, data interface{}) {
	g.C.JSON(httpCode, ResultData{
		Code: code,
		Data: data,
		Msg:  e.GetMsg(code),
	})
}
