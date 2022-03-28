package controller

import (
	"strconv"
	"net/http"

    "ginapp/internal/model/repository"

    "github.com/gin-gonic/gin"
)


type UserController interface {
	GetUsers(c *gin.Context)
	GetUserByUId(c *gin.Context)
	DeleteUserByUId(c *gin.Context)
}


type userController struct {
	ur repository.UserRepository
}


func NewUserController() UserController {
	ur := repository.NewUserRepository()
	return userController{ur}
}


//GET /admin/users
func (uc userController) GetUsers(c *gin.Context) {
	users, err := uc.ur.Select()

	if err != nil {
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}
    c.JSON(200, users)
}


//GET /admin/users/:uid
func (uc userController) GetUserByUId(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	user, err := uc.ur.SelectByUId(uid)

	if err != nil {
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}
	c.JSON(200, user)

}


//DELETE /admin/users/:uid
func (uc userController) DeleteUserByUId(c *gin.Context) {
	uid,_ := strconv.Atoi(c.Param("uid"))
	if uc.ur.DeleteByUId(uid) != nil {
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}
