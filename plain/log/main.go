package log

import (
	"lightning/app/plain/constants"

	"fmt"
	"os"
	"time"
)

func Fatal(s ...any) {
	now := time.Now().Format(time.UnixDate)
	fmt.Println(append([]any{"[FATAL]", now, " -- "}, s...)...)
	os.Exit(1)
}

func Error(s ...any) {
	now := time.Now().Format(time.UnixDate)
	fmt.Println(append([]any{"[ERROR]", now, " -- "}, s...)...)
}

func Info(s ...any) {
	now := time.Now().Format(time.UnixDate)
	fmt.Println(append([]any{"[INFO]", now, " -- "}, s...)...)
}

func Debug(s ...any) {
	if constants.DEBUG {
		now := time.Now().Format(time.UnixDate)
		fmt.Println(append([]any{"[DEBUG]", now, " -- "}, s...)...)
	}
}
