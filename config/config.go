package config

import (
	"log"
	"os"
	"os/user"

	"github.com/spf13/viper"
)

const (
	ConfigFolderName = ".tx"
	ConfigFileName   = "tx"
	ConfigFileExt    = "json"
)

var (
	ConfigDir   string
	HomeDir     string
	TemplateDir string
	AddConfig   AddTemplateConfig
	GenConfig   GenTemplateConfig
)

type Config interface {
	Load() interface{}
}

type AddTemplateConfig struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ignore      []string `json:"ignore"`
	Path        string
}

type GenTemplateConfig struct {
	Template map[string]string `json:"template"`
	Path     string
}

func (c *GenTemplateConfig) Load() GenTemplateConfig {

	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileExt)
	viper.AddConfigPath(c.Path)
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("Error reading config file, ", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&c)

	if err != nil {
		log.Println("Error unmarshalling config file, ", err)
		os.Exit(1)
	}

	return *c
}

func (c *AddTemplateConfig) Load() AddTemplateConfig {
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileExt)
	viper.AddConfigPath(c.Path)
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("Error reading config file, ", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&c)

	if err != nil {
		log.Println("Error unmarshalling config file, ", err)
		os.Exit(1)
	}

	return *c
}

func NewAddTemplateConfig(path string) AddTemplateConfig {
	return AddTemplateConfig{
		Path: path,
	}
}

func NewGenTemplateConfig(path string) GenTemplateConfig {
	return GenTemplateConfig{
		Path: path,
	}
}

func init() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	HomeDir = user.HomeDir
	ConfigDir = user.HomeDir + "/" + ConfigFolderName
	TemplateDir = ConfigDir + "/templates"

	AddConfig = NewAddTemplateConfig("")
}
