package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jasonlvhit/gocron"
	"github.com/showwin/speedtest-go/speedtest"
)

func main() {
	gocron.Every(10).Minutes().Do(task)
	<-gocron.Start()
	// task()
}

func task() {
	dt := time.Now()
	fmt.Println("Current date and time is: ", dt.String())
	// create client
	client := influxdb2.NewClient("http://localhost:8086", "F-QFQpmCL9UkR3qyoXnLkzWj03s6m4eCvYgDl1ePfHBf9ph7yxaSgQ6WN0i9giNgRTfONwVMK1f977r_g71oNQ==")

	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking("3e645a0874fc54c7", "speedtest")
	var speedtestClient = speedtest.New()

	// Use a proxy for the speedtest. eg: socks://127.0.0.1:7890
	// speedtest.WithUserConfig(&speedtest.UserConfig{Proxy: "socks://127.0.0.1:7890"})(speedtestClient)

	// Select a network card as the data interface.
	// speedtest.WithUserConfig(&speedtest.UserConfig{Source: "192.168.1.101"})(speedtestClient)

	user, _ := speedtestClient.FetchUserInfo()
	// Get a list of servers near a specified location
	// user.SetLocationByCity("Detroit")
	// user.SetLocationByCity("Tokyo")
	// user.SetLocation("Osaka", 34.6952, 135.5006)

	serverList, _ := speedtestClient.FetchServers(user)
	// targets, _ := serverList.FindServer([]int{})
	targets, _ := serverList.FindServer([]int{17755})

	for _, s := range targets {
		// Please make sure your host can access this test server,
		// otherwise you will get an error.
		// It is recommended to replace a server at this time
		s.PingTest()
		s.DownloadTest()
		s.UploadTest()

		latency := int64(s.Latency / time.Millisecond)
		fmt.Println(latency)
		fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)
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
