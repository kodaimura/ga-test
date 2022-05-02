package controller

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    
    "ginapp/internal/auth/jwt"
    "ginapp/internal/dto"
    "ginapp/internal/model/repository"
    "ginapp/internal/constants"
)


type LoginController interface {
    LoginPage(c *gin.Context)
    Login(c *gin.Context)
    Logout(c *gin.Context)
}


type loginController struct {
    ur repository.UserRepository
}


func NewLoginController() LoginController {
    ur := repository.NewUserRepository()
    return loginController{ur}
}


//GET /login
func (lc loginController) LoginPage(c *gin.Context) {
    c.HTML(200, "login.html", gin.H{
        "appname": constants.AppName,
    })
}


//POST /login
func (lc loginController) Login(c *gin.Context) {
    ld := &dto.LoginDto{}
    ld.UserName = c.PostForm("username")
    ld.Password = c.PostForm("password")

    user, err := lc.ur.SelectByUserName(ld.UserName)

    if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ld.Password)) != nil{
        c.HTML(401, "login.html", gin.H{
            "appname":constants.AppName,
            "error": "UsernameまたはPasswordが異なります。",
        })
        c.Abort()
        return
    }

    jwtString, err := jwt.GenerateJWT(user.UId)
    if err != nil {
        c.HTML(500, "login.html", gin.H{
            "appname": constants.AppName,
            "error": "ログインに失敗しました。",
        })
        c.Abort()
        return
    }

    c.SetCookie(jwt.JwtKeyName, jwtString, int(jwt.JwtExpires), "/", constants.HostName, false, true)
    c.Redirect(303, "/")
}


//GET /logout
func (lc loginController) Logout(c *gin.Context) {
    c.SetCookie(jwt.JwtKeyName, "", 0, "/", constants.HostName, false, true)
    c.Redirect(303, "/login")
}
