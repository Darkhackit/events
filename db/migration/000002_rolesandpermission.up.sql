CREATE TABLE roles (
    id SERIAL PRIMARY KEY ,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY ,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE user_roles (
    user_id INT REFERENCES users(id) ON DELETE CASCADE ,
    role_id INT REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE role_permissions (
    role_id INT REFERENCES roles(id) ON DELETE CASCADE ,
    permission_id INT REFERENCES permissions(id) ON DELETE CASCADE
);

ALTER TABLE role_permissions
    ADD CONSTRAINT role_permission_unique UNIQUE (role_id, permission_id);
