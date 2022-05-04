package controller

import (
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
    ur := repository.NewUserRepository()
    return signupController{ur}
}


//GET /signup
func (sc signupController) SignupPage(c *gin.Context) {
    c.HTML(200, "signup.html", gin.H{
        "appname": constants.Appname,
    })
}


//POST /signup
func (sc signupController) Signup(c *gin.Context) {
    sd := &dto.SignupDto{} 
    sd.Username = c.PostForm("username")
    sd.Password = c.PostForm("password")

    if _, err := sc.ur.SelectByUsername(sd.Username); err == nil {
        c.HTML(409, "signup.html", gin.H{
            "appname": constants.Appname,
            "error": "Usernameが既に使われています。",
        })
        c.Abort()
        return
    }

    hashed, _ := bcrypt.GenerateFromPassword([]byte(sd.Password), bcrypt.DefaultCost)
    sd.Password = string(hashed)

    if sc.ur.Signup(sd) != nil {
        c.HTML(500, "signup.html", gin.H{
            "appname": constants.Appname,
            "error": "登録に失敗しました。",
        })
        c.Abort()
        return
    }

    c.Redirect(303, "/login")
}
