
-- juste un test, sans doute à supprimer

create or replace view view_qgis_chautre as
    select
        ch.titre                    as titre,
        c.codeinsee||p.code         as code_parcelle11,
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
