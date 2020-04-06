
-- Vente de plaquettes, depuis un lieu de stockage   
create table venteplaq (
    id                      serial primary key,
    id_client               int not null references acteur(id),
    id_fournisseur          int not null references acteur(id),
    puht                    numeric not null,
    tva                     numeric not null,
    datevente               date not null,
    -- facture
    numfacture              varchar(255),
    datefacture             date,
    facturelivraison        boolean not null default false,
    facturelivraisonpuht    numeric,
    facturenotes            boolean not null default false,
    notes                   text
);
create index venteplaq_id_client_idx on venteplaq(id_client);
