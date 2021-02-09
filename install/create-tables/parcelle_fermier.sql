
-- id_parcelle : IdParcelle de la base SCTL
-- id_fermier  : IdExploitant de la base SCTL

create table parcelle_fermier (
    id_parcelle             int not null references parcelle(id),
    id_fermier              int not null,
    primary key(id_parcelle, id_fermier)
);
create index parcelle_fermier_id_parcelle_idx on parcelle_fermier(id_parcelle);
create index parcelle_fermier_id_fermier_idx on parcelle_fermier(id_fermier);
