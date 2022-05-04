package controller

import (
    "github.com/gin-gonic/gin"
    
    "ginapp/internal/auth/jwt"
    "ginapp/internal/model/repository"
    "ginapp/internal/constants"
)


type IndexController interface {
    IndexPage(c *gin.Context)
}


type indexController struct {
    ur repository.UserRepository
}


func NewIndexController() IndexController {
    ur := repository.NewUserRepository()
    return indexController{ur}
}


//GET /
func (ic indexController) IndexPage(c *gin.Context) {
    username, err := jwt.ExtractUsername(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    c.HTML(200, "index.html", gin.H{
        "appname": constants.Appname,
        "username": username,
    })
}

