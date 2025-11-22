create table if not exists users
(
    id       bigint generated always as identity primary key,
    login    text not null unique,
    password text not null
);
