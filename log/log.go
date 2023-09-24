package log

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var std = log.New(os.Stderr, "", 0)

func Debugf(format string, v ...any) {
	printMessage(LevelDebug, fmt.Sprintf(format, v...))
}

func Infof(format string, v ...any) {
	printMessage(LevelInfo, fmt.Sprintf(format, v...))
}

func Warningf(format string, v ...any) {
	printMessage(LevelWarning, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...any) {
	printMessage(LevelError, fmt.Sprintf(format, v...))
}

func Criticalf(format string, v ...any) {
	printMessage(LevelCritical, fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...any) {
	printMessage(LevelCritical, fmt.Sprintf(format, v...))
	os.Exit(1)
}

type severity string

const (
	LevelDebug    severity = "DEBUG"
	LevelInfo     severity = "INFO"
	LevelWarning  severity = "WARNING"
	LevelError    severity = "ERROR"
	LevelCritical severity = "CRITICAL"
)

type message struct {
	Severity severity `json:"severity"`
	Message  string   `json:"message"`
}

func printMessage(sev severity, msg string) {
	m := message{
		Severity: sev,
		Message:  msg,
	}
	b, _ := json.Marshal(m)
	_ = std.Output(3, string(b))
}
