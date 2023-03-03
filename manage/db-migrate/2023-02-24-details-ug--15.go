/*
*
*****************************************************************************

	    https://github.com/bdlarzac/chantiers/issues/15

		Modifier la table ug et ré-écrire l'import depuis ug.csv pour avoir des données plus détaillées.
		Principale contrainte : les ugs existantes doivent conserver leurs ids actuels.

		Intégration : commit

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
	//    create_table_essence_2023_02_24(ctx)
	//    create_table_ug_essence_2023_02_24(ctx)
	//    fill_table_essence_2023_02_24(ctx)
	fill_table_ug_essence_2023_02_24(ctx)
	fmt.Println("Migration effectuée : 2023-02-24-details-ug--15")
}

/*
*

	Remplit la table ug_essence
	en utilisant les ugs déjà en base (pour conserver les ids actuels)
	Les essences sont stockées actuellement dans la colonne type_peuplement

*
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
		//break
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
