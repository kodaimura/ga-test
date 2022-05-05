package logger

import (
    "log"
    "io"
    "os"
    "time"
    "runtime"

    "github.com/gin-gonic/gin"
)


var logF *log.Logger
var logE *log.Logger
var logW *log.Logger
var logI *log.Logger
var logD *log.Logger

var logfolder = "log/"
var logfile = "app.log"
var logpath = logfolder + logfile
var file *os.File


func init() {
    if _, err := os.Stat(logfolder); err != nil {
        os.Mkdir(logfolder, 0777)
    }

    if _, err := os.Stat(logpath); err == nil {
        t := time.Now()
        format := "2006-01-02-15-04-05"
        fname := logfolder + "app~" + t.Format(format) + ".log"
        if err := os.Rename(logpath, fname); err != nil {
            log.Panic(err)
        }
    }

    file, err := os.Create(logpath); 

    if err != nil {
        log.Panic(err)
    }


    logF = log.New(
        io.MultiWriter(os.Stdout, file),
        "[FATAL]",
        log.LstdFlags,
    )

    logE = log.New(
        io.MultiWriter(os.Stdout, file),
        "[ERROR]",
        log.LstdFlags,
    )

    logW = log.New(
        io.MultiWriter(os.Stdout, file),
        "[WARN]",
        log.LstdFlags,
    )

    logI = log.New(
        io.MultiWriter(os.Stdout, file),
        "[INFO]",
        log.LstdFlags,
    )

    logD = log.New(
        os.Stdout,
        "[DEBUG]",
        log.LstdFlags,
    )
}


func SetAccessLogger () {
    gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
}


func LogFatal(msg string) {
    logF.Fatal("Msg:", msg)
}


func LogError(msg string) {
    pc, file, line, _ := runtime.Caller(1)
    f := runtime.FuncForPC(pc)
    logE.Println("\n", "File:", file, "Line:", line, "\n",
     "Func:", f.Name(), "Msg:", msg,
    )
}


func LogWarn(msg string) {
    pc, _, _, _ := runtime.Caller(1)
    f := runtime.FuncForPC(pc)
    logW.Println(f.Name(), msg)
}


func LogInfo(msg string) {
    pc, _, _, _ := runtime.Caller(1)
    f := runtime.FuncForPC(pc)
    logI.Println(f.Name(), msg)
}


func LogDebug(msg string) {
    pc, _, _, _ := runtime.Caller(1)
    f := runtime.FuncForPC(pc)
    logD.Println(f.Name(), msg)
}


