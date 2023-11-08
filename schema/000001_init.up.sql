CREATE TABLE links
(
    id serial not null unique,
    original_url varchar(255) not null unique
);