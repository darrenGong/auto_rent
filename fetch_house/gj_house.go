package fetchHouse

import (
	"uframework/log"
	"errors"
	"strings"
	"fmt"
	//	"time"
)

var (
	TypeUrl = "fang1/"
	MaxNum = 20
)

type GJHouse struct {
	Url     string
	AreaUrl string
}

func (gj GJHouse) GetHouse(chanHouse chan <- []*House) error {
	cityMap, err := GetAllSite(gj.Url)
	if err != nil {
		uflog.ERRORF("Failed to get all city[url:%s]", gj.Url)
		return errors.New("Failed to get all city")
	}

	chanHouses := make(chan []*House)
	chanUrl := make(chan string)
	go gj.RoutineAreaHouse(chanUrl, chanHouses)
	for _, url := range *cityMap {
		chanUrl <- url
		houses, ok := <-chanHouses
		if ok && houses != nil {
			chanHouse <- houses
		}
	}

	close(chanHouse)
	return nil
}

func (gj GJHouse) RoutineAreaHouse(chanUrl <-chan string, chanHouses chan <- []*House) {
	for {
		url := <-chanUrl

		gj.GetAreaHouse(url, chanHouses)
	}
}

func (gj GJHouse) GetAreaHouse(url string, chanHouses chan <- []*House) {
	fmt.Printf("Start url: %s\n", url)
	nodes, err := ApiGet(url)
	if err != nil {
		uflog.ERRORF("Failed to get area house[url:%s]", url)
		chanHouses <- nil
		return
	}
	childNodes := nodes.Find("div")
	childNodes = childNodes.Find(".f-list-item dl")
	houses := make([]*House, MaxNum)
	for i := 0; i < childNodes.Length() && i < MaxNum; i++ {
		house := new(House)

		divNodes := childNodes.Eq(i)
		dtNodes := divNodes.Find("dt")
		childdivNodes := dtNodes.Find("div")

		aNodes := childdivNodes.Find("a")
		house.Id = aNodes.Attr("href")
		house.Url = url + strings.Split(house.Id, "/")[2]

		imgNodes := childdivNodes.Find("img")
		house.HasImage = imgNodes.Attr("src") != ""

		ddNodes := divNodes.Find(".title")
		aNodes = ddNodes.Find("a")
		house.Name = aNodes.Text()

		ddNodes = divNodes.Find(".size")
		spanNodes := ddNodes.Find("span")
		house.Way = spanNodes.Eq(0).Text()
		house.HouseType = spanNodes.Eq(2).Text()
		house.Orientation = spanNodes.Eq(6).Text()

		ddNodes = divNodes.Find(".address")
		spanNodes = ddNodes.Find("span")
		aNodes = spanNodes.Find("a")
		house.Location = aNodes.Eq(0).Text() + "|" + aNodes.Eq(2).Text()

		ddNodes = divNodes.Find(".info")
		priceNodes := ddNodes.Find(".price")
		house.Price = priceNodes.Text()
		hourNodes := ddNodes.Find(".time")
		house.DataTime = hourNodes.Text()

		houses = append(houses, house)
	}

	fmt.Printf("End Url:%s\n", url)
	chanHouses <- houses
	fmt.Printf("Transfer Url:%s\n", url)
}