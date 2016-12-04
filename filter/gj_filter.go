package filter

import (
	"auto_rent/fetch_house"
	"strings"
	"strconv"
	"uframework/log"
)

type GJFilter struct {
	Price       uint32
	Type        uint32
	Orientation uint32
	Way         uint32
}

func (gj *GJFilter) ValidPrice(house *fetchHouse.House, price uint32) bool {
	strHousePrice := strings.TrimRight(house.Price, "元/月")
	housePrice, err := strconv.ParseUint(strHousePrice, 10, 32)
	if err != nil {
		uflog.ERRORF("Failed to parse price[Id:%s, err:%s]", house.Id, err.Error())
		return false
	}
	housePrice = uint32(housePrice)
	return true
}

func (gj *GJFilter) ValidType(house *fetchHouse.House, houseType uint32) bool {
	return true
}

func (gj *GJFilter) ValidOrientation(house *fetchHouse.House, orientation uint32) bool {
	return true
}

func (gj *GJFilter)	ValidWay(house *fetchHouse.House, way uint32) bool {
	return true
}

