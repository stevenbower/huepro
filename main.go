package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	hue "github.com/collinux/gohue"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
var config = flag.String("c", "huepro.conf", "The config file to use")

/*
func register() {
	locators, _ := hue.DiscoverBridges(false)
	locator := locators[0] // find the first locator
	deviceType := "huepro"

	// remember to push the button on your hue first
	bridge, _ := locator.CreateUser(deviceType)
	fmt.Printf("registered new device => %+v\n", bridge)
}
*/

type Config struct {
	IpAddr string
	Token  string
}

func main() {
	flag.Parse()

	raw, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}

	var cfg Config
	json.Unmarshal(raw, &cfg)

	bridge, err := hue.NewBridge(cfg.IpAddr)
	if err != nil {
		panic(err)
	}

	err = bridge.Login(cfg.Token)
	if err != nil {
		panic(err)
	}

	prometheus.MustRegister(NewHueCollector("", bridge))

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
