package routes

func SetUpRoutes() {
	r.GET("/tasks", handler.GetTasks)
}
