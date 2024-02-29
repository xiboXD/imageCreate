package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"image-Designer/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

var (
	mongoClient *mongo.Client
)

func init() {
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	mongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Ping the MongoDB to check the connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
}

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

	// Store data in MongoDB
	err = storeDataInMongoDB(result)
	if err != nil {
		log.Println("Error storing data in MongoDB:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    result,
		"code":    200,
	})
}

func storeDataInMongoDB(data []string) error {
	collection := mongoClient.Database("cats").Collection("generated_img")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, bson.M{"data": data})
	if err != nil {
		return err
	}
	return nil
}
