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
	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/model"
	"fmt"
	"path"
	"strings"
)

func Migrate_2023_02_24_details_ug__15(ctx *ctxt.Context) {
	create_table_essence_2023_02_24(ctx)
	create_table_ug_essence_2023_02_24(ctx)
	fill_table_essence_2023_02_24(ctx)
	fill_table_ug_essence_2023_02_24(ctx)
	create_table_typo_2023_02_24(ctx)
	fill_table_typo_2023_02_24(ctx)
	alter_table_ug_2023_02_24(ctx)
	refill_table_ug_2023_02_24(ctx)
	fmt.Println("Migration effectuée : 2023-02-24-details-ug--15")
}

/*
   Remplit la table ug avec les nouvelles informations.
       Utilise ug.csv du PSG 1, le même fichier qu'utilisé par install.FillUG().
       Ne gère pas les lignes mal formées non traitées par install.FillUG().
   Vérifie l'association code ug <-> id ug sur les données de prod existantes (select dans table ug)
   pour être sûr que ces associations ne changent pas.
*/

func refill_table_ug_2023_02_24(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	//
	// 1 - prepare vérification association code ug => id ug
	//
	stmt_select, err := db.Preparex("select id from ug where code=$1")
	defer stmt_select.Close()
	if err != nil {
		panic(err)
	}
	//
	// 2 - prepare insert
	//
	stmt_update, err := db.Prepare(`
	    update ug set(
            code_typo,
            coupe,
            annee_intervention,
            psg_suivant
        ) = ($1,$2,$3,$4) where id=$5
    `)
	defer stmt_update.Close()
	if err != nil {
		panic(err)
	}
	//
	// 3 - boucle à partir de ug.csv
	//
	filename := path.Join(install.GetDataDir(), "ug.csv")
	ugs_csv, err := tiglib.CsvMap(filename, ',')
	if err != nil {
		panic(err)
	}
	for _, ug_csv := range ugs_csv {
		codeUG_csv := ug_csv["PG"]
		if codeUG_csv == "" || codeUG_csv == "0" {
			continue
		}
		var idUG_db int
		err = stmt_select.Get(&idUG_db, codeUG_csv)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("codeUG_csv = %s, idUG_db = %d\n",codeUG_csv, idUG_db)
		// fmt.Printf("ug_csv = %+v\n",ug_csv)
		//
		var coupe, annee_intervention, psg_suivant string
		coupe = ug_csv["Coupe"]
		if coupe == "0" {
			coupe = ""
		}
		annee_intervention = ug_csv["Annee_intervention"]
		if annee_intervention == "0" {
			annee_intervention = ""
		}
		psg_suivant = ug_csv["PSG_suivant"]
		if psg_suivant == "0" {
			psg_suivant = ""
		}
		// Pour les 2 prochains tests : les "mauvaises" lignes de ug.csv ne sont pas gérées
		// (mais n'ont pas été gérées non plus dans l'import initial
		// => pas de pb par rapport aux données de prod)
		if len(annee_intervention) > 4 {
			annee_intervention = ""
		}
		if len(psg_suivant) > 4 {
			psg_suivant = ""
		}
		//
		_, err = stmt_update.Exec(
			ug_csv["code_typo"],
			coupe,
			annee_intervention,
			psg_suivant,
			idUG_db,
		)
		if err != nil {
			panic(err)
		}
	}
}

func alter_table_ug_2023_02_24(ctx *ctxt.Context) {
	// Nouvelle structure de la table ug :
	// id                 | integer               |
	// code               | character varying(10) |
	// surface_sig        | numeric               |
	// code_typo          | character(1)          |
	// coupe              | character varying(20) |
	// annee_intervention | character varying(4)  |
	// psg_suivant        | character varying(4)  |
	db := ctx.DB
	var err error
	queries := []string{
		`alter table ug drop column type_coupe`,
		`alter table ug drop column previsionnel_coupe`,
		`alter table ug drop column type_peuplement`,
		//
		`alter table ug add column code_typo char(1)`,
		`alter table ug add column coupe varchar(20)`,
		`alter table ug add column annee_intervention varchar(4)`,
		`alter table ug add column psg_suivant varchar(4)`,
		//
		`create index ug_code_typo_idx on ug(code_typo)`,
	}
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
}

func create_table_typo_2023_02_24(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error

	query = `drop table if exists typo`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `
        create table typo (
            code        char(2) not null,
            nom         varchar(255) not null
        )
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `CREATE INDEX typo_idx ON typo(code);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func fill_table_typo_2023_02_24(ctx *ctxt.Context) {
	db := ctx.DB
	filename := path.Join(install.GetDataDir(), "typo.csv")
	typos, _ := tiglib.CsvMap(filename, ';')
	stmt, err := db.Prepare("insert into typo(code, nom) values($1,$2)")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	for _, typo := range typos {
		_, err = stmt.Exec(typo["code"], typo["nom"])
		if err != nil {
			panic(err)
		}
	}
}

/*
Remplit la table ug_essence en utilisant les ugs déjà en base (pour conserver les ids actuels)
Les essences sont stockées actuellement dans la colonne type_peuplement
*/
func fill_table_ug_essence_2023_02_24(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	// map contenu de la colonne type_peuplement => code essence stocké en base
	// Selon Jean Culié, les autres codes sont des erreurs de saisie
	codeMap := map[string]string{
		"Al":                  "AL",
		"Al pars":             "AL",
		"Al. ss tage":         "AL",
		"CH":                  "CN",
		"CH pars":             "CN",
		"CHa":                 "CT",
		"Cd":                  "CD",
		"Cdre":                "CD",
		"Ch":                  "CN",
		"Cha":                 "CT",
		"D pars":              "DG",
		"Douglas":             "DG",
		"Er":                  "ER",
		"Er Ch":               "ER", // Verifier si = érable ou érable champêtre
		"Er de Mpt":           "EM",
		"Er de Mtp":           "EM",
		"Er pars":             "ER",
		"F":                   "FL",
		"F divers":            "FL",
		"H":                   "HT",
		"H pars":              "HT",
		"PN":                  "PN",
		"PN pars":             "PN",
		"PS":                  "PS",
		"PS pars":             "PS",
		"Plantation feuillus": "FL",
	}
	ugs := []*model.UG{}
	err = db.Select(&ugs, `select * from ug`)
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare("insert into ug_essence(id_ug,code_essence,epars) values($1,$2,$3)")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	for _, ug := range ugs {
		strEssences := ug.TypePeuplement
		essences := strings.Split(strEssences, "+")
		for _, ess := range essences {
			ess = strings.TrimSpace(ess)
			codeEss, ok := codeMap[ess]
			if !ok {
				continue
			}
			epars := false
			if strings.Index(ess, "pars") != -1 {
				epars = true
			}
			_, err = stmt.Exec(ug.Id, codeEss, epars)
			if err != nil {
				panic(err)
			}
		}
	}
}

func fill_table_essence_2023_02_24(ctx *ctxt.Context) {
	db := ctx.DB
	filename := path.Join(install.GetDataDir(), "essence.csv")
	essences, _ := tiglib.CsvMap(filename, ';')
	stmt, err := db.Prepare("insert into essence(code, nom, nomlong) values($1,$2,$3)")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	for _, essence := range essences {
		_, err = stmt.Exec(essence["code"], essence["nom"], essence["nomlong"])
		if err != nil {
			panic(err)
		}
	}
}

func create_table_ug_essence_2023_02_24(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `drop table if exists ug_essence`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `
        create table ug_essence (
            id_ug           int not null,
            code_essence    char(2) not null,
            epars           boolean not null default false
        )
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `CREATE INDEX ug_essence_ug_idx ON ug_essence(id_ug);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `CREATE INDEX ug_essence_essence_idx ON ug_essence(code_essence);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func create_table_essence_2023_02_24(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `drop table if exists essence`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `
        create table essence (
            code        char(2) not null,
            nom         varchar(255) not null,
            nomlong     varchar(255)
        )
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `CREATE INDEX essence_idx ON essence(code);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}
