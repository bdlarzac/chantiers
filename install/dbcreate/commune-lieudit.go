/******************************************************************************
    Initialisation commune et lieudit.
    Code servant à initialiser la base, pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du
    @license    GPL
    @history    2019-11-05 06:06:04+01:00, Thierry Graff : Creation from a split
********************************************************************************/
package dbcreate

import (
	"fmt"
	"path"
	"strings"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
)

// *********************************************************
// Prend les données dans install/data/commune.csv
// A la différence des autres données venant de la SCTL,
// commune.csv est versionné (car modifié, ajout de colonne nom_court)
func FillCommune() {
	fmt.Println("Remplit table commune à partir de commune.csv")
	dirCsv := GetDataDir()
	filename := path.Join(dirCsv, "commune.csv")
	records, err := tiglib.CsvMap(filename, ';')
	// insert db
	ctx := ctxt.NewContext()
	db := ctx.DB
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare("insert into commune(id,nom,nomcourt) values($1, $2, $3)")
	for _, v := range records {
        _, err = stmt.Exec(v["id"], v["nom"], v["nom_court"])
        if err != nil {
            panic(err)
        }
	}
}

// *********************************************************
// @param   versionSCTL ex "2020-12-23" - voir commentaire de install-bdl.go
func FillLieudit(versionSCTL string) {
	fmt.Println("Remplit table lieudit à partir de LieuDit.csv")
	
	ctx := ctxt.NewContext()
	
	dirCsv := GetSCTLDataDir(ctx, versionSCTL)
	filename := path.Join(dirCsv, "LieuDit.csv")
	records, err := tiglib.CsvMap(filename, ';')
	if err != nil {
		panic(err)
	}
	// insert db
	db := ctx.DB
	stmt, err := db.Prepare("insert into lieudit(id,nom) values($1, $2)")
	for _, v := range records {
        _, err = stmt.Exec(v["IdLieuDit"], v["Libelle"])
        if err != nil {
            panic(err)
        }
	}
}

// *********************************************************
// @param   versionSCTL ex "2020-12-23" - voir commentaire de install-bdl.go
func FillLiensCommuneLieudit(versionSCTL string) {
	table := "commune_lieudit"
	fmt.Println("Remplit table " + table + " à partir de SubdivCadastre.csv")

	ctx := ctxt.NewContext()
	
	dirCsv := GetSCTLDataDir(ctx, versionSCTL)
	filename := path.Join(dirCsv, "SubdivCadastre.csv")
	//
	records, err := tiglib.CsvMap(filename, ';') // N = 2844 pour base 2018
	if err != nil {
		panic(err)
	}
	// remove doublons
	var k string
	var v [2]string
	uniques := make(map[string][2]string) // N = 433 pour base 2018
	for _, record := range records {
		idC := record["IdCommune"]
		idLD := record["IdLieuDit"]
		k = idC + "-" + idLD
		v[0] = idC
		v[1] = idLD
		uniques[k] = v
	}
	// insert db
	db := ctx.DB
	stmt, err := db.Prepare("insert into " + table + "(id_commune,id_lieudit) values($1, $2)")
    if err != nil {
        panic(err)
    }
	for _, unique := range uniques {
		idC := unique[0]
		idLD := unique[1]
        _, err := stmt.Exec(idC, idLD)
        if err != nil {
            if idLD == "267" {
                continue; // Les Mares, bug sctl
            }
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
