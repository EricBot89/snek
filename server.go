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

func serve_snek() error {
	endpoint := NewEndpoint()

	endpoint.AddHandler("JOIN", handleJoin)
	endpoint.AddHandler("KEY", handleKey)
	endpoint.AddHandler("SYNC", handleSync)
	return endpoint.Listen()
}

type Handler func(*bufio.ReadWriter)

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

func (e *Endpoint) Listen() error {
	var err error
	e.listner, err = net.Listen("tcp", ServerPort) // again, look into udp for this
	if err != nil {
		return err
	}
	log.Println("listning on port " + ServerPort + " @ " + e.listner.Addr().String())
	for {
		conn, err := e.listner.Accept()
		if err != nil {
			log.Println("Failed connection attempt:", err)
			continue
		}
		go e.handleTCP(conn)
	}
}

func (e *Endpoint) handleTCP(conn net.Conn) {
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
		handleCommand(rw)
	}
}

func handleJoin(rw *bufio.ReadWriter) {
	mssg, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read string from stream", readErr)
	}
	mssg = strings.Trim(mssg, "\n ")
	log.Println(mssg)
	_, writeErr := rw.WriteString("All Good\n")
	if writeErr != nil {
		log.Println("Failed to write to steam", writeErr)
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Flush failed.", flushErr)
	}
}

func handleSync(rw *bufio.ReadWriter) {
	_, writeErr := rw.WriteString("All Good\n")
	if writeErr != nil {
		log.Println("Failed to write to steam", writeErr)
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Flush failed.", flushErr)
	}
}

func handleKey(rw *bufio.ReadWriter) {
	var keyPress termbox.Event
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&keyPress)
}
