
-- Table de liens utilisé par plaq, bspied, chautre

create table chantier_fermier (
    type_chantier           varchar(7), -- "plaq" ou "bspied" ou "chautre"
    id_chantier             int not null,
    id_fermier              int not null references fermier(id),
    primary key(type_chantier, id_chantier, id_fermier)
);
CREATE INDEX chantier_fermier_chantier_idx ON chantier_fermier (type_chantier, id_chantier);
CREATE INDEX chantier_fermier_fermier_idx ON chantier_fermier (type_chantier, id_fermier);
