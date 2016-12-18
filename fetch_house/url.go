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





func (city58 *CITY58Url) GeneratorUrl(url string) (string, error) {
	return "", nil
}