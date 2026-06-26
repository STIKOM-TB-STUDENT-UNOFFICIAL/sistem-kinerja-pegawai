package utils

import (
	"fmt"
	"time"
)

func Log(log string) {
	t := time.Now()
	message := fmt.Sprintf("[%s] %s", t.Format(time.DateTime), log)
	fmt.Println(message)
}
