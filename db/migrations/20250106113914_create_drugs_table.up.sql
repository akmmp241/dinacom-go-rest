CREATE TABLE drugs
(
    id          INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    brand_name  VARCHAR(255)                NOT NULL,
    name        VARCHAR(255)                NOT NULL,
    price       INT                         NOT NULL CHECK (price >= 0),
    description TEXT                        NOT NULL,
    image_url   TEXT                        NOT NULL
);