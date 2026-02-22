package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Log  *LogSettings `mapstructure:"log"`
	Port int          `mapstructure:"port"`
	Env  string       `mapstructure:"env"`
	Addr string       `mapstructure:"addr"`
}

type LogSettings struct {
	Level string `mapstructure:"level"`
}

func LoadConfig() (*Config, error) {
	// Setup flags
	pflag.Int("port", 8080, "Port to listen on")
	pflag.String("env", "development", "Environment (e.g., development, production)")
	pflag.String("addr", "", "Address to bind to")
	pflag.String("log.level", "info", "Log level")
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvPrefix("HYPRSHIP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	// Removed SetConfigType so it can automatically detect json, toml, yaml, etc based on extension
	viper.AddConfigPath(".")
	viper.AddConfigPath("./apps/api")

	// Set defaults
	viper.SetDefault("port", 8080)
	viper.SetDefault("env", "development")
	viper.SetDefault("addr", "")
	viper.SetDefault("log.level", "info")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
