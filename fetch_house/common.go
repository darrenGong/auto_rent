package fetchHouse

import (
	"time"
	"github.com/opesun/goquery"
	"uframework/log"
	"auto_rent/http_request"
	"errors"
	"strings"
)

var (
	TIMEOUT = 10 * time.Second

	STATUS = 200
)

const (
	ALLRENT = iota        //整租
	TOGETHERRENT        //合租
	SHORTRENT            //短租
)

const (
	GJPLAT = "GJ"            //赶集
	CITY58PLAT = "58CITY"        //58同城
)

type Price struct {
	MinPrice float32
	MaxPrice float32
}

type Filter struct {
	PlatType    uint32
	City        string
	Area        string
	Price       Price
	HouseType   uint32
	Orientation uint32
	Way         uint32
}

type House struct {
	Id          string
	DataTime    string
	HasImage    bool
	Name        string
	Price       string
	Url         string
	Location    string
	HouseType   string // 三室二厅二卫
	Orientation string // 朝向
	Way         string // 整租

	PlatType    string
}

type AreaHouses struct {
	Area   string
	Houses []*House
}

type HouseInterface interface {
	GetHouse(chanAreaHouse chan <- *AreaHouses) error
}

func GetHouseInterface(platType string, webUrl *WebUrl) (HouseInterface, error) {
	switch platType {
	case GJPLAT:
		return GJHouse{Url: webUrl.Url,
			AreaUrl: webUrl.AreaUrl}, nil
	case CITY58PLAT:
		return CITY58House{Url: webUrl.Url,
			AreaUrl: webUrl.AreaUrl}, nil
	}

	return nil, errors.New("Invalid plat type: " + platType)
}

func ApiGet(url string) (goquery.Nodes, error) {
	httpReq := httpRequest.HttpRequest{
		Timeout: TIMEOUT,
		Url:     url,
	}

	httpRes, err := httpReq.ApiGet(nil)
	if err != nil {
		uflog.ERRORF("Failed to fetch all city[url:%s, err:%v]", httpReq.Url, err)
		return nil, err
	}

	if httpRes.StatusCode != STATUS {
		uflog.ERRORF("Status code is error[url:%s, status:%d]", httpReq.Url, httpRes.StatusCode)
		return nil, errors.New("Failed to get")
	}

	nodes, err := goquery.Parse(httpRes.Body)
	if err != nil {
		uflog.ERRORF("Failed to goquery[url:%s, err: %v]", httpReq.Url, err)
		return nil, err
	}

	return nodes, nil
}

func GetValidUrl(url string) string {
	if !strings.HasPrefix(url, "http://") {
		return "http://" + url
	}
	return url
}

func GetUrlArea(url string) string {
	url = GetValidUrl(url)
	urls := strings.Split(url, ".")
	if len(urls) != 3 {
		uflog.ERRORF("Invalid url: %s", url)
		return ""
	}

	preUrls := strings.Split(urls[0], "//")
	if len(preUrls) != 2 {
		uflog.ERRORF("Invalid url: %s", urls[0])
		return ""
	}

	return preUrls[1]
}