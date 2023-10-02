package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
	"time"
)

type Config struct {
	ServerConfig *ServerConfig `yaml:"ServerConfig"`
	KafkaConfigs []KafkaConfig `yaml:"KafkaConfigs"`
	MetaConfig   *MetaConfig   `yaml:"MetaConfig"`
}

type ServerConfig struct {
	Addr                  string `yaml:"addr"`
	Port                  int    `yaml:"port"`
	ReadTimeout           int    `yaml:"read_timeout"`
	WriteTimeout          int    `yaml:"write_timeout"`
	IdleTimeout           int    `yaml:"idle_timeout"`
	MaxConnsPerIP         int    `yaml:"max_conns_per_ip"`
	MaxIdleWorkerDuration int    `yaml:"max_idle_worker_duration"`
	MaxRequestBodySize    int    `yaml:"max_request_body_size"`
}

type KafkaConfig struct {
	Cluster        string        `yaml:"cluster"`
	Brokers        string        `yaml:"brokers"`
	ClientId       string        `yaml:"client_id"`
	BufferSize     int           `yaml:"buffer_size"`
	MaxRetries     int           `yaml:"max_retries"`
	FlushBytes     int           `yaml:"flush_bytes"`
	FlushMaxSize   int           `yaml:"flush_max_size"`
	FlushSize      int           `yaml:"flush_size"`
	FlushFrequency time.Duration `yaml:"flush_frequency"`
}

type MetaConfig struct {
	XMLFile            string `yaml:"xml_file"`
	RefreshIntervalSec int    `yaml:"refresh_interval_sec"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
}

type Meta struct {
	Apps []App `xml:"app"`
}

type App struct {
	Id     int    `xml:"id"`
	BlCode string `xml:"blcode"`
	Key    string `xml:"key"`
	Kafka  string `xml:"kafka"`
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	// 读文件
	filePath := "/tmp/a.txt"
	writeFile(filePath)
	appendFile(filePath)
	readFile(filePath)

	// 读yaml
	println("=====================================")
	ymlPath := "golang_study/basic_grammar/io.yml"
	readYaml(ymlPath)

	// 读xml
	println("=====================================")
	xmlPath := "golang_study/basic_grammar/io.xml"
	readXml(xmlPath)

	// 读String
	println("=====================================")
	readString()

}

func readString() {
	// 创建了一个 strings.Reader 并以每次 8 字节的速度读取它的输出
	r := strings.NewReader("Hello, Reader!")
	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
}

func readFile(filePath string) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("Failed to read file, err: %s", err.Error())
		panic(err)
	}
	fmt.Println(string(bytes))
}

func writeFile(filePath string) {
	d := []byte("haha")
	if err := os.WriteFile(filePath, d, 0755); err != nil {
		logrus.Errorf("Failed to write file, err: %s", err.Error())
		panic(err)
	}
}

func appendFile(filePath string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0755)
	defer file.Close()
	if err != nil {
		logrus.Errorf("Failed to append file, err: %s", err.Error())
		panic(err)
	}

	//写入文件时，使用带缓存的 *Writer
	writer := bufio.NewWriter(file)
	for i := 0; i < 5; i++ {
		_, err := writer.WriteString("haha\n")
		if err != nil {
			logrus.Errorf("Failed to append file, err: %s", err.Error())
			panic(err)
		}
	}
	_, err = writer.WriteString("haha\nhaha1,haha2\nhaha3,haha4\nhaha5")
	if err != nil {
		logrus.Errorf("Failed to append file, err: %s", err.Error())
	}
	//Flush将缓存的文件真正写入到文件中
	err = writer.Flush()
	if err != nil {
		logrus.Errorf("Failed to flush file, err: %s", err.Error())
	}
}

func readXml(filePath string) {
	// 读取xml配置文件 绑定到结构体
	file, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("Failed to read xml file, err: %s", err.Error())
		panic(err)
	}
	var meta Meta
	err = xml.Unmarshal(file, &meta)
	if err != nil {
		logrus.Errorf("Failed to unmarshal xml file, err: %s", err.Error())
		panic(err)
	}
	println(len(meta.Apps))
	for _, app := range meta.Apps {
		println(app.Id)
	}
}

func readYaml(filePath string) {
	// 读取yaml配置文件 绑定到结构体
	file, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("Failed to read config file, err: %s", err.Error())
		panic(err)
	}
	var conf Config
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		logrus.Errorf("Failed to read config file, err: %s", err.Error())
		panic(err)
	}
	println(conf.ServerConfig.Port)
	println(conf.ServerConfig.MaxConnsPerIP)
	println(len(conf.KafkaConfigs))
}
