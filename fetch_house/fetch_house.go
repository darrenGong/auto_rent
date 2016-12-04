package fetchHouse

import (
	"uframework/log"
)

var (
	chanAreaHouse = make(chan *AreaHouses)
)

func FetchHouse(Config *Config) (map[string][]*House, error) {
	for key, value := range Config.PlatUrl {
		houseInterface, err := GetHouseInterface(key, &value)
		if err != nil {
			uflog.ERRORF("Failed to get house interface [type:%s]", key)
			return nil, err
		}
		go houseInterface.GetHouse(chanAreaHouse)
	}

	houseMaps := make(map[string][]*House)
	for {
		areaHouse, ok := <-chanAreaHouse
		if !ok {
			uflog.WARN("All cities data have got")
			break
		}
		houseMaps[areaHouse.Area] = areaHouse.Houses
	}

	return houseMaps, nil
}
