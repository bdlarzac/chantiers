/*
https://github.com/bdlarzac/chantiers/issues/18

Suite aux nouvelles demandes pour bilans sylviculture :
A partir de data/manage/ug.csv, ajouter dans la table ug les colonnes :
- Volume_stock_ou_recouvrement
- Intensite_prelevement
- Amenagement_divers
Profiter de la migration pour réparer les erreurs de la colonne Coupe

Intégration : commit 6cfba5d

@copyright  BDL, Bois du Larzac
@license    GPL
@history    2023-05-19 12:30:11+02:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"fmt"
	"path"
)

func Migrate_2023_05_19_ug_again__18(ctx *ctxt.Context) {
	fix_coupe_2023_05_19(ctx)
	alter_table_ug_2023_05_19(ctx)
	refill_table_ug_2023_05_19(ctx)
	fmt.Println("---\nMigration effectuée : 2023-05-19-ug-again--18")
}

func fix_coupe_2023_05_19(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	queries := []string{
		`update ug set coupe='Dégagement' where coupe='degagement'`,
		`update ug set coupe='Dégagement' where coupe='Degagement'`,
		`update ug set coupe='Dégagement' where coupe='Dgagement'`,
		`update ug set coupe='ESP1' where coupe='Ei'`,
		`update ug set coupe='ESP1' where coupe='ESP 1'`,
		`update ug set coupe='ESP1' where coupe='E1'`,
		`update ug set coupe='ESP2' where coupe='E2'`,
		`update ug set coupe='ESP3' where coupe='ESP 3'`,
	}
	for _, query := range queries {
		fmt.Println(query)
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}

}

func alter_table_ug_2023_05_19(ctx *ctxt.Context) {
	// Nouvelle structure de la table ug :
	// id                           | integer                |
	// code                         | character varying(10)  |
	// surface_sig                  | numeric                |
	// code_typo                    | character(1)           |
	// coupe                        | character varying(20)  |
	// annee_intervention           | character varying(4)   |
	// psg_suivant                  | character varying(4)   |
	// volume_stock_ou_recouvrement | character varying(255) |
	// intensite_prelevement        | character varying(255) |
	// amenagement_divers           | character varying(255) |
	db := ctx.DB
	var err error
	queries := []string{
		`alter table ug add column volume_stock_ou_recouvrement varchar(255)`,
		`alter table ug add column intensite_prelevement varchar(255)`,
		`alter table ug add column amenagement_divers varchar(255)`,
	}
	for _, query := range queries {
		fmt.Println(query)
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
}

/*
   Remplit la table ug avec les nouvelles informations.
   Utilise ug.csv du PSG 1, le même fichier qu'utilisé par install.FillUG().
   Ne gère pas les lignes mal formées non traitées par install.FillUG().
   Vérifie l'association code ug <-> id ug sur les données de prod existantes (select dans table ug)
   pour être sûr que ces associations ne changent pas.
*/

func refill_table_ug_2023_05_19(ctx *ctxt.Context) {
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
            volume_stock_ou_recouvrement,
            intensite_prelevement,
            amenagement_divers
        ) = ($1,$2,$3) where id=$4
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
		var volume_stock_ou_recouvrement, intensite_prelevement, amenagement_divers string
		volume_stock_ou_recouvrement = ug_csv["Volume_stock_ou_recouvrement"]
		intensite_prelevement = ug_csv["Intensite_prelevement"]
		amenagement_divers = ug_csv["Amenagement_divers"]
		_, err = stmt_update.Exec(
			volume_stock_ou_recouvrement,
			intensite_prelevement,
			amenagement_divers,
			idUG_db,
		)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Table ug remplie avec volume_stock_ou_recouvrement, intensite_prelevement, amenagement_divers")
}
