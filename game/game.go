package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"tictactoe/stats"
)

type Game struct {
	ID          string
	Board       [][]*string
	Mu          sync.Mutex
	Player1     Player
	Player2     Player
	CurrentTurn string
	EmptyCell   string
}

type Player struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Wins   int    `json:"wins"`
}

type Dimension struct {
	N int `json:"n"`
}

type Statistics struct {
	Players map[string]int `json:"players"`
}

var (
	Games      = make(map[string]*Game)
	GamesQueue []*Game
	GamesMu    sync.Mutex
)

func CreateGame(w http.ResponseWriter, r *http.Request) {
	var dimension Dimension

	err := json.NewDecoder(r.Body).Decode(&dimension)
	if err != nil {
		dimension.N = 3 // default
	}

	n := dimension.N

	game := &Game{
		ID:          fmt.Sprintf("G_%d", len(Games)+1),
		Board:       make([][]*string, n),
		EmptyCell:   "_",
		CurrentTurn: "",
	}
	for i := range game.Board {
		game.Board[i] = make([]*string, n)
		for j := range game.Board[i] {
			game.Board[i][j] = &game.EmptyCell
		}
	}

	GamesMu.Lock()
	Games[game.ID] = game
	GamesQueue = append(GamesQueue, game)
	GamesMu.Unlock()

	fmt.Printf("New %dx%d game created [%s]\n", n, n, game.ID)

	w.Header().Set("Content-Type", "text/plain") //change type to txt plain
	json.NewEncoder(w).Encode(game)
}

func CreateDefaultGame() {
	game := &Game{
		ID:          fmt.Sprintf("G_%d", len(Games)+1),
		Board:       make([][]*string, 3),
		EmptyCell:   "_",
		CurrentTurn: "",
	}
	for i := range game.Board {
		game.Board[i] = make([]*string, 3)
		for j := range game.Board[i] {
			game.Board[i][j] = &game.EmptyCell
		}
	}

	Games[game.ID] = game
	GamesQueue = append(GamesQueue, game)
}

func ChangeEmptyCell(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("gameID")

	var newSymbol struct {
		Symbol string `json:"symbol"`
	}

	err := json.NewDecoder(r.Body).Decode(&newSymbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GamesMu.Lock()
	game, exists := Games[gameID]
	GamesMu.Unlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	game.Mu.Lock()
	defer game.Mu.Unlock()

	oldSymbol := game.EmptyCell
	game.EmptyCell = newSymbol.Symbol
	for i := range game.Board {
		for j := range game.Board[i] {
			if *game.Board[i][j] == oldSymbol {
				game.Board[i][j] = &game.EmptyCell
			}
		}
	}

	fmt.Printf("Game character changed to %s [%s]\n", game.EmptyCell, game.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func (g *Game) getBoard(w http.ResponseWriter, r *http.Request) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g.Board)
}

func GetBoard(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("gameID")

	GamesMu.Lock()
	game, exists := Games[gameID]
	GamesMu.Unlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	game.getBoard(w, r)
}

func (g *Game) checkWinner() (string, bool) {
	n := len(g.Board)
	rows := make([]map[string]int, n)
	cols := make([]map[string]int, n)
	diag1 := make(map[string]int)
	diag2 := make(map[string]int)

	for i := 0; i < n; i++ {
		rows[i] = make(map[string]int)
		cols[i] = make(map[string]int)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			cell := *g.Board[i][j]
			if cell != g.EmptyCell {
				rows[i][cell]++
				cols[j][cell]++
				if i == j {
					diag1[cell]++
				}
				if i+j == n-1 {
					diag2[cell]++
				}

				if rows[i][cell] == n || cols[j][cell] == n || diag1[cell] == n || diag2[cell] == n {
					return cell, true
				}
			}
		}
	}

	return "", false
}

func (g *Game) MakeMove(w http.ResponseWriter, r *http.Request, playerNum string) {
	var move struct {
		Row int `json:"row"`
		Col int `json:"col"`
	}

	err := json.NewDecoder(r.Body).Decode(&move)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.Mu.Lock()
	defer g.Mu.Unlock()

	if g.CurrentTurn == "" {
		if playerNum == "Player1" {
			g.CurrentTurn = "Player1"
		} else if playerNum == "Player2" {
			g.CurrentTurn = "Player2"
		} else {
			http.Error(w, "Invalid player", http.StatusBadRequest)
			return
		}
	}

	var playerSymbol string
	var currentPlayer *Player

	if playerNum == "Player1" {
		if g.CurrentTurn != "Player1" {
			http.Error(w, "It's not your turn", http.StatusBadRequest)
			return
		}
		playerSymbol = g.Player1.Symbol
		currentPlayer = &g.Player1
	} else if playerNum == "Player2" {
		if g.CurrentTurn != "Player2" {
			http.Error(w, "It's not your turn", http.StatusBadRequest)
			return
		}
		playerSymbol = g.Player2.Symbol
		currentPlayer = &g.Player2
	} else {
		http.Error(w, "Invalid player", http.StatusBadRequest)
		return
	}

	n := len(g.Board)

	if move.Row < 0 || move.Row >= n || move.Col < 0 || move.Col >= n || *g.Board[move.Row][move.Col] != g.EmptyCell {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	g.Board[move.Row][move.Col] = &playerSymbol

	winnerSymbol, hasWinner := g.checkWinner()
	if hasWinner {
		winnerName := g.getPlayerNameBySymbol(winnerSymbol)
		fmt.Printf("%s is the winner [%s] (Symbol: %s)\n", winnerName, g.ID, winnerSymbol)
		stats.StatisticsData.Players[currentPlayer.Name]++
		stats.SaveStatistics()
		g.resetBoard()
	} else {

		if g.CurrentTurn == "Player1" {
			g.CurrentTurn = "Player2"
		} else {
			g.CurrentTurn = "Player1"
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (g *Game) getPlayerNameBySymbol(symbol string) string {
	if g.Player1.Symbol == symbol {
		return g.Player1.Name
	} else if g.Player2.Symbol == symbol {
		return g.Player2.Name
	}
	return ""
}

func (g *Game) resetBoard() {
	for i := range g.Board {
		for j := range g.Board[i] {
			g.Board[i][j] = &g.EmptyCell
		}
	}
}

func ResetBoard(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("gameID")

	GamesMu.Lock()
	game, exists := Games[gameID]
	GamesMu.Unlock()

	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	game.resetBoard()
}
