
-- id n'est pas serial (auto increment) car utilise IdCommune de la base SCTL

create table commune (
    id                      int not null primary key,
    nom                     varchar(255) not null,
    nomcourt                varchar(255) not null
);
