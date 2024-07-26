package player

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tictactoe/game"
	"tictactoe/stats"
)

type Player game.Player

type Players stats.Statistics

func ConfigurePlayer(w http.ResponseWriter, r *http.Request) {
	var player Player

	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game.GamesMu.Lock()
	defer game.GamesMu.Unlock()

	if len(game.GamesQueue) == 0 {
		game.CreateDefaultGame()
	}
	currentGame := game.GamesQueue[0]

	currentGame.Mu.Lock()
	defer currentGame.Mu.Unlock()

	if currentGame.Player1.Name == "" { //pointer la player pe game
		currentGame.Player1 = game.Player(player)
		fmt.Printf("%s joined to game[%s]\n", player.Name, currentGame.ID)
	} else if currentGame.Player2.Name == "" {
		currentGame.Player2 = game.Player(player)
		fmt.Printf("%s joined to game[%s]\n", player.Name, currentGame.ID)
		// start game
		game.GamesQueue = game.GamesQueue[1:]
	} else {
		http.Error(w, "Both players are already set", http.StatusBadRequest)
		return
	}

	if _, exists := stats.StatisticsData.Players[player.Name]; !exists {
		stats.StatisticsData.Players[player.Name] = 0
	}

	stats.SaveStatistics()

	w.WriteHeader(http.StatusOK)
}

func ChangePlayerSymbol(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gameID := r.URL.Query().Get("gameID")
	game.GamesMu.Lock()
	currentGame, exists := game.Games[gameID]
	game.GamesMu.Unlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	currentGame.Mu.Lock()
	defer currentGame.Mu.Unlock()

	var oldSymbol string
	if currentGame.Player1.Name == request.Name {
		oldSymbol = currentGame.Player1.Symbol
		currentGame.Player1.Symbol = request.Symbol
	} else if currentGame.Player2.Name == request.Name {
		oldSymbol = currentGame.Player2.Symbol
		currentGame.Player2.Symbol = request.Symbol
	} else {
		http.Error(w, "Player not found in this game", http.StatusNotFound)
		return
	}

	//update matrix
	for i := range currentGame.Board {
		for j := range currentGame.Board[i] {
			if *currentGame.Board[i][j] == oldSymbol {
				currentGame.Board[i][j] = &request.Symbol
			}
		}
	}

	fmt.Printf("Player %s changed symbol to %s in game[%s]\n", request.Name, request.Symbol, currentGame.ID)
	stats.SaveStatistics()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
