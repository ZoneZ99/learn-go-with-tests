package iteration

import (
	"testing"
)

func TestRepeat(test *testing.T) {
	repeated := Repeat("a", 5)
	expected := "aaaaa"

	if repeated != expected {
		test.Errorf("expected %q but got %q", expected, repeated)
	}
}

func BenchmarkRepeat(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		Repeat("a", 5)
	}
}