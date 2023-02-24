/*
*****************************************************************************

	Initialisation parcelle.
	Code servant à initialiser la base, pas utilisé en fonctionnement normal.

	@copyright  BDL, Bois du Larzac
	@license    GPL
	@history    2019-11-07 07:54:17+01:00, Thierry Graff : Creation from a split

*******************************************************************************
*/
package install

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

// *********************************************************
// @param   versionSCTL ex "2020-12-23" - voir commentaire de install-bdl.go
func FillParcelle(ctx *ctxt.Context, versionSCTL string) {
	table := "parcelle"
	csvname := "Parcelle.csv"
	fmt.Println("Remplit table parcelle à partir de " + csvname)

	dirCsv := GetSCTLDataDir(versionSCTL)
	filename := path.Join(dirCsv, csvname)

	records, err := tiglib.CsvMap(filename, ';')

	db := ctx.DB

	// 2 propriétaires possibles : SCTL ou GFA
	// ATTENTION : les ids de SCTL et GFA sont récupérés à partir du nom
	// Si le nom change, ce code plante
	// Voir appli/install/dbcreate/acteur.go, AddActeurBDL() et AddActeurSCTL()
	var idSCTL, idGFA int // colonne id
	query := "select id from acteur where nom=$1"
	err = db.QueryRow(query, "SCTL").Scan(&idSCTL)
	if err != nil {
		panic(err)
	}
	err = db.QueryRow(query, "GFA Larzac").Scan(&idGFA)
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
			v["PARCELLE11"],
			idProprio,
			surface)
		if _, err = db.Exec(sql); err != nil {
			panic(err)
		}
	}
}

// *********************************************************
// @param   versionSCTL ex "2020-12-23" - voir commentaire de install-bdl.go
func FillLiensParcelleLieudit(ctx *ctxt.Context, versionSCTL string) {
	table := "parcelle_lieudit"
	fmt.Println("Remplit table " + table + " à partir de Parcelle.csv")

	dirCsv := GetSCTLDataDir(versionSCTL)
	filename := path.Join(dirCsv, "Parcelle.csv")

	// pour lire le csv, on ne peut pas utiliser du code générique tiglib.CsvMap()
	// car certaines lignes contiennent des \n (champ Observation de la table)
	// => à ignorer
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	var records []map[string]string
	var tmp []string
	var nCols int
	var colNames []string
	i := -1
	for scanner.Scan() {
		line := scanner.Text()
		i++
		if i == 0 {
			// récupère le nom des colonnes
			colNames = strings.Split(line, ";")
			nCols = len(colNames)
			continue
		}
		// remplit une ligne de données
		tmp = strings.Split(line, ";")
		if len(tmp) != nCols {
			// ligne foireuse, la ligne précédente contient un \n dans le champ Observation
			continue
		}
		var record = make(map[string]string, nCols)
		for idx, field := range colNames {
			record[field] = tmp[idx]
		}
		records = append(records, record)
	}

	// insert db
	db := ctx.DB
	sql := fmt.Sprintf("truncate table %s", table)
	if _, err = db.Exec(sql); err != nil {
		panic(err)
	}
	for _, record := range records {
		idP := record["IdParcelle"]
		idLD := record["IdLieuDit"]
		sql := fmt.Sprintf("insert into %s(id_parcelle,id_lieudit) values(%s, '%s')", table, idP, idLD)
		if _, err = db.Exec(sql); err != nil {
			panic(err)
		}
	}

}
