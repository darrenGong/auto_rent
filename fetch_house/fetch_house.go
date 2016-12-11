package fetchHouse

import (
	"uframework/log"
)

func FetchHouse(Config *Config) (map[string]map[string][]*House, error) {
	chanAreaHouse := make(chan *AreaHouses)
	for key, value := range Config.PlatUrl {
		houseInterface, err := GetHouseInterface(key, &value)
		if err != nil {
			uflog.ERRORF("Failed to get house interface [type:%s]", key)
			return nil, err
		}
		go houseInterface.GetHouse(chanAreaHouse)
	}

	houseMaps := make(map[string]map[string][]*House)
	for {
		areaHouse, ok := <-chanAreaHouse
		if !ok {
			uflog.WARN("All cities data have got")
			break
		}
		houseMaps[areaHouse.City] = areaHouse.AreaHouses
	}

	return houseMaps, nil
}
