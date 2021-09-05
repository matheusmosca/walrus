package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/matheusmosca/walrus/config"
	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/usecases"
	"github.com/matheusmosca/walrus/domain/vos"
	pb "github.com/matheusmosca/walrus/proto"
	"github.com/matheusmosca/walrus/repositories/memory/topics"
	rpc "github.com/matheusmosca/walrus/rpc/v1"
)

func main() {
	config, err := config.Load()
	if err != nil {
		logrus.Fatal("could not load environment config")
	}

	logger := logrus.New()
	logEntry := logrus.NewEntry(logger).WithFields(logrus.Fields{
		"app_name":    config.AppName,
		"environment": config.Environment,
	})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		logEntry.Fatalf("could not listen in host %s and port %s", config.Host, config.Port)
	}

	server := grpc.NewServer()

	repository := topics.NewMemoryRepository(make(map[vos.TopicName]entities.Topic))

	usecase := usecases.New(repository)
	rpcMethods := rpc.New(usecase, logEntry)

	pb.RegisterWalrusServer(server, rpcMethods)

	logEntry.Infof("starting server on %s:%s", config.Host, config.Port)
	server.Serve(lis)
}
