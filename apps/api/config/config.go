package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Log      *LogSettings    `mapstructure:"log"`
	Otel     *OtelConfig     `mapstructure:"otel"`
	Identity *IdentityConfig `mapstructure:"identity"`
	Port     int             `mapstructure:"port"`
	Env      string          `mapstructure:"env"`
	Addr     string          `mapstructure:"addr"`
}

type IdentityConfig struct {
	Password     PasswordConfig     `mapstructure:"password"`
	Phone        PhoneConfig        `mapstructure:"phone"`
	Verification VerificationConfig `mapstructure:"verification"`
	Lockout      LockoutConfig      `mapstructure:"lockout"`
	Sessions     SessionsConfig     `mapstructure:"sessions"`
	Reset        ResetConfig        `mapstructure:"reset"`
}

type PasswordConfig struct {
	MinLength     int  `mapstructure:"min_length"`
	MaxLength     int  `mapstructure:"max_length"`
	RequireUpper  bool `mapstructure:"require_upper"`
	RequireLower  bool `mapstructure:"require_lower"`
	RequireNumber bool `mapstructure:"require_number"`
	RequireSymbol bool `mapstructure:"require_symbol"`
	HistoryCount  int  `mapstructure:"history_count"`
	ExpireDays    int  `mapstructure:"expire_days"`
}

type PhoneConfig struct {
	Required bool `mapstructure:"required"`
}

type VerificationConfig struct {
	EmailRequired bool `mapstructure:"email_required"`
	PhoneRequired bool `mapstructure:"phone_required"`
}

type LockoutConfig struct {
	MaxAttempts     int `mapstructure:"max_attempts"`
	DurationMinutes int `mapstructure:"duration_minutes"`
}

type SessionsConfig struct {
	MaxConcurrent int  `mapstructure:"max_concurrent"`
	TTLMinutes    int  `mapstructure:"ttl_minutes"`
	IPRequired    bool `mapstructure:"ip_required"`
}

type ResetConfig struct {
	TokenTTLMinutes int `mapstructure:"token_ttl_minutes"`
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

	_ = viper.BindPFlags(pflag.CommandLine)

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

	// Identity defaults
	viper.SetDefault("identity.password.min_length", 12)
	viper.SetDefault("identity.password.max_length", 128)
	viper.SetDefault("identity.password.require_upper", true)
	viper.SetDefault("identity.password.require_lower", true)
	viper.SetDefault("identity.password.require_number", true)
	viper.SetDefault("identity.password.require_symbol", true)
	viper.SetDefault("identity.password.history_count", 5)
	viper.SetDefault("identity.password.expire_days", 0)
	viper.SetDefault("identity.phone.required", false)
	viper.SetDefault("identity.verification.email_required", true)
	viper.SetDefault("identity.verification.phone_required", false)
	viper.SetDefault("identity.lockout.max_attempts", 10)
	viper.SetDefault("identity.lockout.duration_minutes", 15)
	viper.SetDefault("identity.sessions.max_concurrent", 10)
	viper.SetDefault("identity.sessions.ttl_minutes", 1440)
	viper.SetDefault("identity.sessions.ip_required", false)
	viper.SetDefault("identity.reset.token_ttl_minutes", 30)

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
