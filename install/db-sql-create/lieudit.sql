
-- id n'est pas serial (auto increment) car utilise IdLieuDit de la base SCTL

create table lieudit (
    id                      int primary key,
    nom                     varchar(255) not null
);
create index lieudit_nom_idx on lieudit(nom);
