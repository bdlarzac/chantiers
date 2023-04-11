/*
*****************************************************************************

	MAJ de la table qgis_export, qui peut être utilisée par qgis pour afficher les chantiers

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-03-31 13:41:35+02:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	//"fmt"
)

func QGisUpdate(db *sqlx.DB) (err error) {
    
    table := "qgis_chautre"
    
	query := "drop table if exists " + table //  + " cascade"
	if _, err = db.Exec(query); err != nil {
	    return werr.Wrapf(err, "Erreur query: " + query)
	}
	//
	query = `
        create table ` + table + `(
            code_parcelle11 char(11) not null,
            titre                   varchar(255) not null,
            datecontrat             date not null,
            volumerealise           numeric not null
        )
    `
	if _, err = db.Exec(query); err != nil {
	    return werr.Wrapf(err, "Erreur query: " + query)
	}
	//
	query = `
        insert into ` + table + `(code_parcelle11, titre, datecontrat, volumerealise)
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
    `
	if _, err = db.Exec(query); err != nil {
	    return werr.Wrapf(err, "Erreur query: " + query)
	}
	//
    return nil
}

                                             

/* 
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
*/