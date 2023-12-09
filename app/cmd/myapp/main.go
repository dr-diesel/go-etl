package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

type AppMessage struct {
	KafkaMessage   kafka.Message
	DecodedMessage DbMessage
}

type DbMessage struct {
	RemoteAddr       net.IP
	TimeStampMilli   uint64
	Method           string
	Url              string
	ResourceId       uint64
	BytesSent        uint64
	RequestTimeMilli uint64
	ResponseStatus   uint16
	CacheStatus      string
}

const (
	kafkaBrokers = "localhost:9092"
	topic        = "http_log"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	kchan := make(chan kafka.Message)
	cchan := make(chan AppMessage)
	k2chan := make(chan kafka.Message)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokers},
		Topic:   topic,
		GroupID: "anonymizer5",
		//Partition: 0,
	})
	ctx, cancel := context.WithCancel(context.Background())
	ctxc, cancelc := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumeKafkaWorker(r, ctx, kchan)
		close(kchan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		decodeMessageWorker(kchan, cchan)
		close(cchan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		processMessageHttpWorker(cchan, k2chan)
		close(k2chan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		commitKafkaWorker(r, ctxc, k2chan)
	}()

	<-sig
	go func() {
		var i int
		for range time.Tick(1 * time.Second) {
			i++
			log.Printf("Shutting down... %d\n", i)

			if i >= 30 {
				log.Fatal("30s Fatal")
			}
		}
	}()
	go func() {
		<-sig
		log.Fatal("Double sig...")
	}()

	cancel()
	cancelc()

	wg.Wait()

	log.Println("Closing kafka.Reader...")
	r.Close()

}
