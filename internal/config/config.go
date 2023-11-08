package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultGRPCPort               = "4545"
	defaultHTTPPort               = "8080"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
)

type (
	Config struct {
		Postgres PostgresConfig
		HTTP     HTTPConfig
		GRPC     GRPCConfig
	}

	PostgresConfig struct {
		Username string
		Password string
		Port     string
		Host     string
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	HTTPConfig struct {
		Host               string
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegaBytes int           `mapstructure:"maxHeaderMegaBytes"`
	}

	GRPCConfig struct {
		Host string
		Port string `mapstructure:"port"`
	}
)

func InitConfig(configPath string) (*Config, error) {
	setDefaults()

	if err := parseConfigFile(configPath); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("grpc", &cfg.GRPC); err != nil {
		return err
	}

	return viper.UnmarshalKey("http", &cfg.HTTP)
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.GRPC.Host = os.Getenv("GRPC_HOST")
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}

func setDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("grpc.port", defaultGRPCPort)
	viper.SetDefault("http.maxHeaderMegaBytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.readTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPRWTimeout)
}
