/*
*****************************************************************************

	Change les codes parcelle : passe de code à 6 caractères (ex 0C0001) à code à 11 caractères (ex 120820C0001)

	Voir https://github.com/bdlarzac/chantiers/issues/11
	Intégration : commit

	@copyright  BDL, Bois du Larzac
	@license    GPL
	@history    2023-01-16 15:47:07+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package main

import (
	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"fmt"
	"path"
	"strconv"
)

func Migrate_2023_01_16_fix_parcelle(ctx *ctxt.Context) {
	versionSCTL := "2023-01-11"
	alterTableParcelle_2023_01_16(ctx)
	fillTableParcelle_2023_01_16(ctx, versionSCTL)
	alterTableCommune_2023_01_16(ctx)
	fillTableCommune_2023_01_16(ctx, versionSCTL)
	fillLiensParcelleUG_2023_01_16(ctx, versionSCTL)
	fmt.Println("Migration effectuée : 2023-01-16-fix-parcelle")
}

// Fonctions auxiliaires, commencent par des minuscules

func alterTableParcelle_2023_01_16(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `alter table parcelle add column id_commune int`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `alter table parcelle add constraint parcelle_id_commune_fkey foreign key(id_commune) references commune(id) not valid`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func fillTableParcelle_2023_01_16(ctx *ctxt.Context, versionSCTL string) {
	var dirCsv, filename, csvname string
	var err error
	// Charge parcelles
	csvname = "Parcelle.csv"
	dirCsv = install.GetSCTLDataDir(versionSCTL)
	filename = path.Join(dirCsv, csvname)
	parcelles, _ := tiglib.CsvMap(filename, ';')
	// update table parcelle
	db := ctx.DB
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare("update parcelle set id_commune=$1 where id=$2")
	for _, parcelle := range parcelles {
		idCommune, _ := strconv.Atoi(parcelle["IdCommune"])
		idParcelle, _ := strconv.Atoi(parcelle["IdParcelle"])
		_, err = stmt.Exec(idCommune, idParcelle)
		if err != nil {
			panic(err)
		}
	}
}

func alterTableCommune_2023_01_16(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `alter table commune add column codeinsee char(5) not null default ''`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `alter table commune alter column codeinsee drop default`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func fillTableCommune_2023_01_16(ctx *ctxt.Context, versionSCTL string) {
	var dirCsv, filename, csvname string
	csvname = "commune.csv"
	dirCsv = install.GetDataDir()
	filename = path.Join(dirCsv, csvname)
	communes, _ := tiglib.CsvMap(filename, ';')
	db := ctx.DB
	stmt, err := db.Prepare("update commune set codeinsee=$1 where id=$2")
	for _, commune := range communes {
		_, err = stmt.Exec(commune["code_insee"], commune["id"])
		if err != nil {
			panic(err)
		}
	}
}

func fillLiensParcelleUG_2023_01_16(ctx *ctxt.Context, versionSCTL string) {
	var dirCsv, filename, csvname string
	var err error
	db := ctx.DB
	var query string

	// Prépare code parcelle 11 => id parcelle
	parcelle11_id := make(map[string]int, 2351) // 2351 = nb de parcelles en base
	query = "select concat(c.codeinsee, p.code) as code11, p.id as idParcelle from commune as c,parcelle as p where c.id=p.id_commune"
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var idParcelle int
	var code11 string
	for rows.Next() {
		err := rows.Scan(&code11, &idParcelle)
		if err != nil {
			panic(err)
		}
		parcelle11_id[code11] = idParcelle
	}

	// Prépare code ug => id ug
	ug_code_id := make(map[string]int, 659) // 659 = nb de ugs en base
	query = "select id,code from ug"
	rows, err = db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var idUG int
	var codeUG string
	for rows.Next() {
		err := rows.Scan(&idUG, &codeUG)
		if err != nil {
			panic(err)
		}
		ug_code_id[codeUG] = idUG
	}

	// Charge UGs
	csvname = "ug.csv"
	dirCsv = install.GetDataDir()
	filename = path.Join(dirCsv, csvname)
	ugs, err := tiglib.CsvMap(filename, ',') // []map[string]string
	if err != nil {
		panic(err)
	}

	// Vide table parcelle_ug
	query = `truncate table parcelle_ug`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	// Remplit trable parcelle_ug en utilisant code à 11 caractères
	stmt, err := db.Prepare("insert into parcelle_ug(id_parcelle, id_ug) values($1,$2)")
	for _, ug := range ugs {
		codeUG = ug["PG"]
		if codeUG == "" || codeUG == "0" {
			continue
		}
		code11 = ug["ID_PARCELLE_11"]
		idUG = ug_code_id[codeUG]
		idParcelle = parcelle11_id[code11]
		_, err = stmt.Exec(idParcelle, idUG)
		if err != nil {
			//panic(err)
			//fmt.Printf("ERREUR insert parcelle_ug : codeUG = %s - code11 = %s - idParcelle=%d\n",codeUG, code11, idParcelle)
		}
	}
}
