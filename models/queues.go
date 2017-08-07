package models

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func QueueExportTask(task ExportTask, channel *amqp.Channel) error {
	queue, err := channel.QueueDeclare(
		"descriptor_exporter_queue", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.WithError(err).Error("cannot define export task queue")
		return err
	}

	bytes, err := json.Marshal(task)

	if err != nil {
		log.WithError(err).Error("cannot marshall export task")
		return err
	}

	return channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         bytes,
		})
}

type ExportTask struct {
	Name              string
	EntityDescriptors []string
}
