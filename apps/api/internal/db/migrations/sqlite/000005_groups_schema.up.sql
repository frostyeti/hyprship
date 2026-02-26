CREATE TABLE "groups" (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    name_upcase TEXT NOT NULL UNIQUE,
    email TEXT,
    email_upcase TEXT,
    description TEXT,
    image_uri TEXT,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at INTEGER NOT NULL,
    updated_at INTEGER
);

CREATE TABLE groups_users (
    group_id TEXT NOT NULL REFERENCES "groups"(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE groups_admins (
    group_id TEXT NOT NULL REFERENCES "groups"(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE groups_roles (
    group_id TEXT NOT NULL REFERENCES "groups"(id) ON DELETE CASCADE,
    role_id TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, role_id)
);
