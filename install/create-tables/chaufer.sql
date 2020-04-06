
-- chantier chauffage fermier
create table chaufer (
    id                      serial primary key,
    id_fermier              int not null references acteur(id),
    id_ug                   int not null references ug(id),
    datechantier            date not null,
    exploitation            typexploitation not null,
    essence                 typessence not null,
    volume                  numeric not null,
    unite                   typeunite not null,
    notes                   text
);
create index chaufer_id_fermier_idx on chaufer(id_fermier);
create index chaufer_id_ug_idx on chaufer(id_ug);
