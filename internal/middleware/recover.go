package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.JSON(http.StatusInternalServerError, fmt.Errorf("%v", err))
				log.Printf("Error: %+v", err)
			}
		}()
		ctx.Next()
	}
}
