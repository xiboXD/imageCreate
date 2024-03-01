package handler

import (
	"github.com/gin-gonic/gin"
	"image-Designer/internal/service"
	"net/http"
	"time"
)

func SubmitResultHandler(c *gin.Context) {
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

	// Wait for 30 seconds before querying the result
	time.Sleep(30 * time.Second)

	result, err := service.Result(id)
	attempts := 0
	for err != nil && attempts < 2 {
		// Wait for another 30 seconds if result is not available
		time.Sleep(30 * time.Second)
		result, err = service.Result(id)
		attempts++
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"data":    make([]string, 0),
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    result,
		"code":    200,
	})
}
