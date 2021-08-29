package main

import (
	"context"

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

	res, err := client.Publish(context.Background(), &pb.PublishRequest{
		Message: &pb.Message{
			Topic:     "account_created",
			CreatedBy: "example_pub",
			Body:      []byte("some body message"),
		},
	})
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("publisher response: %v", res)
}
