/******************************************************************************
    Initialisation acteurs et rôles
    Code servant à initialiser la base, pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-11-08 08:45:03+01:00, Thierry Graff : Creation from a split
********************************************************************************/
package initialize

import (
	"fmt"
	"path"
	"strings"
    "bufio"
    "os"
	
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
    "golang.org/x/text/encoding/charmap"
    "golang.org/x/text/transform"
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


// *********************************************************
// Ajoute les acteurs saisis dans un fichier csv pour importer 
// les acteurs de BDL au moment du démarrage de la base
func AddActeursInitiaux() {
	table := "acteur"
	csvfile := "acteurs-bdl-bastien.csv"
	fmt.Println("Remplit", table, " à partir de", csvfile)
	dirCsv := getPrivateDir()
	filename := path.Join(dirCsv, csvfile)
	// conversion utf8
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println(csvfile, "inexistant - ne remplit pas les acteurs")
        return
    }
    defer file.Close()
    decodingReader := transform.NewReader(file, charmap.Windows1252.NewDecoder())	
    scanner := bufio.NewScanner(decodingReader)
    lines := []string{}
    for scanner.Scan() {
        lines = append(lines, strings.TrimSpace(scanner.Text()))
    }
    // lecture csv
    sep := ";"
    fields := strings.Split(lines[0], sep)
    nfields := len(fields)
    var csv []map[string]string
    for i, line := range(lines){
        if i == 0 {
            continue
        }
        tmp := strings.Split(line, sep)
        if tmp[0] == ""{
            continue
        }
        current := make(map[string]string, nfields)
        for j, field := range(tmp) {
            current[fields[j]] = field
        }
        csv = append(csv, current)
    }
    
	ctx := ctxt.NewContext()
	db := ctx.DB
	id := int(0)
	n := int(0)
	query := `insert into acteur(
        nom,
        prenom,
        adresse1,
        adresse2,
        cp,
        ville,
        actif,
        notes
        )values($1,$2,$3,$4,$5,$6,$7,$8) returning id`
    actif := true
    for _, line := range(csv){
        if line["actif"] == "oui" {
            actif = true
        } else {
            actif = false
        }
        err := db.QueryRow(
            query,
            line["nom"],
            line["prenom"],
            line["adresse1"],
            line["adresse2"],
            line["cp"],
            line["ville"],
            actif,
            line["notes"]).Scan(&id)
        if err != nil {
            panic(err)
        }
        n++
    }
	fmt.Println(n, "lignes insérées")
}

