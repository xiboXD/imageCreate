package handler

import (
	"image-Designer/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SubmitResultHandler(c *gin.Context) {
	requestMsg := c.Query("message")
	id, err := service.Submit(requestMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"id":      "",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// Wait for 15 seconds before querying the result
	time.Sleep(15 * time.Second)

	result, err := service.Result(id)
	attempts := 0
	for err != nil && attempts < 2 {
		// Wait for another 15 seconds if result is not available
		time.Sleep(15 * time.Second)
		result, err = service.Result(id)
		attempts++
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"data":    make([]string, 0),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    result,
		"code":    http.StatusOK,
	})
}

// func BatchSubmitHandler(c *gin.Context) {
// 	// Declare a struct to hold the request body
// 	var requestBody struct {
// 		Messages []string `json:"messages"`
// 	}

// 	// Bind the request body to the struct
// 	if err := c.BindJSON(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Invalid request body",
// 			"code":    http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	// Define slice to store aggregated data
// 	var aggregatedData []string

// 	// Iterate over each message
// 	for _, message := range requestBody.Messages {
// 		// Submit the message and get the ID
// 		id, err := service.Submit(message)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"message": err.Error(),
// 				"data":    aggregatedData,
// 				"code":    http.StatusInternalServerError,
// 			})
// 			return
// 		}

// 		// Wait for 30 seconds before querying the result
// 		time.Sleep(30 * time.Second)

// 		// Query the result
// 		result, err := service.Result(id)

// 		// Retry if result is not available
// 		attempts := 0
// 		for err != nil && attempts < 2 {
// 			// Wait for another 30 seconds if result is not available
// 			time.Sleep(30 * time.Second)
// 			result, err = service.Result(id)
// 			attempts++
// 		}

// 		// If error still exists after retrying
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"message": err.Error(),
// 				"data":    aggregatedData,
// 				"code":    http.StatusInternalServerError,
// 			})
// 			return
// 		}

// 		// Append result to aggregated data
// 		aggregatedData = append(aggregatedData, result...)
// 	}

// 	// Return success response with aggregated data
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Success",
// 		"data":    aggregatedData,
// 		"code":    http.StatusOK,
// 	})
// }

func BatchSubmitHandler(c *gin.Context) {
	// Declare a struct to hold the request body
	var requestBody struct {
		Messages []string `json:"messages"`
	}

	// Bind the request body to the struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// Define map to store aggregated data
	aggregatedData := make(map[string][]string)

	// Iterate over each message
	for _, message := range requestBody.Messages {
		// Submit the message and get the ID
		id, err := service.Submit(message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"data":    aggregatedData,
				"code":    http.StatusInternalServerError,
			})
			return
		}

		// Wait for 30 seconds before querying the result
		time.Sleep(30 * time.Second)

		// Query the result
		result, err := service.Result(id)

		// Retry if result is not available
		attempts := 0
		for err != nil && attempts < 2 {
			// Wait for another 30 seconds if result is not available
			time.Sleep(30 * time.Second)
			result, err = service.Result(id)
			attempts++
		}

		// If error still exists after retrying
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"data":    aggregatedData,
				"code":    http.StatusInternalServerError,
			})
			return
		}

		// Set the result in the map with the message as key
		aggregatedData[message] = result // Assuming there's only one result for each message
	}

	// Return success response with aggregated data
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    aggregatedData,
		"code":    http.StatusOK,
	})
}
