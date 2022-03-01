package server

import (
	"fmt"
	"liteq/queue/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedLiteQServer // embed UnimplementedLiteQServer
}

// start grpc server
func StartServer(port uint16) {
	log.Println("Starting server on port", port)
	address := fmt.Sprintf("0.0.0.0:%d", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf(err.Error())
	}
	grpcserver := grpc.NewServer()
	proto.RegisterLiteQServer(grpcserver, new(server))

	fmt.Printf("GRPC server listening on %s\n", address)
	grpcserver.Serve(lis)
}
