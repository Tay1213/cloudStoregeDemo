package main

import (
	"cloudStoregeDemo/models"
	"cloudStoregeDemo/pkg/setting"
	"cloudStoregeDemo/router"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

func init() {
	setting.Setup()
	models.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := router.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := "0.0.0.0" + fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	server.ListenAndServe()

	//r := gin.Default()
	//
	//db, _ := gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/cloud_storage?charset=utf8&parseTime=True&loc=Local")
	//defer db.Close()
	//
	//db.AutoMigrate(&User{})
	//var user1 User
	////user1.Id = 1
	//db.Where("email = ?", "222").Find(&user1)
	//fmt.Printf("user1:%#v\n", user1)
	//
	//db.Find(&user1)
	//fmt.Printf("user1:%#v\n", user1)
	//
	//user:= r.Group("/user")
	//{
	//	user.POST("/login", func(context *gin.Context) {
	//		var u1 User
	//		if err :=context.ShouldBindBodyWith(&u1, binding.JSON); err!=nil{
	//
	//		}
	//		fmt.Printf("user1:%#v\n", u1)
	//		var rd entity.ResultData
	//		context.JSON(http.StatusOK, rd.Success(u1, "成功"))
	//		fmt.Printf("rd:%#v\n", rd)
	//	})
	//}
	//
	//err := r.Run()
	//if err != nil {
	//	fmt.Printf("err:%#v\n", err)
	//}
}
