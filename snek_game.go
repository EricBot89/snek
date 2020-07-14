package main

import (
	"flag"
	"log"
)

const title = "SNEK"

func main() {
	serve := flag.Bool("serve", false, "flag to set server mode")
	connect := flag.Bool("c", false, "connect to a snek server")
	ip := flag.String("ip", "", "ip of multi-snek server to connect to")
	flag.Parse()

	server := *serve

	multiplayer := *connect
	server_ip := *ip

	if server {
		serverErr := serve_snek()
		if serverErr != nil {
			log.Println("Error:", serverErr)
		}
	}

	if multiplayer {
		log.Println("should connect on", server_ip)
		clientErr := joinServer(server_ip)
		if clientErr != nil {
			log.Println("Error:", clientErr)
		}
	}
}
