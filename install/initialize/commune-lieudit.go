/******************************************************************************
    Initialisation commune et lieudit.
    Code pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du
    @license    GPL
    @history    2019-11-05 06:06:04+01:00, Thierry Graff : Creation from a split
********************************************************************************/
package initialize

import (
	"fmt"
	"path"
	"strings"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
)

// *********************************************************
func FillCommuneOuLieudit(table string) {
	fmt.Println("Remplit table " + table + " à partir de " + table + ".csv")
	dirCsv := getDataDir()
	filename := path.Join(dirCsv, table+".csv")

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

	for _, v := range records {
		sql := fmt.Sprintf("insert into %s(id,nom) values(%s, '%s')", table, v["id"], strings.Replace(v["nom"], "'", `''`, -1))
		if _, err = tx.Exec(sql); err != nil {
			panic(err)
		}
	}
}

// *********************************************************
func FillLiensCommuneLieudit() {
	table := "commune_lieudit"
	fmt.Println("Remplit table " + table + " à partir de SubdivCadastre.csv")
	dirCsv := getDataDir()
	filename := path.Join(dirCsv, "SubdivCadastre.csv")

	records, err := tiglib.CsvMap(filename, ';') // N = 2844
	if err != nil {
		panic(err)
	}

	// remove doublons
	var k string
	var v [2]string
	uniques := make(map[string][2]string) // N = 433
	for _, record := range records {
		idC := record["IdCommune"]
		idLD := record["IdLieuDit"]
		k = idC + "-" + idLD
		v[0] = idC
		v[1] = idLD
		uniques[k] = v
	}

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

	for _, unique := range uniques {
		idC := unique[0]
		idLD := unique[1]
		sql := fmt.Sprintf("insert into %s(id_commune,id_lieudit) values(%s, '%s')", table, idC, idLD)
		if _, err = tx.Exec(sql); err != nil {
			panic(err)
		}
	}
}

// *********************************************************
func FillLieuditMot() {
	type ld struct {
		id  int
		nom string
	}
	corres := make(map[string][]ld)
	var id int
	var name string
	ignore := []string{"LE", "LA", "LES", "DE", "DU", "D'", "DES", "DEL", "ET", "L'"}
	ctx := ctxt.NewContext()
	db := ctx.DB

	rows, err := db.Query("select id, nom from lieudit")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	stmt, err := db.Prepare("INSERT INTO lieudit_mot(mot,id,nom) VALUES($1,$2,$3)")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		parts := strings.Split(name, " ")
		for _, part := range parts {
			if tiglib.InArrayString(part, ignore) {
				continue
			}
			corres[part] = append(corres[part], ld{id, name})
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for mot, lds := range corres {
		for _, ld := range lds {
			_, err = stmt.Exec(mot, ld.id, ld.nom)
			if err != nil {
				panic(err)
			}
		}
	}
}
