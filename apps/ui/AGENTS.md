# AGENTS - ui

The `@hypership/ui` project is the frontend application for hypership.

It primarily talks to the api project, but may support plugins that run
as other docker services as at later date.

## ✅ ALWAYS

## Configuration

- Use node-config

## Logging

- Use pino for logging on the server side
- Use console.log for logging on the client side

### MISE

- Use the SKILL **mise** for running mise commands.
- Use mise to run go, cast, bun and other tools found in mise.toml.
- Always start a session by reviewing what is installed by mise.

### CAST

- Use the SKILL **cast** for running or updating/adding cast tasks.
- Use cast for common project commands such as build, test, lint, fmt,
  clean, install, audit, etc

#### Cast Tasks

TODO: create tasks that mirror tasks in package.json

### Linting

- Use oxlint for formatting and linting.
- Use mise and cast to run the lint task.
- Fix all lint errors before task/feature is complete.

### Testing

**Always** write and run tests when code is changed, added, or removed before finishing
the current task.

**Always** fix
broken tests before finishing the current task.

For e2e tests spin up the api project including testcontainers for

When writing & running tests:

- Use **mise** and **cast** to run tests or **mise** and **bun**. e.g. `mise exec cast test:unit`
- Write tests using vitest.
- Write e2e tests with playwright and testcontainers for javascript.

## ⚠️ ASK FIRST

- Ask permission to add new dependencies

## 🚫 NEVER DO

- Never store secrets in the repo.
- Never violate owasp top 10.

## Techstack

Use bun, sveltekit, svelte, tailwind, chadcn-svelte, vite, vitest,
oxlint, playwright, hono, pino, node-config

### Modules

- sveltekit
- svelte
- tailwind
- chadcn-svelte
- any @frostyeti modules.
- oxlint
- vite
- playwrite

## Project Structure

The is the directory structure for the apps/hsctl project.

This project is part of a javascript workspace and ties to the package.json in the root directory. This to enable package sharing and make it
easier to invoke tasks from the root as needed.

UPDATE this section as needed.

```text

package.json
```

Shared code is placed at `packages` directory. if you are the
current directory the relative path would be `../../packages`.
