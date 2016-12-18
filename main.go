package main

import (
	"auto_rent/fetch_house"
	"flag"
	"log"
	"auto_rent/filter"
	"os"

	logger "github.com/Sirupsen/logrus"
)

var (
	configPath = flag.String("c", "F:\\go-dev\\src\\auto_rent\\config\\config.json", "Configuration, json format")
)

func Init() {
    file, err := os.OpenFile("log/log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
    if err != nil {
        logger.WithField("LOGFILE", "log").Fatalf("Failed to open[%s]", "log")
    }

    logger.SetLevel(logger.DebugLevel)
    logger.SetOutput(file)
    logger.SetFormatter(&logger.JSONFormatter{})
    //logger.SetFormatter(&logger.TextFormatter{})
}

func main() {
	flag.Parse()
	Init()

	Config := fetchHouse.Config{}
	if err := fetchHouse.ParseConfig(*configPath, &Config); err != nil {
		log.Fatalf("Failed to parse config[%s]", *configPath)
	}

	Filter := filter.Filter{
		IdServices: make(map[string]*filter.Service),
		CityServices: make(map[string]map[string]*filter.Service),
	}
	Filter.Run(&Config)
}
