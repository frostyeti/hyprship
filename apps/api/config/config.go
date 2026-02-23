package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Log  *LogSettings `mapstructure:"log"`
	Otel *OtelConfig  `mapstructure:"otel"`
	Port int          `mapstructure:"port"`
	Env  string       `mapstructure:"env"`
	Addr string       `mapstructure:"addr"`
}

type OtelConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	Exporters []string `mapstructure:"exporters"`
}

type LogSettings struct {
	Level      string `mapstructure:"level"`
	Exporter   string `mapstructure:"exporter"`
	Format     string `mapstructure:"format"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

func LoadConfig() (*Config, error) {
	// Setup flags
	pflag.Int("port", 8080, "Port to listen on")
	pflag.String("env", "development", "Environment (e.g., development, production)")
	pflag.String("addr", "", "Address to bind to")
	pflag.String("log.level", "info", "Log level")
	pflag.String("log.exporter", "stdout", "Log exporter")
	pflag.String("log.format", "text", "Log format")
	pflag.Bool("otel.enabled", false, "Enable OpenTelemetry")
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
	viper.SetDefault("log.exporter", "stdout")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("otel.enabled", false)

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
