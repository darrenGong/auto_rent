package fetchHouse

import (
	"uframework/log"
	"fmt"
)

var (
	TotalMaxNum = 1024 * 10

	chanHouse = make(chan []*House, TotalMaxNum)
)

func FetchHouse(Config *Config) ([]*House, error) {
	for key, value := range Config.PlatUrl {
		houseInterface, err := GetHouseInterface(key, &value)
		if err != nil {
			uflog.ERRORF("Failed to get house interface [type:%s]", key)
			return nil, err
		}
		go houseInterface.GetHouse(chanHouse)
	}

	houseArray := make([]*House, TotalMaxNum)
	for {
		houses, ok := <-chanHouse
		if !ok {
			uflog.WARN("All cities data have got")
			break
		}
		for _, house := range houses {
			houseArray = append(houseArray, house)
		}
	}

	fmt.Println("Done ... ")
	for _, house := range houseArray {
		fmt.Println(house)
	}
	return houseArray, nil
}