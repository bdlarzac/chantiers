
create table parcelle (
    id                      int primary key,
    id_proprietaire         int not null references acteur(id),
    code                    char(6) not null,
    surface                 numeric
);
