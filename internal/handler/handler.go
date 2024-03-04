package handler

import (
	"github.com/gin-gonic/gin"
	"image-Designer/internal/service"
	"net/http"
)

func GenerateHandler(c *gin.Context) {
	requestMsg := c.Query("message")
	id, err := service.Submit(requestMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"id":      "",
			"code":    "500",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Request submitted",
		"id":      id,
		"code":    200,
	})
}

func ResultHandler(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "ID cannot be empty, operation failed",
			"data":    make([]string, 0),
			"code":    500,
		})
		return
	}
	result, err := service.Result(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"data":    make([]string, 0),
			"code":    500,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Operation successful",
		"data":    result,
		"code":    200,
	})
}
