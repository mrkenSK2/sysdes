package service
 
import (
	"crypto/sha256"
    "encoding/hex"
    "net/http"
    "strconv"
    "regexp"
 
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
	database "todolist.go/db"
)

const MIN_PW_LEN = 8
 
func NewUserForm(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "new_user_form.html", gin.H{"Title": "Register user"})
}

func hash(pw string) []byte {
    const salt = "todolist.go#"
    h := sha256.New()
    h.Write([]byte(salt))
    h.Write([]byte(pw))
    return h.Sum(nil)
}

func RegisterUser(ctx *gin.Context) {
    // フォームデータの受け取り
    username := ctx.PostForm("username")
    password := ctx.PostForm("password")
    re_enter_password := ctx.PostForm("re_enter_password")
    switch {
        case username == "":
            ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Usernane is not provided", "Username": username})
        case password == "":
            ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password is not provided", "Password": password})
        case re_enter_password == "":
            ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Confirm password is not provided", "Re_enter_Password": re_enter_password})
    }
    
    // DB 接続
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

    // 重複チェック
    var duplicate int
    err = db.Get(&duplicate, "SELECT COUNT(*) FROM users WHERE name=?", username)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    if duplicate > 0 {
        ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Username is already taken", "Username": username, "Password": password, "Re_enter_Password": re_enter_password})
        return
    }

    // check whether confirm password matches
    if password != re_enter_password{
        ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password are not matching", "Username": username, "Password": password, "Re_enter_Password": re_enter_password})
        return
    }

    errStmt := ""
    badPwFlag := 0
    // check password
    if len(password) < MIN_PW_LEN{
        badPwFlag = 1
        errStmt = "password should be minimum " +strconv.Itoa(MIN_PW_LEN) + " characters. "
        //ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "password should be minimum " +strconv.Itoa(MIN_PW_LEN) + " characters", "Username": username, "Password": password, "Re_enter_Password": re_enter_password})
        //return
    }

    if !(check_regexp(`[a-z]`, password) && check_regexp(`[A-Z]`, password) && check_regexp(`[0-9]`, password)){
        badPwFlag = 1
        errStmt = errStmt + "password must contain at least one lowercase letter, one uppercase letter, and one number"
    }
    if badPwFlag == 1{
        ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": errStmt, "Username": username, "Password": password, "Re_enter_Password": re_enter_password})
        return
    }

 
    // DB への保存
    result, err := db.Exec("INSERT INTO users(name, password) VALUES (?, ?)", username, hash(password))
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
 
    // 保存状態の確認
    id, _ := result.LastInsertId()
    var user database.User
    err = db.Get(&user, "SELECT id, name, password FROM users WHERE id = ?", id)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "HOME"})
}

func check_regexp(reg, str string) bool{
    return regexp.MustCompile(reg).Match([]byte(str))
}

func LoginForm(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Register user"})
}

const userkey = "user"
 
func Login(ctx *gin.Context) {
    username := ctx.PostForm("username")
    password := ctx.PostForm("password")
 
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
 
    // ユーザの取得
    var user database.User
    err = db.Get(&user, "SELECT id, name, password FROM users WHERE name = ?", username)
    if err != nil {
        ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"Title": "Login", "Username": username, "Error": "No such user"})
        return
    }
 
    // パスワードの照合
    if hex.EncodeToString(user.Password) != hex.EncodeToString(hash(password)) {
        ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"Title": "Login", "Username": username, "Error": "Incorrect password"})
        return
    }
 
    // セッションの保存
    session := sessions.Default(ctx)
    session.Set(userkey, user.ID)
    session.Save()
 
    ctx.Redirect(http.StatusFound, "/list")
}

func LoginCheck(ctx *gin.Context) {
    if sessions.Default(ctx).Get(userkey) == nil {
        ctx.Redirect(http.StatusFound, "/user/login")
        ctx.Abort()
    } else {
        ctx.Next()
    }
}

func Logout(ctx *gin.Context) {
    session := sessions.Default(ctx)
    session.Clear()
    session.Options(sessions.Options{Path: "/", MaxAge: -1})
    session.Save()
    ctx.Redirect(http.StatusFound, "/")
}

func UserCheck(ctx *gin.Context) {
    user_id := sessions.Default(ctx).Get(userkey)
    id, _ := strconv.Atoi(ctx.Param("id"))
    db, err := database.GetConnection()
    var pair database.Ownership
    err = db.Get(&pair, "SELECT user_id, task_id FROM ownership WHERE user_id = ? AND task_id = ?", user_id, id)
	if err != nil{
		//Error(http.StatusForbidden, err.Error())(ctx)
        //ctx.Redirect(http.StatusFound, "/")//仮おき
        ctx.HTML(http.StatusOK, "no_permission.html", gin.H{"Title": "No permission"})
		ctx.Abort()
	} else {
        ctx.Next()
    }
}

func LogoutAndDelete(ctx *gin.Context) {
    session := sessions.Default(ctx)
    session.Clear()
    session.Options(sessions.Options{Path: "/", MaxAge: -1})
    session.Save()
    ctx.Redirect(http.StatusFound, "/user/delete")
}
func DeleteUser(ctx *gin.Context) {
    //var task database.Task
    user_id := sessions.Default(ctx).Get(userkey)
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    
    db.Exec("DELETE FROM tasks WHERE id IN (SELECT task_id FROM ownership WHERE user_id = ?", user_id)
    db.Exec("DELETE FROM users WHERE id = ?", user_id)
    db.Exec("DELETE FROM ownership WHERE user_id = ?", user_id)
    
    /*tx := db.MustBegin()
    result, err := tx.Exec("INSERT INTO tasks (title) VALUES (?)", title)
    if err != nil {
        tx.Rollback()
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    taskID, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    _, err = tx.Exec("INSERT INTO ownership (user_id, task_id) VALUES (?, ?)", userID, taskID)
    if err != nil {
        tx.Rollback()
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    tx.Commit()*/
    
    /*err = db.Get(&pair, "SELECT user_id, task_id FROM ownership WHERE user_id = ? AND task_id = ?", user_id, id)
	if err != nil{
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
    }*/
    session := sessions.Default(ctx)
    session.Clear()
    session.Options(sessions.Options{Path: "/", MaxAge: -1})
    session.Save()
    ctx.Redirect(http.StatusFound, "/")
}