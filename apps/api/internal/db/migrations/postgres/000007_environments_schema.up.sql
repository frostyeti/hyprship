CREATE TABLE environments (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    name_upcase VARCHAR(256) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_environments_project_id_name_upcase ON environments (project_id, name_upcase);

CREATE TABLE config_file_names (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    name_upcase VARCHAR(256) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_config_file_names_project_id_name_upcase ON config_file_names (project_id, name_upcase);

CREATE TABLE config_files (
    id UUID PRIMARY KEY,
    config_file_name_id UUID NOT NULL REFERENCES config_file_names(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_config_files_name_env ON config_files (config_file_name_id, environment_id);

CREATE TABLE env_variable_names (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    name_upcase VARCHAR(256) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_env_variable_names_project_id_name_upcase ON env_variable_names (project_id, name_upcase);

CREATE TABLE env_variables (
    id UUID PRIMARY KEY,
    env_variable_name_id UUID NOT NULL REFERENCES env_variable_names(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_env_variables_name_env ON env_variables (env_variable_name_id, environment_id);

CREATE TABLE secret_names (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    name_upcase VARCHAR(256) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_secret_names_project_id_name_upcase ON secret_names (project_id, name_upcase);

CREATE TABLE secrets (
    id UUID PRIMARY KEY,
    secret_name_id UUID NOT NULL REFERENCES secret_names(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_secrets_name_env ON secrets (secret_name_id, environment_id);

CREATE TABLE certificate_names (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    name_upcase VARCHAR(256) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_certificate_names_project_id_name_upcase ON certificate_names (project_id, name_upcase);

CREATE TABLE certificates (
    id UUID PRIMARY KEY,
    certificate_name_id UUID NOT NULL REFERENCES certificate_names(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    certificate_data TEXT NOT NULL,
    private_key_secret TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_certificates_name_env ON certificates (certificate_name_id, environment_id);

CREATE TABLE ssh_key_names (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    name_upcase VARCHAR(256) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_ssh_key_names_project_id_name_upcase ON ssh_key_names (project_id, name_upcase);

CREATE TABLE ssh_keys (
    id UUID PRIMARY KEY,
    ssh_key_name_id UUID NOT NULL REFERENCES ssh_key_names(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    public_key TEXT NOT NULL,
    private_key_secret TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ix_ssh_keys_name_env ON ssh_keys (ssh_key_name_id, environment_id);
