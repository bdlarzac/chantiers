/*
*****************************************************************************

	Ajout de plaq.notes

	Intégration : commit ce96026 - 2021-11-10 18:28:05

	@copyright  BDL, Bois du Larzac
	@license    GPL
	@history    2021-11-10 17:04:32+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2021_11_10_note_plaq(ctx *ctxt.Context) {
	db := ctx.DB
	query := `alter table plaq add column notes text default ''`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration effectuée : 2021-11-10-note-chantier-plaquette")
}
