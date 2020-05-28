CREATE SCHEMA IF NOT EXISTS book_downloader;

USE book_downloader;

CREATE TABLE BookStatus
(
    id         INT                 NOT NULL UNIQUE AUTO_INCREMENT,
    enumerator VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

INSERT INTO BookStatus (enumerator)
VALUES ('pending'),
       ('finished'),
       ('error');

CREATE TABLE Book
(
    id         INT                 NOT NULL UNIQUE AUTO_INCREMENT,
    book_key   CHAR(36)            NOT NULL UNIQUE,
    name       VARCHAR(200)        NOT NULL,
    mask       VARCHAR(200)        NOT NULL,
    status_id  INT                 NOT NULL,
    created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME                     DEFAULT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY fk_status (status_id) REFERENCES BookStatus (id)
);

