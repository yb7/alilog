package alilog

import (
	"fmt"
)

func _debug(format string, v ...interface{}) {
	if ALI_INTERNAL_DEBUG {
		fmt.Printf("[SLS] "+format, v...)
	}
}
