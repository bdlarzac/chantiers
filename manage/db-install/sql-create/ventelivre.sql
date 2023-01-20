
-- Livraison de plaquettes à un client
create table ventelivre (
    id                      serial primary key,
    id_vente                int not null references venteplaq(id),
    id_livreur              int not null references acteur(id),
    id_conducteur           int not null references acteur(id),
    id_proprioutil          int not null references acteur(id),
    datelivre               date not null,
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
create index ventelivre_id_vente_idx on ventelivre(id_vente);
create index ventelivre_id_livreur_idx on ventelivre(id_livreur);
create index ventelivre_id_conducteur_idx on ventelivre(id_conducteur);
create index ventelivre_id_proprioutil_idx on ventelivre(id_proprioutil);
