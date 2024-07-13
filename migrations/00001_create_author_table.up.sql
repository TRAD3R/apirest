CREATE SEQUENCE author_id_seq INCREMENT BY 1 MINVALUE 1 START 1;
CREATE TABLE author
(
    id          INTEGER     NOT NULL PRIMARY KEY DEFAULT nextval('author_id_seq'),
    name        VARCHAR(30) NOT NULL,
    phonenumber VARCHAR(12) NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NULL
);