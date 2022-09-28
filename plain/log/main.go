package log

import (
	"fmt"
	"os"
)

func Fatal(s string, e error) {
	fmt.Println("[FATAL]", s, e)
	os.Exit(1)
}

func Error(s string, e error) {
	fmt.Println("[ERROR]", s, e)
}

func Info(s string) {
	fmt.Println("[INFO]", s)
}
