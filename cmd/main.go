package main

import (
	"time"

	"github.com/matheusmosca/walrus"
)

func main() {
	dispatcher := walrus.NewDispatcher()

	dispatcher.Activate()

	pub1 := walrus.NewPublisher(dispatcher)
	pub2 := walrus.NewPublisher(dispatcher)

	subscriber1 := walrus.NewSubscriber(1, dispatcher)
	subscriber2 := walrus.NewSubscriber(2, dispatcher)
	subscriber1.Subscribe()
	subscriber2.Subscribe()

	time.Sleep(time.Second * 1)

	for {
		msg1 := "{\"name\": \"matheus\"}"
		msg2 := "{\"age\": 10}"
		go pub1.Publish([]byte(msg1))
		go pub2.Publish([]byte(msg2))

		go func() {
			time.Sleep(time.Second * 6)
			subscriber := walrus.NewSubscriber(5, dispatcher)

			subscriber.Subscribe()
		}()

		time.Sleep(time.Second * 5)
	}
}
