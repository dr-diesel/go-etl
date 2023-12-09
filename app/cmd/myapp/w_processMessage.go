package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	clickhouseDriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/segmentio/kafka-go"
)

type kafkaCommit struct {
	topic     string
	partition int
	offset    int64
}

// from kafka-go :
// offsetStash holds offsets by topic => partition => offset.
type kafkaCommitStash map[string]map[int]int64

// merge updates the offsetStash with the offsets from the provided messages.
func (o kafkaCommitStash) merge(c kafkaCommit) {
	offsetsByPartition, ok := o[c.topic]
	if !ok {
		offsetsByPartition = map[int]int64{}
		o[c.topic] = offsetsByPartition
	}

	if offset, ok := offsetsByPartition[c.partition]; !ok || c.offset > offset {
		offsetsByPartition[c.partition] = c.offset
	}
}

func zeroLastOctet(ip net.IP) net.IP {
	ip = ip.To4()
	ip[3] = 0 // Zero out the last octet
	return ip
}

func processMessageWorker(cchan <-chan AppMessage, k2chan chan<- kafka.Message) {
	fmt.Println("Connecting ClickHouse...")
	ctx := context.Background()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:           false,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		//Protocol:        clickhouse.HTTP,
	})

	if err != nil {
		log.Fatalln("Unable to connect ClickHouse")
	}

	defer conn.Close()

	//conn.Exec(context.Background(), "TRUNCATE example")

	ticker := time.NewTicker(10 * time.Second)

	var batch clickhouseDriver.Batch
	f_resetbatch := func() error {
		var err error
		batch, err = conn.PrepareBatch(ctx, "INSERT INTO http_log (timestamp, resource_id, bytes_sent, request_time_milli, response_status, cache_status, method, remote_addr, url)")
		return err
	}
	if err := f_resetbatch(); err != nil {
		log.Println("Err batch:", err)
		return
	}

	commitStash := make(kafkaCommitStash, 0)

	for {
		select {
		case m, ok := <-cchan:
			if ok {
				log.Println("Adding to batch...", m.DecodedMessage)

				dbr := m.DecodedMessage

				err := batch.Append(
					time.UnixMilli(int64(dbr.TimeStampMilli)),
					dbr.ResourceId,
					dbr.BytesSent,
					dbr.RequestTimeMilli,
					dbr.ResponseStatus,
					dbr.CacheStatus,
					dbr.Method,
					zeroLastOctet(dbr.RemoteAddr),
					dbr.Url,
				)
				if err != nil {
					log.Println(err)
				}

				commit := kafkaCommit{topic: "http_log", partition: m.KafkaMessage.Partition, offset: m.KafkaMessage.Offset}
				commitStash.merge(commit)

			} else {
				//closed input channel - breaking - TODO wait for 1min or close uncommitted kafka messages
				return
			}
		case <-ticker.C:
			{
				r := batch.Rows()
				fmt.Println("FLUSHING TO CLICKHOUSE...")
				if r > 0 {
					//race ?
					batch.Send()
					log.Println("done rows:", r)

					if err := f_resetbatch(); err != nil {
						log.Println("Err batch")
						return
					}

					for topic, commit := range commitStash {
						for partition, offset := range commit {

							var cmsg kafka.Message
							cmsg.Topic = topic
							cmsg.Partition = partition
							cmsg.Offset = offset

							k2chan <- cmsg
						}
					}
				} else {
					log.Println("norows")
				}
			}

		}
	}
}
