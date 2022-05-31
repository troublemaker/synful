package log

import (
	"fmt"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"

	wheel = "|/-\\"
)

var wheelState uint8

func init() {
	wheelState = 0
}

func MoveCursorUp(n int) {
	fmt.Printf("\033[%dF", n)
}

func MoveCursorDown(n int) {
	fmt.Printf("\033[%dE", n)
}

func ColorPrint(text string, color string, a ...interface{}) {
	fmt.Printf(color+text+Reset, a...)
}

func Spinner(color string) {
	wheelState++
	fmt.Printf(color + string(wheel[wheelState%4]) + Reset)
}

func ClearScreen() {
	fmt.Printf("\033[2J\n")
}

func DebugHex(data []byte, line int) {
	fmt.Printf(Red)
	fmt.Println("-----------")
	for i, b := range data {
		fmt.Printf("%02X ", b)
		if (i+1)%line == 0 {
			fmt.Println("")
		}
	}
	fmt.Printf("-----------\n\n")
	fmt.Printf(Reset)
}
