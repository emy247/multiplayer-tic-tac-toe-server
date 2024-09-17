package main

import (
	"tictactoe/server"
	"tictactoe/stats"
)

func main() {

	stats.LoadStatistics()
	server.StartRouter()

	defer stats.SaveStatistics()
}
