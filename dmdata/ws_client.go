package dmdata

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string `json:"type"`
}

type Ping struct {
	Type   string `json:"type"`
	PingId string `json:"pingId"`
}

type Data struct {
	Head        Head   `json:"head"`
	Format      string `json:"format"`
	Compression string `json:"compression"`
	Encoding    string `json:"encoding"`
	Body        string `json:"body"`
}

type Head struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

func StartWebSocket(ctx context.Context, url string, callback func(head Head, xml string)) {
	log.Printf("connecting to %s", url)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, b, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s\n", b)

			// JSON parse
			message := Message{}
			if err := json.Unmarshal(b, &message); err != nil {
				log.Println("parse:", err)
				return
			}

			// Message handling
			if message.Type == "ping" {
				ping := Ping{}
				if err := json.Unmarshal(b, &ping); err != nil {
					log.Println("ping parse:", err)
					return
				}

				pong := Ping{Type: "pong", PingId: ping.PingId}
				log.Printf("pong: %#v", pong)
				if err := c.WriteJSON(pong); err != nil {
					log.Println("pong write:", err)
					return
				}
			}

			if message.Type == "data" {
				data := Data{}
				if err := json.Unmarshal(b, &data); err != nil {
					log.Println("data parse:", err)
					return
				}

				if data.Format != "xml" {
					log.Println("invalid format:", data.Format)
					return
				}

				if data.Compression == "gzip" && data.Encoding == "base64" {
					// BASE64 デコードの後 gzip 展開
					gzippedBytes, err := base64.StdEncoding.DecodeString(data.Body)
					if err != nil {
						log.Println("base64 decode:", err)
						return
					}

					reader, err := gzip.NewReader(bytes.NewBuffer(gzippedBytes))
					if err != nil {
						log.Println("gzip read:", err)
						return
					}
					defer reader.Close()

					gunzippedBytes := bytes.Buffer{}
					if _, err = gunzippedBytes.ReadFrom(reader); err != nil {
						log.Println("gzip read:", err)
						return
					}

					callback(data.Head, string(gunzippedBytes.Bytes()))
				}

				if data.Compression == "" && data.Encoding == "utf-8" {
					callback(data.Head, data.Body)
				}
			}

		}
	}()

	log.Println("connected.")

	func() {
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	log.Println("closed.")
}
