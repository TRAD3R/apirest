-- Up
CREATE SEQUENCE author_id_seq INCREMENT BY 1 MINVALUE 1 START 1;
CREATE TABLE author
(
    id          INTEGER     NOT NULL PRIMARY KEY DEFAULT nextval('author_id_seq'),
    name        VARCHAR(30) NOT NULL,
    phonenumber VARCHAR(12) NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NULL
);
CREATE SEQUENCE post_id_seq INCREMENT BY 1 MINVALUE 1 START 1;
CREATE TABLE post
(
    id INTEGER NOT NULL PRIMARY KEY DEFAULT nextval('post_id_seq'),
    subject        VARCHAR(255),
    body TEXT DEFAULT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NULL,
    author_id      INTEGER,
    CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES author(id) ON DELETE CASCADE
);

-- Down
DROP TABLE post;
DROP SEQUENCE post_id_seq;
DROP TABLE author;
DROP SEQUENCE author_id_seq;