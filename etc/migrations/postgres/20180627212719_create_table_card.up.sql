-- table
create table card (
    id varchar(40) not null,
    name varchar(255) not null,
    types varchar(128)[] not null,
    costs varchar(5)[] not null default '{}',
    number_cost numeric(14,4) not null default 0,
    id_external varchar(100) not null,
    id_rarity varchar(40) not null,
    id_set varchar(40) not null,
    id_asset varchar(40) not null,
    rate numeric(14,4) not null default 0,
    rate_votes numeric(10) not null default 0,
    rules varchar(512)[] not null default '{}',
    order_external varchar(15),
    artist varchar(255),
    flavor varchar(1024),
    data jsonb,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_card primary key (id)
);

-- index
create index ix_card_name on card (name);
create index ix_card_idexternal on card (id_external);
create index ix_card_orderexternal on card (order_external);
create index ix_card_type on card (types);

-- foreign key
alter table card add constraint
      fk_card_set foreign key (id_set) references set (id);
alter table card add constraint
      fk_card_rarity foreign key (id_rarity) references rarity (id);
