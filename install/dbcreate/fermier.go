/******************************************************************************
    Initialisation acteurs et rôles
    Code servant à initialiser la base, pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-11-08 08:45:03+01:00, Thierry Graff : Creation from a split
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
// Remplit les acteurs à partir d'un export de la base SCTL
// @param   versionSCTL ex "2020-12-23" - voir commentaire de install-bdl.go
func FillFermier(versionSCTL string) {
	table := "fermier"
	fmt.Println("Remplit " + table + " à partir de Exploita.csv")
	
	ctx := ctxt.NewContext()
	
	dirCsv := GetSCTLDataDir(ctx, versionSCTL)
	filename := path.Join(dirCsv, "Exploita.csv")
	records, err := tiglib.CsvMap(filename, ';')

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

    var adresse, cp, ville, tel, email string
    
	for _, v := range records {
		if v["Agricole"] != "1" {
			// Importer que les agricoles
			continue
		}
        if PRIVACY{
            adresse, cp, ville, tel, email = "", "", "", "", ""
        } else {
            adresse = strings.Replace(v["AdresseExp"], "'", `''`, -1)
            cp := v["CPExp"]
            if len(cp) > 5 {
                cp = cp[:5] // fix une typo dans la base SCTL
            }
            ville = strings.Replace(v["VilleExp"], "'", `''`, -1)
            tel = v["Telephone"]
            email = v["Mail"]
        }
		query := `insert into %s(
            id,
            nom,
            prenom,
            adresse,
            cp,
            ville,
            tel,
            email
            ) values(%s,'%s','%s','%s','%s','%s','%s','%s')`
		sql := fmt.Sprintf(
			query,
			table,
			v["IdExploitant"],
			strings.Replace(v["NOMEXP"], "'", `''`, -1),
			strings.Replace(v["Prenom"], "'", `''`, -1),
			adresse,
			cp,
			ville,
			tel,
			email)
		if _, err = tx.Exec(sql); err != nil {
			panic(err)
		}
	}
}

// *********************************************************
// Remplit les liens parcelle - exploitant à partir d'un export de la base SCTL
// @param   versionSCTL ex "2020-12-23" - voir commentaire de install-bdl.go
func FillLiensParcelleFermier(versionSCTL string) {
	table := "parcelle_fermier"
	fmt.Println("Remplit table " + table + " à partir de Subdivision.csv")
	
	ctx := ctxt.NewContext()
	
	dirCsv := GetSCTLDataDir(ctx, versionSCTL)
	filename := path.Join(dirCsv, "Subdivision.csv")

	records, err := tiglib.CsvMap(filename, ';') // N = 2844
	if err != nil {
		panic(err)
	}
	// remove doublons
	var k string
	var v [2]string
	uniques := make(map[string][2]string) // N = 433
	for _, record := range records {
		idP := record["IdParcelle"]
		idE := record["IdExploitant"]
		k = idP + "-" + idE
		v[0] = idP
		v[1] = idE
		uniques[k] = v
	}

	// insert db
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

	n := 0
	for _, unique := range uniques {
		idP := unique[0]
		idE := unique[1]
		sql := fmt.Sprintf("insert into %s(id_parcelle,id_fermier) values(%s, %s)", table, idP, idE)
		if _, err = tx.Exec(sql); err != nil {
			n++
			continue
		}
	}
	fmt.Printf("  %d associations pas enregistrées (bugs SCTL)\n", n)
}
