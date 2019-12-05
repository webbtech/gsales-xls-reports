package error

import (
	"errors"
	"fmt"
)

var errFoo = errors.New("Error 1")

func x() error {
	return fmt.Errorf("adding more context: %w", errFoo)
}

func main() {
	fmt.Printf("effFoo Type %T\n", errFoo)
	e := x()
	fmt.Printf("e from main %+v\n", e)
	if errors.Is(e, errFoo) { // Magical it works
		// handle gracefully
		fmt.Println("WORKS")
	}
}
