package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jinzhu/configor"
)

// LoadConfig 加载配置文件
func LoadConfig(config interface{}, path string) {
	if os.Getenv("CONFIGOR_ENV") == "" {
		os.Setenv("CONFIGOR_ENV", "local")
	}
	err := configor.Load(config, path)
	if err != nil {
		panic(err)
	}
	PrintConfig(config)
}

func PrintConfig(config interface{}) {
	var bytes []byte
	bytes, _ = json.Marshal(config)
	fmt.Println(string(bytes))
}
