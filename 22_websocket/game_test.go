package poker_test

import (
	"../22_websocket"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestGame_Start(t *testing.T) {
	var dummyPlayerStore = &poker.StubPlayerStore{}

	t.Run("schedules Alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(5, ioutil.Discard)

		cases := []poker.ScheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(7, ioutil.Discard)

		cases := []poker.ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})
}

func TestGame_Finish(t *testing.T) {
	dummyBlindAlerter := &poker.SpyBlindAlerter{}
	store := &poker.StubPlayerStore{}
	game := poker.NewGame(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner)
}

func checkSchedulingCases(cases []poker.ScheduledAlert, t *testing.T, blindAlerter *poker.SpyBlindAlerter) {
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(blindAlerter.Alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.Alerts)
			}

			got := blindAlerter.Alerts[i]
			poker.AssertScheduledAlert(t, got, want)
		})
	}
}
