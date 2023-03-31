
-- juste un test, sans doute à supprimer

create or replace view view_qgis_chautre as
    select
        c.codeinsee||p.code         as code_parcelle11,
        ch.titre                    as titre,
        ch.essence                  as essence,
        ch.datecontrat              as date,
        ch.volumerealise            as quantite,
        ch.unite                    as unite,
        ch.puht*ch.volumerealise    as prixht
    from parcelle           "p",
         commune            "c",
         chantier_parcelle  "cp",
         chautre            "ch"
    where p.id_commune = c.id
        and cp.type_chantier='chautre'
        and cp.id_chantier=ch.id
        and cp.id_parcelle=p.id
;



-- essai avec une table => marche
create table qgis_chautre(
    code_parcelle11 char(11) not null primary key,
    titre                   varchar(255) not null,
    datecontrat             date not null,
    volumerealise           numeric not null
);

insert into qgis_chautre(code_parcelle11, titre, datecontrat, volumerealise)
    select
        c.codeinsee||p.code         as code_parcelle11,
        ch.titre                    as titre,
        ch.datecontrat              as date,
        ch.volumerealise            as quantite
    from parcelle           "p",
         commune            "c",
         chantier_parcelle  "cp",
         chautre            "ch"
    where p.id_commune = c.id
        and cp.type_chantier='chautre'
        and cp.id_chantier=ch.id
        and cp.id_parcelle=p.id
;
