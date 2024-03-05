package handler

import (
	"encoding/json"
	"fmt"
	"image-Designer/internal/service"
	"net/http"
	"os"
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

	timestamp := time.Now().Format("20060102_150405")

	// Return success response with timestamp
	c.JSON(http.StatusOK, gin.H{
		"message": "Processing request. Results will be available shortly.",
		"data":    timestamp,
		"code":    http.StatusOK,
	})

	// Process messages in the background
	go processMessages(requestBody.Messages, timestamp)
}

func processMessages(messages []string, timestamp string) {
	// Define map to store aggregated data
	aggregatedData := make(map[string][]string)

	// Iterate over each message
	for _, message := range messages {
		// Submit the message and get the ID
		id, err := service.Submit(message)
		if err != nil {
			fmt.Printf("Error submitting message %s: %v\n", message, err)
			continue
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
			fmt.Printf("Error getting result for message %s: %v\n", message, err)
			continue
		}

		// Set the result in the map with the message as key
		aggregatedData[message] = result // Assuming there's only one result for each message
	}

	// Save aggregated data to a JSON file
	saveDataToFile(aggregatedData, timestamp)
}

func saveDataToFile(data map[string][]string, timestamp string) {

	filename := fmt.Sprintf("output/output_%s.json", timestamp) // Output folder path

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	fmt.Printf("Data saved to file: %s\n", filename)
}

func GetOutputByIDHandler(c *gin.Context) {
	id := c.Query("id")
	filename := fmt.Sprintf("output/output_%s.json", id)

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// If the file does not exist, return a message indicating that it is being generated
		c.JSON(http.StatusOK, gin.H{
			"message": "File is being generated. Please wait.",
			"code":    http.StatusAccepted,
		})
		return
	}

	// If the file exists, read its contents
	file, err := os.Open(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error reading file.",
			"code":    http.StatusInternalServerError,
		})
		return
	}
	defer file.Close()

	var data map[string][]string
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error decoding JSON.",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// Return the contents of the file
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    data,
		"code":    http.StatusOK,
	})
}

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
