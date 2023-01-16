
-- mesure d'humidité
create table humid (
    id                      serial primary key,
    id_tas                  int not null references tas(id),
    valeur                  numeric not null,
    datemesure              date not null,
    notes                   text
);
