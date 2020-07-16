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
		serverErr := s.serve_snek()
		if serverErr != nil {
			log.Println("Error:", serverErr)
		}
	}

	if multiplayer {
		log.Println("should connect on", *ip)
		c := NewClient(*name, *ip, *port)
		clientErr := c.join_server()
		if clientErr != nil {
			log.Println("Error:", clientErr)
		}
		c.play_snek()
	}
}
