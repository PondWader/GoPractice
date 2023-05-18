package utils

import "fmt"

func Info(a ...any) {
	a = append([]any{Grey("[") + Cyan("INFO") + Grey("]")}, a...)
	fmt.Println(a...)
}

func Error(a ...any) {
	a = append([]any{Grey("[") + Red("ERROR") + Grey("]")}, a...)
	fmt.Println(a...)
}
