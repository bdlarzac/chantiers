
-- Rangement plaquettes dans un hangar à plaquettes
create table plaqrange (
    id                      serial primary key,
    id_chantier             int not null references plaq(id),
    id_tas                  int not null references tas(id),
    id_rangeur              int not null references acteur(id),
    id_conducteur           int not null references acteur(id),
    id_proprioutil          int not null references acteur(id),
    daterange               date not null,
    typecout                char(1) not null, -- G (global) ou D (détail)
    -- coût global
    glprix                  numeric,
    gltva                   numeric,
    gldatepay               date,
    -- concerne le conducteur
    conheure                numeric,
    coprixh                 numeric,
    cotva                   numeric,
    codatepay               date,
    -- concerne l'outil
    ouprix                  numeric,
    outva                   numeric,
    oudatepay               date,
    notes                   text
);
create index plaqrange_id_chantier_idx on plaqrange(id_chantier);
create index plaqrange_id_rangeur_idx on plaqrange(id_rangeur);
create index plaqrange_id_conducteur_idx on plaqrange(id_conducteur);
create index plaqrange_id_proprioutil_idx on plaqrange(id_proprioutil);
