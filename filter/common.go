package filter

import (
	"auto_rent/fetch_house"
	"errors"
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
		}
	case fetchHouse.CITY58PLAT:
		return City58Filter{
			Price: 			service.Price,
			Type: 			service.Type,
			Orientation: 	service.Orientation,
			Way: 			service.Way,
		}
	}

	return nil, errors.New("Unkown plat type")
}

func SetBIT(val *uint32, bit uint8) error {
	if bit > 31 {
		return errors.New("Out max bit number")
	}

	*val |= 1 << bit
	return nil
}

func UnSetBIT(val *uint32, bit uint8) {
	if 0 == bit || bit >= 32 {
		return errors.New("Out max bit number")
	}

	*val |= 0 << bit
	return nil
}
