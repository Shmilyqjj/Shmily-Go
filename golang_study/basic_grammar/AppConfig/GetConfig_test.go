package AppConfig

import (
	"encoding/xml"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

// -c /home/shmily/Projects/MyProjects/GolangProjects/Shmily-Go/golang_study/basic_grammar/AppConfig/config.yml
func GetYamlConfig() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	//Init config
	confPath := "config.yml"
	if len(os.Args) > 1 {
		if os.Args[1] == "-c" {
			confPath = os.Args[2]
		}
	}
	file, err := os.ReadFile(confPath)
	if err != nil {
		logrus.Errorf("Failed to read config file '%s', err: %s", confPath, err.Error())
		panic(err)
	}
	var appConfig Config
	err = yaml.Unmarshal(file, &appConfig)
	if err != nil {
		logrus.Errorf("Failed to init config, err: %s", err.Error())
		panic(err)
	}
	logrus.Infof("Config inited: %+v", appConfig)
}

func GetXmlConfig() error {
	confPath := "config.xml"
	file, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}
	var users Users
	err = xml.Unmarshal(file, &users)
	if err != nil {
		return err
	}
	logrus.Infof("Config inited: %+v", users)
	return nil
}

func TestGetConfig(t *testing.T) {
	GetYamlConfig()
}

func TestGetXmlConfig(t *testing.T) {
	if err := GetXmlConfig(); err != nil {
		panic(err)
	}
}
