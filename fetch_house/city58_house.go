package fetchHouse

type CITY58House struct {
	Url     string
	AreaUrl string
}

func (city58 CITY58House) GetHouse(chanHouse chan <- []*House) error {
	return nil
}