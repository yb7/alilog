package alilog

import (
  "os"
  "log"
)

var stdDebug = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
var stdInfo = log.New(os.Stdout, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)
var stdWarning = log.New(os.Stdout, "WARN:  ", log.Ldate|log.Ltime|log.Lshortfile)
var stdError = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
