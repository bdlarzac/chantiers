/*
*****************************************************************************

	Pour les chantiers, remplace les titres créés par le programme par des titres saisis par l'utilisateur.

	Intégration : commit 4b4bf34

	@copyright  BDL, Bois du Larzac
	@license    GPL
	@history    2023-02-20 05:43:23+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"fmt"
	"strconv"
	"time"
)

func Migrate_2023_02_21_num_facture(ctx *ctxt.Context) {
	create_table_2023_02_21(ctx)
	fill_num_vides_chautre_2023_02_21(ctx)
	fill_num_vides_venteplaq_2023_02_21(ctx)
	fmt.Println("Migration effectuée : 2023-02-21-num-facture")
}

func create_table_2023_02_21(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `
        create table facture (
            annee   char(4) not null,
            lastnum int not null
        )
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func fill_num_vides_chautre_2023_02_21(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	var query string
	type rowStruct struct {
		Id          int
		DateContrat time.Time
	}
	rows := []*rowStruct{}
	query = `select id,datecontrat from chautre where numfacture=''`
	err = db.Select(&rows, query)
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare(`update chautre set numfacture=$1 WHERE id=$2`)
	for _, row := range rows {
		annee := strconv.Itoa(row.DateContrat.Year())
		numFacture, err := model.NouveauNumeroFacture(db, annee)
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec(numFacture, row.Id)
		if err != nil {
			panic(err)
		}
	}
}

func fill_num_vides_venteplaq_2023_02_21(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	var query string
	type rowStruct struct {
		Id        int
		DateVente time.Time
	}
	rows := []*rowStruct{}
	query = `select id,datevente from venteplaq where numfacture=''`
	err = db.Select(&rows, query)
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare(`update venteplaq set numfacture=$1 WHERE id=$2`)
	for _, row := range rows {
		annee := strconv.Itoa(row.DateVente.Year())
		numFacture, err := model.NouveauNumeroFacture(db, annee)
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec(numFacture, row.Id)
		if err != nil {
			panic(err)
		}
	}

}
