package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User represents a user in the database
type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
	Age  int    `json:"age" gorm:"not null"`
}

// Database variables
var db *sql.DB
var gormDB *gorm.DB

// Initialize the databases
func initDB() {
	var err error
	// Direct SQL connection
	connStr := "user=postgres password=Adilek2003alimgozha dbname=golab3 host=localhost sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// GORM connection
	gormDB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate for GORM
	gormDB.AutoMigrate(&User{})
}

// Get Users with optional filtering and sorting
func getUsersDirect(c *gin.Context) {
	var users []User
	ageFilter := c.Query("age")
	order := c.Query("order")

	var query string
	var args []interface{}

	if ageFilter != "" {
		query = "SELECT id, name, age FROM users WHERE age = $1"
		args = append(args, ageFilter)
	} else {
		query = "SELECT id, name, age FROM users"
	}

	// Sorting
	if order == "asc" {
		query += " ORDER BY name ASC"
	} else if order == "desc" {
		query += " ORDER BY name DESC"
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}

// Create User
func createUserDirect(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO users(name, age) VALUES($1, $2)", newUser.Name, newUser.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

// Update User
func updateUserDirect(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", updatedUser.Name, updatedUser.Age, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedUser.ID = uint(id)
	c.IndentedJSON(http.StatusOK, updatedUser)
}

// Delete User
func deleteUserDirect(c *gin.Context) {
	idStr := c.Param("id")
	_, err := db.Exec("DELETE FROM users WHERE id = $1", idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Get Users with GORM
func getUsersGORM(c *gin.Context) {
	var users []User
	if err := gormDB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// Create User with GORM
func createUserGORM(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := gormDB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

// Update User with GORM
func updateUserGORM(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := gormDB.Model(&User{}).Where("id = ?", id).Updates(updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedUser.ID = uint(id)
	c.IndentedJSON(http.StatusOK, updatedUser)
}

// Delete User with GORM
func deleteUserGORM(c *gin.Context) {
	idStr := c.Param("id")
	var user User
	if err := gormDB.First(&user, idStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := gormDB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Get Users with GORM and Pagination
func getUsersGORM_p(c *gin.Context) {
	var users []User
	var count int64
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default to 10 items per page
	}

	offset := (page - 1) * pageSize

	if err := gormDB.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gormDB.Model(&User{}).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"page":     page,
		"pageSize": pageSize,
		"total":    count,
		"users":    users,
	})
}

func main() {
	initDB()
	router := gin.Default()

	// Direct SQL routes
	/*router.GET("/users/direct", getUsersDirect)
	router.POST("/users/direct", createUserDirect)
	router.PUT("/users/direct/:id", updateUserDirect)
	router.DELETE("/users/direct/:id", deleteUserDirect)*/

	// GORM routes
	router.GET("/users/gorm", getUsersGORM)
	router.POST("/users/gorm", createUserGORM)
	router.PUT("/users/gorm/:id", updateUserGORM)
	router.DELETE("/users/gorm/:id", deleteUserGORM)

	router.GET("/users/gorm", getUsersGORM_p)

	router.Run(":8080") // Start the server on port 8080
}
