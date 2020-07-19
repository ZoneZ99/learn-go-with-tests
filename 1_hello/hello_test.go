package hello

import "testing"

func TestHello(test *testing.T) {

	assertCorrectMessage := func(test *testing.T, got, want string) {
		test.Helper()
		if got != want {
			test.Errorf("got %q want %q", got, want)
		}
	}

	test.Run("saying hello to people", func(test *testing.T) {
		got := Hello("Chris", "")
		want := "Hello, Chris"
		assertCorrectMessage(test, got, want)
	})

	test.Run("say 'Hello, World' when an empty string is supplied", func(test *testing.T) {
		got := Hello("", "")
		want := "Hello, World"
		assertCorrectMessage(test, got, want)
	})

	test.Run("in Spanish", func(test *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(test, got, want)
	})

	test.Run("in French", func(test *testing.T) {
		got := Hello("Paul", "French")
		want := "Bonjour, Paul"
		assertCorrectMessage(test, got, want)
	})
} 