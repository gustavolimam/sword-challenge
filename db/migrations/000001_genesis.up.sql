-- Create roles tables
CREATE TABLE IF NOT EXISTS roles
(
    id   BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Create users table
CREATE TABLE IF NOT EXISTS users
(
    id       BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    role_id  BIGINT       NOT NULL REFERENCES roles
);

-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks
(
    id             BIGINT           NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id        BIGINT           NOT NULL REFERENCES users,
    summary        VARBINARY(10012) NOT NULL, # 2500*(max char size in UTF-8) + IV
    completed_date TIMESTAMP
);

-- Create tokens table
CREATE TABLE IF NOT EXISTS tokens
(
    uuid         VARCHAR(128) UNIQUE NOT NULL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users,
    created_date TIMESTAMP           NOT NULL
);

-- Insert initial values for users and roles
INSERT INTO roles (name)
VALUES ('tech')
ON DUPLICATE KEY UPDATE name=name;
INSERT INTO roles (name)
VALUES ('manager')
ON DUPLICATE KEY UPDATE name=name;

SET @admin_role_id = (SELECT id
                      FROM roles
                      WHERE name = 'manager'
                      LIMIT 1);
SET @tech_role_id = (SELECT id
                     FROM roles
                     WHERE name = 'tech'
                     LIMIT 1);

INSERT INTO users (username, role_id)
VALUES ('gustavo lima', @user_role_id);

INSERT INTO users (username, role_id)
VALUES ('thiago ', @admin_role_id);