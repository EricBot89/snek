package main

import (
	"bufio"
)

const (
	ServerPort = ":8080"
)

type Snek_Server struct {
	endpoint *Endpoint
	game     *Game
	players  []string
	port     string
}

func NewServer(port string) Snek_Server {
	return Snek_Server{port: port}
}

func (server *Snek_Server) serve_snek() error {

	game := NewGame()
	server.game = game
	endpoint := NewEndpoint()
	server.endpoint = endpoint
	endpoint.AddHandler("JOIN", handleJoin)
	endpoint.AddHandler("KEY", handleKey)
	endpoint.AddHandler("SYNC", handleSync)
	go server.game.run_snek()
	return server.endpoint.Listen(server.game, server.port)
}

type Handler func(*bufio.ReadWriter, *Game)
