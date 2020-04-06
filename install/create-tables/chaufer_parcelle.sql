
-- lien entre chantier chauffage fermier et parcelle
-- Pour chaque parcelle, on doit préciser s'il s'agit d'une parcelle entière ou pas.
-- S'il ne s'agit pas d'une parcelle entière, il faut préciser la surface concernée par la coupe.
create table chaufer_parcelle (
    id_chaufer              int not null references chaufer(id),
    id_parcelle             int not null references parcelle(id),
    entiere                 boolean not null,
    surface                 numeric,
    primary key(id_chaufer, id_parcelle)
);
create index chaufer_parcelle_id_chaufer_idx on chaufer_parcelle(id_chaufer);
