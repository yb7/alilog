package alilog

import "fmt"

var enableDebug = false
func _debug(format string, v... interface{}) {
  if enableDebug {
    fmt.Printf(format, v...)
  }
}
