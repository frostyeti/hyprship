# Configuration

The application uses Viper to manage its configuration. Configuration can be provided via:
1. Environment variables
2. Configuration files (`config.yaml`)
3. Command-line flags (to be implemented)

## Environment Variables

All environment variables must be prefixed with `HYPRSHIP_` to be parsed by the application. Dots in configuration paths map to underscores `_` in environment variable names.

| Configuration Path | Environment Variable | Type   | Default Value | Description                             |
|--------------------|----------------------|--------|---------------|-----------------------------------------|
| `env`              | `HYPRSHIP_ENV`       | string | `development` | Environment mode (`development`, `prod`) |
| `port`             | `HYPRSHIP_PORT`      | int    | `8080`        | Port on which the API listens           |
| `addr`             | `HYPRSHIP_ADDR`      | string | `""`          | Address to bind to (overrides port)     |
| `log.level`        | `HYPRSHIP_LOG_LEVEL` | string | `info`        | Logging level                           |
| `log.exporter`     | `HYPRSHIP_LOG_EXPORTER` | string | `stdout`   | Log exporter (`stdout`, `stderr`, `file`, `rolling-file`, `tint`) |
| `log.format`       | `HYPRSHIP_LOG_FORMAT`   | string | `text`     | Log format (`text`, `json`) |
| `log.filename`     | `HYPRSHIP_LOG_FILENAME` | string | `""`       | Target file for `file` and `rolling-file` exporters |
| `log.max_size`     | `HYPRSHIP_LOG_MAX_SIZE` | int    | `0`        | Max file size in MB for `rolling-file` exporter |
| `log.max_backups`  | `HYPRSHIP_LOG_MAX_BACKUPS` | int | `0`        | Max backups for `rolling-file` exporter |
| `log.max_age`      | `HYPRSHIP_LOG_MAX_AGE`  | int    | `0`        | Max age in days for `rolling-file` backups |
| `log.compress`     | `HYPRSHIP_LOG_COMPRESS` | bool   | `false`    | Compress `rolling-file` backups |
| `otel.enabled`     | `HYPRSHIP_OTEL_ENABLED` | bool   | `false`    | Enable OpenTelemetry tracing and metrics |
| `otel.exporters`   | `HYPRSHIP_OTEL_EXPORTERS`| []string| `[]`      | List of OTEL exporters (`stdout`, `otlp`, `prometheus`, `zipkin`) |
| `identity.password.min_length` | `HYPRSHIP_IDENTITY_PASSWORD_MIN_LENGTH` | int | `12` | Minimum password length |
| `identity.password.max_length` | `HYPRSHIP_IDENTITY_PASSWORD_MAX_LENGTH` | int | `128` | Maximum password length |
| `identity.password.require_upper` | `HYPRSHIP_IDENTITY_PASSWORD_REQUIRE_UPPER` | bool | `true` | Require at least one uppercase letter |
| `identity.password.require_lower` | `HYPRSHIP_IDENTITY_PASSWORD_REQUIRE_LOWER` | bool | `true` | Require at least one lowercase letter |
| `identity.password.require_number` | `HYPRSHIP_IDENTITY_PASSWORD_REQUIRE_NUMBER` | bool | `true` | Require at least one number |
| `identity.password.require_symbol` | `HYPRSHIP_IDENTITY_PASSWORD_REQUIRE_SYMBOL` | bool | `true` | Require at least one symbol |
| `identity.password.history_count` | `HYPRSHIP_IDENTITY_PASSWORD_HISTORY_COUNT` | int | `5` | Number of previous passwords to remember |
| `identity.password.expire_days` | `HYPRSHIP_IDENTITY_PASSWORD_EXPIRE_DAYS` | int | `0` | Number of days before a password expires (0 = never) |
| `identity.phone.required` | `HYPRSHIP_IDENTITY_PHONE_REQUIRED` | bool | `false` | Require a phone number on registration |
| `identity.verification.email_required` | `HYPRSHIP_IDENTITY_VERIFICATION_EMAIL_REQUIRED` | bool | `true` | Require email verification before login |
| `identity.verification.phone_required` | `HYPRSHIP_IDENTITY_VERIFICATION_PHONE_REQUIRED` | bool | `false` | Require phone verification before login |
| `identity.lockout.max_attempts` | `HYPRSHIP_IDENTITY_LOCKOUT_MAX_ATTEMPTS` | int | `10` | Maximum failed login attempts before lockout |
| `identity.lockout.duration_minutes` | `HYPRSHIP_IDENTITY_LOCKOUT_DURATION_MINUTES` | int | `15` | Lockout duration in minutes |
| `identity.sessions.max_concurrent` | `HYPRSHIP_IDENTITY_SESSIONS_MAX_CONCURRENT` | int | `10` | Maximum concurrent sessions per user |
| `identity.sessions.ttl_minutes` | `HYPRSHIP_IDENTITY_SESSIONS_TTL_MINUTES` | int | `1440` | Session time-to-live in minutes (24 hours) |
| `identity.sessions.ip_required` | `HYPRSHIP_IDENTITY_SESSIONS_IP_REQUIRED` | bool | `false` | Require IP address match for session resumption |
| `identity.reset.token_ttl_minutes` | `HYPRSHIP_IDENTITY_RESET_TOKEN_TTL_MINUTES` | int | `30` | Password reset token time-to-live in minutes |

## Configuration Struct

The configuration is parsed into the following Go structure in `apps/api/config/config.go`:

```go
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
```
