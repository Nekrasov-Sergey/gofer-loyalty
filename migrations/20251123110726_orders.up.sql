create table if not exists orders
(
    id          bigint generated always as identity primary key,
    number      text        not null unique,
    user_id     bigint      not null references users (id) on delete cascade,
    uploaded_at timestamptz not null
);
