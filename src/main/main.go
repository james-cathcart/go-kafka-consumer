package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func main() {

	fmt.Println("Initializing consumer...")

	consumer, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": "localhost",
			"group.id":          "myGroup",
			"auto.offset.reset": "earliest",
		},
	)

	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	consumer.Subscribe("myTopic",
		nil,
	)

	fmt.Println("Beginning message polling...")

	run := true
	for run == true {

		polledEvent := consumer.Poll(0)
		switch event := polledEvent.(type) {
		case *kafka.Message:
			fmt.Printf(
				"%% Message on %s:\n%s\n",
				event.TopicPartition,
				string(event.Value),
			)
		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", event)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", event)
			run = false
		}
	}

}
