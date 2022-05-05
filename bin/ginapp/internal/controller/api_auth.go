package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"ginapp/internal/dto"
	"ginapp/internal/model/repository"
	"ginapp/internal/constants"
	"ginapp/internal/auth/jwt"
	"ginapp/internal/pkg/logger"

)

/*
login.goのJSONレスポンスver
POST /api/login に対して jwtトークンをレスポンス
*/

type AuthController interface {
	Login(c *gin.Context)
	Signup(c *gin.Context)
	ChangePassword(c *gin.Context)
	ChangeUsername(c *gin.Context)
	GetProfile(c *gin.Context)
	DeleteAccount(c *gin.Context)
}


type authController struct {
	ur repository.UserRepository
}


func NewAuthController() AuthController {
	ur := repository.NewUserRepository()
	return authController{ur}
}


//POST /api/login
func (ac authController) Login(c *gin.Context) {
	ld := &dto.LoginDto{}
	c.BindJSON(&ld)

	user, err := ac.ur.SelectByUsername(ld.Username)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ld.Password)) != nil{
		c.JSON(401, gin.H{"error": http.StatusText(401)})
		c.Abort()
		return
	}

	jwtString, err := jwt.GenerateJWT(user.UId)
	if err != nil {
		logger.LogError(err.Error())
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}
	c.SetCookie(jwt.JwtKeyname, jwtString, int(jwt.JwtExpires), "/", constants.Hostname, false, true)
	c.JSON(200 , gin.H{jwt.JwtKeyname: jwtString})
}


//POST /api/signup
func (ac authController) Signup(c *gin.Context) {
	sd := &dto.SignupDto{} 
	c.BindJSON(&sd)

	if _, err := ac.ur.SelectByUsername(sd.Username); err == nil {
		c.JSON(409, gin.H{"error": http.StatusText(409)})
		c.Abort()
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(sd.Password), bcrypt.DefaultCost)
	sd.Password = string(hashed)

	if err := ac.ur.Signup(sd); err != nil {
		logger.LogError(err.Error())
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}

    c.JSON(201, gin.H{})
}


//PUT[POST] /api/password
func (ac authController) ChangePassword(c *gin.Context) {
	uid, _ := jwt.ExtractUId(c)

	var body map[string]interface{}
	c.BindJSON(&body)
	newPassword := body["password"].(string)
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil || ac.ur.UpdatePasswordByUId(uid, string(hashed)) != nil {
		logger.LogError(err.Error())
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}


//PUT[POST] /api/username
func (ac authController) ChangeUsername(c *gin.Context) {
	uid, _ := jwt.ExtractUId(c)
	var body map[string]interface{}
	c.BindJSON(&body)
	newUsername := body["username"].(string)

	if err := ac.ur.UpdateUsernameByUId(uid, newUsername); err != nil {
		logger.LogError(err.Error())
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}


//GET /api/profile
func (ac authController) GetProfile(c *gin.Context) {
	uid, _ := jwt.ExtractUId(c)
	user, err := ac.ur.SelectByUId(uid)

	if err != nil {
		logger.LogError(err.Error())
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}

	c.JSON(200, user)
}


//DELETE /api/account
func (ac authController) DeleteAccount(c *gin.Context) {
	uid, _ := jwt.ExtractUId(c)

	if err := ac.ur.DeleteByUId(uid); err != nil {
		logger.LogError(err.Error())
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}
