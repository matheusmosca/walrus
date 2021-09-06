package main

import (
	"context"
	"io"

	pb "github.com/matheusmosca/walrus/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:3000", opts...)
	if err != nil {
		logrus.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewWalrusClient(conn)

	req := &pb.SubscribeRequest{
		Topic: "account_created",
	}
	stream, err := client.Subscribe(context.Background(), req)
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Infof("subscriber got %v", in)
	}
}
