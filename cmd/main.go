package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/AlexSamarskii/delivery_service/internal/config"
	"github.com/AlexSamarskii/delivery_service/internal/middleware"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.Load()

	rout := gin.Default()
	rout.Use(middleware.Recover())

	rout.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK,
			gin.H{"message": "Ok"},
		)
	})

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatalln("Failed to listen:", err)
		}

		s := grpc.NewServer()

		log.Fatal(s.Serve(lis))
	}()

	rout.Run(fmt.Sprintf(":%s", cfg.Port))
}
