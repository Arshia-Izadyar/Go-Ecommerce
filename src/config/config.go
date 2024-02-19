package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Logger   LoggerConfig
	Otp      OtpConfig
	JWT      JWTConfig
	Cors     CorseConfig
}
type CorseConfig struct {
	AllowOrigins string
}
type ServerConfig struct {
	Port    int
	RunMode string
}

type JWTConfig struct {
	Secret                     string
	RefreshSecret              string
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
}

type PostgresConfig struct {
	Host            string
	User            string
	Password        string
	DbName          string
	SslMode         string
	Port            int
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host               string
	Port               int
	Password           string
	Db                 int
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	PoolTimeout        int
	IdleCheckFrequency int
}
type OtpConfig struct {
	Digits     int
	ExpireTime time.Duration
	Limiter    time.Duration
}

type LoggerConfig struct {
	FilePath string
	Encoding string
	Level    string
	Logger   string
}

func findPath(name string) string {
	switch name {
	case "dev":
		return "../config/config-development.yml"
	case "docker":
		return "../config/config-docker.yml"
	case "prod":
		return "../config/config-production.yml"
	default:
		return "../config/config-development.yml"
	}
}

func loadConfig(fileName, filePath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(fileName)
	v.SetConfigType("yml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New(fmt.Sprintf("file %s not found", fileName))
		}
		return nil, err
	}
	return v, nil
}

func parsConfig(v *viper.Viper) (cfg *Config, err error) {
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return
}

func GetConfig() *Config {
	fileName := findPath("dev")
	v, err := loadConfig(fileName, "yml")
	if err != nil {
		return nil
	}
	config, err := parsConfig(v)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
