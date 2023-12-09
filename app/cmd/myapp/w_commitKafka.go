package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func commitKafkaWorker(r *kafka.Reader, ctxc context.Context, k2chan <-chan kafka.Message) {
	for m := range k2chan {
		log.Println("Commiting:", m.Topic, m.Partition, m.Offset)
		r.CommitMessages(ctxc, m)
	}
}
