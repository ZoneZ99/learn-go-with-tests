package poker

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGETPlayers(t *testing.T) {
	dummyGame := &GameSpy{}
	store := StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		winCalls: nil,
		league:   nil,
	}
	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	dummyGame := &GameSpy{}
	store := StubPlayerStore{
		scores:   map[string]int{},
		winCalls: nil,
		league:   nil,
	}
	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("it records wins on POST", func(t *testing.T) {
		request := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusAccepted)
		AssertPlayerWin(t, &store, "Pepper")
	})
}

func TestLeague(t *testing.T) {
	dummyGame := &GameSpy{}

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{
			scores:   nil,
			winCalls: nil,
			league:   wantedLeague,
		}
		server := mustMakePlayerServer(t, &store, dummyGame)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		AssertStatus(t, response.Code, http.StatusOK)
		AssertLeague(t, got, wantedLeague)

		const jsonContentType = "application/json"

		AssertContentType(t, response, jsonContentType)
	})
}

func TestGame(t *testing.T) {
	dummyGame := &GameSpy{}

	t.Run("GET /playGame returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{}, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("start a playGame with 3 players and declare Ruth the winner", func(t *testing.T) {
		dummyPlayerStore := &StubPlayerStore{}
		game := &GameSpy{}
		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		writeWsMessage(t, ws, "3")
		writeWsMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		AssertGameStartedWith(t, game, 3)
		AssertGameFinishedWith(t, game, winner)
	})

	t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
		dummyPlayerStore := &StubPlayerStore{}
		wantedBlindAlert := "Blind is 100"
		winner := "Ruth"

		game := &GameSpy{
			BlindAlert: []byte(wantedBlindAlert),
		}
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		writeWsMessage(t, ws, "3")
		writeWsMessage(t, ws, winner)

		ten := 10 * time.Millisecond
		time.Sleep(ten)

		AssertGameStartedWith(t, game, 3)
		AssertGameFinishedWith(t, game, winner)
		within(t, ten, func() {
			AssertWebsocketGotMsg(t, ws, wantedBlindAlert)
		})
	})
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league League) {
	t.Helper()
	league, err := NewLeague(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return
}

func newGameRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return req
}

func mustMakePlayerServer(t *testing.T, store PlayerStore, game Game) *PlayerServer {
	server, err := NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}

func writeWsMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}
