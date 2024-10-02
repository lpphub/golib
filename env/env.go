package env

import (
	"gopkg.in/yaml.v3"
	"os"
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

func LoadConf(filename, subConf string, s interface{}) {
	var path string
	path = filepath.Join(GetRootPath(), subConf, filename)

	if yamlFile, err := os.ReadFile(path); err != nil {
		panic(filename + " read error: " + err.Error())
	} else if err = yaml.Unmarshal(yamlFile, s); err != nil {
		panic(filename + " unmarshal error: " + err.Error())
	}
}
