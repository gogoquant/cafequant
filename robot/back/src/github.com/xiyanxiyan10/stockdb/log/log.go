package log

import (
	"fmt"
	"os"
	"time"
)

const (
	InfoLog    = "[ INFO  ]"
	SuccessLog = "[SUCCESS]"
	ErrorLog   = "[ ERROR ]"
	FatalLog   = "[ FATAL ]"
	RequestLog = "[REQUEST]"
)

func Log(level string, msgs ...interface{}) {
	fmt.Printf(
		"[%s] %9s %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		level,
		fmt.Sprint(msgs...),
	)
	if level == FatalLog {
		os.Exit(1)
	}
}
