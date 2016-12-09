package filter

import (
	"strings"
	"path"
	"io/ioutil"
	"uframework/log"
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
	AreaServices map[string]map[string]*Service // key: Area subKey: Id
}

func (f *Filter) Print() {
	for area, idMap := range f.AreaServices {
		for _, service := range idMap {
			fmt.Printf("AreaLen:%d, IdLen:%d, %s: %v\n", len(f.AreaServices), len(idMap), area, *service)
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
		uflog.ERRORF("Failed to fetch house[err:%s]", err.Error())
		}
		f.HandleService(areaHouseMap)
		time.Sleep(config.FetchDuration * time.Second)
	}

	<- quitChan
	return nil
}

func (f *Filter) HandleService(areaHouseMap map[string][]*fetchHouse.House) error {
	sm.Lock()
	defer sm.Unlock()

	for area, serviceMap := range f.AreaServices {
		houses := areaHouseMap[area]
		if houses == nil {
			uflog.ERRORF("Area[%s] have not info on rent house", area)
		}
		if err := AnalysisData(serviceMap, houses); err != nil {
			uflog.ERRORF("Failed to analysis rent house and services, err:%s", err.Error())
			continue
		}
	}
	return nil
}

func (f *Filter) watcherService(dirPath string, quitChan <-chan struct{}) {
	if err := f.OnScanServiceDir(dirPath); err != nil {
		uflog.ERRORF("Failed to scan dir[path:%s, err:%s]", dirPath, err.Error())
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		uflog.ERRORF("Failed to new watcher[path:%s, err:%s]", dirPath, err.Error())
		return
	}
	defer watcher.Close()

	if err := watcher.Add(dirPath); err != nil {
		uflog.ERRORF("Failed to add watcher[path:%s, err:%s]", dirPath, err.Error())
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			f.handleFileChange(event)
		case <-quitChan:
			uflog.INFOF("Watcher service routine was force exits !")
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
		uflog.ERRORF("Failed to handle change file[path:%s, err:%s]", event.Name, err.Error())
	}
	return err
}

func (f *Filter) newServiceCreate(file string) error {
	var err error
	if f.IsValidFilePath(file) {
		uflog.INFOF("Create new file[%s]", file)
		err = f.LoadService(file)
	} else {
		uflog.WARNF("Invalid file[%s], ignore it", file)
	}

	return err
}

func (f *Filter) oldServiceChange(file string) error {
	var err error
	if f.IsValidFilePath(file) {
		uflog.INFOF("Change old file[%s]", file)
		err = f.LoadService(file)
	} else {
		uflog.WARNF("Invalid file[%s], ignore it", file)
	}

	return err
}

func (f *Filter) oldServiceRemove(file string) error {
	var err error
	if f.IsValidFilePath(file) {
		uflog.INFOF("Remove old file[%s]", file)
		err = f.RemoveService(file)
	} else {
		uflog.WARNF("Invalid file[%s], ignore it", file)
	}

	return err
}

func (f *Filter) OnScanServiceDir(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		uflog.ERRORF("Failed to open dir[%s]", dirPath)
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fName := file.Name()
			fPath := path.Join(dirPath, fName)

			if f.IsValidFilePath(fPath) {
				if err := f.LoadService(fPath); err != nil {
					uflog.ERRORF("Failed to lioad service[path:%s, err:%s]",
						fPath, err.Error())
				}
			} else {
				uflog.ERRORF("Invalid file[%s]", fPath)
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
		uflog.INFOF("No exist service[%s], create it now", service.Id)
	} else {
		f.OnServiceChange(service)
		uflog.INFOF("Exist service[%s], update it now", service.Id)
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
	if _, ok := f.AreaServices[service.Area]; !ok {
		f.AreaServices[service.Area] = make(map[string]*Service)
	}
	f.AreaServices[service.Area][service.Id] = service

	f.Print()
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

	delete(f.AreaServices[service.Area], id)
	delete(f.IdServices, id)

	f.Print()
	return nil
}

func (f *Filter) OnServiceChange(service *Service) error {
	sm.Lock()
	defer sm.Unlock()

	f.IdServices[service.Id] = service
	f.AreaServices[service.Area][service.Id] = service

	f.Print()
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

func AnalysisData(idServiceMap map[string]*Service, houses []*fetchHouse.House) error {
	for _, service := range idServiceMap {
		housesMap, err := GetValidHouse(service, houses)
		if err != nil {
			uflog.ERRORF("Failed to get valid house[err:%s]", err)
			continue
		}

		// send data to target
		SendHouseToTarget(service, housesMap)
	}
	return nil
}

func GetValidHouse(service *Service, houses []*fetchHouse.House) (*map[string]string, error) {
	housesMap := make(map[string]string)
	for _, house := range houses {
		houseFilter, err := GetPlatInterface(house.PlatType, service)
		if err != nil {
			uflog.ERRORF("Failed to get Plat interface[err:%s, platType:%s]",
				err.Error(), house.PlatType)
			continue
		}
		status := GetMatchStatus(service)
		matchStatus := FilterMatchStatus(houseFilter, house)
		if status == matchStatus {
			uflog.INFOF("Successful match[service:%v, house:%v]", *house, *service)

			houseString := fmt.Sprintf("Congratulations! The following house:\n%s|%s|%s|%s|%s|%s|%s",
				house.Url, house.Name, house.Location, house.Price, house.Way, house.HouseType, house.Orientation)

			housesMap[house.Id] = houseString
		}
	}
	return &housesMap, nil
}

func SendHouseToTarget(service *Service, houseMap *map[string]string) error {
	for _, houseFormat := range *houseMap {
		fmt.Println("Have match data[%s]", houseFormat)
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
