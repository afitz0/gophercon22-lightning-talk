package log

import (
	"os"
)

func Fatal(s string, e error) {
	os.Exit(1)
}

func Error(s string, e error) {}
