package main

import (
	"flag"
	"liteq/queue"
	"liteq/server"
)

var (
	gport = flag.Uint("port", 8080, "GRPC TCP port to listen on")
)

func main() {
	flag.Parse()

	// init Queue
	queue.Q = queue.NewQueue()

	// grpc server
	server.StartServer(uint16(*gport))
}
