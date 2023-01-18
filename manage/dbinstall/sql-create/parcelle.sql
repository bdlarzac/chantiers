
-- id n'est pas serial (auto increment) car utilise les ids de la base SCTL
create table parcelle (
    id                      int primary key,
    id_proprietaire         int not null references acteur(id),
    code                    char(11) not null, -- = code INSEE commune + code à 6 car.
    surface                 numeric
);
