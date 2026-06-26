package utils

import (
	"fmt"
	"time"
)

func PanicLog(log string) {
	t := time.Now()
	message := fmt.Sprintf("[%s] %s", t.Format(time.DateTime), log)
	panic(message)
}
