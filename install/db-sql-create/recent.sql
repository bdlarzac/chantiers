
-- Pour stocker les dernières pages visitées
-- Table contenant toujours max nb-recent lignes (cf config)
-- pas d'index car table toujours très petite
create table recent (
    url                     varchar(255) not null unique,
    label                   varchar(255) not null,
    --dateVisite              date not null
    dateVisite              timestamp not null
);
