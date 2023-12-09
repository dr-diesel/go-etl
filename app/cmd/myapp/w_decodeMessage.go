package main

import (
	"log"
	"net"

	"capnproto.org/go/capnp/v3"
	"g.o/~/anonymizer/cmd/myapp/http_log_schema"
	"github.com/segmentio/kafka-go"
)

func decodeMessageWorker(kchan <-chan kafka.Message, cchan chan<- AppMessage) {
	log.Println("Running worker decodeKafka...")
	for m := range kchan {
		log.Print("Decoding message:...")

		msg, err := capnp.Unmarshal(m.Value)
		if err != nil {
			log.Println("[DECODER] Unable to unmarshal, dropping message:", err)
			continue
		}

		h, err := http_log_schema.ReadRootHttpLogRecord(msg)
		if err != nil {
			log.Println("[DECODER] Unable to ReadRootHttLogRecord, dropping message:", err)
			continue
		}

		dmsg := DbMessage{}

		ipstr, err := h.RemoteAddr()
		if err != nil {
			log.Println("[DECODER] Unable to read RemoteAddr(), dropping message: ", err)
			continue
		}

		ip := net.ParseIP(ipstr).To4()
		dmsg.RemoteAddr = ip

		tsmilli := h.TimestampEpochMilli()
		if tsmilli == 0 {
			log.Println("[DECODER] Unable to read TimestampEpochMilli(), dropping message: ", err)
			continue
		}
		dmsg.TimeStampMilli = tsmilli

		resource_id := h.ResourceId()
		if tsmilli == 0 {
			log.Println("[DECODER] Unable to read ResourceId(), dropping message: ", err)
			continue
		}
		dmsg.ResourceId = resource_id

		request_time_milli := h.RequestTimeMilli()
		if request_time_milli == 0 {
			log.Println("[DECODER] Unable to read RequestTimeMilli(), dropping message: ", err)
			continue
		}
		dmsg.RequestTimeMilli = request_time_milli

		response_status := h.ResponseStatus()
		if response_status == 0 {
			log.Println("[DECODER] Unable to read ResponseStatus(), dropping message: ", err)
			continue
		}
		dmsg.ResponseStatus = response_status

		cache_status, err := h.CacheStatus()
		if err != nil {
			log.Println("[DECODER] Unable to read CacheStatus(), dropping message: ", err)
			continue
		}
		dmsg.CacheStatus = cache_status

		method, err := h.Method()
		if err != nil {
			log.Println("[DECODER] Unable to read Method(), dropping message: ", err)
			continue
		}
		dmsg.Method = method

		url, err := h.Url()
		if err != nil {
			log.Println("[DECODER] Unable to read Url(), dropping message:", err)
			continue
		}
		dmsg.Url = url

		bytes_sent := h.BytesSent()
		if bytes_sent == 0 {
			log.Println("[DECODER] Unable to read BytesSent(), dropping message: ", err)
			continue
		}
		dmsg.BytesSent = bytes_sent

		log.Print("Sending message...")
		cchan <- AppMessage{
			KafkaMessage:   m,
			DecodedMessage: dmsg,
		}
		log.Println("...done")
	}
}
