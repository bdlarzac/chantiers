/*
*****************************************************************************
    
    https://github.com/bdlarzac/chantiers/issues/15
    
	Modifier la table ug et ré-écrire l'import depuis ug.csv pour avoir des données plus détaillées.

	Intégration : commit 

	@copyright  BDL, Bois du Larzac
	@license    GPL
	@history    2023-02-24 14:43:07+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/generic/tiglib"
	"fmt"
	"path"
)

func Migrate_2023_02_24_details_ug__15(ctx *ctxt.Context) {
//    create_table_essence_2023_02_24(ctx)
//    fill_table_essence_2023_02_24(ctx)
    update_ug_2023_02_24(ctx)
	fmt.Println("Migration effectuée : 2023-02-24-details-ug--15")
}

func update_ug_2023_02_24(ctx *ctxt.Context) {
	filename := path.Join(install.GetDataDir(), "ug.csv")
    
}

func fill_table_essence_2023_02_24(ctx *ctxt.Context) {
	db := ctx.DB
	filename := path.Join(install.GetDataDir(), "essence.csv")
	essences, _ := tiglib.CsvMap(filename, ';')
	stmt, err := db.Prepare("insert into essence(code, nom, nomlong) values($1,$2,$3)")
	for _, essence := range essences {
		_, err = stmt.Exec(essence["code"], essence["nom"], essence["nomlong"])
		if err != nil {
			panic(err)
		}
	}
}

func create_table_essence_2023_02_24(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `
        create table essence (
            id          serial primary key,
            code        varchar(10) not null,
            nom         varchar(255) not null,
            nomlong     varchar(255)
        )
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `CREATE INDEX essence_idx ON essence (id);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}
