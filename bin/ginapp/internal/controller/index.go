package controller

import (
    "github.com/gin-gonic/gin"
    
    "ginapp/pkg/jwtauth"
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


func (ic indexController) IndexPage(c *gin.Context) {
    username, err := jwtauth.ExtractUserName(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    c.HTML(200, "index.html", gin.H{
        "appname": constants.AppName,
        "username": username,
    })
}

