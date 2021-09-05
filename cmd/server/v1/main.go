package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/matheusmosca/walrus/domain/usecases"
	pb "github.com/matheusmosca/walrus/proto"
	rpc "github.com/matheusmosca/walrus/rpc/v1"
)

// TODO pick those values through env
const (
	host = "localhost"
	port = 3000
)

func main() {
	usecase := usecases.New()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logrus.Fatal("could not listen in host %s and port %d", host, port)
	}

	server := grpc.NewServer()

	rpcMethods := rpc.New(usecase)

	pb.RegisterWalrusServer(server, rpcMethods)

	logrus.Infof("starting server on %s:%d", host, port)
	server.Serve(lis)
}
