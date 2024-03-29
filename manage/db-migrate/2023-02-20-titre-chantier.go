/*
Pour les chantiers, remplace les titres créés par le programme par des titres saisis par l'utilisateur.

Intégration : commit

@copyright  BDL, Bois du Larzac
@license    GPL
@history    2023-02-20 05:43:23+01:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/model"
	"fmt"
	"strings"
)

func Migrate_2023_02_20_titre_chantier(ctx *ctxt.Context) {
	alter_tables_2023_02_20(ctx)
	update_plaq_2023_02_20(ctx)
	update_chautre_2023_02_20(ctx)
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

//
// Met des valeurs par défaut aux titres - plaq
//

func update_plaq_2023_02_20(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	ids := []int{}
	err = db.Select(&ids, "select id from plaq")
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare(`update plaq set titre=$1 WHERE id=$2`)
	for _, id := range ids {
		ch, err := model.GetPlaq(db, id)
		if err != nil {
			panic(err)
		}
		ch.ComputeLieudits(db)
		_, err = stmt.Exec(plaq_titreParDefaut_2023_02_20(ch), id)
		if err != nil {
			panic(err)
		}
	}
}

// Ancienne fonction plaq.String()
func plaq_titreParDefaut_2023_02_20(ch *model.Plaq) string {
	if len(ch.Lieudits) == 0 {
		panic("Erreur dans le code - Les lieux-dits d'un chantier plaquettes doivent être calculés avant d'appeler TitreParDefaut()")
	}
	res := ""
	var noms []string
	for _, ld := range ch.Lieudits {
		noms = append(noms, ld.Nom)
	}
	res += strings.Join(noms, " - ")
	res += " " + tiglib.DateFr(ch.DateDebut)
	return res
}

//
// Met des valeurs par défaut aux titres - chautre
//

func update_chautre_2023_02_20(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	ids := []int{}
	err = db.Select(&ids, "select id from chautre")
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare(`update chautre set titre=$1 WHERE id=$2`)
	for _, id := range ids {
		ch, err := model.GetChautre(db, id)
		if err != nil {
			panic(err)
		}
		ch.ComputeAcheteur(db)
		_, err = stmt.Exec(chautre_titreParDefaut_2023_02_20(ch), id)
		if err != nil {
			panic(err)
		}
	}
}

// Ancienne fonction chautre.String()
func chautre_titreParDefaut_2023_02_20(ch *model.Chautre) string {
	if ch.Acheteur == nil {
		panic("Erreur dans le code - L'acheteur d'un chantier autre valorisation doit être calculé avant d'appeler String()")
	}
	return model.ValoMap[ch.TypeValo] + " " + ch.Acheteur.String() + " " + tiglib.DateFr(ch.DateContrat)
}

//
// Met des valeurs par défaut aux titres - chaufer
//

func update_chaufer_2023_02_20(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	ids := []int{}
	err = db.Select(&ids, "select id from chaufer")
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare(`update chaufer set titre=$1 WHERE id=$2`)
	for _, id := range ids {
		ch, err := model.GetChaufer(db, id)
		if err != nil {
			panic(err)
		}
		ch.ComputeFermier(db)
		_, err = stmt.Exec(chaufer_titreParDefaut_2023_02_20(ch), id)
		if err != nil {
			panic(err)
		}
	}
}

// Ancienne fonction chaufer.String()
func chaufer_titreParDefaut_2023_02_20(ch *model.Chaufer) string {
	if ch.Fermier == nil {
		panic("Erreur dans le code - Le fermier d'un chantier chauffage fermier doit être calculé avant d'appeler String()")
	}
	return ch.Fermier.String() + " " + tiglib.DateFr(ch.DateChantier)
}
