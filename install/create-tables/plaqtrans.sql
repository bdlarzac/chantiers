
-- transport plaquette
create table plaqtrans (
    id                      serial primary key,
    id_chantier             int not null references plaq(id),
    id_tas                  int not null references tas(id),
    id_transporteur         int not null references acteur(id),
    id_conducteur           int not null references acteur(id),
    id_proprioutil          int not null references acteur(id),
    datetrans               date not null,
    qte                     numeric not null,
--    pourcentperte           numeric not null, -- différence bois sec - bois vert
    typecout                char(1) not null, -- G (global) ou C (camion) ou T (tracteur+benne)
    -- coût global
    glprix                  numeric,
    gltva                   numeric,
    gldatepay               date,
	-- coût conducteur
    conheure                numeric,              
    coprixh                 numeric,
    cotva                   numeric,
    codatepay               date,
	-- coût transport camion
	cankm                   numeric,
    caprixkm                numeric,
    catva                   numeric,
    cadatepay               date,
	-- coût transport tracteur + benne
	tbnbenne                int,
    tbduree                 numeric,
    tbprixh                 numeric,
    tbtva                   numeric,
    tbdatepay               date,
    --
    notes                   text
);
create index plaqtrans_id_chantier_idx on plaqtrans(id_chantier);
create index plaqtrans_id_transporteur_idx on plaqtrans(id_transporteur);
create index plaqtrans_id_conducteur_idx on plaqtrans(id_conducteur);
create index plaqtrans_id_proprioutil_idx on plaqtrans(id_proprioutil);
