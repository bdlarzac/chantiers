
create table bspied_parcelle (
    id_bspied               int not null references bspied(id),
    id_parcelle             int not null references parcelle(id),
    primary key(id_bspied, id_parcelle)
);
create index bspied_parcelle_id_bspied_idx on bspied_parcelle(id_bspied);
