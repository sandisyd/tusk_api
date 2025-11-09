package main

import (
	"net/http"
	"tusk/config"
	"tusk/controller"
	"tusk/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// db connection
	db := config.DatabaseConnection()
	db.AutoMigrate(&models.User{}, &models.Task{})
	config.CreateOwnerAccount(db)

	// controller
	userCt := controller.UserController{DB: db}
	taskCt := controller.TaskController{DB: db}
	// Router

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		
		c.JSON(http.StatusOK, "Welcome to my API")
	})
	// route account
	router.POST("/users/login", userCt.Login)
	router.POST("/users/register", userCt.Register)
	router.DELETE("/users/:id", userCt.Delete)
	router.GET("/users/employee", userCt.GetListUsersEmployee)
	
	// route task
	router.POST("/task", taskCt.CreateTask)
	router.DELETE("/task/:id", taskCt.Delete)
	router.PATCH("/task/:id/patch-sub", taskCt.PatchData)
	router.PATCH("/task/:id/rejected", taskCt.RejectedTask)


	router.Static("/attachment", "./attachment")
	router.Run("localhost:8080")
}