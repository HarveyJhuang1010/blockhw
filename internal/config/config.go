package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var self Config

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return path.Join(filepath.Dir(d), "..")
}

func NewConfig() *Config {
	viper.AddConfigPath(path.Join(rootDir(), "config"))
	viper.AddConfigPath("/blockhw/config")
	viper.SetConfigName("common")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("failed to read config", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	if err := viper.Unmarshal(&self); err != nil {
		panic(errors.Wrap(err, "failed to marshal config"))
	}

	return &self
}

type Config struct {
	Database *DatabaseConfig `mapstructure:"database"`
	Rest     *RestConfig     `mapstructure:"rest"`
	Ethereum *Ethereum       `mapstructure:"ethereum"`
	Worker   *WorkerConfig   `mapstructure:"worker"`
	Nats     *NatsConfig     `mapstructure:"nats"`
}

func GetConfig() *Config {
	return &self
}

type RestConfig struct {
	ListenAddress             string   `mapstructure:"listen_address"`
	ListenPort                int      `mapstructure:"listen_port" validate:"gte=1,lte=65535"`
	AllowOrigins              []string `mapstructure:"allow_origins"`
	AllowHeaders              []string `mapstructure:"allow_headers"`
	ExposeHeaders             []string `mapstructure:"expose_headers"`
	AllowMethods              []string `mapstructure:"allow_methods"`
	RateLimitIntervalSeconds  int      `mapstructure:"rate_limit_interval_seconds"`
	RateLimitRequestPerSecond int      `mapstructure:"rate_limit_requests_per_second"`
}

type DatabaseConfig struct {
	Dialect      string `mapstructure:"dialect"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host" validate:"ipv4"`
	Port         int    `mapstructure:"port" validate:"gte=1,lte=65535"`
	DBName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifetime  int    `mapstructure:"max_life_time"`
}

func (cfg *DatabaseConfig) IsMySQL() bool {
	return strings.EqualFold(cfg.Dialect, "mysql")
}

func (cfg *DatabaseConfig) IsPostgreSQL() bool {
	return strings.EqualFold(cfg.Dialect, "postgres") ||
		strings.EqualFold(cfg.Dialect, "postgresql")
}

func (cfg *DatabaseConfig) Open() gorm.Dialector {
	if cfg.IsMySQL() {
		return mysql.Open(cfg.DSN())
	}

	if cfg.IsPostgreSQL() {
		return postgres.New(postgres.Config{
			DSN:                  cfg.DSN(),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		})
	}

	panic("no match dialect")
}

func (cfg *DatabaseConfig) DSN() string {
	if cfg.IsMySQL() {
		//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&time_zone=UTC",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		)
	}

	if cfg.IsPostgreSQL() {
		//dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			cfg.Host,
			cfg.User,
			cfg.Password,
			cfg.DBName,
			cfg.Port,
		)
	}

	return ""
}

type Ethereum struct {
	Endpoint string `mapstructure:"endpoint"`
}

type WorkerConfig struct {
	StartNumber uint64 `mapstructure:"start_number"`
	DelayMinute int    `mapstructure:"delay_minute"`
	MaxWorkers  int64  `mapstructure:"max_workers"`
}

type NatsConfig struct {
	Host       string  `mapstructure:"host" validate:"ipv4"`
	ClientPort int     `mapstructure:"client_port" validate:"gte=1,lte=65535"`
	Password   *string `mapstructure:"password"`
}

func (cfg *NatsConfig) GetURL() string {
	return fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.ClientPort)
}
