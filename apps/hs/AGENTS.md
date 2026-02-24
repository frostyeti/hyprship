# AGENTS - hs

The hsctl is the hyprship control command line app for managing the
the hypership-api.

TUI capabilities may be added using charm or other go modules.

hsctl should enable commands that mirror routes and capabilities of the
`api` project.  

It should also be able to generate the hashed value of a password
or api key for the app for the purposes putting that value into
sql directly.

The cli may add additional flags and commands to make it easier to upload/download files, edit files, bulk upload, etc.

Commands and their associated apis that are implementet **MUST** be tracked in the `apps/hs/docs/commands.md` file after they have been
implemented.

## ✅ ALWAYS

## Configuration

Enable yaml, toml, environment variables, and read environment values using `github.com/frostyeti/go/dotenv` module.  Use the viper module.

For yaml and toml files may be in the current directory, in the user's config hs config folder e.g. ~/hyprship/hs.yaml, or in the os config folder c:\ProgramData\hyprship\hs.yaml or /etc/hypership/hs.toml.

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

Primarily write integration tests for the hs app that spins up
the API project using testcontainers and postgres and then
run code against the api.

For e2e tests make exec calls using github.com/frostyeti/go/exec
that call the cli after testcontainers and postgres are spun up.

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

## Commands

- Write commands using the cobra module.
- Store, write, get configuration values using viper.

## ⚠️ ASK FIRST

- Ask permission to add new dependencies
- Prompt the user before making changes to the db design.

## 🚫 NEVER DO

- Never store secrets in the repo.

## Techstack

Use golang and gin for the API.

### Modules

- github.com/stretchr/testify/assert
- github.com/frostyeti/go/dotenv
- github.com/frostyeti/go/exec
- github.com/frostyeti/go/env
- go.yaml.in/yaml/v4
- github.com/spf13/viper
- github.com/rs/zerolog
- github.com/spf13/cobra
- github.com/go-playground/validator/v10
- github.com/xuri/excelize/v2

## Project Structure

The is the directory structure for the apps/hsctl project.

UPDATE this section as needed.

```text
cmd/             # commands
main.go          # entry point

```

Shared code is placed at `internal` directory. if you are the current directoy
the relative path would be `../../internal`.
