/******************************************************************************
    Initialisation acteurs et rôles
    Code pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du
    @license    GPL
    @history    2019-11-08 08:45:03+01:00, Thierry Graff : Creation from a split
********************************************************************************/
package initialize

import (
	"fmt"
	"path"
	"strings"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
)

// Valeur du champ id_sctl pour l'acteur SCTL
// valeur fixe venant de la base SCTL
const SCTL_ID_SCTL = 28

// *********************************************************
func FillActeur() {
	table := "acteur"
	fmt.Println("Remplit " + table + " acteur à partir de Exploita.csv")
	dirCsv := getPrivateDir()
	filename := path.Join(dirCsv, "Exploita.csv")

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
		if v["Agricole"] != "1" {
			// Importer que les agricoles
			// continue
			// @todo - fait planter FillLiensParcelleExploitant
			// Pour ça, vérifier sur une carte
		}
		cp := v["CPExp"]
		if len(cp) > 5 {
			cp = cp[:5] // fix une typo dans la base SCTL
		}
		query := `insert into %s(
            id_sctl,
            nom,
            prenom,
            adresse1,
            cp,
            ville,
            tel,
            email,
            fournisseur,
            actif
            ) values(%s,'%s','%s','%s','%s','%s','%s','%s',%t,%t)`
		sql := fmt.Sprintf(
			query,
			table,
			v["IdExploitant"],
			strings.Replace(v["NOMEXP"], "'", `''`, -1),
			strings.Replace(v["Prenom"], "'", `''`, -1),
			strings.Replace(v["AdresseExp"], "'", `''`, -1),
			cp,
			strings.Replace(v["VilleExp"], "'", `''`, -1),
			v["Telephone"],
			v["Mail"],
			false,
			true)
		if _, err = tx.Exec(sql); err != nil {
			panic(err)
		}
	}
}

// *********************************************************
// Ajoute un acteur BDL
// Pas présent dans la base SCTL mais nécessaire au fonctionnement
// du logiciel BDL (car c'est un fournisseur)
func AddActeurBDL() {
	ctx := ctxt.NewContext()
	db := ctx.DB
	query := `insert into acteur(
        nom,
        prenom,       
        adresse1,     
        cp,
        ville,
        tel,
        email,
        fournisseur,  
        actif
        )values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		"BDL",
		"Bois du Larzac",
		"Montredon",
		"12230",
		"La Roque-Sainte-Marguerite",
		"05 65 62 13 39",
		"lesboisdularzac@larzac.org",
		true,
		true).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Crée acteur BDL, id = %d\n", id)
}

// *********************************************************
// Ajoute un acteur GFA
// Pas présent dans la base SCTL mais nécessaire au fonctionnement
// du logiciel BDL (car c'est un propriétaire)
func AddActeurGFA() {
	ctx := ctxt.NewContext()
	db := ctx.DB
	query := `insert into acteur(
        nom,
        proprietaire,  
        actif
        )values($1,$2,$3) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		"GFA",
		true,
		true).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Crée acteur GFA, id = %d\n", id)
}

// *********************************************************
// Remplit le champ proprietaire de l'acteur SCTL
func FillProprietaire() {
	fmt.Printf("Set proprietaire=true pour l'acteur SCTL, id_sctl = %d\n", SCTL_ID_SCTL)
	ctx := ctxt.NewContext()
	db := ctx.DB
	query := `update acteur set proprietaire=$1 where id_sctl=$2`
	_, err := db.Exec(query, true, SCTL_ID_SCTL)
	if err != nil {
		panic(err)
	}
}

// *********************************************************
func FillLiensParcelleExploitant() {
	table := "parcelle_exploitant"
	fmt.Println("Remplit table " + table + " à partir de Subdivision.csv")
	dirCsv := getDataDir()
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

	//n := 0
	for _, unique := range uniques {
		idP := unique[0]
		idE := unique[1]
		sql := fmt.Sprintf("insert into %s(id_parcelle,id_sctl_exploitant) values(%s, %s)", table, idP, idE)
		if _, err = tx.Exec(sql); err != nil {
			//n++
			//fmt.Printf("idP=%s - idE=%s\n", idP, idE)
			panic(err)
		}
	}
	//fmt.Printf("%d associations pas enregistrées (non-agricoles)\n", n)

}
