
-- stockage = hangar à plaquettes
create table stockage (
    id                      serial primary key,
    nom                     varchar(255) not null,
    archived                boolean not null default false
);
create index stockage_archived_idx on stockage(archived);
