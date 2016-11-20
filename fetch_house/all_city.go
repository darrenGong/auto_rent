package fetchHouse

import (
	"uframework/log"
	"github.com/opesun/goquery"
)

func GetAllSite(url string) (*map[string]string, error) {
	nodes, err := ApiGet(url)
	if err != nil {
		uflog.ERRORF("Failed to get url[%s]", url)
		return nil, err
	}

	return ParseNodes(nodes)
}

func ParseNodes(nodes goquery.Nodes) (*map[string]string, error) {
	childNodes := nodes.Find(".all-city dl")
	childNodes = childNodes.Find("a")

	cityMap := make(map[string]string)
	for i := 0; i < childNodes.Length(); i++ {
		cityMap[childNodes.Eq(i).Text()] = childNodes.Eq(i).Attr("href")
	}

	return &cityMap, nil
}
