package main

import (
	"github.com/gin-gonic/gin"
	"image-Designer/internal/handler"
)

func main() {
	r := gin.Default()

	// 设置静态文件服务
	r.Static("/static", "./static")

	// 设置提交处理程序路由
	r.GET("/submit", handler.SubmitResultHandler)

	// 设置根路由，返回 HTML 页面
	r.GET("/", func(c *gin.Context) {
		c.File("static/image_submission.html")
	})

	// 启动服务器
	err := r.Run(":9999")
	if err != nil {
		panic("Server starts fail: " + err.Error())
	}
}
