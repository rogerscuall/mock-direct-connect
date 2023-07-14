package main

import (
	"log"
	"os"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

func NewCustomLogger(l string) *CustomLogger {
	clog := &CustomLogger{
		debug:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		info:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warning: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		error:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	clog.SetLogLevel(l)
	return clog
}

func (c *CustomLogger) SetLogLevel(l string) {
	switch l {
	case "DEBUG":
		c.logLevel = DEBUG
	case "INFO":
		c.logLevel = INFO
	case "WARNING":
		c.logLevel = WARNING
	case "ERROR":
		c.logLevel = ERROR
	default:
		c.logLevel = INFO
	}
}

func (c *CustomLogger) Debug(v ...interface{}) {
	if c.logLevel <= DEBUG {
		c.debug.Println(v...)
	}
}

func (c *CustomLogger) Info(v ...interface{}) {
	if c.logLevel <= INFO {
		c.info.Print(v...)
	}
}

func (c *CustomLogger) Warning(v ...interface{}) {
	if c.logLevel <= WARNING {
		c.warning.Println(v...)
	}
}

func (c *CustomLogger) Error(v ...interface{}) {
	if c.logLevel <= ERROR {
		c.error.Println(v...)
	}
}

type CustomLogger struct {
	logLevel LogLevel
	debug    *log.Logger
	info     *log.Logger
	warning  *log.Logger
	error    *log.Logger
}
