/******************************************************************************

    Ajoute un lien entre les chantiers (plaquettes et autres valorisations)

    Intégration : commit a2d9d68

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2023-01-04 08:53:52+01:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
//	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
	"fmt"
)

/** 
    1 - Ajout des liens chantier - parcelle pour plaq et chautre => modification de chaufer_parcelle
        - Copie la table chaufer_parcelle dans la table chantier_parcelle
          From https://stackoverflow.com/questions/2974057/move-data-from-one-table-to-another-postgresql-edition
        - Supprime la table chaufer_parcelle
        - Ajoute une colonne type_chantier à chantier_parcelle
        - Renomme colonne id_chaufer en id_chantier dans chantier_parcelle
        - Remplit type_chantier avec "chaufer" dans toutes les lignes de chantier_parcelle
        - Ajoute des index à la table chantier_parcelle
    2 - Transformation de la liaison chaufer - ug, de 1-n à n-n
        - Transfère les id_ug de la table chaufer dans chantier_ug
        - Supprime la colonne chaufer.id_ug
**/
func Migrate_2023_01_23_chantier_parcelle(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	//
	// 1 - chantier_parcelle
	//
	query = `
	    CREATE TABLE chantier_parcelle AS
        WITH moved_rows AS (
            DELETE FROM chaufer_parcelle a
            RETURNING a.*
        )
        SELECT * FROM moved_rows`
	_, err = db.Exec(query)
	if err != nil {
	    panic(err)
	}
	//
	query = `DROP TABLE chaufer_parcelle`
	_, err = db.Exec(query)
	if err != nil {
	    panic(err)
	}
	//
	query = `alter table chantier_parcelle add column type_chantier varchar(7)`
	_, err = db.Exec(query)
	if err != nil {
	    panic(err)
	}
	//
	query = `alter table chantier_parcelle rename column id_chaufer to id_chantier`
	_, err = db.Exec(query)
	if err != nil {
	    panic(err)
	}
	//
	query = `update chantier_parcelle set type_chantier='chaufer'`
	_, err = db.Exec(query)
	if err != nil {
	    panic(err)
	}
	//
	query = `CREATE INDEX chantier_parcelle_chantier_idx ON chantier_parcelle (type_chantier, id_chantier)`
	_, err = db.Exec(query)
	if err != nil {
	    panic(err)
	}
	//
	query = `CREATE INDEX chantier_parcelle_parcelle_idx ON chantier_parcelle (type_chantier, id_parcelle)`
	_, err = db.Exec(query)
	if err != nil {
	    panic(err)
	}
	//
	// 2 - chaufer.id_ug
	//
	type ligne struct {
		IdChantier    int `db:"id"`
		IdUg          int `db:"id_ug"`
	}
    rows := []*ligne{}
	err = db.Select(&rows, "select id, id_ug from chaufer")
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare(`insert into chantier_ug(type_chantier,id_chantier,id_ug) values($1,$2,$3)`)
	if err != nil {
		panic(err)
	}
	for _, row := range rows {
        _, err = stmt.Exec("chaufer", row.IdChantier, row.IdUg)
        if err != nil {
            panic(err)
        }
	}
	//
	query = `alter table chaufer drop column id_ug`
	_, err = db.Exec(query)
	if err != nil {
	    panic(err)
	}
	
	//
	fmt.Println("Migration effectuée : 2023-01-23-chantier-parcelle")
}
