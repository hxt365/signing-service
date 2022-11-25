CREATE TABLE progresses
(
    id        INT AUTO_INCREMENT,
    key_id    INT UNIQUE,
    committed int,
    PRIMARY KEY (id)
);

CREATE TABLE records
(
    id    INT AUTO_INCREMENT,
    value VARCHAR(1024),
    PRIMARY KEY (id)
);

CREATE TABLE secret_keys
(
    id         INT AUTO_INCREMENT,
    identifier INT,
    value      VARCHAR(256),
    PRIMARY KEY (id)
);

CREATE TABLE signatures
(
    id        INT AUTO_INCREMENT,
    record_id INT UNIQUE,
    key_id    INT,
    value     VARCHAR(256),
    PRIMARY KEY (id)
);