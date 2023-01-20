
create table parcelle_ug (
    id_parcelle             int not null references parcelle(id),
    id_ug                   int not null references ug(id),
    primary key(id_ug, id_parcelle)
);
create index parcelle_ug_id_parcelle_idx on parcelle_ug(id_parcelle);
create index parcelle_ug_id_ug_idx on parcelle_ug(id_ug);
