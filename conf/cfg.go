package conf

// config.go 读取全局参数

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var HostName string
var NodesPreFix string = "/configService/nodes/"
var PathPreFix string = "/configService/path/"
var ConfigPreFix string = "/configService/config/"
var ZkConfSlice []string = []string{"172.16.17.100:2181", "172.16.17.101:2181", "172.16.17.102:2181"}

// Config 全局配置参数
type Config struct {
	Server struct {
		Addr         string `yaml:"addr"`
		BasePath     string `yaml:"basePath"`
		RegisterName string `yaml:"registerName"`
	} `yaml:"server"`

	Redis struct {
		Mode     string   `yaml:"mode"` // 模式
		Password string   `yaml:"pass"` // 密码
		Addr     []string `yaml:"host"` // 集群地址
	} `yaml:"redis"`

	Runtime struct {
		PartInterval int `yaml:"partInterval"`
		Interval     int `yaml:"interval"`
		LogLevel     int `yaml:"logLevel"` // 日志输出级别 0:不输出 1:只输出panic及fault 2:输出Error 3:输出Warn 4:输出Info 5:输出Debug
	} `yaml:"runtime"`

	Zookeeper struct {
		Addr []string `yaml:"addrs"`
	} `yaml:"zookeeper"`

	Service struct {
		Debug string `yaml:"debug"`
	} `yaml:"service"`

	Logger struct {
		Filename   string `yaml:"filename"`   //文件名, 默认"../log/foo.log"
		MaxSize    int    `yaml:"maxSize"`    //文件最大容量，超过会自动拆分，默认为100M
		MaxBackups int    `yaml:"maxBackups"` //保留的文件个数，超过个数，最旧的文件会被删除, 默认10个
		MaxDays    int    `yaml:"maxDays"`    //保留的日志天数，默认14天
		Compress   bool   `yaml:"compress"`   //是否压缩文件，默认压缩
	} `yaml:"logger"`
}

var Cfg Config //全局配置句柄

func init() {
	// 设置配置文件路径
	var err error
	HostName, err = os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		panic(err)

	}
	err = ParseYaml("../config/configService.yaml", &Cfg)
	if err != nil {
		err = ParseYaml("../config/conf.yaml", &Cfg)
		if err != nil {
			fmt.Println("no config file")
		}
	}
}

func (c *Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return ""
	}
	return out.String()
}

// ParseYaml 读取配置参数
func ParseYaml(path string, cfg any) error {
	//加载配置文件
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	//反序列化配置文件数据结构
	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		return err
	}
	return nil
}
