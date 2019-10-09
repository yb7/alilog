package alilog

import (
  "os"
  "log"
)

var stdTrace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime)
var stdDebug = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
var stdInfo = log.New(os.Stdout, "INFO:  ", log.Ldate|log.Ltime)
var stdWarning = log.New(os.Stdout, "WARN:  ", log.Ldate|log.Ltime)
var stdError = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
