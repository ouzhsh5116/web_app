package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局配置变量
var Conf = new(AppConfig)

type AppConfig struct {
	Mode         string `mapstructure:"mode"`
	Port         int    `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

func Init(fileName string) (err error) {
	// 读取的路径默认是以项目名开头的了，不用再另外写项目名，项目启动默认在项目目录下！！！
	// 读取方式 1.
	viper.SetConfigFile(fileName)

	// 读取方式 2.
	// viper.SetConfigName("config")       // 配置文件名称(无扩展名)
	// viper.SetConfigType("yaml")  // 如果配置文件的名称中没有扩展名，则需要配置此项，专门用于配置远程文件类型
									//可以不写
	// viper.AddConfigPath("conf") // 查找配置文件所在的路径

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	}

	// 将配置文件解析到对应的结构体中
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(fmt.Errorf("unmarshal to Conf failed, err:%v", err))
	}

	// 监控配置文件是否修改
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("夭寿啦~配置文件被人修改啦...")
		if err := viper.Unmarshal(&Conf); err != nil {
			fmt.Printf("ReadInConfig failed, err: %v", err)
		}
	})
	return err
}
