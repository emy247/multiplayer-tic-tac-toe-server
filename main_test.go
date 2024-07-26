package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"tictactoe/game"
	"tictactoe/player"
	//"tictactoe/player"
)

func TestNewGame(t *testing.T) {

	http.HandleFunc("/newgame", game.CreateGame)
	router := http.DefaultServeMux

	//creare joc 3x3
	body := map[string]int{"n": 3}

	//transformare in json
	jsonBody, _ := json.Marshal(body)

	//creare request
	req := httptest.NewRequest("POST", "/newgame", bytes.NewReader(jsonBody))

	//
	w := httptest.NewRecorder()

	//
	router.ServeHTTP(w, req)

	// pana aici face testul si trece pass (daca n exista si e mai mare ca 0)

	if w.Code != http.StatusOK { //test daca status e ok
		t.Fatalf("expected status ok, got %v", w.Code)
	}

}

func TestConfigurePlayer(t *testing.T) {

	http.HandleFunc("/configureplayer", player.ConfigurePlayer)
	router := http.DefaultServeMux

	// body player
	playerBody := map[string]string{"name": "Cristi", "symbol": "A"}
	playerBody2 := map[string]string{"name": "Marius", "symbol": "M"}

	jsonBody, _ := json.Marshal(playerBody)
	jsonBody2, _ := json.Marshal(playerBody2)

	req := httptest.NewRequest("POST", "/configureplayer?gameID=G_1", bytes.NewReader(jsonBody))
	req2 := httptest.NewRequest("POST", "/configureplayer?gameID=G_1", bytes.NewReader(jsonBody2))

	// recorder pt raspuns
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	router.ServeHTTP(w2, req2)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status ok, got %v", w.Code)
	}
}

func TestMakeMove(t *testing.T) {
	// handler pentru mutari
	http.HandleFunc("/player1/move", func(w http.ResponseWriter, r *http.Request) {
		gameID := r.URL.Query().Get("gameID")
		game.GamesMu.Lock()
		g, exists := game.Games[gameID]
		game.GamesMu.Unlock()
		if !exists {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		g.MakeMove(w, r, "Player1")
	})
	http.HandleFunc("/player2/move", func(w http.ResponseWriter, r *http.Request) {
		gameID := r.URL.Query().Get("gameID")
		game.GamesMu.Lock()
		g, exists := game.Games[gameID]
		game.GamesMu.Unlock()
		if !exists {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		g.MakeMove(w, r, "Player2")
	})
	router := http.DefaultServeMux

	// mutari pana un player castiga
	moves := []struct {
		player string
		row    int
		col    int
	}{
		{"player1", 0, 0}, // X
		{"player2", 0, 1}, // 0
		{"player1", 1, 1}, // X
		{"player2", 0, 2}, // 0
		{"player1", 2, 2}, // X
	}

	for _, move := range moves {
		moveBody := map[string]int{"row": move.row, "col": move.col}
		jsonBody, _ := json.Marshal(moveBody)
		req := httptest.NewRequest("POST", "/"+move.player+"/move?gameID=G_1", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status ok, got %v", w.Code)
		}
	}
}
