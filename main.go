package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jasonlvhit/gocron"
	"github.com/showwin/speedtest-go/speedtest"
)

func main() {
	// env var for the interval check
	interval := os.Getenv("CHECK_INTERVAL")
	intVar, err := strconv.Atoi(interval)
	if err != nil {
		log.Panic(err)
	}

	gocron.Every(uint64(intVar)).Minutes().Do(task)
	<-gocron.Start()
}

func task() {
	org := os.Getenv("INFLUXDB_ORG")
	bucket := os.Getenv("INFLUXDB_BUCKET")
	database := os.Getenv("INFLUXDB")
	token := os.Getenv("INFLUXDB_TOKEN")
	dt := time.Now()
	fmt.Println("Current date and time is: ", dt.String())
	// create client
	client := influxdb2.NewClient(database, token)

	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)
	var speedtestClient = speedtest.New()
	user, _ := speedtestClient.FetchUserInfo()

	serverList, _ := speedtestClient.FetchServers(user)
	targets, _ := serverList.FindServer([]int{17755}) // using detroit, MI

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest()
		s.UploadTest()

		// convert latency string to int64
		latency := int64(s.Latency / time.Millisecond)

		// print results to console
		fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)

		// post results to influxdb
		p := influxdb2.NewPointWithMeasurement("stat").
			AddTag("unit", "speed").
			AddField("upload", s.ULSpeed).
			AddField("download", s.DLSpeed).
			AddField("latency", latency).
			SetTime(time.Now())
		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			panic(err)
		}
	}
}
