CREATE TABLE environments (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    description TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_environments_project_id_name_upcase ON environments (project_id, name_upcase);

CREATE TABLE config_file_names (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_config_file_names_project_id_name_upcase ON config_file_names (project_id, name_upcase);

CREATE TABLE config_files (
    id TEXT PRIMARY KEY,
    config_file_name_id TEXT NOT NULL,
    environment_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (config_file_name_id) REFERENCES config_file_names(id) ON DELETE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_config_files_name_env ON config_files (config_file_name_id, environment_id);

CREATE TABLE env_variable_names (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_env_variable_names_project_id_name_upcase ON env_variable_names (project_id, name_upcase);

CREATE TABLE env_variables (
    id TEXT PRIMARY KEY,
    env_variable_name_id TEXT NOT NULL,
    environment_id TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (env_variable_name_id) REFERENCES env_variable_names(id) ON DELETE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_env_variables_name_env ON env_variables (env_variable_name_id, environment_id);

CREATE TABLE secret_names (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_secret_names_project_id_name_upcase ON secret_names (project_id, name_upcase);

CREATE TABLE secrets (
    id TEXT PRIMARY KEY,
    secret_name_id TEXT NOT NULL,
    environment_id TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (secret_name_id) REFERENCES secret_names(id) ON DELETE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_secrets_name_env ON secrets (secret_name_id, environment_id);

CREATE TABLE certificate_names (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_certificate_names_project_id_name_upcase ON certificate_names (project_id, name_upcase);

CREATE TABLE certificates (
    id TEXT PRIMARY KEY,
    certificate_name_id TEXT NOT NULL,
    environment_id TEXT NOT NULL,
    certificate_data TEXT NOT NULL,
    private_key_secret TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (certificate_name_id) REFERENCES certificate_names(id) ON DELETE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_certificates_name_env ON certificates (certificate_name_id, environment_id);

CREATE TABLE ssh_key_names (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_ssh_key_names_project_id_name_upcase ON ssh_key_names (project_id, name_upcase);

CREATE TABLE ssh_keys (
    id TEXT PRIMARY KEY,
    ssh_key_name_id TEXT NOT NULL,
    environment_id TEXT NOT NULL,
    public_key TEXT NOT NULL,
    private_key_secret TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (ssh_key_name_id) REFERENCES ssh_key_names(id) ON DELETE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_ssh_keys_name_env ON ssh_keys (ssh_key_name_id, environment_id);
