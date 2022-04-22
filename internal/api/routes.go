package api

import "github.com/sword-challenge/internal/api/middlewares"

func (r Router) RegisterRoutes() {
	r.registerUserRoutes()

}

func (r Router) registerUserRoutes() {
	userRoute := r.base.Group("/user")

	userRoute.POST("/login", r.userS.Login)
}

func (r Router) registerTaskRoutes() {
	taskRoute := r.base.Group("/task", middlewares.AuthenticatedUser(r.userS))

	taskRoute.GET("/tasks", r.taskS.GetTasks)
	taskRoute.PUT("/tasks/:id", r.taskS.UpdateTask, middlewares.MustReceiveID)
	taskRoute.DELETE("/tasks/:id", r.taskS.DeleteTask, middlewares.MustReceiveID)
	taskRoute.POST("/tasks", r.taskS.CreateTask, middlewares.MustReceiveID)
}
