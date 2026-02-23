# AGENTS - api

## ✅ ALWAYS DO

### MISE

- Use the SKILL **mise** for running mise commands.
- Use mise to run go, cast, bun and other tools found in mise.toml.
- Always start a session by reviewing what is installed by mise.

### CAST

- Use the SKILL **cast** for running or updating/adding cast tasks.
- Use cast for common project commands such as build, test, lint, fmt,
  clean, install, audit, etc

#### Cast Tasks

- **build**
- **test:unit**
- **test:integration**
- **test:e2e**
- **test:all**
- **lint**
- **fmt**
- **clean**
- **install**
- **tidy** 
- **audit** - scan for vulnerabilities

### Interfaces

- Use the SKILL `golang-interfaces` to writing and implement interfaces.

### Linting

- Use mise and cast to run the lint task.
- Fix all lint errors before task/feature is complete.

### Testing

**Always** use the skill **golang-testing** when writing, running, or updating tests.

**Always** write and run tests when code is changed, added, or removed before finishing
the current task.

**Always** fix
broken tests before finishing the current task.

When writing & running tests:

- Use **mise** and **cast** to run tests. e.g. `mise exec cast test:unit`
- Write integration tests in a seperate file and use the `// +build integration` directive.
- Write e2e tests for the api under `apps/api/test/e2e` and use the `// +build e2e` directive.
- Write api unit tests using `net/http/httptest`.
- Write .http files for testing/triaging issues in the `apps/api/test/http` folder.
- Use asserts from github.com/stretchr/testify/assert
- Generate stubs or fakes instead of mocks as needed.
- Use testcontainers for go for integration tests for services like
  databases, redis, search, open telemetry recievers, etc.

### Database Design

- Use the skill **db-common-type-system** when working with database schema.
- Always use snakecase for database identifiers.
- Always use pluralized names for collections such as tables and views.
- Support implementations for SQLITE, POSTGRESS, MYSQL, and MSSQL.
- Update the `apps/api/docks/database-schema.md` when the schema is updated. Use the
  db-common-type-system. Include, tables, views, indexes & constraints.
- Use the db common type system when generating a plan.
- Translate the db common type system to the correct db driver type when implemeting the plan.
- Provide the user with recommendations for the db schema design.
- Provide the user with recommendations for updates to the db-common-type-system skill as needed.
- Provide the user with recommendations for missing indexes and constraints.
- Follow the field name conventions found in the db-common-type-system.
- Write integrations tests against the db implementation for go.
- With the exception of sqlite, use testcontainers for writing db integration tests.

### Respostories

Use the name `stores` for repositories and store for `respository`.  e.g. `UserStore` instead of `UserRespository`.

Use the golang-interfaces skill for designing interfaces for repositories.

### Configuration

**Always** update `apps/api/docs/configuration.md` when adding, changing, or removing configuration settings in `apps/api/config/config.go`.
Use the prefix `HYPRSHIP` for all environment variables with `viper`.

### API Design

**Always** use the skill `rest-api-design` when working on REST APIS, including planning
and implementing rest apis.

**Always** update the doc `apps/api/docs/routes.md` when routes are added, changed, removed.

**Always** return an envelope response.

```golang
type ApiError struct {
   Code string
   Target *string
   Message string
   InnerError *ApiError
   Details []ApiErrorDetail
}

//  used for form field validation errors 
type ApiErrorDetail struct {
  // in this case the code is name of the field
  Code *string 
  Message string 
  Target string
}

type Response struct {
   Ok bool  `json:"ok"`
   Error *Error `json:"error, omitempty"`
   TraceId *string
   SpanId *string
   Value any
   NextLink *string `json:"@nextLink, omitempty"`
}
```

- implement paging for collections form the rest-api-design skill.
- implement filtering for collections from the *rest-api-design* skill.
- implement sorting for collections from the *rest-api-design* skill.
- implement expand when all properties are not returned, but optionally may be returned.
- Always use rate limiting.
- Write unit and e2e tests for api endpoints.
- For collections enable batch operations for bulk imports and exports.
- Imports and exports must support json, yaml, excel and csv.
  - csv imports/exports that allow more than one colleciton must support zip files
    with multiple csv files.
  - excel exports with more than one collection must used one sheet per collection.
  - json/yaml with more than one collection must have a top level property per collection.
  - imports must be atomic and errors including the row number and must be provided for each error.

## ⚠️ ASK FIRST

- Ask permission to add new dependencies

## 🚫 NEVER DO

- Never store secrets in the repo.

## Techstack

Use golang and gin for the API.  UPDATE the list of modules as needed

### Modules

- github.com/stretchr/testify/assert
- github.com/frostyeti/go/dotenv
- github.com/frostyeti/go/exec
- github.com/frostyeti/go/env
- go.yaml.in/yaml/v4
- github.com/spf13/viper
- log/slog
- github.com/samber/slog-gin
- github.com/lmittmann/tint
- github.com/open-telemetry/opentelemetry-go
- github.com/gin-gonic/gin
- net/http/httptest
- golang.org/x/time/rate

### Logging (slog)

- Always use the structured logging package `log/slog`.
- Avoid high-allocation debug or info logging statements when not needed. **Always** check if the level is enabled before logging high cost statements:
  ```go
  if logger.Enabled(ctx, slog.LevelDebug) {
      slog.Debug("costly debug message", "data", getCostlyData())
  }
  ```
- Use `slog.Group` to structure related log attributes.

### OpenTelemetry (Telemetry)

- Only trace and measure when OpenTelemetry is enabled (`cfg.Otel.Enabled == true`).
- Ensure tracing or metrics spans are only initialized if an exporter that supports it is configured. (e.g. `prometheus` does not support tracing for Go Otel).
- Add OpenTelemetry collectors/instrumentation for databases and datastores as they are implemented.
- Follow the OpenTelemetry Semantic Conventions (see [semconv](https://opentelemetry.io/docs/concepts/semantic-conventions/)).

### Cryptography

- **Hashing:** Use the PBKDF2 hasher in `apps/api/crypto` for internal application passwords, API keys, or any secrets that only need to be verified and never retrieved in plaintext. A good hint for the database schema is that any column ending with `_digest` should use the hasher.
- **Encryption:** Use the AES-GCM encryption driver in `apps/api/crypto` for secrets that are stored for external apps and systems. These secrets must be encrypted at rest so they can be decrypted when sent to the service.

## Project Structure

The is the directory structure for the apps/api project.

UPDATE this section as needed.

```text
claims/                  # consts for claim names
  permissions/           # consts for permissions
config/                  # vipr configuration
internal/
  core/                  # shared models and interfaces
  crypto/                # hashing and encryption drivers
  db/                    # database schema and tools
    actions/
    migrations/
    seed/
    migrator.go
    seeder.go
  models/                # data models
  routes/                # api routes may go here
    v1/                  # version 1 routes go 
  validators/            # validators
  stores/                # repositories
    mssql/               # microsoft sql server impl for stores
    mysql/               # mysql/maria db impl for stores
    pg/                  # postgres db impl for stores
    sqlite/              # sqlite db impl for stores
  svc/                   # services may go here or in child folders
  telemetry/             # logging and otel goes here
  test/                  # integration and e2e testing tools
```

Shared code is placed at `../../internal`. If you are in the
project root then it is just the `internal` folder.
