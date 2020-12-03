
-- chantier plaquette, champs généraux
create table plaq (
    id                      serial primary key,
    datedeb                 date not null,
    datefin                 date not null,
    surface                 numeric not null,
    granulo                 typegranulo not null,
    exploitation            typexploitation not null,
    essence                 typessence not null,
    fraisrepas              numeric,
    fraisreparation         numeric
);
