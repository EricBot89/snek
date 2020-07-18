package main

//SnekServer struct for a snek server
type SnekServer struct {
	endpoint *Endpoint
	game     *Game
	players  []string
	port     string
}

//NewServer starts a snek server on a port
func NewServer(port string) SnekServer {
	return SnekServer{port: port}
}

func (server *SnekServer) serveSnek() error {

	game := NewGame()
	server.game = game
	endpoint := NewEndpoint()
	server.endpoint = endpoint
	endpoint.AddHandler("JOIN", handleJoin)
	endpoint.AddHandler("KEY", handleKey)
	endpoint.AddHandler("SYNC", handleSync)
	endpoint.AddHandler("QUIT", handleQuit)
	endpoint.AddHandler("", handleDC)
	go server.game.runSnek()
	return server.endpoint.Listen(server.game, server.port)
}
