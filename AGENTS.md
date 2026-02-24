# Agents

Hyprship enables start ups or home labs to ship code faster
and manage infra.

## ✅ Always DO

- Use short concise language.
- Prefer bullet points in order of important over paragraphs.
- When a design such as a database table is missing field or can
  be greatly improved, suggest an updated design and prompt the user for input.
- Run linters and formatters before completing the task and fix lint errors.
- Run tests before completing the task and fix tests.  
- Format commits following the Conventional Commits specification.
  - Structure: `<type>[optional scope]: <description>`
  - Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
  - Scopes: Indicate project component e.g. `feat(api)`, `fix(ui)`
  - Link GitHub issues in the description or footer (e.g. `#1`, `Resolves #1`, `Fixes #1`)
  

## ⚠️ ASK FIRST

- Ask permission before adding new dependencies unless dependencies are included
  as part of the tech stack.

## 🚫 Never DO

- expose or commit secrets in the repository.

## High Level Project Structure

```text
apps/                              # the directory where the application is stored.
  api/                             # backend api
    main.go                        # main go file for backendapi


  hs/                              # hyprship cli util. calls api and other utilities
    cmd/                           # commands
      users.go                     # subcommands(s) for user
    main.go                        # main entry point

  ui/                              # frontend/ui goes here
    hs-ui/                         # folder for site
    libs/                          # package folder for shared 
    AGENTS.md                      # ui context
    package.json                   # private package.json for monorepo
    tsconfig.json                  # ts config

   docs/                           # astro/starlight docs site
backlog                            # the folder that is used by the backlog.md tool
 archive/                          # archived tasks
 completed/                        # completed tasks
 decisions/                        # decision records
 docs/                             # documents
   database-schema.md              # keep the database schema documentation here 
   routes.md                       # keep routes documentation here.
 drafts/                           # drafts
 milestones/                       # milestones
 tasks/                            # tasks
 config.yml                        # configuration for backlog-md
eng/                               # engineering folder. for scripts, ci/cd stuff, etc.
  scripts/                         # scripts
  bin/                             # scripts that should be available on the path
  compose/                         # compose projects
    api/                           # api compose project
      compose.yaml                 # compose file
      .env                         # dot env file for variables
      castfile                     # castfile should include up/down commands
    mssql/
      compose.yaml
      .env
      castfile 
    mysql/
      compose.yaml
      .env
      castfile 
    pg/
      compose.yaml
      .env
      castfile 
    sqlite/
      compose.yaml
      .env
      castfile 
    ui/                             #
      compose.yaml
      .env
      castfile 
internal/                          # internal packages
  routes/                          # routes folder
    auth/                          # auth routes
    users/                         # user routes
    .../                           # create route folders as needed
    shared.go                      # shared route logic
  validation/                      # validation
  stores/                          # stores.  use stores/store instead of the name repositories/
    mssql/                         # mssql impl folder
    mysql/                         # pg impl folder
    pg/                            # pg impl folder
    sqlite/                        # sqlite impl folder
    users.go                    
    roles.go
  db/                              # related things like migrations, seeds, etc.
    migrations/                    # migrations.
      mssql/
      mysql/
      pg/
      sqlite/
    seeds/
      users.yaml   
      roles.yaml
    migrator.go
    seeder.go
AGENTS.md
castfile
go.mod
mise.toml
README.md
```

## API project

The api project's is found at `apps/api` and should be built as hyprship-api.



