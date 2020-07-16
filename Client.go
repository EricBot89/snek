package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"log"
	"net"
	"strings"
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
	game GameData
}

func NewClient(name string, ip string, port string) Snek_Client {
	return Snek_Client{
		name: name,
		ip:   ip,
		port: port,
	}
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
	response = strings.Trim(response, "\n")
	if response == "JOINED" {
		return nil
	}
	return errors.New("Player already joined with that name")
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
	_, writeErr = rw.WriteString(client.name + "\n")
	if writeErr != nil {
		log.Println("Failed to write to stream", writeErr)
		return writeErr
	}
	enc := gob.NewEncoder(rw)
	err := enc.Encode(key)
	if err != nil {
		log.Println("Failed to encode termbox event")
	}
	flushErr := rw.Flush()
	if flushErr != nil {
		log.Println("Failed to flush", flushErr)
		return flushErr
	}
	return nil
}

func (client *Snek_Client) request_game() error {
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
	response = strings.Trim(response, "\n")
	if response == "DEAD" {
		return errors.New("You are Dead")
	}
	var game GameData
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&game)
	if err != nil {
		log.Println("Failed to sync game")
	}
	client.game = game
	return nil
}

func (client *Snek_Client) quit_game() error {
	rw, openErr := Open(client.ip + client.port)
	if openErr != nil {
		log.Println("Failed to connect", openErr)
		return openErr
	}

	_, writeErr := rw.WriteString("QUIT\n")
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
	return nil
}

func Open(addr string) (*bufio.ReadWriter, error) {
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
					client.quit_game()
					break loop
				}
				client.send_key(event)
			case termbox.EventError:
				panic(event.Err)
			}

		default:
			err := client.request_game()
			if err != nil {
				client.quit_game()
				log.Println(err)
			}
			draw(&client.game)
			time.Sleep(15 * time.Millisecond)
		}
	}
}
