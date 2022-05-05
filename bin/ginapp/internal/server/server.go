package server

import (
    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"

    "ginapp/internal/constants"
    "ginapp/internal/controller"
    "ginapp/internal/auth/jwt"
    "ginapp/internal/pkg/logger"
)


func Run() {
    logger.SetAccessLogger()
    r := router()
    r.Run(constants.Port)
}


func router() *gin.Engine {
    r := gin.Default()
    
    //TEMPLATE
    r.LoadHTMLGlob("web/template/*.html")

    //STATIC
    r.Static("/css", "web/static/css")
    r.Static("/js", "web/static/js")

    
    //ROOT
    lc := controller.NewLoginController()
    sc := controller.NewSignupController()
    r.GET("/login", lc.LoginPage)
    r.POST("/login", lc.Login)
    r.GET("/logout", lc.Logout)
    r.GET("/signup", sc.SignupPage)
    r.POST("/signup", sc.Signup)


    ic := controller.NewIndexController()
    //ROOT (Authentication required)
    auth := r.Group("/")
    auth.Use(jwt.JwtAuthMiddleware())
    {
        auth.GET("/", ic.IndexPage)
    }


    ac := controller.NewAuthController()
    //API
    api := r.Group("/api")
    {
        api.POST("/login", ac.Login)
        api.POST("/signup", ac.Signup)
    
        //API (Authentication required)
        apiAuth := api.Group("/")
        apiAuth.Use(jwt.JwtApiAuthMiddleware())
        {
            apiAuth.GET("/profile", ac.GetProfile)
            apiAuth.PUT("/username", ac.ChangeUsername)
            apiAuth.POST("/username", ac.ChangeUsername)
            apiAuth.PUT("/password", ac.ChangePassword)
            apiAuth.POST("/password", ac.ChangePassword)
            apiAuth.DELETE("/account", ac.DeleteAccount)
        }
    }

/*  
    //ADMIN
    admin := r.Group("/admin")
    uc := controller.NewUserController(ur)
    {
        admin.GET("/users", uc.GetUsers)
        admin.GET("/users/:uid", uc.GetUser)
        admin.DELETE("/users/:uid", uc.DeleteUser)
    } 
 */
    
    return r
}
