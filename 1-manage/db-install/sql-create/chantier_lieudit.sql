
-- Table de liens utilisé par plaq, chautre

create table chantier_lieudit (
    type_chantier           varchar(7), -- "plaq" ou "chautre"
    id_chantier             int not null,
    id_lieudit              int not null references lieudit(id),
    primary key(type_chantier, id_chantier, id_lieudit)
);
CREATE INDEX chantier_lieudit_chantier_idx ON chantier_lieudit (type_chantier, id_chantier);
CREATE INDEX chantier_lieudit_lieudit_idx ON chantier_lieudit (type_chantier, id_lieudit);
