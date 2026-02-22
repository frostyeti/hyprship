---
name: db-common-type-system
description: Common type system for mapping types across different database engines like MSSQL, SQLite, PostgreSQL, etc. Used in design documents such as tickets or feature specifications to ensure consistent database design and implementation.
license: MIT
---

## When to Use this Skill

- Designing a new database schema for an application
- Refactoring an existing database schema for better consistency and maintainability
- Creating database design documentation for a project
- Establishing conventions for a team or organization on how to design databases
- Writing database-related tickets or feature specifications that require defining new tables, columns, or data types
- Generating Classes/Structs/Interfaces in application code that map to database tables and columns.

## ✅ ALWAYS DO

- Update `.agents/instructions/database-design.md` relative to the root project directory 
  when the user request new general rules for databases to be followed or
  when the agent notices common patterns.
- Update `docs/design/database-schema.md` relative to the root project directory when the user requests a database design document or when the agent notices common patterns that should be documented for future reference.
- Use a commmon type system that gets mapped to each database engine e.g. mssql, sqlite, postgres etc. 
- Convert lists that has a common type system to markdown tables for better readability and documentation.
- Add ansi sql queries to provide specific use cases for the common type system when relevant.

## Common Type System

These are common types that are mapped to types found in the database.

| Type         | SQLite  | MSSQL            | MySQL             | PostgreSQL       | Notes                                            |
|--------------|---------|------------------|-------------------|------------------|--------------------------------------------------|
| **uuid**     | TEXT    | UNIQUEIDENTIFIER | CHAR(36)          | UUID             | use uuid7                                        |
| **text**     | TEXT    | NVARCHAR         | VARCHAR           | TEXT             | size via parens e.g. `text(256)`                 |
| **bin**   | BLOB    | VARBINARY        | VARBINARY         | BYTEA            |                                                  |
| **blob**   | BLOB    | VARBINARY        | VARBINARY         | BYTEA            |                                                  |
| **boolean**  | INTEGER | BIT              | TINYINT(1)        | BOOLEAN          | 0/1 in SQLite                                    |
| **bit**      | INTEGER | BIT              | BIT               | BIT              |                                                  |
| **datetime** | INTEGER | DATETIME2        | DATETIME          | TIMESTAMP        | SQLite stores `*_at` as unix epoch (UTC)         |
| **date**     | TEXT    | DATE             | DATE              | DATE             | ISO 8601 in SQLite                               |
| **time**     | TEXT    | TIME             | TIME              | TIME             | ISO 8601 in SQLite                               |
| **i16**      | INTEGER | SMALLINT         | SMALLINT          | SMALLINT         |                                                  |
| **i32**      | INTEGER | INT              | INT               | INTEGER          |                                                  |
| **int**      | INTEGER | INT              | INT               | INTEGER          | alias for i32                                    |
| **i64**      | INTEGER | BIGINT           | BIGINT            | BIGINT           |                                                  |
| **i128**     | BLOB    | BINARY(16)       | BINARY(16)        | NUMERIC(39,0)    |                                                  |
| **u16**      | INTEGER | INT              | SMALLINT UNSIGNED | INTEGER          | no unsigned in SQLite/MSSQL/PG                   |
| **u32**      | INTEGER | BIGINT           | INT UNSIGNED      | BIGINT           |                                                  |
| **u64**      | INTEGER | DECIMAL(20,0)    | BIGINT UNSIGNED   | NUMERIC(20,0)    |                                                  |
| **json**     | TEXT    | NVARCHAR(MAX)    | JSON              | JSONB            |                                                  |
| **double**   | REAL    | FLOAT            | DOUBLE            | DOUBLE PRECISION | precision via parens e.g. `double(10,3)`         |
| **decimal**  | NUMERIC | DECIMAL          | DECIMAL           | NUMERIC          | precision/scale via parens e.g. `decimal(10,2)`  |
| **duration** | INTEGER | BIGINT           | BIGINT            | INTERVAL         | stored as milliseconds in SQLite/MSSQL/MySQL     |

If a type has a size it will use parens.  e.g. `text(256)`. 

If a type has precision and scale it will parens with a comman to delimit the other value. e.g.  `double(10, 3)`

### Common Attributes

- **pk**  primary key
- **fk** foreign **key
- **nil** or **?** nullable
- **ix** index
- **uniq**  unique
- **auto++** auto increment
- **asc** or **desc** for indexes to specify sort order
- **df()** default e.g. df(0) = default of 0.
- **\#** notes a comment
- **cascade()** means cascade on fk
  - cascase(delete) = cascade on delete
  - cascade(update) = cascade on update
  - cascade(delete,update) = cascade on delete,update.

### Field Conventions

- `*_at` date time fields
- `*_on` date fields
- `*_count` i32 field
- `is_*` bit/boolean field
- `*_digest` is only for hashed auth secrets (passwords, API keys). Use `password` or `secret` for encrypted external credentials.
- `_upcase` uppercase text used for case insensitive search

### Table Format

```markdown
| Field          | Type         | Attributes                | Notes                          |
|----------------|--------------|---------------------------|--------------------------------|
| id             | uuid         | pk                        | Primary key                    |
| name           | text(256)    | nil, uniq, ix             | Name of the entity             |
| description    | text         | nil                       | Optional description           |
| created_at     | datetime     |                           | Timestamp of creation (UTC)    |
| updated_at     | datetime     | nil                       | Timestamp of last update (UTC) |
| expires_at     | datetime     | nil                       | Timestamp of expiration (UTC)  |
| last_used_at   | datetime     | nil                       | Timestamp of last use (UTC)    |
```


extended version with target database types:

```markdown
| Field          | Common Type    | SQLite Type  | MSSQL Type         | MySQL Type      | PostgreSQL Type  | Attributes      | Notes                           |
|----------------|--------------- |--------------|--------------------|-----------------|------------------|-----------------|---------------------------------|
| id             | uuid           | TEXT         | UNIQUEIDENTIFIER   | CHAR(36)        | UUID             | pk              | Primary key                     | 
| name           | text(256)      | TEXT         | NVARCHAR(256)      | VARCHAR(256)    | TEXT             | nil, uniq       | Name of the entity              |
| description    | text           | TEXT         | NVARCHAR(MAX)      | TEXT            | TEXT             | nil             | Optional description            |
| created_at     | datetime       | INTEGER      | DATETIME2          | DATETIME        | TIMESTAMP        |                 | Timestamp of creation (UTC)     |
| updated_at     | datetime       | INTEGER      | DATETIME2          | DATETIME        | TIMESTAMP        | nil             | Timestamp of last update (UTC)  |
| expires_at     | datetime       | INTEGER      | DATETIME2          | DATETIME        | TIMESTAMP        | nil             | Timestamp of expiration (UTC)   |
| last_used_at   | datetime       | INTEGER      | DATETIME2          | DATETIME        | TIMESTAMP        | nil             | Timestamp of last use (UTC)     |
```

if you need to note additional contraints that can't easily be captured in the attributes column add it below the table with a reference to the field name. e.g.
```
- ix_table_name (use_id, name) asc 
```

## SQLite Preferences

### Timestamps

- Store all `*_at` fields as **unix epoch integers (UTC)**.
- Use INTEGER columns for timestamps (`created_at`, `updated_at`, `expires_at`, `last_used_at`, etc.).
- Compare and sort timestamps numerically in SQL.
- Convert to/from RFC3339 strings in application code.

#### Rationale

- Numeric comparisons are faster and simpler than text comparisons.
- INTEGER storage is smaller than RFC3339 TEXT.
- Avoids driver scan issues with TEXT → time.Time conversions.

#### Example

SQLite schema:

```sql
CREATE TABLE example (
    id TEXT PRIMARY KEY,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NULL
);
```

Go conversion:

```go
createdAt := time.Unix(createdAtEpoch, 0).UTC()
updatedAt := sql.NullTime{Time: time.Unix(updatedAtEpoch, 0).UTC(), Valid: true}
```
