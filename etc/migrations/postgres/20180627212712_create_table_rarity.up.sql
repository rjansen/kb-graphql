--table
create table rarity (
    id varchar(40) not null,
    name varchar(255) not null,
    alias varchar(15) not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_rarity primary key(id)
);

-- index
create index ix_rarity_name on rarity (name);
create index ix_rarity_alias on rarity (alias);

-- data
insert into rarity values (encode(digest('Rarity#Common', 'sha1'), 'hex'), 'Common', 'C');
insert into rarity values (encode(digest('Rarity#Uncommon', 'sha1'), 'hex'), 'Uncommon', 'U');
insert into rarity values (encode(digest('Rarity#Rare', 'sha1'), 'hex'), 'Rare', 'R');
insert into rarity values (encode(digest('Rarity#Mythic Rare', 'sha1'), 'hex'), 'Mythic Rare', 'M');
insert into rarity values (encode(digest('Rarity#Special', 'sha1'), 'hex'), 'Special', 'S');
insert into rarity values (encode(digest('Rarity#Land', 'sha1'), 'hex'), 'Land', 'L');
insert into rarity values (encode(digest('Rarity#Promo', 'sha1'), 'hex'), 'Promo', 'P');
insert into rarity values (encode(digest('Rarity#Bonus', 'sha1'), 'hex'), 'Bonus', 'B');

