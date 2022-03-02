package server

import (
	"context"
	"liteq/queue/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func Dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	sv := grpc.NewServer()
	proto.RegisterLiteQServer(sv, &server{})
	go func() {
		if err := sv.Serve(listener); err != nil {
			log.Fatal(err.Error())
		}
	}()

	return func(c context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}
}
