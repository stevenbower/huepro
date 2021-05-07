package main

import (
	"fmt"
	hue "github.com/collinux/gohue"
	log "github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"regexp"
)

type hueCollector struct {
	bridge     *hue.Bridge
	counts     *prometheus.GaugeVec
	brightness *prometheus.GaugeVec
}

func NewHueCollector(namespace string, bridge *hue.Bridge) prometheus.Collector {
	c := hueCollector{
		bridge: bridge,
		counts: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "hue_lights",
				Name:      "count",
				Help:      "Count of Hue lights",
			},
			[]string{
				"state",
				"reachable",
			},
		),
		brightness: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "hue_lights",
				Name:      "brightness",
				Help:      "Brightness by light",
			},
			[]string{
				"name",
				//"room",
			},
		),
	}

	return c
}

func (c hueCollector) Describe(ch chan<- *prometheus.Desc) {
	c.counts.Describe(ch)
	c.brightness.Describe(ch)
}

func (c hueCollector) Collect(ch chan<- prometheus.Metric) {
	c.counts.Reset()

	lights, err := c.bridge.GetAllLights()
	if err != nil {
		log.Errorf("Failed to update lights: %v", err)
		return
	}

	//groups, err := c.bridge.GetGroups()
	//if err != nil {
	//log.Errorf("Failed to update groups: %v", err)
	//return
	//}
	//fmt.Printf("Groups: %+v\n", groups)
	nameRe := regexp.MustCompile("[^a-zA-Z0-9_]")

	c.counts.With(prometheus.Labels{"state": "on", "reachable": "yes"}).Set(0)
	c.counts.With(prometheus.Labels{"state": "off", "reachable": "yes"}).Set(0)
	c.counts.With(prometheus.Labels{"state": "off", "reachable": "no"}).Set(0)
	for _, light := range lights {
		name := nameRe.ReplaceAllString(light.Name, "_")

		st := "off"
		if light.State.On {
			st = "on"
		}

		rc := "no"
		if light.State.Reachable {
			rc = "yes"
		}

		c.counts.With(prometheus.Labels{"state": st, "reachable": rc}).Inc()

		if !light.State.Reachable {
			c.brightness.With(prometheus.Labels{"name": name}).Set(-1.0)
		} else if light.State.On {
			c.brightness.With(prometheus.Labels{"name": name}).Set(float64(light.State.Bri))
		} else {
			c.brightness.With(prometheus.Labels{"name": name}).Set(0.0)
		}

		fmt.Printf("Light: %s %s\n", name, light.State.Reachable)
	}
	c.counts.Collect(ch)
	c.brightness.Collect(ch)
}
