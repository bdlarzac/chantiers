
-- chantier autres valorisations
create table chautre (
    id                      serial primary key,
    id_acheteur             int not null references acteur(id),
    typevalo                typevalorisation not null,
    typevente               typevente not null,
    datecontrat             date not null,
    exploitation            typexploitation not null,
    essence                 typessence not null,
    volume                  numeric not null,
--    volumecontrat           numeric not null,
--    volumerealise           numeric not null,
    unite                   typeunite not null,
    puht                    numeric not null,
    tva                     numeric not null,
    datefacture             date,
    numfacture              varchar(255),
    notes                   text
);
create index chautre_id_acheteur_idx on chautre(id_acheteur);
