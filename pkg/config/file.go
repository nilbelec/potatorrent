package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const baseURLDefault = "https:" + "//" + "pc" + "tmi" + "x.com"
const intervalInMinutesDefault = 60
const portDefault = 8080

type Configuration struct {
	BaseURL           string `json:"baseURL"`
	IntervalInMinutes int    `json:"intervalInMinutes"`
	DownloadFolder    string `json:"downloadFolder"`
	Port              int    `json:"port"`
}

type ConfigFile struct {
	sync.Mutex
	filename      string
	configuration *Configuration
}

func NewConfigFile(filename string) *ConfigFile {
	c := &ConfigFile{filename: filename, configuration: defaults()}
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

func defaults() *Configuration {
	path, _ := filepath.Abs(".")
	return &Configuration{
		BaseURL:           baseURLDefault,
		DownloadFolder:    path,
		Port:              portDefault,
		IntervalInMinutes: intervalInMinutesDefault,
	}
}

func (c *ConfigFile) GetConfiguration() *Configuration {
	c.load()
	return c.configuration
}

func (c *ConfigFile) SaveConfiguration(cfg *Configuration) error {
	err := c.validate(cfg)
	if err != nil {
		return err
	}
	c.configuration = cfg
	return c.persist()
}

func (c *ConfigFile) validate(cfg *Configuration) error {
	cfg.BaseURL = strings.TrimSuffix(cfg.BaseURL, "/")
	_, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return errors.New("La URL no es válida")
	}
	cfg.DownloadFolder = strings.TrimSuffix(cfg.DownloadFolder, "/")
	if _, err := os.Stat(cfg.DownloadFolder); os.IsNotExist(err) {
		return errors.New("El directorio introducido no existe")
	}
	return nil
}

func (c *ConfigFile) BaseURL() string {
	c.load()
	return c.configuration.BaseURL
}

func (c *ConfigFile) IntervalInMinutes() int {
	c.load()
	return c.configuration.IntervalInMinutes
}

func (c *ConfigFile) DownloadFolder() string {
	c.load()
	return c.configuration.DownloadFolder
}

func (c *ConfigFile) Port() string {
	c.load()
	return strconv.Itoa(c.configuration.Port)
}

func (c *ConfigFile) exists() (bool, error) {
	_, err := os.Stat(c.filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (c *ConfigFile) persist() error {
	b, err := json.MarshalIndent(&c.configuration, "", " ")
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
	return json.Unmarshal(b, &c.configuration)
}
