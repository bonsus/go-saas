package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App   AppConfig `yaml:"app"`
	DB    DBConfig  `yaml:"db"`
	JWT   JWT       `yaml:"jwt"`
	File  File      `yaml:"file"`
	S3    S3        `yaml:"s3"`
	Redis Redis     `yaml:"redis"`
}
type AppConfig struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
}
type DBConfig struct {
	Host           string                 `yaml:"host"`
	Port           string                 `yaml:"port"`
	User           string                 `yaml:"user"`
	Password       string                 `yaml:"password"`
	Name           string                 `yaml:"name"`
	TZ             string                 `yaml:"timezone"`
	ConnectionPool DBConnectionPoolConfig `yaml:connection_pool"`
}
type DBConnectionPoolConfig struct {
	MaxIdleConnection     uint8 `yaml:"max_idle_connection"`
	MaxOpenConnection     uint8 `yaml:"max_open_connection"`
	MaxLifetimeConnection uint8 `yaml:"max_lifetime_connection"`
	MaxIdletimeConnection uint8 `yaml:"max_idletime_connection"`
}
type JWT struct {
	Key        string `yaml:"key"`
	ExpireTime int    `yaml:"expire_time"`
	Type       string `yaml:"type"`
}
type File struct {
	TempDir string `yaml:"temp_dir"`
}
type S3 struct {
	Dir        string `yaml:"dir"`
	Domain     string `yaml:"domain"`
	BucketName string `yaml:"bucket_name"`
	Endpoint   string `yaml:"endpoint"`
	Region     string `yaml:"region"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
}
type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// var Cfg Config
var Cfg *Config = &Config{}

func LoadConfig(filename string) (err error) {
	configByte, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	return yaml.Unmarshal(configByte, &Cfg)
}
func GetConfig() *Config {
	return Cfg
}
