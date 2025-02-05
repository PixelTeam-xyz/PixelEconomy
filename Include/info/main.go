package info

import (
	"fmt"
	"os"
	"strings"
)

func Info(msgs ...any) {
	fmt.Println("\033[34;1m[ INFO ]:\033[0m " + join(msgs...))
}

func Infof(format string, args ...any) {
	fmt.Println("\033[34;1m[ INFO ]:\033[0m " + fmt.Sprintf(format, args...))
}

func Debug(msgs ...any) {
	fmt.Println("\033[36;1m[ DEBUG ]:\033[0m " + join(msgs...))
}

func Debugf(format string, args ...any) {
	fmt.Println("\033[36;1m[ DEBUG ]:\033[0m " + fmt.Sprintf(format, args...))
}

func Warn(msgs ...any) {
	fmt.Println("\033[33;1m[ WARN ]:\033[0m " + join(msgs...))
}

func Warnf(format string, args ...any) {
	fmt.Println("\033[33;1m[ WARN ]:\033[0m " + fmt.Sprintf(format, args...))
}

func Error(msgs ...any) {
	fmt.Fprintln(os.Stderr, "\033[31;1m[ ERROR ]:\033[0m "+join(msgs...))
}

func Errorf(format string, args ...any) {
	fmt.Fprintln(os.Stderr, "\033[31;1m[ ERROR ]:\033[0m "+fmt.Sprintf(format, args...))
}

func Fatal(msgs ...any) {
	fmt.Fprintln(os.Stderr, "\033[35;1m[ FATAL ]:\033[0m "+join(msgs...))
}

func Fatalf(format string, args ...any) {
	fmt.Fprintln(os.Stderr, "\033[35;1m[ FATAL ]:\033[0m "+fmt.Sprintf(format, args...))
}

func join(msgs ...any) string {
	var str strings.Builder
	for _, msg := range msgs {
		str.WriteString(fmt.Sprintf("%v", msg))
	}
	return str.String()
}
