
-- chantier bois sur pied
create table bspied (
    id                      serial primary key,
    id_acheteur             int not null references acteur(id),
    id_lieudit              int not null references lieudit(id),
    id_ug                   int not null references ug(id),
    datecontrat             date not null,
    exploitation            typexploitation not null,
    essence                 typessence not null,
    nsterecontrat           numeric not null,
    nsterecoupees           numeric,
    prixstere               numeric not null,
    tva                     numeric not null,
    datefacture             date,
    numfacture              varchar(255),
    notes                   text                                                                                  
);
create index bspied_id_acheteur_idx on bspied(id_acheteur);
create index bspied_id_ug_idx on bspied(id_ug);
