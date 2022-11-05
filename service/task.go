package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

    // Get query parameter
    kw := ctx.Query("kw")
    kw_check := ctx.Query("kw_check")
    done := ctx.Query("done")
    notdone := ctx.Query("notdone")
    //kw_check := ctx.("kw_check")
    //if checkedtxt == "on" {}else{}
    

	// Get tasks in DB
	var tasks []database.Task

    if done=="on" && kw_check!="on"{
        db.Select(&tasks, "SELECT * FROM tasks WHERE is_done LIKE 1")
    }else if notdone=="on" && kw_check!="on"{
        db.Select(&tasks, "SELECT * FROM tasks WHERE is_done LIKE 0")
    }else {
        switch{
            case kw != "":
                err = db.Select(&tasks, "SELECT * FROM tasks WHERE title LIKE ?", "%" + kw + "%")
	        default:
                err = db.Select(&tasks, "SELECT * FROM tasks") // Use DB#Select for multiple entries
	    }
    }
    if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks, "Kw": kw})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Render task
	//ctx.String(http.StatusOK, task.Title)  // Modify it!!
	ctx.HTML(http.StatusOK, "task.html", task)
}

// registration
func NewTaskForm(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "form_new_task.html", gin.H{"Title": "Task registration"})
}

func RegisterTask(ctx *gin.Context) {
    // Get task title
    title, exist := ctx.GetPostForm("title")
    if !exist {
        Error(http.StatusBadRequest, "No title is given")(ctx)
        return
    }
	detail, exist := ctx.GetPostForm("detail")
    if !exist {
        Error(http.StatusBadRequest, "Problem in detail")(ctx)
        return
    }
    // Get DB connection
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
	// Create new data with given title on DB
    result, err := db.Exec("INSERT INTO tasks (title, detail) VALUES (?, ?)", title, detail)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
	 // Render status
	 path := "/list"  // デフォルトではタスク一覧ページへ戻る
	 if id, err := result.LastInsertId(); err == nil {
		 path = fmt.Sprintf("/task/%d", id)   // 正常にIDを取得できた場合は /task/<id> へ戻る
	 }
	 ctx.Redirect(http.StatusFound, path)
}

func EditTaskForm(ctx *gin.Context) {
    // ID の取得
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Get DB connection
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    // Get target task
    var task database.Task
    err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Render edit form
    ctx.HTML(http.StatusOK, "form_edit_task.html",
        gin.H{"Title": fmt.Sprintf("Edit task %d", task.ID), "Task": task})
}

func UpdateTask(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Get task title
    title, exist := ctx.GetPostForm("title")
    if !exist {
        Error(http.StatusBadRequest, "No title is given")(ctx)
        return
    }
	is_done_bool, exist := ctx.GetPostForm("is_done")
    if !exist {
        Error(http.StatusBadRequest, "Problem in is_done")(ctx)
        return
    }
	is_done, err := strconv.ParseBool(is_done_bool)
	if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
	detail, exist := ctx.GetPostForm("detail")
    if !exist {
        Error(http.StatusBadRequest, "problem in detail")(ctx)
        return
    }
    // Get DB connection
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
	// update task
    _, err = db.Exec("UPDATE tasks SET title=?,is_done=?,detail=? WHERE id=?", title, is_done, detail, id)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
	ctx.Redirect(http.StatusFound, "/list")
}

func DeleteTask(ctx *gin.Context) {
    //var task database.Task
    // ID の取得
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Get DB connection
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    // Delete the task from DB
    _,err = db.Exec("DELETE FROM tasks WHERE id=?", id)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    // Redirect to /list
    ctx.Redirect(http.StatusFound, "/list")
}
