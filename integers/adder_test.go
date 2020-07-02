package integers

import (
	"testing"
	"fmt"
)

func TestAdder(test *testing.T) {
	sum := Add(2, 2)
	expected := 4

	if sum != expected {
		test.Errorf("expected '%d' but got '%d'", expected, sum)
	}
}

func ExampleAdd() {
	sum := Add(1,5)
	fmt.Println(sum)
}