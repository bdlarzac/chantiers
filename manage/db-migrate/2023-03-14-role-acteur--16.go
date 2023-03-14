/*
*
*****************************************************************************

	    https://github.com/bdlarzac/chantiers/issues/15

		Modifier la table ug et ré-écrire l'import depuis ug.csv pour avoir des données plus détaillées.
		Principale contrainte : les ugs existantes doivent conserver leurs ids actuels.

		Intégration : commit 95e55d4

		@copyright  BDL, Bois du Larzac
		@license    GPL
		@history    2023-02-24 14:43:07+01:00, Thierry Graff : Creation

*******************************************************************************
*
*/
package main

import (
//	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
//	"bdl.local/bdl/model"
	"fmt"
)

func Migrate_2023_03_14_role_acteur__16(ctx *ctxt.Context) {
	alter_table_acteur_2023_03_14(ctx)
	fmt.Println("Migration effectuée : 2023-03-14-role-acteur--16")
}

func alter_table_acteur_2023_03_14(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	queries := []string{
		`alter table acteur rename column proprietaire to isproprietaire`,
		`alter table acteur rename column fournisseur to isfournisseur`,
		`alter table acteur rename column actif to isactif`,
		//
		`alter table ug add column isclientplaq boolean`,
		`alter table ug add column isclientchauffage boolean`,
		`alter table ug add column isclientautre boolean`,
		`alter table ug add column isabatteur boolean`,
		`alter table ug add column isdebardeur boolean`,
		`alter table ug add column isdechiqueteur boolean`,
		`alter table ug add column isbroyeur boolean`,
		`alter table ug add column istransporteur boolean`,
		`alter table ug add column isproprioutil boolean`,
		`alter table ug add column is boolean`,
		`alter table ug add column is boolean`,
		`alter table ug add column is boolean`,
		`alter table ug add column is boolean`,
		`alter table ug add column is boolean`,
		`alter table ug add column is boolean`,
	}
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
}

