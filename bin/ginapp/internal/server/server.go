package server

import (
    "log"
    "io"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"

    "ginapp/internal/controller"
    "ginapp/pkg/jwtauth"
)


func Run() {
    setLogger()
    r := router()
    r.Run()
}


func setLogger () {
    logfolder := "log"
    logfile := "log/app.log"

    if _, err := os.Stat(logfolder); err != nil {
        os.Mkdir(logfolder, 0777)
    }

    if _, err := os.Stat(logfile); err == nil {
        t := time.Now()
        format := "2006-01-02-15-04-05"
        fname := "log/~" + t.Format(format) + ".log"
        if err := os.Rename(logfile, fname); err != nil {
            log.Panic(err)
        }
    }

    f, err := os.Create(logfile); 

    if err != nil {
        log.Panic(err)
    }

    gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
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
    auth.Use(jwtauth.JwtAuthMiddleware())
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
        apiAuth.Use(jwtauth.JwtApiAuthMiddleware())
        {
            apiAuth.GET("/profile", ac.GetProfile)
            apiAuth.PUT("/username", ac.ChangeUserName)
            apiAuth.POST("/username", ac.ChangeUserName)
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
