CREATE TABLE environments (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    project_id UNIQUEIDENTIFIER NOT NULL,
    name NVARCHAR(256) NOT NULL,
    name_upcase NVARCHAR(256) NOT NULL,
    description NVARCHAR(MAX),
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_environments_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_environments_project_id_name_upcase ON environments (project_id, name_upcase);

CREATE TABLE config_file_names (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    project_id UNIQUEIDENTIFIER NOT NULL,
    name NVARCHAR(256) NOT NULL,
    name_upcase NVARCHAR(256) NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_config_file_names_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_config_file_names_project_id_name_upcase ON config_file_names (project_id, name_upcase);

CREATE TABLE config_files (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    config_file_name_id UNIQUEIDENTIFIER NOT NULL,
    environment_id UNIQUEIDENTIFIER NOT NULL,
    content NVARCHAR(MAX) NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_config_files_name FOREIGN KEY (config_file_name_id) REFERENCES config_file_names(id) ON DELETE CASCADE,
    CONSTRAINT fk_config_files_env FOREIGN KEY (environment_id) REFERENCES environments(id)
);

CREATE UNIQUE INDEX ix_config_files_name_env ON config_files (config_file_name_id, environment_id);

CREATE TABLE env_variable_names (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    project_id UNIQUEIDENTIFIER NOT NULL,
    name NVARCHAR(256) NOT NULL,
    name_upcase NVARCHAR(256) NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_env_variable_names_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_env_variable_names_project_id_name_upcase ON env_variable_names (project_id, name_upcase);

CREATE TABLE env_variables (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    env_variable_name_id UNIQUEIDENTIFIER NOT NULL,
    environment_id UNIQUEIDENTIFIER NOT NULL,
    value NVARCHAR(MAX) NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_env_variables_name FOREIGN KEY (env_variable_name_id) REFERENCES env_variable_names(id) ON DELETE CASCADE,
    CONSTRAINT fk_env_variables_env FOREIGN KEY (environment_id) REFERENCES environments(id)
);

CREATE UNIQUE INDEX ix_env_variables_name_env ON env_variables (env_variable_name_id, environment_id);

CREATE TABLE secret_names (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    project_id UNIQUEIDENTIFIER NOT NULL,
    name NVARCHAR(256) NOT NULL,
    name_upcase NVARCHAR(256) NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_secret_names_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_secret_names_project_id_name_upcase ON secret_names (project_id, name_upcase);

CREATE TABLE secrets (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    secret_name_id UNIQUEIDENTIFIER NOT NULL,
    environment_id UNIQUEIDENTIFIER NOT NULL,
    value NVARCHAR(MAX) NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_secrets_name FOREIGN KEY (secret_name_id) REFERENCES secret_names(id) ON DELETE CASCADE,
    CONSTRAINT fk_secrets_env FOREIGN KEY (environment_id) REFERENCES environments(id)
);

CREATE UNIQUE INDEX ix_secrets_name_env ON secrets (secret_name_id, environment_id);

CREATE TABLE certificate_names (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    project_id UNIQUEIDENTIFIER NOT NULL,
    name NVARCHAR(256) NOT NULL,
    name_upcase NVARCHAR(256) NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_certificate_names_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_certificate_names_project_id_name_upcase ON certificate_names (project_id, name_upcase);

CREATE TABLE certificates (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    certificate_name_id UNIQUEIDENTIFIER NOT NULL,
    environment_id UNIQUEIDENTIFIER NOT NULL,
    certificate_data NVARCHAR(MAX) NOT NULL,
    private_key_secret NVARCHAR(MAX),
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_certificates_name FOREIGN KEY (certificate_name_id) REFERENCES certificate_names(id) ON DELETE CASCADE,
    CONSTRAINT fk_certificates_env FOREIGN KEY (environment_id) REFERENCES environments(id)
);

CREATE UNIQUE INDEX ix_certificates_name_env ON certificates (certificate_name_id, environment_id);

CREATE TABLE ssh_key_names (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    project_id UNIQUEIDENTIFIER NOT NULL,
    name NVARCHAR(256) NOT NULL,
    name_upcase NVARCHAR(256) NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_ssh_key_names_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_ssh_key_names_project_id_name_upcase ON ssh_key_names (project_id, name_upcase);

CREATE TABLE ssh_keys (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    ssh_key_name_id UNIQUEIDENTIFIER NOT NULL,
    environment_id UNIQUEIDENTIFIER NOT NULL,
    public_key NVARCHAR(MAX) NOT NULL,
    private_key_secret NVARCHAR(MAX) NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    CONSTRAINT fk_ssh_keys_name FOREIGN KEY (ssh_key_name_id) REFERENCES ssh_key_names(id) ON DELETE CASCADE,
    CONSTRAINT fk_ssh_keys_env FOREIGN KEY (environment_id) REFERENCES environments(id)
);

CREATE UNIQUE INDEX ix_ssh_keys_name_env ON ssh_keys (ssh_key_name_id, environment_id);
