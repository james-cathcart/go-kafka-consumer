package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro/v2"
	"io/ioutil"
	"kafka-test/dto"
	"os"
)

var (
	codec *goavro.Codec
)

func init() {

	schema, err := ioutil.ReadFile("dto/book-schema.avsc")
	if err != nil {
		panic(err)
	}

	codec, err = goavro.NewCodec(string(schema))
	if err != nil {
		panic(err)
	}

}

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

	consumer.Subscribe("bookTopic",
		nil,
	)

	fmt.Println("Beginning message polling...")

	run := true
	for run == true {

		polledEvent := consumer.Poll(0)
		switch event := polledEvent.(type) {
		case *kafka.Message:

			nativeMap, _, err := codec.NativeFromBinary(event.Value)
			if err != nil {
				panic(err)
			}

			nativeBook := dto.StringMapToUser(nativeMap.(map[string]interface{}))
			SaveBook(nativeBook)

			fmt.Printf(
				"%% Message on %s: %s by %s\n",
				event.TopicPartition,
				//string(event.Value),
				nativeBook.Title,
				fmt.Sprintf("%s, %s", nativeBook.Author.LastName, nativeBook.Author.FirstName),
			)
		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", event)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", event)
			run = false
		}
	}

}

func SaveBook(book *dto.Book) {
	// Arbitrary mock save function
	fmt.Printf(
		"\nSaving book: %s, by %s\n",
		book.Title,
		fmt.Sprintf("%s, %s", book.Author.LastName, book.Author.FirstName),
	)
}
