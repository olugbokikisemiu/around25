package main

import (
	"fmt"
	"time"

	"github.com/around-project/config"
	"github.com/around-project/handler"
	"github.com/around-project/router"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Printf("failed to load config data: %v", err)
	}

	var con config.Config

	if err := envconfig.Process("", &con); err != nil {
		fmt.Printf("failed to process config data: %v", err)
	}

	h := handler.NewHandler(time.Duration(con.OrderTTL) * time.Second)
	router := router.NewRouter(h)

	go func() {
		removeExpiredCron(con, h)
	}()

	r := gin.Default()
	orderPath := r.Group("/location")

	orderPath.POST("/:order_id/now", router.MakeOrder)
	orderPath.GET("/:order_id", router.GetOrder)
	orderPath.DELETE("/:order_id", router.DeleteOrder)

	if err := r.Run(fmt.Sprintf(":%s", con.ServerPort)); err != nil {
		panic(err)

	}
}

func removeExpiredCron(con config.Config, h *handler.Handler) {

	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(int(con.OrderTTL)).Second().SingletonMode().Do(h.RemoveExpiredOrder)
	if err != nil {
		fmt.Printf("error setting up runner removeExpiredOrder: %v", err)
	}

	scheduler.StartBlocking()
}
