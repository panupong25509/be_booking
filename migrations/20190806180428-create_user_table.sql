
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id              VARCHAR(36)                 NOT NULL,
    username        VARCHAR(255)                NOT NULL,
    password        VARCHAR(255)                NOT NULL,
    fname           VARCHAR(255)                NOT NULL,
    lname           VARCHAR(255)                NOT NULL,
    organization    VARCHAR(255)                NOT NULL,
    email           VARCHAR(255)                NOT NULL,
    role            VARCHAR(255)                NOT NULL,
    create_at       datetime                    NOT NULL,
    update_at       datetime                    NOT NULL,
    PRIMARY KEY (id)
);
-- +migrate Down
DROP TABLE users;
