
-- chantier plaquette, champs généraux
create table plaq (
    id                      serial primary key,
    id_lieudit              int not null references lieudit(id),
    id_fermier              int not null references acteur(id),    
    id_ug                   int not null references ug(id),
    datedeb                 date not null,
    datefin                 date not null,
    surface                 numeric not null,
    granulo                 typegranulo not null,
    exploitation            typexploitation not null,
    essence                 typessence not null,
    fraisrepas              numeric,
    fraisreparation         numeric
);
create index plaq_id_ug_idx on plaq(id_ug);