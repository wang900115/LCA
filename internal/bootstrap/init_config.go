package bootstrap

import (
	"os"

	"github.com/spf13/viper"
)

// func NewConfig() *viper.Viper {
// 	envConf := filepath.Join("config", "app.yaml")
// 	return getConfig(filepath.Join(GetAppPath(), envConf))
// }

// func getConfig(path string) *viper.Viper {
// 	conf := viper.New()
// 	conf.SetConfigFile(path)
// 	err := conf.ReadInConfig()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return conf
// }

// func GetAppPath() string {
// 	file, _ := exec.LookPath(os.Args[0])
// 	path, _ := filepath.Abs(file)
// 	index := strings.LastIndex(path, string(os.PathSeparator))

// 	return path[:index]
// }

func NewConfig() *viper.Viper {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/app.yaml"
	}
	return getConfig(path)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	if err := conf.ReadInConfig(); err != nil {
		panic(err)
	}
	return conf
}
