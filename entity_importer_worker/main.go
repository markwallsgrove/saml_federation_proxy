package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/markwallsgrove/saml_federation_proxy/models"
	"github.com/streadway/amqp"
)

func setChannelQOS(channel *amqp.Channel) error {
	return channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
}

func createQueue(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare(
		"ingest_queue", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
}

func main() {
	conn, err := amqp.Dial(os.Getenv("QUEUE_CONN"))
	if err != nil {
		log.Fatal("cannot connect to queue", err)
		return
	}

	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("cannot create channel in queue", err)
		return
	}

	if err = setChannelQOS(channel); err != nil {
		log.Fatal("cannot set qos on channel", err)
		return
	}

	queue, err := createQueue(channel)
	if err != nil {
		log.Fatal("cannot create queue")
		return
	}

	entityDescriptors, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		log.Fatal("cannot consume queue", err)
		return
	}

	for delivery := range entityDescriptors {
		var entityDescriptor models.EntityDescriptor

		if err = json.Unmarshal(delivery.Body, &entityDescriptor); err != nil {
			fmt.Println("cannot get body of queue message", err, delivery.MessageId)
			continue
		}

		fmt.Println("entity descriptor", entityDescriptor)
		delivery.Ack(false)
	}
}
