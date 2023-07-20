/*
	MAJ des tables utilisées par qgis pour afficher les chantiers

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-03-31 13:41:35+02:00, Thierry Graff : Creation
*/

/*
-- test mais qgis ne fonctionne pas avec les vues
create or replace view view_qgis_chautre as
    select
        c.codeinsee||p.code         as code_parcelle11,
        ch.titre                    as titre,
        ch.essence                  as essence,
        ch.datecontrat              as date,
        ch.volumerealise            as quantite,
        ch.unite                    as unite
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

package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

func QGisUpdate(db *sqlx.DB) (err error) {
	table := "qgis_chantier"
	query := "drop table if exists " + table
	if _, err = db.Exec(query); err != nil {
		return werr.Wrapf(err, "Erreur query: "+query)
	}
	//
	query = `
        create table ` + table + `(
            code_parcelle11         char(11) not null,
            titre                   varchar(255) not null,
            typechantier            varchar(7) not null,
            datechantier            date not null,
            essence                 char(2),
            quantite                numeric not null,
            unite                   char(2)
        )
    `
	if _, err = db.Exec(query); err != nil {
		return werr.Wrapf(err, "Erreur query: "+query)
	}
	//
	// plaq
	//
	query = `
        insert into ` + table + `(code_parcelle11, titre, typechantier, datechantier, essence, quantite, unite)
            select
                c.codeinsee||p.code,
                ch.titre,
                'plaq',
                ch.datedeb,
                ch.essence,
                0,
                'MA'
            from parcelle           "p",
                 commune            "c",
                 chantier_parcelle  "cp",
                 plaq               "ch"
            where p.id_commune = c.id
                and cp.type_chantier='plaq'
                and cp.id_chantier=ch.id
                and cp.id_parcelle=p.id
    `
	if _, err = db.Exec(query); err != nil {
		return werr.Wrapf(err, "Erreur query: "+query)
	}
	//
	// chautre
	//
	query = `
        insert into ` + table + `(code_parcelle11, titre, typechantier, datechantier, essence, quantite, unite)
            select
                c.codeinsee||p.code,
                ch.titre,
                'chautre',
                ch.datecontrat,
                ch.essence,
                ch.volumerealise,
                ch.unite
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
		return werr.Wrapf(err, "Erreur query: "+query)
	}
	//
	// chaufer
	//
	query = `
        insert into ` + table + `(code_parcelle11, titre, typechantier, datechantier, essence, quantite, unite)
            select
                c.codeinsee||p.code,
                ch.titre,
                'chaufer',
                ch.datechantier,
                ch.essence,
                ch.volume,
                unite
            from parcelle           "p",
                 commune            "c",
                 chantier_parcelle  "cp",
                 chaufer            "ch"
            where p.id_commune = c.id
                and cp.type_chantier='chaufer'
                and cp.id_chantier=ch.id
                and cp.id_parcelle=p.id
    `
	if _, err = db.Exec(query); err != nil {
		return werr.Wrapf(err, "Erreur query: "+query)
	}
	//
	return nil
}
