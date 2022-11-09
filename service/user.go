package service
 
import (
	"crypto/sha256"
    "net/http"
    "strconv"
    "regexp"
 
    "github.com/gin-gonic/gin"
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
    ctx.JSON(http.StatusOK, user)
}

func check_regexp(reg, str string) bool{
    return regexp.MustCompile(reg).Match([]byte(str))
}