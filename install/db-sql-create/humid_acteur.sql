
create table humid_acteur (
    id_humid                int not null references humid(id),
    id_acteur               int not null references acteur(id),
    primary key(id_humid, id_acteur)
);
create index humid_acteur_id_humid_idx on humid_acteur(id_humid);
create index humid_acteur_id_acteur_idx on humid_acteur(id_acteur);
