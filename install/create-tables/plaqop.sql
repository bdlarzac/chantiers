
create table plaqop (
    id                      serial primary key,
    id_chantier             int not null references plaq(id),
    id_acteur               int not null references acteur(id),
    typop                   typeop not null,
    dateop                  date not null,
    qte                     numeric not null,
    unite                   typeunite not null, -- JO, HE ou MA
    puht                    numeric not null,
    tva                     numeric not null,
    datepay                 date,
    notes                   text
);
create index plaqop_id_chantier_idx on plaqop(id_chantier);
create index plaqop_typop_idx on plaqop(typop);
