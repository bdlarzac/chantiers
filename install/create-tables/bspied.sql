
-- chantier bois sur pied
create table bspied (
    id                      serial primary key,
    id_acheteur             int not null references acteur(id),
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
