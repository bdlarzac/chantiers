/*
	    Ajouter date paiement dans chautre et venteplaq

	    Bizarre, si on met pas de défaut dans les colone date, le scan donne une erreur
	    "sql: Scan error on column index" (...) "unsupported Scan, storing driver.Value type" (...)

		Voir https://github.com/bdlarzac/chantiers/issues/22
		Intégration : commit

		@copyright  BDL, Bois du Larzac
		@license    GPL
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2023_05_26_date_paiement__22(ctx *ctxt.Context) {
	alterTableChautre_2023_05_26(ctx)
	alterTableVentePlaq_2023_05_26(ctx)
	fmt.Println("Migration effectuée : 2023-05-26-date-paiement--22")
}

func alterTableVentePlaq_2023_05_26(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `alter table venteplaq add column datepaiement date not null default '0001-01-01'`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func alterTableChautre_2023_05_26(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `alter table chautre add column datepaiement date not null default '0001-01-01'`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}
