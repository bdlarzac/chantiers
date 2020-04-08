/******************************************************************************
    Initialisation parcelle.
    Code pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du
    @license    GPL
    @history    2019-11-07 07:54:17+01:00, Thierry Graff : Creation from a split
********************************************************************************/
package initialize

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"fmt"
	"path"
	"strconv"
)

// *********************************************************
/**
    ATTENTION : code à remplacer, parcelle doit être remplie à partir de la base Access SCTL
    Le remplissage à partir du csv est provisoire
**/
func FillParcelle() {
	table := "parcelle"
	csvname := "Parcelle_Corr.csv"
	fmt.Println("Remplit table parcelle à partir de " + csvname)
	dirCsv := getDataDir()
	filename := path.Join(dirCsv, csvname)

	records, err := tiglib.CsvMap(filename, ';')

	ctx := ctxt.NewContext()
	db := ctx.DB

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// 2 propriétaires possibles : SCTL ou GFA
	var idSCTL, idGFA int // colonne id
	query := "select id from acteur where nom=$1"
	err = db.QueryRow(query, "SCTL").Scan(&idSCTL)
	if err != nil {
		panic(err)
	}
	err = db.QueryRow(query, "GFA").Scan(&idGFA)
	if err != nil {
		panic(err)
	}

	var idProprio int
	for _, v := range records {
		surface, err := strconv.ParseFloat(v["SURFACE"], 32)
		if err != nil {
			panic(err)
		}
		surface /= 10000 // m2 -> ha
		if v["SCTL"] == "1" {
			idProprio = idSCTL
		} else {
			idProprio = idGFA
		}
		sql := fmt.Sprintf(
			"insert into %s(id,code,id_proprietaire,surface) values(%s, '%s', '%d', %7.4f)",
			table,
			v["IdParcelle"],
			v["PARCELLE"],
			idProprio,
			surface)
		if _, err = tx.Exec(sql); err != nil {
			panic(err)
		}
	}
}

// *********************************************************
func FillLiensParcelleLieudit() {
	table := "parcelle_lieudit"
	fmt.Println("Remplit table " + table + " à partir de Parcelle_Corr.csv")
	dirCsv := getDataDir()
	filename := path.Join(dirCsv, "Parcelle_Corr.csv")

	records, err := tiglib.CsvMap(filename, ';')
	if err != nil {
		panic(err)
	}

	// remove doublons
	/*
	   var k string
	   var v [2]string
	   uniques := make(map[string][2]string) // N =
	   for _, record := range records{
	       idC := record["IdParcelle"]
	       idLD := record["IdLieuDit"]
	       k = idC + "-" + idLD
	       v[0] = idC
	       v[1] = idLD
	       uniques[k] = v
	   }
	*/

	// insert db
	ctx := ctxt.NewContext()
	db := ctx.DB

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, record := range records {
		idC := record["IdParcelle"]
		idLD := record["IdLieuDit"]
		sql := fmt.Sprintf("insert into %s(id_parcelle,id_lieudit) values(%s, '%s')", table, idC, idLD)
		if _, err = tx.Exec(sql); err != nil {
			panic(err)
		}
	}

}
