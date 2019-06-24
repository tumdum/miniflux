drop table tokens;

create table sessions (
    id text not null,
    data jsonb not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    primary key(id)
);
