package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Subscription API
// @version 1.0
// @description API for managing subscriptions
// @host localhost:8080
// @BasePath /
var db *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	initDB()

	r := gin.Default()

	r.POST("/subscriptions", createSubscription)
	r.GET("/subscriptions", getSubscriptions)
	r.GET("/subscriptions/:id", getSubscription)
	r.PUT("/subscriptions/:id", updateSubscription)
	r.DELETE("/subscriptions/:id", deleteSubscription)
	r.GET("/subscriptions/total", getTotalCost)

	r.Run()
}

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = db.AutoMigrate(&Subscription{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type TotalCostResponse struct {
	TotalCost float64 `json:"total_cost"`
}

// @Summary Create a subscription
// @Description Create a new subscription
// @Accept json
// @Produce json
// @Param subscription body Subscription true "Subscription"
// @Success 201 {object} Subscription
// @Failure 400 {object} ErrorResponse
// @Router /subscriptions [post]
func createSubscription(c *gin.Context) {
	var subscription Subscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(400, ErrorResponse{Message: err.Error()})
		return
	}
	db.Create(&subscription)
	log.Printf("Subscription created: %+v\n", subscription)
	c.JSON(201, subscription)
}

// @Summary Get all subscriptions
// @Description Get all subscriptions
// @Produce json
// @Success 200 {array} Subscription
// @Router /subscriptions [get]
func getSubscriptions(c *gin.Context) {
	var subscriptions []Subscription
	db.Find(&subscriptions)
	c.JSON(200, subscriptions)
}

// @Summary Get a subscription by ID
// @Description Get a subscription by its ID
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} Subscription
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [get]
func getSubscription(c *gin.Context) {
	id := c.Param("id")
	var subscription Subscription
	if err := db.First(&subscription, "id = ?", id).Error; err != nil {
		log.Println("Subscription not found:", err)
		c.JSON(404, ErrorResponse{Message: "Subscription not found!"})
		return
	}
	c.JSON(200, subscription)
}

// @Summary Update a subscription by ID
// @Description Update an existing subscription by its ID
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body Subscription true "Updated Subscription"
// @Success 200 {object} Subscription
// @Failure 404 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /subscriptions/{id} [put]
func updateSubscription(c *gin.Context) {
	id := c.Param("id")
	var subscription Subscription

	if err := db.First(&subscription, "id = ?", id).Error; err != nil {
		log.Println("Subscription not found:", err)
		c.JSON(404, ErrorResponse{Message: "Subscription not found!"})
		return
	}

	if err := c.ShouldBindJSON(&subscription); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(400, ErrorResponse{Message: err.Error()})
		return
	}

	db.Save(&subscription)
	log.Printf("Subscription updated: %+v\n", subscription)
	c.JSON(200, subscription)
}

// @Summary Delete a subscription by ID
// @Description Delete a subscription by its ID
// @Param id path int true "Subscription ID"
// @Success 204 {object} nil
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [delete]
func deleteSubscription(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Subscription{}, id).Error; err != nil {
		log.Println("Subscription not found:", err)
		c.JSON(404, ErrorResponse{Message: "Subscription not found!"})
		return
	}
	log.Printf("Subscription deleted: %s\n", id)
	c.JSON(204, nil)
}

// @Summary Get total cost of subscriptions
// @Description Get the total cost of subscriptions for a user and service name
// @Produce json
// @Param user_id query string true "User ID"
// @Param service_name query string true "Service Name"
// @Success 200 {object} TotalCostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/total [get]
func getTotalCost(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	var total float64
	err := db.Model(&Subscription{}).
		Where("user_id = ? AND service_name = ?", userID, serviceName).
		Select("SUM(cost)").Scan(&total).Error

	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Internal server error"})
		return
	}

	c.JSON(200, TotalCostResponse{TotalCost: total})
}
