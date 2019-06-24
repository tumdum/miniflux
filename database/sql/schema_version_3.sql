create table tokens (
    id text not null,
    value text not null,
    -- created_at timestamp with time zone not null default now(),
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    primary key(id, value)
);
