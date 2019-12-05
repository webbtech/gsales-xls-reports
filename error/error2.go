package error

import (
	"errors"
	"fmt"
)

type myError struct {
	err  string
	more string
}

func (e *myError) Error() string {
	return fmt.Sprintf("%s: %s", e.more, e.err)
}

func x2() error {
	return fmt.Errorf("adding more context: %w", &myError{"error", "more"})
}

func main2() {
	e := x2()
	fmt.Printf("e in main2 %+v\n", e)

	var err *myError
	if ok := errors.As(e, &err); ok {
		// handle gracefully
		fmt.Println(err.more)
		fmt.Printf("err %+v\n", err)
	}
}
