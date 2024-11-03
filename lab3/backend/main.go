package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	var tasks = []Task{
		{ID: "1", Title: "Task 1", Description: "First task", Status: "Pending"},
		{ID: "2", Title: "Task 2", Description: "Second task", Status: "In Progress"},
		{ID: "3", Title: "Task 3", Description: "Third task", Status: "Completed"},
	}

	r := gin.Default()

}
