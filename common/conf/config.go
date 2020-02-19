package conf

import (
	"fmt"
	"github.com/jinzhu/configor"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type MqttConfig struct {
	BrokerAddr       string `yaml:"broker_addr"`
	QoS              byte    `yaml:"qos"`
	ConnectTimeoutMS int    `yaml:"connect_timeout_ms"`
	RequestTimeoutMS int    `yaml:"request_timeout_ms"`
}

func LoadConfig(dest interface{}, defaultConfig string) {
	file := os.Getenv("CONFIG_FILE")
	if file == "" {
		file = defaultConfig
	}
	err := configor.Load(dest, file)
	if err != nil {
		log.Panicf("Failed to load config file %s: %s", file, err)
	}
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
