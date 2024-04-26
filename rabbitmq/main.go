package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	role string
)

func main() {
	flag.StringVar(&role, "role", "r", "sender or receiver")
	flag.Parse()

	conn, err := amqp.Dial("amqp://guest:guest@192.168.56.101:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"log",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	switch role {
	case "sender", "s":
		body := "Hello World!"
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		failOnError(err, "Failed to publish a message")
	case "receiver", "r":
		consumer, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		failOnError(err, "Failed to register a consumer")
		forever := make(chan bool)
		go func() {
			for d := range consumer {
				log.Printf("Receive a message:\n%s\n", d.Body)
			}
		}()
		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	default:
		fmt.Println("sender or receiver")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
