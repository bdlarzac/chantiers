
create table parcelle_lieudit (
    id_parcelle             int not null references parcelle(id),
    id_lieudit              int not null references lieudit(id),
    primary key(id_parcelle, id_lieudit)
);
create index parcelle_lieudit_id_parcelle_idx on parcelle_lieudit(id_parcelle);
create index parcelle_lieudit_id_lieudit_idx on parcelle_lieudit(id_lieudit);
