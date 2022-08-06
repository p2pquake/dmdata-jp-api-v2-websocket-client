package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/p2pquake/dmdata-jp-api-v2-websocket-client/dmdata"
	"github.com/spf13/cobra"
)

// ビルド時に設定
var (
	Version = "develop"
	Commit  = "unknown"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "dmdata-jp-api-v2-websocket-client",
	Short:   "DMDATA.JP (Project DM-D.S.S) API v2 の非公式クライアント実装",
	Example: "DMDATA_JP_API_KEY=<API_KEY> ./dmdata-jp-api-v2-websocket-client -k -c telegram.earthquake -c eew.warning",
	Version: fmt.Sprintf("%s (commit %s, built at %s)", Version, Commit, Date),
	Run:     run,
}

var keepConnection bool
var classifications []string

func Execute() error {
	rootCmd.Flags().BoolVarP(&keepConnection, "keep-existing-connections", "k", false, "Keep existing connections")
	rootCmd.Flags().StringSliceVarP(&classifications, "classifications", "c", []string{"telegram.earthquake"}, "Retreive classifications")

	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Println("DMDATA.JP API v2 WebSocket client (unofficial)")

	apiKey := os.Getenv("DMDATA_JP_API_KEY")
	outputDir := "./xml"

	if apiKey == "" {
		log.Fatalln("DMDATA_JP_API_KEY is not defined")
	}

	c := dmdata.V2Client{ApiKey: apiKey}

	if !keepConnection {
		lr, err := c.ListSocket("open")
		if err != nil {
			log.Fatalln("ListSocket error:", err)
		}

		for _, v := range lr.Items {
			err = c.CloseSocket(v.Id)
			if err != nil {
				log.Fatalln("CloseSocket error:", err)
			}
		}
	}

	sr, err := c.StartSocket(classifications, nil, "P2PQuakeWSV2Client 0.2")
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
