package sigmarsolver

import "fmt"

var log = false

func Log(args ...any) {
	if log {
		fmt.Println(args...)
	}
}
func Logf(format string, args ...any) {
	if log {
		fmt.Printf(format, args...)
	}
}
