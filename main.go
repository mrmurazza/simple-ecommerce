package main

import (
	"ecommerce/domain/order"
	orderImpl "ecommerce/domain/order/impl"
	userImpl "ecommerce/domain/user/impl"
	"ecommerce/handler"
	"ecommerce/pkg/auth"
	"ecommerce/pkg/database"
	"fmt"

	"github.com/gin-gonic/gin"
	cron "github.com/robfig/cron/v3"
)

func main() {
	r := gin.Default()

	database.InitDatabase()

	// init auth & user service & repo
	authSvc := auth.InitAuthService()
	userRepo := userImpl.NewRepo(database.DB)
	userSvc := userImpl.NewService(userRepo, authSvc)

	// ini order service & repo
	orderRepo := orderImpl.NewRepo(database.DB)
	orderSvc := orderImpl.NewService(orderRepo, userSvc)

	// init handler
	userApiHandler := handler.NewUserApiHandler(userSvc)
	orderApiHandler := handler.NewOrderApiHandler(orderSvc)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", userApiHandler.Login)
		v1.POST("/user", userApiHandler.CreateUser)

		authorized := v1.Group("", authSvc.AuthenticateMiddleware())
		authorized.POST("/orders", orderApiHandler.Order)
		authorized.GET("/orders", orderApiHandler.GetOrderHistories)
		authorized.GET("/products", orderApiHandler.GetAllProducts)

		admin := v1.Group("", authSvc.AuthenticateAdmin())
		admin.GET("/admin/orders", orderApiHandler.GetAllOrders)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080

	initWorker(orderSvc)
}

func initWorker(orderSvc order.Service) {
	// Cron
	c := cron.New()
	c.AddFunc("@midnight", func() {
		err := orderSvc.RemindPendingOrder()

		if err != nil {
			fmt.Print("Error running jobs remind pending order")
			// 	logger.Infof("Error Publish By Cron err : %v", err)
		}
	})

	c.Start()
}
