drop table if exists users;

CREATE TABLE users (
    id  INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(128) NOT NULL,
    email VARCHAR(256) NOT NULL,
    password longtext NOT NULL,
    PRIMARY KEY(`id`)
);

