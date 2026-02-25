CREATE TABLE "groups" (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    name_upcase TEXT NOT NULL UNIQUE,
    email TEXT,
    email_upcase TEXT,
    description TEXT,
    image_uri TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

CREATE TABLE groups_users (
    group_id UUID NOT NULL REFERENCES "groups"(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE groups_admins (
    group_id UUID NOT NULL REFERENCES "groups"(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE groups_roles (
    group_id UUID NOT NULL REFERENCES "groups"(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, role_id)
);
