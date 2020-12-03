
-- Table de liens utilisé par plaq, bspied, chautre

create table chantier_ug (
    type_chantier           varchar(7), -- "plaq" ou "bspied" ou "chautre"
    id_chantier             int not null,
    id_ug                   int not null references ug(id),
    primary key(type_chantier, id_chantier, id_ug)
);
CREATE INDEX chantier_ug_idx ON chantier_ug (type_chantier, id_chantier);
