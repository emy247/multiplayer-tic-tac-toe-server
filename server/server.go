package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tictactoe/game"
	"tictactoe/player"
	"tictactoe/stats"
)

func StartRouter() {
	http.HandleFunc("/newgame", game.CreateGame)
	http.HandleFunc("/configureplayer", player.ConfigurePlayer)
	http.HandleFunc("/player1/move", func(w http.ResponseWriter, r *http.Request) {
		gameID := r.URL.Query().Get("gameID")

		game.GamesMu.Lock()
		g, exists := game.Games[gameID] // ]'g' type *game.Game
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

	http.HandleFunc("/changeemptycell", game.ChangeEmptyCell)
	http.HandleFunc("/board", game.GetBoard)
	http.HandleFunc("/resetboard", game.ResetBoard)
	http.HandleFunc("/statistics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats.StatisticsData.Players)
	})

	http.HandleFunc("/changeplayersymbol", player.ChangePlayerSymbol)

	fmt.Println("Server is running on port 5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
