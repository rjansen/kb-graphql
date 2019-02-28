-- table
create table type (
    id varchar(40) not null,
    name varchar(255) not null,
    data jsonb,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_type primary key (id)
);

-- index
create index ix_type_name on type (name);

-- data
insert into type (id, name) values (encode(digest('Type#Artifact', 'sha1'), 'hex'), 'Artifact');
insert into type (id, name) values (encode(digest('Type#Basic', 'sha1'), 'hex'), 'Basic');
insert into type (id, name) values (encode(digest('Type#Conspiracy', 'sha1'), 'hex'), 'Conspiracy');
insert into type (id, name) values (encode(digest('Type#Creature', 'sha1'), 'hex'), 'Creature');
insert into type (id, name) values (encode(digest('Type#Eaturecray', 'sha1'), 'hex'), 'Eaturecray');
insert into type (id, name) values (encode(digest('Type#Enchantment', 'sha1'), 'hex'), 'Enchantment');
insert into type (id, name) values (encode(digest('Type#Ever', 'sha1'), 'hex'), 'Ever');
insert into type (id, name) values (encode(digest('Type#Host', 'sha1'), 'hex'), 'Host');
insert into type (id, name) values (encode(digest('Type#Instant', 'sha1'), 'hex'), 'Instant');
insert into type (id, name) values (encode(digest('Type#Land', 'sha1'), 'hex'), 'Land');
insert into type (id, name) values (encode(digest('Type#Legendary', 'sha1'), 'hex'), 'Legendary');
insert into type (id, name) values (encode(digest('Type#Ongoing', 'sha1'), 'hex'), 'Ongoing');
insert into type (id, name) values (encode(digest('Type#Phenomenon', 'sha1'), 'hex'), 'Phenomenon');
insert into type (id, name) values (encode(digest('Type#Plane', 'sha1'), 'hex'), 'Plane');
insert into type (id, name) values (encode(digest('Type#Planeswalker', 'sha1'), 'hex'), 'Planeswalker');
insert into type (id, name) values (encode(digest('Type#Scariest', 'sha1'), 'hex'), 'Scariest');
insert into type (id, name) values (encode(digest('Type#Scheme', 'sha1'), 'hex'), 'Scheme');
insert into type (id, name) values (encode(digest('Type#See', 'sha1'), 'hex'), 'See');
insert into type (id, name) values (encode(digest('Type#Snow', 'sha1'), 'hex'), 'Snow');
insert into type (id, name) values (encode(digest('Type#Sorcery', 'sha1'), 'hex'), 'Sorcery');
insert into type (id, name) values (encode(digest('Type#Summon', 'sha1'), 'hex'), 'Summon');
insert into type (id, name) values (encode(digest('Type#Tribal', 'sha1'), 'hex'), 'Tribal');
insert into type (id, name) values (encode(digest('Type#Vanguard', 'sha1'), 'hex'), 'Vanguard');
insert into type (id, name) values (encode(digest('Type#World', 'sha1'), 'hex'), 'World');
insert into type (id, name) values (encode(digest('Type#You''ll', 'sha1'), 'hex'), 'You''ll');
