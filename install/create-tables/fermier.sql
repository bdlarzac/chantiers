
-- Représente une exploitation agricole
-- Table remplie à partir de la base SCTL

create table fermier (
    id                      int not null unique default 0, -- IdExploitant de la base SCTL
    nom                     varchar(255) not null,
    prenom                  varchar(255) not null default '',
    adresse                 varchar(255) not null default '',
    cp                      char(5) not null default '',
    ville                   varchar(255) not null default '',
    tel                     varchar(15) not null default '',
    email                   varchar(255) not null default ''
);
