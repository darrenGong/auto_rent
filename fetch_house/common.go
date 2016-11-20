package fetchHouse

import (
	"time"
	"github.com/opesun/goquery"
	"uframework/log"
	"auto_rent/http_request"
	"errors"
)

var (
	TIMEOUT = 60 * time.Second

	STATUS = 200

	GJURL = "http://www.ganji.com/index.htm"
	GJCITYURL = "http://%s.ganji.com/fang1/"
	CITY58URL = "http://www.58.com/changecity.aspx"
	CITY2URL = "http://%s.58.com/zufang/"
)

const (
	ALLRENT = iota    //整租
	TOGETHERRENT    //合租
	SHORTRENT        //短租
)

type Price struct {
	MinPrice float32
	MaxPrice float32
}

type Filter struct {
	PlatType	uint32
	City        string
	Area        string
	Price       Price
	HouseType   uint32
	Orientation uint32
	Way         uint32
}

type House struct {
	Id          string
	DataTime    time.Duration
	HasImage    bool
	Name        string
	Price       string
	Url         string
	Location    string
	HouseType   string // 三室二厅二卫
	Orientation string // 朝向
	Way         uint32 // 整租
}

func GetUrl(platType uint32) string {
	switch platType {
	case 10001:
		return GJCITYURL
	case 10002:
		return CITY2URL
	}

	return GJCITYURL
}

func GetPriceType(price float32) string {
	price = uint32(price)
	switch price {
	case 800...1500:
		return "p2"
	case 1500...2000:
		return "p3"
	case 2000...3000:
		return "p4"
	case 3000...5000:
		return "p5"
	case 5000...6500:
		return "p6"
	case 6500...8000:
		return "p7"
	case 8000...1000000000:
		return "p8"
	}

	return "p1"	//800以下
}

func GetOrientation(orientation uint32) string {
	switch orientation {
	case 2:	//south
		return "j2"
	case 3: //west
		return "j3"
	case 4:	//north
		return "j4"
	case 5:	//sn
		return "j5"
	case 6:	//ew
		return "j6"
	}
	return "j1"	//1 east
}

func GetRentWay(way uint32) string {
	switch way {
	case 2:	//合租
		return "a3"
	}
	return "m1"	//1 整租
}

func GetHouseType(houseType uint32) string {
	switch houseType {
	case 2:
		return "h2"
	case 3:
		return "h3"
	case 4:
		return "h4"
	case 5:
		return "h5"
	case 6:	//五室以上
		return "h6"
	}
	return "h1"
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