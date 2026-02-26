CREATE TABLE projects (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    name NVARCHAR(256) NOT NULL,
    name_upcase NVARCHAR(256) NOT NULL,
    slug NVARCHAR(256) NOT NULL UNIQUE,
    description NVARCHAR(MAX),
    is_active BIT NOT NULL DEFAULT 1,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2
);

CREATE INDEX ix_projects_name_upcase ON projects (name_upcase);
CREATE INDEX ix_projects_slug ON projects (slug);

CREATE TABLE projects_groups (
    project_id UNIQUEIDENTIFIER NOT NULL,
    group_id UNIQUEIDENTIFIER NOT NULL,
    permissions BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (project_id, group_id),
    CONSTRAINT fk_projects_groups_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_projects_groups_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);
