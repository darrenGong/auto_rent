package filter

import (
	"auto_rent/fetch_house"
	"strings"
	"strconv"

	logger "github.com/Sirupsen/logrus"
)

var (
	houseTypeMap = map[uint32]string{
		1: "h1",
		2: "h2",
		3: "h3",
		4: "h4",
		5: "h5",
		6: "h6",
	}
	MAXTYPE = uint32(6)
)

type GJFilter struct {
	Price       string	// 代表价格范围
	Type        string
	Orientation string
	Way         string
}

func (gj GJFilter)GetPriceType(price uint32) string {
	switch {
	case 800 <= price && price <= 1500:
		return "p2"
	case 1500 <= price && price <= 2000:
		return "p3"
	case 2000 <= price && price <= 3000:
		return "p4"
	case 3000 <= price && price <= 5000:
		return "p5"
	case 5000 <= price && price <= 6500:
		return "p6"
	case 6500 <= price && price <= 8000:
		return "p7"
	case 8000 <= price && price <= 1000000000:
		return "p8"
	}

	return "p1"    //800以下
}

func (gj GJFilter) ValidPrice(house *fetchHouse.House) bool {
	strHousePrice := strings.TrimRight(house.Price, "元/月")
	housePrice64, err := strconv.ParseUint(strHousePrice, 10, 32)
	if err != nil {
		logger.WithField("Filter", "House").Errorf("Failed to parse price[Id:%s, err:%s]",
			house.Id, err.Error())
		return false
	}
	housePrice := uint32(housePrice64)
	priceType := gj.GetPriceType(housePrice)
	if priceType == gj.Price {
		logger.WithField("Filter", "House").Infof("Successful match to price[sourcePrice:%s, targetPrice:%s]",
			house.Price, gj.Price)
		return true
	}

	return false
}

func (gj GJFilter) GetHouseType(uPrefixType uint32) string {

	houseType, ok := houseTypeMap[uPrefixType]
	if !ok {
		logger.WithField("Filter", "House").Errorf("House type out of index[%d]", uPrefixType)
		return houseTypeMap[MAXTYPE]
	}

	return houseType
}

func (gj GJFilter) ValidType(house *fetchHouse.House) bool {
	if "" == house.HouseType {
		logger.WithField("Filter", "House").Error("HouseType is equal at empty")
		return false
	}
	preFixHouses := strings.Split(house.HouseType, "室")
	if len(preFixHouses) <= 1 {
		logger.WithField("Filter", "House").Errorf("Invalid house type[%s]", house.HouseType)
		return false
	}

	uPrefixType64, err := strconv.ParseUint(preFixHouses[0], 10, 32)
	if err != nil {
		logger.WithField("Filter", "House").Errorf("Failed to parse prefix house type[%s]", preFixHouses[0])
		return false
	}
	uPrefixType := uint32(uPrefixType64)
	houseType := gj.GetHouseType(uPrefixType)
	if houseType == gj.Type {
		logger.WithField("Filter", "House").Infof("Successful match to house type[sourceType:%s, targetType:%s]",
			house.HouseType, gj.Type)
		return true
	}

	return false
}

func (gj GJFilter) GetOrientation(orientation string) string {
	switch orientation {
	case "南向":    //south
		return "j2"
	case "西向": //west
		return "j3"
	case "北向":    //north
		return "j4"
	case "南北向":    //sn
		return "j5"
	case "东西向":    //ew
		return "j6"
	}
	return "j1"    //1 east
}

func (gj GJFilter) ValidOrientation(house *fetchHouse.House) bool {
	houseOrientation := gj.GetOrientation(house.Orientation)
	if houseOrientation == gj.Orientation {
		logger.WithField("Filter", "House").Infof("Successful match to orientation[sourceOri:%s, targetOri:%s]",
			house.Orientation, gj.Orientation)
		return true
	}
	return false
}

func (gj GJFilter) GetRentWay(way string) string {
	switch way {
	case "合租":    //合租
		return "a3"
	}
	return "m1"    //1 整租
}

func (gj GJFilter)	ValidWay(house *fetchHouse.House) bool {
	rentWay := gj.GetRentWay(house.Way)
	if rentWay == gj.Way {
		logger.WithField("Filter", "House").Infof("Successful match to rent way[sourceWay:%s, targetWay:%s]",
			house.Way, gj.Way)
		return true
	}

	return false
}