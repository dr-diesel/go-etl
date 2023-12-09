package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/segmentio/kafka-go"
)

type loggingTransport struct{}

func (s *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	rbytes, _ := httputil.DumpRequestOut(r, true)
	fmt.Printf("%s\n", rbytes)

	resp, err := http.DefaultTransport.RoundTrip(r)
	// err is returned after dumping the response

	respBytes, _ := httputil.DumpResponse(resp, true)
	//bytes = append(bytes, respBytes...)
	fmt.Printf("%s\n", respBytes)

	return resp, err
}

func processMessageHttpWorker(cchan <-chan AppMessage, k2chan chan<- kafka.Message) {
	log.Println("Running worker processMessageHttp...")

	var totalSent int

	ticker := time.NewTicker(9 * time.Second)

	buffer := make([]AppMessage, 0)

	commitStash := make(kafkaCommitStash, 0)

	/*
		chansend := make(chan []AppMessage)
		go httpSender(chansend)
		defer close(chansend)
	*/

	// Create a writer that caches compressors.
	// For this operation type we supply a nil Reader.
	var zencoder, _ = zstd.NewWriter(nil)

	for {
		select {
		case m, ok := <-cchan:
			if ok {
				log.Println("Adding to batch...", m.DecodedMessage)

				buffer = append(buffer, m)

				commit := kafkaCommit{topic: "http_log", partition: m.KafkaMessage.Partition, offset: m.KafkaMessage.Offset}
				commitStash.merge(commit)

			} else {
				//closed input channel - breaking - TODO wait for 1min or close uncommitted kafka messages
				return
			}
		case <-ticker.C:
			{
				r := len(buffer)
				fmt.Println("CLICKHOUSE FLUSHING TO HTTP...")
				if r > 0 {

					respCode := func() int {
						clickHouseURL := "http://127.0.0.1:8124/" // URL for ClickHouse server
						tableName := "http_log"                   // Replace with your table name

						// Build the URL for the ClickHouse INSERT query
						p := url.Values{}
						p.Add("query", fmt.Sprintf("INSERT INTO %s (timestamp, resource_id, bytes_sent, request_time_milli, response_status, cache_status, method, remote_addr, url) FORMAT TSVRaw", tableName))

						url := clickHouseURL + "?" + p.Encode()

						// Compress the TSV data with Zstd
						tsvData := getTsvData(buffer)
						compressedData := zencoder.EncodeAll(tsvData, make([]byte, 0, len(tsvData)))

						// Create a POST request with the compressed data
						req, err := http.NewRequest("POST", url, bytes.NewReader(compressedData))
						if err != nil {
							fmt.Println("Error creating request:", err)
							return -1
						}
						req.Header.Set("Content-Encoding", "zstd")
						req.Header.Set("Content-Type", "text/tab-separated-values")

						// Send the request
						client := &http.Client{
							//Transport: &loggingTransport{},
						}

						resp, err := client.Do(req)
						if err != nil {
							fmt.Println("Error sending request:", err)
							return -2
						}
						defer resp.Body.Close()

						return resp.StatusCode
					}()

					if respCode != 200 {
						fmt.Println("INVALID RESPONSE CODE: ", respCode)
						continue
					}

					// buffer reset
					buffer = nil
					totalSent += r
					fmt.Println("---\nCLICKHOUSE sent to API ", r, " item(s) Total: ", totalSent, "\n---")

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

/*
func httpSender(c <-chan []AppMessage) {
	for buffer := range c {
		for _, m := range buffer {
			fmt.Println(m.DecodedMessage)
		}
	}
}
*/

func getTsvData(buffer []AppMessage) []byte {

	var lines []string
	for _, m := range buffer {

		dmsg := m.DecodedMessage

		//(timestamp, resource_id, bytes_sent, request_time_milli, response_status, cache_status, method, remote_addr, url)
		sfmt := make([]string, 0)
		args := make([]interface{}, 0)

		sfmt = append(sfmt, "%d")
		args = append(args, dmsg.TimeStampMilli)

		sfmt = append(sfmt, "%d")
		args = append(args, dmsg.ResourceId)

		sfmt = append(sfmt, "%d")
		args = append(args, dmsg.BytesSent)

		sfmt = append(sfmt, "%d")
		args = append(args, dmsg.RequestTimeMilli)

		sfmt = append(sfmt, "%d")
		args = append(args, dmsg.ResponseStatus)

		sfmt = append(sfmt, "%s")
		args = append(args, dmsg.CacheStatus)

		sfmt = append(sfmt, "%s")
		args = append(args, dmsg.Method)

		sfmt = append(sfmt, "%s")
		args = append(args, zeroLastOctet(dmsg.RemoteAddr).String())

		sfmt = append(sfmt, "%s")
		args = append(args, dmsg.Url)

		line := fmt.Sprintf(strings.Join(sfmt, "\t"), args...)
		lines = append(lines, line)
	}
	tsvData := strings.Join(lines, "\n")

	return []byte(tsvData)
}
