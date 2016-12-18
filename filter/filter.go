package filter

import (
	"strings"
	"path"
	"io/ioutil"
	"sync"
	"github.com/fsnotify/fsnotify"
	"auto_rent/fetch_house"
	"fmt"
	"errors"

	logger "github.com/Sirupsen/logrus"
	"time"
)

type SimpleHouse struct {
	Price       string // 价格范围
	Type        string // 3室一厅
	Orientation string // 朝向
	Way         string // 整租
}

type Service struct {
	Id        string
	Username  string
	Email     string
	Number    string
	Platform  uint32
	City      string
	Area      string
	Community string // 社区
	SimpleHouse
}

var (
	FPrefix = "House_"
	FSuffix = ".json"
)

var sm sync.Mutex

type Filter struct {
	IdServices   map[string]*Service            // key: Id
	CityServices map[string]map[string]*Service // key: Area subKey: Id
}

func DebugHouses(areaHousesMap *map[string]map[string][]*fetchHouse.House) {
	logger.WithField("fetch", "house").Debugf("City len:%d", len(*areaHousesMap))
	for city, areaHouses := range *areaHousesMap {
		logger.WithField("fetch", "house").Debugf("Area len:%d", len(areaHouses))
		for area, houses := range areaHouses {
			logger.WithField("fetch", "house").Debugf("city:%s, area:%s, House len:%d, houses:%v",
				city, area, len(houses), houses)

			}
	}
}

func DebugMatchData(service *Service, house *fetchHouse.House) {
	logger.WithField("filter", "house").Debugf(`Start matching service[
	Id: %s
	Username: %s
	Email: %s
	Number: %s
	Platform: %d
	City: %s
	Area: %s
	Community: %s
	Price: %s
	Type: %s
	Orientation: %s 
	Way: %s] and house[
	Id: %s          
	DataTime: %s    
	HasImage: %v
	Name: %s        
	Price: %s       
	Url: %s         
	Location: %s    
	HouseType: %s
	Orientation: %s
	Way: %s
	PlatType: %s]`,
	service.Id, service.Username, service.Email, service.Number, service.Platform, service.City, service.Area, service.Community, service.Price, service.Type, service.Orientation, service.Way,
	house.Id, house.DataTime, house.HasImage, house.Name, house.Price, house.Url, house.Location, house.HouseType, house.Orientation, house.Way, house.PlatType)
}

func (f *Filter) Print() {
	for area, idMap := range f.CityServices {
		for _, service := range idMap {
			fmt.Printf("AreaLen:%d, IdLen:%d, %s: %v\n", len(f.CityServices), len(idMap), area, *service)
		}
	}
}

func (f *Filter) Run(config *fetchHouse.Config) error {
	// watcher dir
	quitChan := make(chan struct{})
	go f.watcherService(config.ServiceDir, quitChan)
	defer close(quitChan)

	// filter houses
	for {
		areaHouseMap, err := fetchHouse.FetchHouse(config)
		if err != nil {
			logger.WithField("filter", "house").Errorf("Failed to fetch house[err:%s]", err.Error())
		}
		DebugHouses(&areaHouseMap)

		f.HandleService(areaHouseMap)
		time.Sleep(time.Duration(config.FetchDuration) * time.Second)
	}

	<- quitChan
	return nil
}

func (f *Filter) HandleService(cityHouseMap map[string]map[string][]*fetchHouse.House) error {
	sm.Lock()
	defer sm.Unlock()

	for city, serviceMap := range f.CityServices {
		areaHouses, ok := cityHouseMap[city]
		if !ok {
			logger.WithField("filter", "house").Errorf("city[%s] have not info on rent house", city)
		}

		if err := AnalysisData(serviceMap, areaHouses); err != nil {
			logger.WithField("filter", "house").Errorf("Failed to analysis rent house and services, err:%s", err.Error())
			continue
		}
	}
	return nil
}

func (f *Filter) watcherService(dirPath string, quitChan <-chan struct{}) {
	if err := f.OnScanServiceDir(dirPath); err != nil {
		logger.WithField("filter", "house").Errorf("Failed to scan dir[path:%s, err:%s]", dirPath, err.Error())
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.WithField("filter", "house").Errorf("Failed to new watcher[path:%s, err:%s]", dirPath, err.Error())
		return
	}
	defer watcher.Close()

	if err := watcher.Add(dirPath); err != nil {
		logger.WithField("filter", "house").Errorf("Failed to add watcher[path:%s, err:%s]", dirPath, err.Error())
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			f.handleFileChange(event)
		case <-quitChan:
			logger.WithField("filter", "house").Infof("Watcher service routine was force exits !")
			break
		}
	}
}

func (f *Filter) handleFileChange(event fsnotify.Event) error {
	logger.WithField("service", "parse").Infof("received an watcher event[%v]", event)
	var err error
	if event.Op & fsnotify.Create == fsnotify.Create {
		err = f.newServiceCreate(event.Name)
	} else if event.Op & fsnotify.Remove == fsnotify.Remove {
		err = f.oldServiceRemove(event.Name)
	} else if event.Op & fsnotify.Write == fsnotify.Write {
		err = f.oldServiceChange(event.Name)
	}

	if err != nil {
		logger.WithField("filter", "house").Errorf("Failed to handle change file[path:%s, err:%s]", event.Name, err.Error())
	}
	return err
}

func (f *Filter) newServiceCreate(file string) error {
	var err error
	if f.IsValidFilePath(file) {
		logger.WithField("filter", "house").Infof("Create new file[%s]", file)
		err = f.LoadService(file)
	} else {
		logger.WithField("filter", "house").Warnf("Invalid file[%s], ignore it", file)
	}

	return err
}

func (f *Filter) oldServiceChange(file string) error {
	var err error
	if f.IsValidFilePath(file) {
		logger.WithField("filter", "house").Infof("Change old file[%s]", file)
		err = f.LoadService(file)
	} else {
		logger.WithField("filter", "house").Warnf("Invalid file[%s], ignore it", file)
	}

	return err
}

func (f *Filter) oldServiceRemove(file string) error {
	var err error
	if f.IsValidFilePath(file) {
		logger.WithField("filter", "house").Infof("Remove old file[%s]", file)
		err = f.RemoveService(file)
	} else {
		logger.WithField("filter", "house").Warnf("Invalid file[%s], ignore it", file)
	}

	return err
}

func (f *Filter) OnScanServiceDir(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		logger.WithField("filter", "house").Errorf("Failed to open dir[%s]", dirPath)
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fName := file.Name()
			fPath := path.Join(dirPath, fName)

			if f.IsValidFilePath(fPath) {
				if err := f.LoadService(fPath); err != nil {
					logger.WithField("filter", "house").Errorf("Failed to lioad service[path:%s, err:%s]",
						fPath, err.Error())
				}
			} else {
				logger.WithField("filter", "house").Errorf("Invalid file[%s]", fPath)
			}
		}
	}
	return nil
}

func (f *Filter) LoadService(file string) error {
	service, err := UnmarshalFile(file)
	if err != nil {
		logger.WithField("service", "parse").Errorf("Failed to unmarshal[file:%s, err:%v]", file, err)
		return err
	}

	sm.Lock()
	_, ok := f.IdServices[service.Id]
	sm.Unlock()

	if !ok {
		f.OnServiceCreate(service)
		logger.WithField("filter", "house").Infof("No exist service[%s], create it now", service.Id)
	} else {
		f.OnServiceChange(service)
		logger.WithField("filter", "house").Infof("Exist service[%s], update it now", service.Id)
	}

	return nil
}

func (f *Filter) RemoveService(file string) error {
	if !f.IsValidFilePath(file) {
		logger.WithField("service", "parse").Errorf("Invalid file[%s]", file)
		return errors.New("Invalid file")
	}
	fName := path.Base(file)
	prefixName := strings.TrimRight(fName, FSuffix)

	logger.WithField("service", "parse").Infof("Successful remove service[file:%s]", fName)
	f.OnServiceRemove(prefixName)

	return nil
}

func (f *Filter) OnServiceCreate(service *Service) error {
	sm.Lock()
	defer sm.Unlock()

	f.IdServices[service.Id] = service
	if _, ok := f.CityServices[service.City]; !ok {
		f.CityServices[service.City] = make(map[string]*Service)
	}
	f.CityServices[service.City][service.Id] = service

	return nil
}

func (f *Filter) OnServiceRemove(id string) error {
	sm.Lock()
	defer sm.Unlock()

	service, ok := f.IdServices[id]
	if !ok {
		logger.WithField("service", "parse").Errorf("Failed to find service[%d] when remove service", id)
		return errors.New("Failed to remove service")
	}

	delete(f.CityServices[service.City], id)
	delete(f.IdServices, id)

	return nil
}

func (f *Filter) OnServiceChange(service *Service) error {
	sm.Lock()
	defer sm.Unlock()

	f.IdServices[service.Id] = service
	f.CityServices[service.City][service.Id] = service

	return nil
}

func (f *Filter) IsValidFilePath(fPath string) bool {
	fName := path.Base(fPath)
	fType := path.Ext(fPath)

	if strings.HasPrefix(fName, FPrefix) && fType == FSuffix {
		return true
	}

	return false
}

func AnalysisData(idServiceMap map[string]*Service, areaHouses map[string][]*fetchHouse.House) error {
	for _, service := range idServiceMap {
		logger.WithField("filter", "house").Infof("Start filter service[%v]", service)

		if "" == service.Area {
			for _, houses := range areaHouses {
				housesMap, err := GetValidHouse(service, houses)
				if err != nil {
					logger.WithField("filter", "house").Errorf("Failed to get valid house[err:%s]", err)
					continue
				}

				// send data to target
				SendHouseToTarget(service, housesMap)
			}
		} else {
			housesMap, err := GetValidHouse(service, areaHouses[service.Area])
			if err != nil {
				logger.WithField("filter", "house").Errorf("Failed to get valid house[err:%s]", err)
				continue
			}

			// send data to target
			SendHouseToTarget(service, housesMap)
		}
	}
	return nil
}

func GetValidHouse(service *Service, houses []*fetchHouse.House) (*map[string]string, error) {
	logger.WithField("filter", "house").Infof("Start range house:[len:%d]", len(houses));
	housesMap := make(map[string]string)
	for _, house := range houses {
		DebugMatchData(service, house)

		houseFilter, err := GetPlatInterface(house.PlatType, service)
		if err != nil {
			logger.WithField("filter", "house").Errorf("Failed to get Plat interface[err:%s, platType:%s]",
				err.Error(), house.PlatType)
			continue
		}
		status := GetMatchStatus(service)
		matchStatus := FilterMatchStatus(houseFilter, house)
		logger.WithField("filter", "house").Debugf("House match status:%d, service match status:%d", matchStatus, status)
		if status == matchStatus {
			logger.WithField("filter", "house").Infof("Successful match[service:%v, house:%v]", *house, *service)

			houseString := fmt.Sprintf("Congratulations! The following house:%s|%s|%s|%s|%s|%s|%s",
				house.Url, house.Name, house.Location, house.Price, house.Way, house.HouseType, house.Orientation)

			housesMap[house.Id] = houseString
		}
	}
	return &housesMap, nil
}

func SendHouseToTarget(service *Service, houseMap *map[string]string) error {
	for _, houseFormat := range *houseMap {
		fmt.Printf("Have match data[%s]\n", houseFormat)
	}
	return nil
}

func GetMatchStatus(service *Service) uint32 {
	matchStatus := uint32(0)
	if service.Price != "" {
		matchStatus |= 1 << 0
	}

	if service.Type != ""  {
		matchStatus |= 1 << 1
	}

	if service.Way != ""  {
		matchStatus |= 1 << 2
	}

	if service.Orientation != ""  {
		matchStatus |= 1 << 3
	}

	return matchStatus
}

func FilterMatchStatus(houseFilter HouseFilter, house *fetchHouse.House) uint32 {
	var matchStatus = uint32(0)

	if bValidPrice := houseFilter.ValidPrice(house); bValidPrice {
		SetBIT(&matchStatus, 1)
	}

	if bValidType := houseFilter.ValidType(house); bValidType {
		SetBIT(&matchStatus, 2)
	}

	if bValidWay := houseFilter.ValidWay(house); bValidWay {
		SetBIT(&matchStatus, 3)
	}

	if bValidOrientation := houseFilter.ValidOrientation(house); bValidOrientation {
		SetBIT(&matchStatus, 4)
	}

	return matchStatus
}
