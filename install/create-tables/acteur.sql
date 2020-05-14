
-- Note : on devrait avoir id_sctl unique
-- pour que parcelle_exploitant.id_sctl_exploitant puisse référencer id_sctl comme fk
-- Mais si on met id_sctl unique, on ne peut plus insérer de nouveaux acteurs sans id_sctl.
-- Donc supprimé unique pour id_sctl et supprimé la fk dans parcelle_exploitant.
-- Mais pas satisfaisant

create table acteur (
    id                      serial primary key,
--  id_sctl                   int unique,
    id_sctl                 int not null default 0,
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
    fournisseur             boolean not null default false,
    actif                   boolean not null default true,
    notes                   text not null default ''
);
create index acteur_id_sctl_idx on acteur(id_sctl);
create index acteur_fournisseur_idx on acteur(fournisseur);
