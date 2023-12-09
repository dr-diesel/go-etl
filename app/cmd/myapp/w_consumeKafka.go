package main

import (
	"context"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

func consumeKafkaWorker(r *kafka.Reader, ctx context.Context, kchan chan<- kafka.Message) {
	log.Println("Running worker consumeKafka...")
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Kafka context canceled")
				return
			} else {
				log.Printf("Error reading message: %v\n", err)
			}
			return
		}
		log.Printf("Received message: %s (%d:%d)\n", msg.Topic, msg.Partition, msg.Offset)
		kchan <- msg
		log.Printf("...sent message: %s (%d:%d)\n", msg.Topic, msg.Partition, msg.Offset)

	}
}
