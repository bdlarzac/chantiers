
-- transport plaquette
create table plaqtrans (
    id                      serial primary key,
    id_chantier             int not null references plaq(id),
    id_tas                  int not null references tas(id),
    id_transporteur         int not null references acteur(id),
    datetrans               date not null,
    qte                     numeric not null,
    typecout                char(1) not null, -- G (global) ou C (camion) ou T (tracteur+benne)
    -- coût global
    glprix                  numeric,
    gltva                   numeric,
    gldatepay               date,
	-- coût transporteur
    trnheure                numeric,              
    trprixh                 numeric,
    trtva                   numeric,
    trdatepay               date,
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
