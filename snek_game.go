package main

import (
	"flag"
	"fmt"
	"log"
)

const title = "SNEK"

func main() {
	serve := flag.Bool("s", false, "serve snek")
	connect := flag.Bool("c", false, "connect to a snek server")
	ip := flag.String("ip", "127.0.0.1", "ip")
	port := flag.String("p", ":8080", "port")
	name := flag.String("n", "PLAYER", "player name")
	flag.Parse()

	server := *serve
	multiplayer := *connect

	if server {
		s := NewServer(*port)
		serverErr := s.serveSnek()
		if serverErr != nil {
			log.Println("Error:", serverErr)
		}
	}

	if multiplayer {
		log.Println("Connecting to snek server on ", *ip)
		c := NewClient(*name, *ip, *port)
		clientErr := c.joinServer()
		if clientErr != nil {
			log.Println("Error:", clientErr)
			return
		}
		end := c.playSnek()
		log.Println("Game Ended: " + end)
	}
	if !server && !multiplayer {
		fmt.Println("Usage of snek_game:")
		flag.PrintDefaults()
	}
}
