package main

import (
	"fmt"
	"net/http"

	"github.com/adilalimgozha/Go/tree/main/lab3/backend/models"
	"github.com/adilalimgozha/Go/tree/main/lab3/backend/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var users = make(map[string]string)
var tasks = []models.Task{
	{ID: "1", Title: "Task 1", Description: "First task", Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", Status: "Completed"},
}

func register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Store user (In-memory for simplicity)
	users[user.Username] = user.Password

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// User login
func login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check credentials
	storedPassword, exists := users[user.Username]
	if !exists || storedPassword != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, models.JWT{Token: token})
}

// Middleware to protect routes
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Print the token for debugging
		fmt.Println("Token received:", tokenStr)

		// Validate token
		token, err := utils.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Protected route to get all tasks
func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

// Protected route to create a new task
func createTask(c *gin.Context) {
	var newTask models.Task

	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tasks = append(tasks, newTask)

	c.JSON(http.StatusCreated, newTask)
}

// Protected route to update a task
func updateTask(c *gin.Context) {
	id := c.Param("id")

	var updatedTask models.Task

	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = updatedTask
			c.JSON(http.StatusOK, updatedTask)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

// Protected route to delete a task
func deleteTask(c *gin.Context) {
	id := c.Param("id")

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.POST("/register", register)
	r.POST("/login", login)

	// Protected routes group
	authGroup := r.Group("/auth")
	authGroup.Use(authMiddleware())
	{
		authGroup.GET("/tasks", getTasks)          // Get all tasks
		authGroup.POST("/tasks", createTask)       // Create a new task
		authGroup.PUT("/tasks/:id", updateTask)    // Update an existing task
		authGroup.DELETE("/tasks/:id", deleteTask) // Delete a task
	}

	r.Run(":8080")
}
