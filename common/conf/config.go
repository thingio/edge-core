package conf

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/configor"
	"os"
	"strconv"
	"strings"
	"time"
)

type MqttConfig struct {
	BrokerAddr     string        `json:"broker_addr" yaml:"broker_addr"`
	ConnectTimeout time.Duration `json:"connect_timeout" yaml:"connect_timeout"`
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout"`
}

type DBConfig struct {
	File    string        `json:"file" yaml:"file"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

func LoadConfig(dest interface{}, defaultConfig string) {
	file := os.Getenv("CONFIG_FILE")
	if file == "" {
		file = defaultConfig
	}

	fmt.Printf("Loading config: %s\n", file)
	if _, err := os.Stat(file); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err := configor.Load(dest, file)
	if err != nil {
		fmt.Printf("Fail to load config file %s: %s\n", file, err)
		os.Exit(1)
	}


	data, _ := json.Marshal(dest)
	fmt.Printf("Config loaded: %s\n", data)
}

const (
	EnvSep = ","
)

func LoadEnvs(envs map[string]interface{}) error {
	var err error
	for env, target := range envs {
		value := os.Getenv(env)
		if value == "" {
			continue
		}
		switch target.(type) {
		case *string:
			*(target.(*string)) = value
		case *[]string:
			values := strings.Split(value, EnvSep)
			result := make([]string, 0)
			for _, v := range values {
				if v != "" {
					result = append(result, v)
				}
			}
			if len(result) != 0 {
				*(target.(*[]string)) = result
			}
		case *int:
			*(target.(*int)), err = strconv.Atoi(value)
		case *bool:
			*(target.(*bool)), err = strconv.ParseBool(value)
		case *int64:
			*(target.(*int64)), err = strconv.ParseInt(value, 10, 64)
		case *float32:
			if v, err := strconv.ParseFloat(value, 32); err == nil {
				*(target.(*float32)) = float32(v)
			}
		case *float64:
			*(target.(*float64)), err = strconv.ParseFloat(value, 64)
		case *time.Duration:
			*(target.(*time.Duration)), err = time.ParseDuration(value)
		default:
			return fmt.Errorf("unsupported env type : %T", target)
		}
		if err != nil {
			return fmt.Errorf("error while loading environments: %v", err)
		}
	}
	return nil
}
