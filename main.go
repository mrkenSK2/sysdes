package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"

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
    
	// prepare session
    store := cookie.NewStore([]byte("my-secret"))
    engine.Use(sessions.Sessions("user-session", store))
	
	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Home)
	engine.GET("/list", service.LoginCheck, service.TaskList)

	engine.GET("/task/sort", service.LoginCheck, service.Sort)

	taskGroup := engine.Group("/task")
    taskGroup.Use(service.LoginCheck)
    {
        //taskGroup.GET("/:id", service.ShowTask)
        taskGroup.GET("/new", service.NewTaskForm)
        taskGroup.POST("/new", service.RegisterTask)
        //taskGroup.GET("/edit/:id", service.EditTaskForm)
        taskGroup.POST("/edit/:id", service.UpdateTask)
        //taskGroup.GET("/delete/:id", service.DeleteTask)
    }

	taskGroupGuard := engine.Group("/task")
    taskGroupGuard.Use(service.LoginCheck, service.UserCheck)
    {
        taskGroupGuard.GET("/:id", service.ShowTask)
        taskGroupGuard.GET("/edit/:id", service.EditTaskForm)
        taskGroupGuard.GET("/delete/:id", service.DeleteTask)
    }
	/*engine.GET("/task/:id", service.ShowTask) // ":id" is a parameter

	// add new task
	engine.GET("/task/new", service.NewTaskForm)
	engine.POST("/task/new", service.RegisterTask)
	
	// edit existing task
	engine.GET("/task/edit/:id", service.EditTaskForm)
    engine.POST("/task/edit/:id", service.UpdateTask)
	
	// delete existing task
    engine.GET("/task/delete/:id", service.DeleteTask)
*/
	// register user
	engine.GET("/user/new", service.NewUserForm)
    engine.POST("/user/new", service.RegisterUser)

	// user login
	engine.GET("/user/login", service.LoginForm)
    engine.POST("/user/login", service.Login)

	engine.GET("/user/edit", service.EditUserForm)
    engine.POST("/user/edit", service.UpdateUser)

	// user logout
	engine.GET("/user/logout", service.Logout)

	// delete user
	engine.GET("/user/delete", service.DeleteUser)

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
