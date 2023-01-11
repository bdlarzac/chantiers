
create table acteur (
    id                      serial primary key,
    nom                     varchar(255) not null,
    prenom                  varchar(255) not null default '',
    adresse1                varchar(255) not null default '',
    adresse2                varchar(255) not null default '',
    cp                      char(5) not null default '',
    ville                   varchar(255) not null default '',
    tel                     varchar(15) not null default '',
    mobile                  varchar(15) not null default '',
    email                   varchar(255) not null default '',
    bic                     char(11) not null default '',
    iban                    char(27) not null default '',
    siret                   char(14) not null default '',
    proprietaire            boolean not null default false,
    fournisseur             boolean not null default false, -- fournisseur de plaquettes (que BDL à priori)
    actif                   boolean not null default true,
    notes                   text not null default ''
);
create index acteur_fournisseur_idx on acteur(fournisseur);
