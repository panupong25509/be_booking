
-- +migrate Up
CREATE TABLE IF NOT EXISTS signs (
    id              INTEGER                     NOT NULL,
    sign_name       VARCHAR(255)                UNIQUE NOT NULL,
    location        VARCHAR(255)                NOT NULL,
    limitdate       VARCHAR(255)                NOT NULL,
    beforebooking   VARCHAR(255)                NOT NULL,
    picture         VARCHAR(255)                NOT NULL,
    created_at       TIMESTAMP                    NOT NULL,
    updated_at       TIMESTAMP                    NOT NULL,
    PRIMARY KEY (id)
);
-- +migrate Down
DROP TABLE signs;

