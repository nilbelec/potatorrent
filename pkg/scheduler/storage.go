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
	filename  string
	schedules schedulesMap
}

type schedulesMap map[string]*ScheduleSearch

// NewSchedulesFile creates a new SchedulesFile to store the schedules
func NewSchedulesFile(filename string) *SchedulesFile {
	f := &SchedulesFile{filename: filename, schedules: make(schedulesMap)}
	e, err := f.exists()
	if err != nil {
		panic(err)
	}
	if !e {
		f.persist()
	}
	f.load()
	return f
}

func (f *SchedulesFile) GetAll() []*ScheduleSearch {
	f.Lock()
	defer f.Unlock()
	r := make([]*ScheduleSearch, 0)
	for _, v := range f.schedules {
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
	delete(f.schedules, id)
	return f.persist()
}

func (f *SchedulesFile) GetAllIDs() ([]string, error) {
	f.Lock()
	defer f.Unlock()
	r := make([]string, 0)
	for k := range f.schedules {
		r = append(r, k)
	}
	return r, nil
}

func (f *SchedulesFile) Save(ss *ScheduleSearch) error {
	f.Lock()
	defer f.Unlock()
	f.schedules[ss.ID] = ss
	return f.persist()
}

func (f *SchedulesFile) Get(id string) *ScheduleSearch {
	f.Lock()
	defer f.Unlock()
	return f.schedules[id]
}

func (f *SchedulesFile) exists() (bool, error) {
	_, err := os.Stat(f.filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (f *SchedulesFile) persist() error {
	b, err := json.MarshalIndent(&f.schedules, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.filename, b, 0644)
}

func (f *SchedulesFile) load() error {
	b, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &f.schedules)
}
