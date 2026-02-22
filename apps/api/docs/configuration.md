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

## Configuration Struct

The configuration is parsed into the following Go structure in `apps/api/config/config.go`:

```go
type Config struct {
	Log  *LogSettings `mapstructure:"log"`
	Port int          `mapstructure:"port"`
	Env  string       `mapstructure:"env"`
	Addr string       `mapstructure:"addr"`
}

type LogSettings struct {
	Level string `mapstructure:"level"`
}
```
