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

//Endpoint TCP struct for handling resquest/response from client
type Endpoint struct {
	listner net.Listener
	handler map[string]Handler

	m sync.RWMutex
}

//Handler handles incoming tcp requests
type Handler func(*bufio.ReadWriter, *Game) bool

//NewEndpoint endpoint constructor
func NewEndpoint() *Endpoint {

	return &Endpoint{
		handler: map[string]Handler{},
	}
}

//AddHandler adds a handler for a command
func (e *Endpoint) AddHandler(name string, f Handler) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

//Listen goroutine to start listning
func (e *Endpoint) Listen(game *Game, port string) error {
	var err error
	e.listner, err = net.Listen("tcp", port)
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
		//goroutine for handling each connection
		go e.handleTCP(conn, game)
	}
}

func (e *Endpoint) handleTCP(conn net.Conn, game *Game) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()
	for {
		cmd, err := rw.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println(err)
		}

		cmd = strings.Trim(cmd, "\n ")
		e.m.RLock()
		handleCommand, ok := e.handler[cmd]
		e.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		done := handleCommand(rw, game)
		if done {
			return
		}
	}
}

func handleJoin(rw *bufio.ReadWriter, game *Game) bool {
	name, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read string from stream", readErr)
	}
	name = strings.Trim(name, "\n ")
	if _, joined := game.Sneks[name]; joined {
		_, writeErr := rw.WriteString("Player with that name already joined\n")
		if writeErr != nil {
			log.Println("Failed to write to steam", writeErr)
		}
		flushErr := rw.Flush()
		if flushErr != nil {
			log.Println("Flush failed.", flushErr)
		}
		return true
	}
	log.Println(name + " Joined Snek")
	game.m.Lock()
	game.Sneks[name] = NewSnek(name)
	game.m.Unlock()
	_, writeErr := rw.WriteString("JOINED\n")
	if writeErr != nil {
		log.Println("Failed to write to steam", writeErr)
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Flush failed.", flushErr)
	}
	return false
}

func handleSync(rw *bufio.ReadWriter, game *Game) bool {
	name, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read player name from stream", readErr)
	}
	name = strings.Trim(name, "\n")
	if game.Sneks[name].Dead {
		_, writeErr := rw.WriteString("DEAD\n")
		if writeErr != nil {
			log.Println("Failed to write to steam", writeErr)
		}
		flushErr := rw.Flush()
		if flushErr != nil {
			log.Println("Flush failed.", flushErr)
		}
		return false
	}
	_, writeErr := rw.WriteString("NOT DEAD\n")
	if writeErr != nil {
		log.Println("Failed to write to steam", writeErr)
	}
	game.m.RLock()
	g := NewGameData(game)
	game.m.RUnlock()
	enc := gob.NewEncoder(rw)
	err := enc.Encode(g)
	if err != nil {
		log.Println("Failed to write to steam", err)
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Flush failed.", flushErr)
	}

	return false
}

func handleKey(rw *bufio.ReadWriter, game *Game) bool {
	var keyPress termbox.Event

	name, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read player name from stream", readErr)
	}
	name = strings.Trim(name, "\n")
	log.Print(name + " Moved Snek")
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&keyPress)
	if err != nil {
		log.Println("unable to decode keypress")
		return false
	}

	switch keyPress.Type {
	case termbox.EventKey:
		game.m.Lock()
		game.Sneks[name].chSnek(keyPress)
		game.m.Unlock()
	default:
		return false
	}

	return false
}

func handleQuit(rw *bufio.ReadWriter, game *Game) bool {
	name, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read player name from stream", readErr)
	}
	name = strings.Trim(name, "\n")
	log.Print(name + " Quit the Game")
	game.m.Lock()
	delete(game.Sneks, name)
	game.m.Unlock()
	_, writeErr := rw.WriteString("QUIT\n")
	if writeErr != nil {
		log.Println("Failed to write to steam", writeErr)
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Flush failed.", flushErr)
	}

	return true
}

func handleDC(rw *bufio.ReadWriter, game *Game) bool {
	log.Println("client disconnected")
	return true
}
