package filter

import (
	"strings"
	"path"
	"io/ioutil"
	"uframework/log"
	"encoding/json"
	"sync"
	"github.com/fsnotify/fsnotify"
	"auto_rent/fetch_house"
	"forward_port/config"
)

type SimpleHouse struct {
	Price       uint32 // 价格范围
	Type        uint32 // 3室一厅
	Orientation uint32 // 朝向
	Way         uint32 // 整租
}

type Service struct {
	Id        string
	Username  string
	Email     string
	Number    string
	Platform  uint32
	City      string
	Area      string
	Community string
	SimpleHouse
}

var (
	FPrefix = "House_"
	FSuffix = ".json"
)

var sm sync.Mutex

type Filter struct {
	IdServices map[string]*Service	// key: Id

}

func (f *Filter) Run(config *fetchHouse.Config) error {
	// watcher dir
	quitChan := make(chan struct{})
	go f.watcherService(config.ServiceDir, quitChan)
	defer close(quitChan)

	// filter houses
	houses, err := fetchHouse.FetchHouse(config)
	if err != nil {
		uflog.ERRORF("Failed to fetch house[err:%s]", err.Error())
	}
	return nil
}

func (f *Filter) FilterHouses(houses *fetchHouse.House) error {
	return nil
}

func (f *Filter) HandleHouse(house *fetchHouse.House) error {
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
	var err error
	if event & fsnotify.Create == fsnotify.Create {
		err = f.newServiceCreate(event.Name)
	} else if event & fsnotify.Remove == fsnotify.Remove {
		err = f.oldServiceRemove(event.Name)
	} else if event & fsnotify.Write == fsnotify.Write {
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
		err = f.LoadService(file)
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
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		uflog.ERRORF("Failed to read file[%s]", file)
		return err
	}
	uflog.INFOF("Read %s bytes from %s", len(bytes), file)

	var service Service
	if err := json.Unmarshal(bytes, &service); err != nil {
		uflog.ERRORF("Failed to unmarshal file[%s]", file)
		return err
	}

	sm.Lock()
	_, ok := f.IdServices[service.Id]
	sm.Unlock()

	if !ok {
		f.OnServiceCreate(&service)
		uflog.INFOF("No exist service[%s], create it now", service.Id)
	} else {
		f.OnServiceChange(&service)
		uflog.INFOF("Exist service[%s], update it now", service.Id)
	}

	return nil
}

func (f *Filter) OnServiceCreate(service *Service) error {
	sm.Lock()
	defer sm.Unlock()

	f.IdServices[service.Id] = service
	return nil
}

func (f *Filter) OnServiceRemove(service *Service) error {
	sm.Lock()
	defer sm.Unlock()

	delete(f.IdServices, service.Id)
	return nil
}

func (f *Filter) OnServiceChange(service *Service) error {
	sm.Lock()
	defer sm.Unlock()

	f.IdServices[service.Id] = service
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