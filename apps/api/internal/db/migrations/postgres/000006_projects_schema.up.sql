CREATE TABLE projects (
    id UUID PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    name_upcase VARCHAR(256) NOT NULL,
    slug VARCHAR(256) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX ix_projects_name_upcase ON projects (name_upcase);
CREATE INDEX ix_projects_slug ON projects (slug);

CREATE TABLE projects_groups (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permissions BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (project_id, group_id)
);
