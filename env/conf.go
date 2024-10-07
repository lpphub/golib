package env

import (
	"github.com/spf13/viper"
	"path/filepath"
)

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
