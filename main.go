package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/p2pquake/dmdata-jp-api-v2-websocket-client/dmdata"
)

type Message struct {
	Type string `json:"type"`
}

type Ping struct {
	Type   string `json:"type"`
	PingId string `json:"pingId"`
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Println("DMDATA.JP API v2 WebSocket client (unofficial)")

	apiKey := os.Getenv("DMDATA_JP_API_KEY")
	outputDir := "./xml"

	if apiKey == "" {
		log.Fatalln("DMDATA_JP_API_KEY is not defined")
	}

	c := dmdata.V2Client{ApiKey: apiKey}

	sr, err := c.StartSocket([]string{"telegram.earthquake"}, nil, "P2PQuakeWSV2Client 0.1")
	if err != nil {
		log.Fatalln("StartSocket error:", err)
	}

	dmdata.StartWebSocket(context.TODO(), sr.WebSocket.URL, func(head dmdata.Head, xml string) {
		fileName := time.Now().Format("20060102150405") + "_" + fmt.Sprintf("%08d", rand.Intn(10000000)) + "_" + head.Type + ".xml"
		err = ioutil.WriteFile(outputDir+"/"+fileName+".tmp", []byte(xml), 0644)
		if err != nil {
			log.Fatalln("Callback WriteFile error:", err)
		}
		err = os.Rename(outputDir+"/"+fileName+".tmp", outputDir+"/"+fileName)
		if err != nil {
			log.Fatalln("Callback Rename error:", err)
		}
	})
}
