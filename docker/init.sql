CREATE SEQUENCE books_id_sequence;

CREATE TABLE books (
    id bigint not null primary key default nextval('books_id_sequence'),
    title varchar not null unique,
    author varchar not null
);

CREATE SEQUENCE comments_id_sequence;

CREATE TABLE comments (
    id bigint not null primary key default nextval('comments_id_sequence'),
    book_id bigint not null references books(id) on delete cascade,
    text varchar not null
);
