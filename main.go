package main

import (
	"github.com/gin-gonic/gin"
	"image-Designer/internal/handler"
	"log"
)

func main() {

	r := gin.Default()
	r.GET("/generate", handler.GenerateHandler)
	r.GET("/result/:id", handler.ResultHandler)
	err := r.Run(":9000")
	if err != nil {
		log.Fatal("Fail:", err)
	}
}