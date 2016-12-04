package filter

import (
	"auto_rent/fetch_house"
	"errors"
)

type HouseFilter interface {
	ValidPrice(house *fetchHouse.House, price uint32) bool
	ValidType(house *fetchHouse.House, houseType uint32) bool
	ValidOrientation(house *fetchHouse.House, orientation uint32) bool
	ValidWay(house *fetchHouse.House, way uint32) bool
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
