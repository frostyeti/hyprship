# Issue #2: Setup Logging and OpenTelemetry for API

This plan outlines the implementation steps to add robust structured logging (using `slog` and `slog-gin`) and OpenTelemetry (metrics and tracing) to the `apps/api` project, resolving Issue #2.

## Phase 1: Modules and Dependencies
Add the following Go dependencies to `apps/api`:
- `log/slog` (standard library, Go 1.21+)
- `github.com/lmittmann/tint`
- `github.com/mattn/go-colorable`
- `github.com/mattn/go-isatty`
- `gopkg.in/natefinch/lumberjack.v2`
- `github.com/samber/slog-gin`
- `go.opentelemetry.io/otel` and related SDK/exporter packages (`go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin`, `go.opentelemetry.io/otel/exporters/stdout/stdouttrace`, `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp`, etc.)

## Phase 2: Configuration Enhancements
Update `apps/api/config/config.go` to support logging and OpenTelemetry settings.

1. **Update `LogSettings` struct:**
```go
type LogSettings struct {
    Level      string `mapstructure:"level"`
    Exporter   string `mapstructure:"exporter"` // stdout, stderr, file, rolling-file, tint
    Format     string `mapstructure:"format"`   // text, json
    Filename   string `mapstructure:"filename"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxBackups int    `mapstructure:"max_backups"`
    MaxAge     int    `mapstructure:"max_age"`
    Compress   bool   `mapstructure:"compress"`
}
```
2. **Add `OtelConfig` struct:**
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
    Exporters []string `mapstructure:"exporters"` // stdout, otlp, prometheus, zipkin
}
```
3. Initialize sensible defaults (`Log.Exporter="stdout"`, `Log.Format="text"`, `Otel.Enabled=false`).

## Phase 3: Logging Implementation
Create a new logging package/initializer (`apps/api/svc/logger.go` or similar) to initialize the default `slog.Logger`.

1. **Configure Slog Handler based on `Log.Exporter`:**
   - `stdout` / `stderr`: Use `os.Stdout` or `os.Stderr`. If `config.Env` is "dev" or "development", promote to `tint`.
   - `tint`: Use `tint.NewHandler` with `go-colorable` and `go-isatty` to ensure Windows colored log support. Format config is ignored.
   - `file`: Write directly to the specified `Log.Filename`.
   - `rolling-file`: Configure `lumberjack.Logger` using `MaxSize`, `MaxBackups`, `MaxAge`, and `Compress`. Write to this logger.
2. **Configure Format (if not `tint`):**
   - Use `slog.NewJSONHandler` if `Format` is "json".
   - Use `slog.NewTextHandler` if `Format` is "text".
3. **Gin Middleware:**
   - Integrate `github.com/samber/slog-gin` as middleware in the gin router configuration.

## Phase 4: OpenTelemetry (Otel) Implementation
Create an OTEL initialization routine (`apps/api/svc/telemetry.go` or similar).

1. **Service Configuration:**
   - Service name: `hyprship-api`.
2. **Exporters & Providers:**
   - Read `Otel.Exporters` list.
   - For traces: Initialize exporters based on values (`stdout`, `otlp`, `zipkin`). (Note: `prometheus` doesn't support tracing in Go).
   - For metrics: Initialize exporters (`prometheus`, `stdout`, `otlp`).
3. **Suggestions for Otel Sampling (Requested by Issue):**
   - Use `ParentBased(AlwaysSample)` for development environments.
   - Use `ParentBased(TraceIDRatioBased(0.1))` for high-traffic production environments to reduce overhead and costs by sampling 10% of new traces.
   - Allow configuration of this ratio via `OtelConfig.SamplingRatio` in future phases.
4. **Gin Middleware:**
   - Add `otelgin` middleware to the Gin engine to auto-trace HTTP endpoints.

## Phase 5: Documentation & Guidelines Updates
1. **Update `AGENTS.md`**:
   - **Logging Section:** Describe how to use `slog`. Add strict guidance to wrap high-allocation debug/info logs in `if logger.Enabled(ctx, slog.LevelDebug) { ... }` checks for performance.
   - **Telemetry Section:** Provide constraints: only create metrics/traces when OTEL is enabled and the relevant exporter is configured (e.g., don't initialize tracing if only `prometheus` is selected). List semantic conventions links. Emphasize adding standard OTEL collectors for databases when new stores are added.
2. **Update `apps/api/docs/configuration.md`**:
   - Add documentation detailing all `LogSettings` and `OtelConfig` environment variables and properties (e.g., `HYPRSHIP_LOG_EXPORTER`, `HYPRSHIP_OTEL_ENABLED`).
