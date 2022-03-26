package controller

import (
    "log"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    
    "ginapp/internal/dto"
    "ginapp/internal/constants"
    "ginapp/internal/model/repository"
)


type SignupController interface {
    SignupPage(c *gin.Context)
    Signup(c *gin.Context)
}


type signupController struct {
    ur repository.UserRepository
}


func NewSignupController() SignupController {
    ur, err := repository.NewUserRepository()
    if err != nil {
        log.Panic(err)
    }
    return signupController{ur}
}


func (sc signupController) SignupPage(c *gin.Context) {
    c.HTML(200, "signup.html", gin.H{
        "appname": constants.AppName,
    })
}


func (sc signupController) Signup(c *gin.Context) {
    sd := &dto.SignupDto{} 
    sd.UserName = c.PostForm("username")
    sd.Password = c.PostForm("password")

    if _, err := sc.ur.SelectByUserName(sd.UserName); err == nil {
        c.HTML(409, "signup.html", gin.H{
            "appname": constants.AppName,
            "error": "Usernameが既に使われています。",
        })
        c.Abort()
        return
    }

    hashed, _ := bcrypt.GenerateFromPassword([]byte(sd.Password), bcrypt.DefaultCost)
    sd.Password = string(hashed)

    if sc.ur.Signup(*sd) != nil {
        c.JSON(500, gin.H{
            "appname": constants.AppName,
            "error": "登録に失敗しました。",
        })
        c.Abort()
        return
    }

    c.Redirect(303, "/login")
}
