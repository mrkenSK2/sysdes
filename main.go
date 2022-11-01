package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Home)
	engine.GET("/list", service.TaskList)
	engine.GET("/task/:id", service.ShowTask) // ":id" is a parameter

	// adding new task
	engine.GET("/task/new", service.NotImplemented)
	engine.POST("/task/new", service.NotImplemented)
	// existing task edit
	engine.GET("/task/edit/:id", service.NotImplemented)
    engine.POST("/task/edit/:id", service.NotImplemented)
	// existing task delete
    engine.GET("/task/delete/:id", service.NotImplemented)

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
