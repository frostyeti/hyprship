CREATE TABLE projects (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    name_upcase VARCHAR(256) NOT NULL,
    slug VARCHAR(256) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at DATETIME NOT NULL,
    updated_at DATETIME
);

CREATE INDEX ix_projects_name_upcase ON projects (name_upcase);
CREATE INDEX ix_projects_slug ON projects (slug);

CREATE TABLE projects_groups (
    project_id CHAR(36) NOT NULL,
    group_id CHAR(36) NOT NULL,
    permissions BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (project_id, group_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES `groups`(id) ON DELETE CASCADE
);
