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

func joinServer(ip string) error {
	rw, openErr := Open(ip + ClientPort)
	if openErr != nil {
		log.Println("Failed to connect", openErr)
		return openErr
	}

	_, writeErr := rw.WriteString("JOIN\n")
	if writeErr != nil {
		log.Println("Failed to write to stream", writeErr)
		return writeErr
	}
	_, writeErr = rw.WriteString("join\n")
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

func requestGame(ip string) error {
	rw, openErr := Open(ip + ClientPort)
	if openErr != nil {
		log.Println("Failed to connect", openErr)
		return openErr
	}

	_, writeErr := rw.WriteString("SYNC\n")
	if writeErr != nil {
		log.Println("Failed to write to stream", writeErr)
		return writeErr
	}
	_, writeErr = rw.WriteString("join\n")
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

func play_snek() {
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
				if event.Key == termbox.KeyArrowUp && s.Dir != "D" {
					s.Dir = "U"
				}
				if event.Key == termbox.KeyArrowDown && s.Dir != "U" {
					s.Dir = "D"
				}
				if event.Key == termbox.KeyArrowLeft && s.Dir != "R" {
					s.Dir = "L"
				}
				if event.Key == termbox.KeyArrowRight && s.Dir != "L" {
					s.Dir = "R"
				}
			case termbox.EventError:
				panic(event.Err)
			}

		default:
			g.game_tick()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
