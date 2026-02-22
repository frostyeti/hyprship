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

## Configuration Struct

The configuration is parsed into the following Go structure in `apps/api/config/config.go`:

```go
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
```
