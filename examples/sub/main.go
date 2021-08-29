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

	stream, err := client.Subscribe(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}
	req := &pb.SubscribeRequest{
		Topic:        "account_created",
		ConsumerName: "example_consumer",
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
		if err = stream.Send(req); err != nil {
			logrus.Fatal(err)
		}

	}
}
