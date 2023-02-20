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
	"fmt"
)

func Migrate_2023_02_20_titre_chantier(ctx *ctxt.Context) {
//    alter_tables_2023_02_20(ctx)
//    update_plaq_2023_02_20(ctx)
//    update_chautre_2023_02_20(ctx)
    update_chaufer_2023_02_20(ctx)
	fmt.Println("Migration effectuée : 2023-02-20-titre-chantier")
}

func alter_tables_2023_02_20(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `alter table plaq add column titre varchar(255) not null default ''`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `alter table chautre add column titre varchar(255) not null default ''`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `alter table chaufer add column titre varchar(255) not null default ''`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

// Met des valeurs par défaut aux titres

func update_plaq_2023_02_20(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	ids := []int{}
	err = db.Select(&ids, "select id from plaq")
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare(`update plaq set titre=$1 WHERE id=$2`)
	for _, id := range(ids){
        ch, err := model.GetPlaq(db, id)
        if err != nil {
            panic(err)
        }
        ch.ComputeLieudits(db)
        _, err = stmt.Exec(ch.String(), id)
        if err != nil {
            panic(err)
        }
	}
}

func update_chautre_2023_02_20(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	ids := []int{}
	err = db.Select(&ids, "select id from chautre")
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare(`update chautre set titre=$1 WHERE id=$2`)
	for _, id := range(ids){
        ch, err := model.GetChautre(db, id)
        if err != nil {
            panic(err)
        }
        ch.ComputeAcheteur(db)
        _, err = stmt.Exec(ch.String(), id)
        if err != nil {
            panic(err)
        }
	}
}

func update_chaufer_2023_02_20(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	ids := []int{}
	err = db.Select(&ids, "select id from chaufer")
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare(`update chaufer set titre=$1 WHERE id=$2`)
	for _, id := range(ids){
        ch, err := model.GetChaufer(db, id)
        if err != nil {
            panic(err)
        }
        ch.ComputeFermier(db)
        _, err = stmt.Exec(ch.String(), id)
        if err != nil {
            panic(err)
        }
	}
}

