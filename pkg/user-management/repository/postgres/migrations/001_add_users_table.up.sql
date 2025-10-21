CREATE TABLE IF NOT EXISTS users
(
    id          BIGSERIAL PRIMARY KEY,
    email       VARCHAR UNIQUE NOT NULL,
    username    VARCHAR        NOT NULL,
    first_name  VARCHAR,
    last_name   VARCHAR,
    role        VARCHAR        NOT NULL,
    is_active   BOOLEAN DEFAULT TRUE,
    is_deleted  BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP      NOT NULL,
    modified_at TIMESTAMP,
    CHECK (role IN ('admin', 'user', 'guest'))
);

CREATE TABLE IF NOT EXISTS user_passwords
(
    user_id       BIGINT PRIMARY KEY,
    password_hash VARCHAR NOT NULL,
    modified_at   TIMESTAMP,
    CONSTRAINT fk_user_password FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refresh_tokens
(
    id         SERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    token      TEXT UNIQUE NOT NULL,
    expires_at BIGINT      NOT NULL,
    revoked    BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP   NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);