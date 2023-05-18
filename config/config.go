package config

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/PondWader/GoPractice/utils"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

type ServerConfiguration struct {
	Port                 int    `default:"25565" yaml:"port"`
	Motd                 string `default:"GoPractice Server" yaml:"motd"`
	OnlineMode           bool   `default:"true" yaml:"online-mode"`
	CompressionThreshold int    `default:"256" yaml:"compression-threshold"`

	DatabaseUser     string `default:"" yaml:"database-user"`
	DatabasePassword string `default:"" yaml:"database-password"`
	DatabaseName     string `default:"" yaml:"database-name"`
	DatabaseHost     string `default:"" yaml:"database-host"`
	DatabasePort     uint16 `default:"3306" yaml:"database-port"`
}

func LoadConfig() ServerConfiguration {
	config := ServerConfiguration{}
	defaults.Set(&config)

	if _, err := os.Stat("server.yml"); err == nil {
		configData, err := ioutil.ReadFile("server.yml")
		if err != nil {
			utils.Error("Failed to read server.yml file:", err)
			os.Exit(1)
			return config
		}
		err = yaml.Unmarshal(configData, &config)
		if err != nil {
			utils.Error("Error parsing server.yml YAML:", err)
			os.Exit(1)
			return config
		}
		return config
	} else if errors.Is(err, os.ErrNotExist) {
		defaultConfig, _ := yaml.Marshal(config)
		ioutil.WriteFile("server.yml", defaultConfig, 0777)
		return config
	} else {
		utils.Error("Failed to load server.yml file info:", err)
		os.Exit(1)
		return config
	}
}
