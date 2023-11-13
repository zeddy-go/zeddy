package errx

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	err := New("123")
	err2 := WrapWithSkip(err, "321", 0)
	fmt.Printf("%#v", err2)
}
