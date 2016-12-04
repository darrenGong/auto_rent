package fetchHouse

import (
	"errors"
	"strings"
	"time"
	"uframework/log"
	//	"time"
	"fmt"
)

var (
	TypeUrl = "fang1/"
	MaxNum  = 20
)

type GJHouse struct {
	Url     string
	AreaUrl string
}

func (gj GJHouse) GetHouse(chanAreaHouse chan<- *AreaHouses) error {
	cityMap, err := GetAllSite(gj.Url)
	if err != nil {
		uflog.ERRORF("Failed to get all city[url:%s]", gj.Url)
		return errors.New("Failed to get all city")
	}

	fmt.Println(*cityMap)
	chanHouse := make(chan *AreaHouses)
	for _, url := range *cityMap {
		go gj.GetAreaHouse(url+TypeUrl, chanHouse)
	}

	for _, url := range *cityMap {
		select {
		case areaHouse := <-chanHouse:
			if areaHouse != nil {
				chanAreaHouse <- areaHouse
			}
		case <-time.After(30 * time.Second):
			uflog.ERRORF("Get url:%s area timeout", url)
		}
	}

	close(chanAreaHouse)
	return nil
}

func (gj GJHouse) GetAreaHouse(url string, chanAreaHouse chan<- *AreaHouses) {
	uflog.DEBUGF("Start url: %s\n", url)
	nodes, err := ApiGet(url)
	if err != nil {
		uflog.ERRORF("Failed to get area house[url:%s]", url)
		chanAreaHouse <- nil
		return
	}

	areaHouse := &AreaHouses{
		Area: GetUrlArea(url),
		Houses: make([]*House, 0),
	}

	childNodes := nodes.Find("div")
	childNodes = childNodes.Find(".f-list-item dl")
	for i := 0; i < childNodes.Length() && i < MaxNum; i++ {
		house := new(House)
		house.PlatType = GJPLAT

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
		locationList := strings.Split(aNodes.Eq(0).Attr("href"), "/")
		if len(locationList) >= 3 {
			house.Location = locationList[2] + "|"
		}
		communityList := strings.Split(aNodes.Eq(2).Attr("href"), "/")
		if len(communityList) >= 3 {
			house.Location += communityList[2]
		}

		ddNodes = divNodes.Find(".info")
		priceNodes := ddNodes.Find(".price")
		house.Price = priceNodes.Text()
		hourNodes := ddNodes.Find(".time")
		house.DataTime = hourNodes.Text()

		areaHouse.Houses = append(areaHouse.Houses, house)
	}

	chanAreaHouse <- areaHouse
	uflog.DEBUGF("End url:%s\n", url)
}
