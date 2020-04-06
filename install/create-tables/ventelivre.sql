
-- Livraison de plaquettes à un client
create table ventelivre (
    id                      serial primary key,
    id_vente                int not null references venteplaq(id),
    id_livreur              int not null references acteur(id),
    datelivre               date not null,
    typecout                char(1) not null, -- G (global) ou D (détail)
    -- coût global
    glprix                  decimal,
    gltva                   decimal,
    -- coût main d'oeuvre
    monheure                decimal,
    moprixh                 decimal,
    motva                   decimal,
    --
    datepay                 date,
    notes                   text
);
create index ventelivre_id_vente_idx on ventelivre(id_vente);
create index ventelivre_id_livreur_idx on ventelivre(id_livreur);
