/*
Ajoute table blocnotes
Intégration: commit 474e894

@copyright  BDL, Bois du Larzac
@license    GPL
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2023_07_21_bloc_notes(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	_, err = db.Exec(`drop table if exists blocnotes`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`create table blocnotes (contenu text not null default '')`)
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare("insert into blocnotes(contenu) VALUES($1)")
	_, err = stmt.Exec("")
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration effectuée : 2023-07-21-bloc-notes")
}
