package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/AlexSamarskii/delivery_service/internal/config"
	"github.com/AlexSamarskii/delivery_service/internal/middleware"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func setupDefaultOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {

	setupDefaultOutput()
	cfg := config.Load()

	server := gin.New()
	server.Use(middleware.Recover(), middleware.Logger())

	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK,
			gin.H{"message": "OK"},
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

	server.Run(fmt.Sprintf(":%s", cfg.Port))
}
