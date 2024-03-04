package main

import (
	"image-Designer/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/generate", handler.SubmitResultHandler)

	r.POST("/batch-generate", handler.BatchSubmitHandler)

	r.GET("/", func(c *gin.Context) {
		c.File("static/image_submission.html")
	})

	err := r.Run(":9999")
	if err != nil {
		panic("Server starts fail: " + err.Error())
	}
}
