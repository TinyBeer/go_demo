package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	r := gin.Default()

	limiter := rate.NewLimiter(200, 200)
	r.GET("/rate", func(ctx *gin.Context) {
		ok := limiter.Allow()
		if ok {
			ctx.JSON(http.StatusOK, nil)
		} else {
			ctx.JSON(http.StatusTooManyRequests, nil)
		}
	})

	r.Run(":8080")
}
