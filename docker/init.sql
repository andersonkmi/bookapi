CREATE SEQUENCE books_id_sequence;

CREATE TABLE books (
    id bigint not null primary key default nextval('books_id_sequence'),
    title varchar not null unique,
    author varchar not null
);
