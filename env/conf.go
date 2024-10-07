package env

import (
	"github.com/spf13/viper"
	"path/filepath"
)

const DefaultRootPath = "."

var (
	rootPath string
)

func SetRootPath(r string) {
	rootPath = r
}

func GetRootPath() string {
	if rootPath != "" {
		return rootPath
	} else {
		return DefaultRootPath
	}
}

func LoadConf(filename string, s interface{}) {
	path := filepath.Join(GetRootPath(), filename)

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(path + " read error: " + err.Error())
	}

	if err := viper.Unmarshal(s); err != nil {
		panic(filename + " unmarshal error: " + err.Error())
	}
}
