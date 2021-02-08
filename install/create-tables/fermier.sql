
-- Représente une exploitation agricole
-- Table remplie à partir de la base SCTL
create table fermier (
    id                      serial primary key,
    id_sctl                 int not null default 0,
    nom                     varchar(255) not null,
    prenom                  varchar(255) not null default '',
    adresse                 varchar(255) not null default '',
    cp                      char(5) not null default '',
    ville                   varchar(255) not null default '',
    tel                     varchar(15) not null default '',
    email                   varchar(255) not null default ''
);
create index fermier_id_sctl_idx on fermier(id_sctl);
