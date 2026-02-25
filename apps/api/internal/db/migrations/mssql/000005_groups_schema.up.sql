CREATE TABLE [groups] (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    name NVARCHAR(256) NOT NULL UNIQUE,
    name_upcase NVARCHAR(256) NOT NULL UNIQUE,
    email NVARCHAR(320),
    email_upcase NVARCHAR(320),
    description NVARCHAR(MAX),
    image_uri NVARCHAR(2048),
    is_active BIT NOT NULL DEFAULT 1,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2
);

CREATE TABLE groups_users (
    group_id UNIQUEIDENTIFIER NOT NULL,
    user_id UNIQUEIDENTIFIER NOT NULL,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES [groups](id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE groups_admins (
    group_id UNIQUEIDENTIFIER NOT NULL,
    user_id UNIQUEIDENTIFIER NOT NULL,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES [groups](id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) 
);

CREATE TABLE groups_roles (
    group_id UNIQUEIDENTIFIER NOT NULL,
    role_id UNIQUEIDENTIFIER NOT NULL,
    PRIMARY KEY (group_id, role_id),
    FOREIGN KEY (group_id) REFERENCES [groups](id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
