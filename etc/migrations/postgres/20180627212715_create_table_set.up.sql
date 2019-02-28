-- table
create table set (
    id varchar(40) not null,
    name varchar(255) not null,
    alias varchar(15) not null,
    asset jsonb,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_set primary key (id)
);

-- index
create index ix_set_name on set (name);
create index ix_set_alias on set (alias);
