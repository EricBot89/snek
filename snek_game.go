package main

import (
	"flag"
	"log"
)

const title = "SNEK"

func main() {
	serve := flag.Bool("serve", false, "flag to set server mode")
	connect := flag.Bool("c", false, "connect to a snek server")
	ip := flag.String("ip", "127.0.0.1", "ip of multi-snek server to connect to or create")
	port := flag.String("p", ":8080", "port to connect on")
	name := flag.String("n", "dorkus", "name to join snek as")
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
}
