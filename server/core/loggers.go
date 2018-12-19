package core

import (
	"log"
	"os"
)

var (
	Error   *log.Logger
	Info    *log.Logger
	Request *log.Logger
)

func InitLoggers() {
	Error = log.New(os.Stderr, "[ERR]", log.LstdFlags|log.Llongfile)
	Info = log.New(os.Stdout, "[INFO]", log.LstdFlags|log.Lshortfile)
	Request = log.New(os.Stdout, "[REQ]", log.LstdFlags)
}
