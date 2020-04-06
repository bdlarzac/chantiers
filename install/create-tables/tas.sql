
-- tas dans un hangar à plaquette
create table tas (
    id                      serial primary key,
    id_stockage             int not null references stockage(id),
    id_chantier             int not null references plaq(id),
    stock                   numeric,
    actif                   boolean not null default true
);
create index tas_id_chantier_idx on tas(id_chantier);
create index tas_actif_idx on tas(actif);
-- pas mis index sur id_stockage parce que à chaque fois
-- utilisé dans un actif and id_stockage=$1 => test d'abord sur actif
-- create index tas_id_stockage_idx on tas(id_stockage);
