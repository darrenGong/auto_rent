package filter

import "auto_rent/fetch_house"

type City58Filter struct {
	Price       string
	Type        string
	Orientation string
	Way         string
}

func (city58 City58Filter) ValidPrice(house *fetchHouse.House) bool {
	return true
}

func (city58 City58Filter) ValidType(house *fetchHouse.House) bool {
	return true
}

func (city58 City58Filter) ValidOrientation(house *fetchHouse.House) bool {
	return true
}

func (city58 City58Filter)	ValidWay(house *fetchHouse.House) bool {
	return true
}
