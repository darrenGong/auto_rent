package main

import (
	"auto_rent/fetch_house"
	"flag"
	"log"
	"auto_rent/filter"
)

var (
	configPath = flag.String("c", "F:\\go-dev\\src\\auto_rent\\config\\config.json", "Configuration, json format")
)

func main() {
	flag.Parse()

	Config := fetchHouse.Config{}
	if err := fetchHouse.ParseConfig(*configPath, &Config); err != nil {
		log.Fatalf("Failed to parse config[%s]", *configPath)
	}

	Filter := filter.Filter{
		IdServices: make(map[string]*filter.Service),
		AreaServices: make(map[string]map[string]*filter.Service),
	}
	Filter.Run(&Config)
}
