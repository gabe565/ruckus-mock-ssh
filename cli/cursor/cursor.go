package cursor

import "fmt"

func MoveUpBeginning(n int) string {
	return fmt.Sprintf("\x1b[%dF", n)
}

func ClearScreenBelow() string {
	return "\x1B[0J"
}
