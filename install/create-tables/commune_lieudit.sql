
create table commune_lieudit (
    id_commune              int not null references commune(id),
    id_lieudit              int not null references lieudit(id),
    primary key(id_commune, id_lieudit)
);
create index commune_lieudit_id_lieudit_idx on commune_lieudit(id_lieudit);
