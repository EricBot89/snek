package main

import (
	"bufio"
	"log"
	"net"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	ClientPort = ":8080"
)

type Snek_Client struct {
	name string
	ip   string
	port string
	game Game
}

func (client *Snek_Client) join_server() error {
	rw, openErr := Open(client.ip + client.port)
	if openErr != nil {
		log.Println("Failed to connect", openErr)
		return openErr
	}

	_, writeErr := rw.WriteString("JOIN\n")
	if writeErr != nil {
		log.Println("Failed to write to stream", writeErr)
		return writeErr
	}
	_, writeErr = rw.WriteString(client.name + "\n")
	if writeErr != nil {
		log.Println("Failed to write to stream", writeErr)
		return writeErr
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Failed to flush", flushErr)
		return flushErr
	}
	response, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read response", readErr)
		return readErr
	}

	log.Println("Read Response" + response)

	return nil
}

func (client *Snek_Client) send_key(key termbox.Event) error {
	rw, openErr := Open(client.ip + client.port)
	if openErr != nil {
		log.Println("Failed to connect", openErr)
		return openErr
	}

	_, writeErr := rw.WriteString("KEY\n")
	if writeErr != nil {
		log.Println("Failed to write to stream", writeErr)
		return writeErr
	}

	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Failed to flush", flushErr)
		return flushErr
	}
	response, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read response", readErr)
		return readErr
	}

	log.Println("Read Response" + response)

	return nil
}

func (client *Snek_Client) requestGame() error {
	rw, openErr := Open(client.ip + client.port)
	if openErr != nil {
		log.Println("Failed to connect", openErr)
		return openErr
	}

	_, writeErr := rw.WriteString("SYNC\n")
	if writeErr != nil {
		log.Println("Failed to write to stream", writeErr)
		return writeErr
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Failed to flush", flushErr)
		return flushErr
	}
	response, readErr := rw.ReadString('\n')
	if readErr != nil {
		log.Println("Failed to read response", readErr)
		return readErr
	}

	log.Println("Read Response" + response)

	return nil
}

func Open(addr string) (*bufio.ReadWriter, error) {
	log.Println("opening connection to snek server at", addr)
	conn, err := net.Dial("tcp", addr) //Look into UDP for this application, TCP might not be the best choice
	if err != nil {
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func (client *Snek_Client) play_snek() {
	err := termbox.Init()

	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

loop:
	for {
		select {
		case event := <-eventQueue:
			switch event.Type {
			case termbox.EventKey:
				if event.Key == termbox.KeyCtrlQ {
					break loop
				}
				client.send_key(event)
			case termbox.EventError:
				panic(event.Err)
			}

		default:
			client.requestGame()
			time.Sleep(5 * time.Millisecond)
		}
	}
}
