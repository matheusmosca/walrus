package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/matheusmosca/walrus/config"
	"github.com/matheusmosca/walrus/domain/usecases"
	pb "github.com/matheusmosca/walrus/proto"
	rpc "github.com/matheusmosca/walrus/rpc/v1"
)

func main() {
	config, err := config.Load()
	if err != nil {
		logrus.Fatal("could not load environment config")
	}

	usecase := usecases.New()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		logrus.Fatal("could not listen in host %s and port %s", config.Host, config.Port)
	}

	server := grpc.NewServer()

	rpcMethods := rpc.New(usecase)

	pb.RegisterWalrusServer(server, rpcMethods)

	logrus.Infof("starting server on %s:%s", config.Host, config.Port)
	server.Serve(lis)
}
