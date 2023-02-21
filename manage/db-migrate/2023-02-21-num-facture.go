/******************************************************************************

    Pour les chantiers, remplace les titres créés par le programme par des titres saisis par l'utilisateur.

    Intégration : commit

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2023-02-20 05:43:23+01:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	// "bdl.local/bdl/generic/tiglib"
	"fmt"
	// "strings"
)

func Migrate_2023_02_21_num_facture(ctx *ctxt.Context) {
    
test(ctx)
    
//    create_table_2023_02_21(ctx)
	fmt.Println("Migration effectuée : 2023-02-21-num-facture")
}

func test(ctx *ctxt.Context){
    _, _ = model.NouveauNumeroAffacture(ctx.DB, "2023")
}

func create_table_2023_02_21(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `
        create table affacture (
            annee   char(4) not null,
            lastnum int not null
        )
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}