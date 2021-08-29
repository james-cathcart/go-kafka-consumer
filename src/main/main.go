package main

import (
    "fmt"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

    fmt.Println("Initializing consumer...")

    consumer, err := kafka.NewConsumer(
        &kafka.ConfigMap {
            "bootstrap.servers": "localhost",
            "group.id": "myGroup",
            "auto.offset.reset": "earliest",
        },
    )

    if err != nil {
        panic(err)
    }
    defer consumer.Close()

    consumer.SubscribeTopics(
        []string{ "myTopic", "^aRegex.*[Tt]opic" },
        nil,
    )

    fmt.Println("Beginning message polling...")
    for {

        msg, err := consumer.ReadMessage(-1)
        if err == nil {
            fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
        } else {
            fmt.Printf("Consumer error: %v (%v)\n", err, msg)
        }

    }
}