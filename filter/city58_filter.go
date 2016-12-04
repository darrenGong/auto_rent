package filter

import "auto_rent/fetch_house"

type City58Filter struct {
	Price       uint32
	Type        uint32
	Orientation uint32
	Way         uint32
}

func (city58 *City58Filter) ValidPrice(house *fetchHouse.House, price uint32) bool {
	return true
}

func (city58 *City58Filter) ValidType(house *fetchHouse.House, houseType uint32) bool {
	return true
}

func (city58 *City58Filter) ValidOrientation(house *fetchHouse.House, orientation uint32) bool {
	return true
}

func (city58 *City58Filter)	ValidWay(house *fetchHouse.House, way uint32) bool {
	return true
}
