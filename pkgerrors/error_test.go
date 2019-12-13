package pkgerrors

import (
	"math"
	"testing"
)

// TestAbs function
func TestAbs(t *testing.T) {
	got := math.Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}

// TestError function
func TestMainError(t *testing.T) {
	main()
}

// TestError2 function
func TestMain2Error(t *testing.T) {
	main2()
	// e := x2()
	// fmt.Printf("e from x2: %+v\n", e)
}

// TestError3 function
func TestMain3Error(t *testing.T) {
	main3()
	// e := x2()
	// fmt.Printf("e from x2: %+v\n", e)
}
