package ArraysAndSlices

import (
	"reflect"
	"testing"
)

func TestSum(test *testing.T) {

	test.Run("collection of 5 numbers", func(test *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			test.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	test.Run("collection of any size", func(test *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			test.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}

func TestSumAll(test *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		test.Errorf("got %v want %v", got, want)
	}
}

func TestSumAllTails(test *testing.T) {

	checkSums := func(test *testing.T, got, want [] int) {
		test.Helper()
		if !reflect.DeepEqual(got, want) {
			test.Errorf("got %v want %v", got, want)
		}
	}

	test.Run("make the sums of some slices", func(test *testing.T) {
		got := SumAllTails([]int{1, 2, 3}, []int{0, 9, 10})
		want := []int{5, 19}
		checkSums(test, got, want)
	})

	test.Run("safely sum empty slices", func(test *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checkSums(test, got, want)
	})

}