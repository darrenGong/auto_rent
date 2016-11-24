package main

import (
	"auto_rent/fetch_house"
	"flag"
	"fmt"
	"log"
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
	fmt.Println(Config)
	fetchHouse.FetchHouse(&Config)
}
