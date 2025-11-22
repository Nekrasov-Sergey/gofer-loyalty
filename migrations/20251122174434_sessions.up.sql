create table if not exists sessions
(
    token      text primary key,
    user_id    bigint      not null references users (id) on delete cascade,
    expires_at timestamptz not null
);
