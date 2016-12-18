package fetchHouse

import (
	"errors"
	"strings"
	"time"
	"uframework/log"
	//	"time"

	logger "github.com/Sirupsen/logrus"
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
	logger.WithField("fetch", "house").Debug(*cityMap)
/*
	cityMap := map[string]string{
		"sh": "http://sh.ganji.com/",
	}
*/
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
		City: GetUrlCity(url),
		AreaHouses: make(map[string][]*House),
	}

	if "" == areaHouse.City {
		logger.WithField("fetch", "house").Errorf("Failed to get url area[url:%s]", url)
		chanAreaHouse <- nil
		return
	}
	fetchNum := GetFetchNum(areaHouse.City, MaxNum)

	goHeavyMap := make(map[string]bool)
	childNodes := nodes.Find("div")
	childNodes = childNodes.Find(".f-list-item dl")
	for i := 0; i < childNodes.Length() && i < fetchNum; i++ {
		house := new(House)
		house.Init()
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
		house.Location = aNodes.Eq(0).Text()
		if aNodes.Eq(2).Text() != "" {
			house.Location += "|" + aNodes.Eq(2).Text()
		}
		locationList := strings.Split(aNodes.Eq(0).Attr("href"), "/")

		ddNodes = divNodes.Find(".info")
		priceNodes := ddNodes.Find(".price")
		house.Price = priceNodes.Text()
		hourNodes := ddNodes.Find(".time")
		house.DataTime = hourNodes.Text()

		if _, ok := goHeavyMap[house.Id]; !ok {
			if len(locationList) >= 3 {
				areaHouse.AreaHouses[locationList[2]] = append(areaHouse.AreaHouses[locationList[2]], house)
			} else {
				areaHouse.AreaHouses[""] = append(areaHouse.AreaHouses[""], house)
			}
			goHeavyMap[house.Id] = true
		} else {
			fetchNum += 1
		}
	}

	chanAreaHouse <- areaHouse
	uflog.DEBUGF("End url:%s\n", url)
}
