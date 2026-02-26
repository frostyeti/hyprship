CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL,
    updated_at DATETIME
);

CREATE INDEX ix_projects_name_upcase ON projects (name_upcase);
CREATE INDEX ix_projects_slug ON projects (slug);

CREATE TABLE projects_groups (
    project_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    permissions INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (project_id, group_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES "groups"(id) ON DELETE CASCADE
);
