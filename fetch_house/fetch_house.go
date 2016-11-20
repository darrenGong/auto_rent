package fetchHouse

var (
	areaUrl = "%s/"
	WayUrl = "%s"
	PriceUrl = "%s"			// combine	example: m1p1/
	HouseTypeUrl = "%s"		// combine HouseTypeUrl WayUrl PriceUrl/
	OrientationUrl = "%s"	// combine HouseTypeUrl OrientationUrl WayUrl PriceUrl/
)

func (filter *Filter) FetchHouse() ([]*House, error) {
	url := GetUrl(filter.PlatType)
	return nil, nil
}