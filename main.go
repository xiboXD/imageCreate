package main

import (
	"github.com/gin-gonic/gin"
	"image-Designer/internal/handler"
)

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/submit", handler.SubmitResultHandler)

	r.GET("/", func(c *gin.Context) {
		c.File("static/image_submission.html")
	})

	err := r.Run(":9999")
	if err != nil {
		panic("Server starts fail: " + err.Error())
	}
}
