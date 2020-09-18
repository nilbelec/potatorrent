package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const baseURLOptionsKey = "baseURL"
const baseURLDefault = "https:" + "//" + "pc" + "tmi" + "x.com"
const downloadFolderOptionsKey = "downloadFolder"
const portOptionsKey = "port"
const portDefault = "8080"

type ConfigFile struct {
	sync.Mutex
	filename string
	options  map[string]string
}

func NewConfigFile(filename string) *ConfigFile {
	c := &ConfigFile{filename: filename, options: defaults()}
	e, err := c.exists()
	if err != nil {
		panic(err)
	}
	if !e {
		c.persist()
	}
	c.load()
	return c
}

func defaults() map[string]string {
	opts := make(map[string]string)
	opts[baseURLOptionsKey] = baseURLDefault
	opts[downloadFolderOptionsKey], _ = filepath.Abs(".")
	opts[portOptionsKey] = portDefault
	return opts
}

func (c *ConfigFile) BaseURL() string {
	if val, ok := c.options[baseURLOptionsKey]; ok {
		return val
	}
	return baseURLDefault
}

func (c *ConfigFile) DownloadFolder() string {
	if val, ok := c.options[downloadFolderOptionsKey]; ok {
		return val
	}
	path, _ := filepath.Abs(".")
	return path
}

func (c *ConfigFile) Port() string {
	if val, ok := c.options[portOptionsKey]; ok {
		return val
	}
	return portDefault
}

func (c *ConfigFile) exists() (bool, error) {
	_, err := os.Stat(c.filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (c *ConfigFile) persist() error {
	b, err := json.MarshalIndent(&c.options, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.filename, b, 0644)
}

func (c *ConfigFile) load() error {
	b, err := ioutil.ReadFile(c.filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &c.options)
}
