
-- Chrgement de plaquettes pour une livraison
create table ventecharge (
    id                      serial primary key,
    id_livraison            int not null references ventelivre(id),
    id_chargeur             int not null references acteur(id),
    id_conducteur           int not null references acteur(id),
    id_proprioutil          int not null references acteur(id),
    id_tas                  int not null references tas(id),
    qte                     decimal not null,
    datecharge              date not null,    
    typecout                char(1) not null, -- G (global) ou D (détail)
    -- coût global
    glprix                  decimal,
    gltva                   decimal,
    gldatepay               date,
    -- coût détaillé, outil
    ouprix                  decimal,
    outva                   decimal,
    oudatepay               date,
    -- coût détaillé, main d'oeuvre
    monheure                decimal,
    moprixh                 decimal,
    motva                   decimal,
    modatepay               date,
    --
    notes                   text
);
create index ventecharge_id_livraison_idx on ventecharge(id_livraison);
create index ventecharge_id_tas_idx on ventecharge(id_tas);
create index ventecharge_id_chargeur_idx on ventecharge(id_chargeur);
create index ventecharge_id_conducteur_idx on ventecharge(id_conducteur);
create index ventecharge_id_proprioutil_idx on ventecharge(id_proprioutil);
