
-- +migrate Up
CREATE TABLE IF NOT EXISTS bookings (
    id              INTEGER                NOT NULL,
    booking_code    VARCHAR(255)       UNIQUE NOT NULL,
    description     VARCHAR(255)       NOT NULL,
    first_date      TIMESTAMP           NOT NULL,
    last_date       TIMESTAMP           NOT NULL,
    sign_id         INTEGER                NOT NULL REFERENCES signs(id),
    applicant_id    VARCHAR(36)        NOT NULL REFERENCES users(id),
    status          VARCHAR(255)       NOT NULL,
    comment         VARCHAR(255)       NOT NULL,
    created_at       TIMESTAMP           NOT NULL,
    updated_at       TIMESTAMP           NOT NULL,
    PRIMARY KEY (id)
);
-- +migrate Down
DROP TABLE bookings;