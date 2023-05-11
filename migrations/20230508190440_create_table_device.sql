-- +goose Up
CREATE TABLE IF NOT EXISTS device
(
    id SERIAL PRIMARY KEY,
    num        BIGINT       NOT NULL,
    mqtt       TEXT         NOT NULL,
    invid      VARCHAR(32)  NOT NULL,
    unit_guid  uuid         NOT NULL,
    msg_id     VARCHAR(32)  NOT NULL,
    text       TEXT         NOT NULL,
    context    VARCHAR(256) NOT NULL,
    class      VARCHAR(32)  NOT NULL,
    level      INT          NOT NULL,
    area       VARCHAR(32)  NOT NULL,
    addr       VARCHAR(256) NOT NULL,
    block      boolean      NOT NULL,
    type       VARCHAR(32)  NOT NULL,
    bit        INT          NOT NULL,
    invert_bit boolean      NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS device;
