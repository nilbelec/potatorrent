package scheduler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"sync"
)

type SchedulesFile struct {
	sync.Mutex
	filename string
}

type schedulesMap map[string]*ScheduleSearch

// NewSchedulesFile creates a new SchedulesFile to store the schedules
func NewSchedulesFile(filename string) *SchedulesFile {
	f := &SchedulesFile{filename: filename}
	e, _ := f.exists()
	if !e {
		f.saveJSON(make(schedulesMap))
	}
	return f
}

func (f *SchedulesFile) GetAll() []*ScheduleSearch {
	f.Lock()
	defer f.Unlock()
	schedules, err := f.readJSON()
	if err != nil {
		return nil
	}
	r := make([]*ScheduleSearch, 0)
	for _, v := range schedules {
		r = append(r, v)
	}
	sort.SliceStable(r, func(i, j int) bool {
		return r[i].LastExecution.After(r[j].LastExecution)
	})
	return r
}

func (f *SchedulesFile) Delete(id string) error {
	f.Lock()
	defer f.Unlock()
	schedules, err := f.readJSON()
	if err != nil {
		return err
	}
	delete(schedules, id)
	f.saveJSON(schedules)
	return nil
}

func (f *SchedulesFile) GetAllIDs() ([]string, error) {
	f.Lock()
	defer f.Unlock()
	r := make([]string, 0)
	schedules, err := f.readJSON()
	if err != nil {
		return nil, err
	}
	for k := range schedules {
		r = append(r, k)
	}
	return r, nil
}

func (f *SchedulesFile) Save(ss *ScheduleSearch) error {
	f.Lock()
	defer f.Unlock()
	schedules, err := f.readJSON()
	if err != nil {
		return err
	}
	schedules[ss.ID] = ss
	return f.saveJSON(schedules)
}

func (f *SchedulesFile) Get(id string) (*ScheduleSearch, error) {
	f.Lock()
	defer f.Unlock()
	var schedules, err = f.readJSON()
	if err != nil {
		return nil, err
	}
	return schedules[id], nil
}

func (f *SchedulesFile) exists() (bool, error) {
	_, err := os.Stat(f.filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (f *SchedulesFile) saveJSON(v schedulesMap) error {
	b, err := json.MarshalIndent(&v, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.filename, b, 0644)
}

func (f *SchedulesFile) readJSON() (schedulesMap, error) {
	b, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return nil, err
	}
	var schedules = schedulesMap{}
	err = json.Unmarshal(b, &schedules)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}
