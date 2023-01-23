/******************************************************************************

    Ajoute un lien entre les chantiers (plaquettes et autres valorisations)

    Intégration : commit

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2023-01-04 08:53:52+01:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
//	"fmt"
)

/** 
    - Copie la table chaufer_parcelle dans la table chantier_parcelle
      From https://stackoverflow.com/questions/2974057/move-data-from-one-table-to-another-postgresql-edition
    - Supprime la table chaufer_parcelle
    - Ajoute une colonne type_chantier à chantier_parcelle
    - Renomme colonne id_chaufer en id_chantier dans chantier_parcelle
    - Remplit type_chantier avec "chaufer" dans toutes les lignes de chantier_parcelle
    - Ajoute des index à la table chantier_parcelle
**/
func Migrate_2023_01_23_chantier_parcelle(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
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
	fmt.Println("Migration effectuée : 2023-01-23-chantier-parcelle")
}
