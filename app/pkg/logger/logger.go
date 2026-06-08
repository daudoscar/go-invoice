package logger

import (
	"fmt"
	"os"
	"time"
)

// ANSI Color Codes untuk terminal
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func format(level, color, msg string) string {
	ts := time.Now().Format("15:04:05")
	return fmt.Sprintf("%s[%s] %-7s %s%s", color, ts, level, msg, Reset)
}

func Info(m string)    { fmt.Println(format("INFO", Blue, m)) }
func Success(m string) { fmt.Println(format("SUCCESS", Green, m)) }
func Warn(m string)    { fmt.Println(format("WARN", Yellow, m)) }

func Error(m string, err error) {
	msg := m
	if err != nil {
		msg = fmt.Sprintf("%s | Error: %v", m, err)
	}
	fmt.Println(format("ERROR", Red, msg))
}

func Fatal(m string, err error) {
	errMsg := m
	if err != nil {
		errMsg = m + " | " + err.Error()
	}
	fmt.Println(format("FATAL", Red, errMsg))
	os.Exit(1)
}
