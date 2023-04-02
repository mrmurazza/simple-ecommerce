package main

import (
	userImpl "ecommerce/domain/user/impl"
	"ecommerce/handler"
	"ecommerce/pkg/auth"
	"ecommerce/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database.InitDatabase()

	// init service & repo
	authSvc := auth.InitAuthService()
	itemRepo := userImpl.NewRepo(database.DB)
	itemSvc := userImpl.NewService(itemRepo, authSvc)

	// init handler
	apiHandler := handler.NewApiHandler(itemSvc)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", apiHandler.Login)
		v1.POST("/user", apiHandler.CreateUser)

		authorized := v1.Group("", authSvc.AuthenticateMiddleware())
		authorized.GET("/check-auth", apiHandler.CheckAuth)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
