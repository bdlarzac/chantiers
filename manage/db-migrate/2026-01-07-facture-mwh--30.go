/*
	    Ajoute les champs pour gérer type de facture map / mwh

	    Intégration: commit 

		@copyright  BDL, Bois du Larzac
		@license    GPL
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2026_01_07_facture_mwh__30(ctx *ctxt.Context) {
	db := ctx.DB
	queries := []string{
	    "alter table venteplaq add column typefacture char(2)            -- 'MA' (map) ou 'MW' (Mwh)",
	    "alter table venteplaq add column facture_mwh_nb numeric not null default 0",
	    "update venteplaq set typefacture='MA'",
	}
	for _, query := range queries {
		stmt, err := db.Prepare(query)
		defer stmt.Close()
		_, err = stmt.Exec()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Migration effectuée : 2026-01-07-facture-mwh")
}
