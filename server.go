package main

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/nsf/termbox-go"
)

const (
	ServerPort = ":8080"
)

type Snek_Server struct {
	endpoint Endpoint
	game     *Game
	players  []string
	port     string
}

func NewServer(port string) Snek_Server {
	return Snek_Server{port: port}
}

func (server *Snek_Server) serve_snek() error {

	game := NewGame()
	server.game = &game
	endpoint := NewEndpoint()
	server.endpoint = *endpoint
	endpoint.AddHandler("JOIN", handleJoin)
	endpoint.AddHandler("KEY", handleKey)
	endpoint.AddHandler("SYNC", handleSync)
	go server.game.run_snek()
	return endpoint.Listen(server.game, server.port)
}

type Handler func(*bufio.ReadWriter, *Game)

type Endpoint struct {
	listner net.Listener
	handler map[string]Handler

	m sync.RWMutex //Maps are not threadsafe, so we need a mutex to control access.
}

func NewEndpoint() *Endpoint {

	return &Endpoint{
		handler: map[string]Handler{},
	}
}

func (e *Endpoint) AddHandler(name string, f Handler) {
	e.m.Lock()          //Lock access to the endpoint to prevent race conditions? something else? check this out
	e.handler[name] = f //Assign handler to the endpoint
	e.m.Unlock()        //unlock access to the endpoint
}

func (e *Endpoint) Listen(game *Game, port string) error {
	var err error
	e.listner, err = net.Listen("tcp", port) // again, look into udp for this
	if err != nil {
		return err
	}
	log.Println("listning on port " + port + " @ " + e.listner.Addr().String())
	for {
		conn, err := e.listner.Accept()
		if err != nil {
			log.Println("Failed connection attempt:", err)
			continue
		}
		go e.handleTCP(conn, game)
	}
}

func (e *Endpoint) handleTCP(conn net.Conn, game *Game) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)) //Grab a readwriter for this connection
	defer conn.Close()                                                      //close the connection when we're all done
	for {
		cmd, err := rw.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("EOF reached, done reading successfully")
				return
			}
			log.Println("Something fucked up", err)
		}

		cmd = strings.Trim(cmd, "\n ")

		log.Println("Read Command " + cmd)

		e.m.RLock()
		handleCommand, ok := e.handler[cmd]
		e.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		handleCommand(rw, game)
	}
}

func handleJoin(rw *bufio.ReadWriter, game *Game) {
	name, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read string from stream", readErr)
	}
	name = strings.Trim(name, "\n ")
	log.Println(name)
	game.Sneks[name] = NewSnek()
	_, writeErr := rw.WriteString("All Good\n")
	if writeErr != nil {
		log.Println("Failed to write to steam", writeErr)
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Flush failed.", flushErr)
	}
}

func handleSync(rw *bufio.ReadWriter, game *Game) {
	enc := gob.NewEncoder(rw)
	err := enc.Encode(*game)
	if err != nil {
		log.Println("Failed to write to steam", err)
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Flush failed.", flushErr)
	}
}

func handleKey(rw *bufio.ReadWriter, game *Game) {
	var keyPress termbox.Event

	name, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read string from stream", readErr)
	}
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&keyPress)
	if err != nil {
		log.Println("unable to decode keypress")
		return
	}

	switch keyPress.Type {
	case termbox.EventKey:
		var s = game.Sneks[name]
		if keyPress.Key == termbox.KeyArrowUp && s.Dir != "D" {
			s.Dir = "U"
		}
		if keyPress.Key == termbox.KeyArrowDown && s.Dir != "U" {
			s.Dir = "D"
		}
		if keyPress.Key == termbox.KeyArrowLeft && s.Dir != "R" {
			s.Dir = "L"
		}
		if keyPress.Key == termbox.KeyArrowRight && s.Dir != "L" {
			s.Dir = "R"
		}
		game.Sneks[name] = s
	default:
		return
	}
}
