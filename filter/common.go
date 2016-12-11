package filter

import (
	"auto_rent/fetch_house"
	"errors"
	"encoding/json"
	"io/ioutil"

	logger "github.com/Sirupsen/logrus"
)

type HouseFilter interface {
	ValidPrice(house *fetchHouse.House) bool
	ValidType(house *fetchHouse.House) bool
	ValidOrientation(house *fetchHouse.House) bool
	ValidWay(house *fetchHouse.House) bool
}

func GetPlatInterface(platType string, service *Service) (HouseFilter, error) {
	switch platType {
	case fetchHouse.GJPLAT:
		return GJFilter{
			Price:    		service.Price,
			Type:    		service.Type,
			Orientation: 	service.Orientation,
			Way: 			service.Way,
		}, nil
	case fetchHouse.CITY58PLAT:
		return City58Filter{
			Price: 			service.Price,
			Type: 			service.Type,
			Orientation: 	service.Orientation,
			Way: 			service.Way,
		}, nil
	}

	return nil, errors.New("Unkown plat type")
}

func SetBIT(val *uint32, bit uint8) error {
	if 0 == bit || bit > 32 {
		return errors.New("Out max bit number")
	}

	*val |= 1 << (bit - 1)
	return nil
}

func UnSetBIT(val *uint32, bit uint8) error {
	if 0 == bit || bit >= 32 {
		return errors.New("Out max bit number")
	}

	*val |= 0 << bit
	return nil
}

func UnmarshalFile(file string) (*Service, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logger.WithField("service", "parse").Errorf("Failed to read file[%s]", file)
		return nil, err
	}
	logger.WithField("service", "parse").Infof("Read %d bytes from %s", len(bytes), file)

	var service Service
	if err := json.Unmarshal(bytes, &service); err != nil {
		logger.WithField("service", "parse").Errorf("Failed to unmarshal file[%s]", file)
		return nil, err
	}

	return &service, nil
}
