package wildcard

import (
	"fmt"
	"testing"
)

func TestCompilePattern(t *testing.T) {
	pattern := CompilePattern("[a-f]")
	fmt.Println(pattern.IsMatch("b"))
}
