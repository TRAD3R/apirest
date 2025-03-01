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