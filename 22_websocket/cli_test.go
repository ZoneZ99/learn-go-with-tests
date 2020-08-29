package poker_test

import (
	"../22_websocket"
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {

	var dummyStdOut = &bytes.Buffer{}
	var dummyBlindAlerter = &poker.SpyBlindAlerter{}
	var dummyPlayerStore = &poker.StubPlayerStore{}

	t.Run("finish game with 'Chris' as winner", func(t *testing.T) {
		in := strings.NewReader("1\nChris wins\n")
		game := &poker.GameSpy{}
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		if game.FinishedWith != "Chris" {
			t.Errorf("expected finish called with 'Chris' but got %q", game.FinishedWith)
		}
	})

	t.Run("record 'Cleo' win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nCleo wins\n")
		game := &poker.GameSpy{}
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		if game.FinishedWith != "Cleo" {
			t.Errorf("expected finish called with 'Cleo' but got %q", game.FinishedWith)
		}
	})

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &poker.GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}

		if game.StartedWith != 7 {
			t.Errorf("wanted Start called with 7 but got %d", game.StartedWith)
		}
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &poker.GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Errorf("game should not have started")
		}

		wantPrompt := poker.PlayerPrompt + poker.BadPlayerInputErrMsg

		poker.AssertMessagesSentToUser(t, stdout, wantPrompt)
	})

	t.Run("do not read beyond the first newline", func(t *testing.T) {
		in := failOnEndReader{
			t:   t,
			rdr: strings.NewReader("1\nChris wins\n hello there"),
		}

		game := poker.NewGame(dummyBlindAlerter, dummyPlayerStore)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()
	})
}

type failOnEndReader struct {
	t   *testing.T
	rdr io.Reader
}

func (m failOnEndReader) Read(p []byte) (n int, err error) {
	n, err = m.rdr.Read(p)

	if n == 0 || err == io.EOF {
		m.t.Fatal("Read to the end when you shouldn't have")
	}

	return n, err
}
