package router

import (
	_ "cloudStoregeDemo/docs"
	"cloudStoregeDemo/router/api"
	v1 "cloudStoregeDemo/router/api/v1"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	user := r.Group("/user")
	{
		user.GET("/name/:name", v1.GetUserByName)
		user.GET("/email/:email", v1.GetUserByEmail)
		user.POST("/login", v1.Login)
		user.POST("/logout", api.Logout(), v1.Logout)
		user.POST("/reg", v1.Reg)
		user.DELETE("/delete/:id", api.Auth(), v1.DeleteUser)
		user.PUT("/update", api.Auth(), v1.UpdateUser)
	}

	file := r.Group("/file", api.Auth())
	{
		file.GET("/get/:id", v1.GetFile)
		file.POST("/getAll", v1.GetFiles)
		file.POST("/update", v1.UpdateFile)
		file.POST("/add", v1.AddFile)
		file.DELETE("/delete/:id", v1.DeleteFile)
	}

	return r
}
