create table if not exists users
(
    id       bigint generated always as identity primary key,
    login    text not null unique,
    password text not null
);

create table if not exists sessions
(
    token      text primary key,
    user_id    bigint      not null references users (id),
    expires_at timestamptz not null
);

do
$$
    begin
        create type order_status as enum ('NEW', 'INVALID', 'PROCESSING', 'PROCESSED');
    exception
        when duplicate_object then null;
    end
$$;

create table if not exists orders
(
    id          bigint generated always as identity primary key,
    number      text           not null unique,
    status      order_status   not null,
    accrual     numeric(20, 2) not null,
    user_id     bigint         not null references users (id),
    uploaded_at timestamptz    not null
);

create table if not exists balances
(
    id        bigint generated always as identity primary key,
    current   numeric(20, 2) not null,
    withdrawn numeric(20, 2) not null,
    user_id   bigint         not null unique references users (id)
);

create table if not exists withdrawals
(
    id           bigint generated always as identity primary key,
    order_number text           not null,
    withdrawn    numeric(20, 2) not null,
    user_id      bigint         not null references users (id),
    processed_at timestamptz    not null
);
