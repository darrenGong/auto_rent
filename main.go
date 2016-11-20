package main

import (

)
import (
	"auto_rent/fetch_house"
	"fmt"
)

func main() {
	cityMap, _ := fetchHouse.GetAllSite("http://www.ganji.com/index.htm")
	fmt.Println(cityMap)
}
