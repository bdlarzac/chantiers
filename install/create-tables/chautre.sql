
-- chantier autres valorisations
create table chautre (
    id                      serial primary key,
    id_acheteur             int not null references acteur(id),
    typevalo                typevalorisation not null,
    datecontrat             date not null,
    exploitation            typexploitation not null,
    essence                 typessence not null,
    volume                  numeric not null,
    unite                   typeunite not null,
    puht                    numeric not null,
    tva                     numeric not null,
    datefacture             date,
    numfacture              varchar(255),
    notes                   text
);
create index chautre_id_client_idx on chautre(id_client);
