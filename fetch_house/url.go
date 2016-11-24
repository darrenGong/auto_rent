package fetchHouse

import (
	"uframework/log"
	"errors"
)

type UrlInterface interface {
	GeneratorUrl(url string) (string, error)
}

type GJUrl struct {
	AreaUrl        string // %s/
	WayUrl         string // %s
	PriceUrl       string // combine	example: WayUrl PriceUrl/
	HouseTypeUrl   string // combine HouseTypeUrl WayUrl PriceUrl/
	OrientationUrl string // combine HouseTypeUrl OrientationUrl WayUrl PriceUrl/
}

type CITY58Url struct {
	AreaUrl        string // %s/
}

func (gjUrl *GJUrl) GeneratorUrl(url string) (string, error) {
	if "" == gjUrl.AreaUrl {
		return "", errors.New("AreaUrl is empty")
	}
	url += gjUrl.AreaUrl
	uflog.DEBUGF("Url: %s", url)

	url += gjUrl.HouseTypeUrl + gjUrl.OrientationUrl + gjUrl.WayUrl + gjUrl.PriceUrl
	return url, nil
}

func (gjUrl *GJUrl) GetPriceType(priceFloat float64) string {
	price := uint32(priceFloat)
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

func (gjUrl *GJUrl) GetOrientation(orientation uint32) string {
	switch orientation {
	case 2:    //south
		return "j2"
	case 3: //west
		return "j3"
	case 4:    //north
		return "j4"
	case 5:    //sn
		return "j5"
	case 6:    //ew
		return "j6"
	}
	return "j1"    //1 east
}

func (gjUrl *GJUrl) GetRentWay(way uint32) string {
	switch way {
	case 2:    //合租
		return "a3"
	}
	return "m1"    //1 整租
}

func (gjUrl *GJUrl) GetHouseType(houseType uint32) string {
	switch houseType {
	case 2:
		return "h2"
	case 3:
		return "h3"
	case 4:
		return "h4"
	case 5:
		return "h5"
	case 6:    //五室以上
		return "h6"
	}
	return "h1"
}

func (city58 *CITY58Url) GeneratorUrl(url string) (string, error) {
	return "", nil
}