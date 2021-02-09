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
    "golang.org/x/text/encoding/charmap"
    "golang.org/x/text/transform"
)

func AddActeursInitiaux() {
    // Attention ordre d'appel important
    // car détermine les ids des 4 premiers acteurs
    // ces ids sont utilisés dans l'appli BDL
    addActeurZero()
    addActeurSCTL()
    addActeurBDL()
    addActeurGFA()
}

// *********************************************************
// Acteur servant à autoriser des clés étrangères optionnelles dans les tables plaqtrans et plaqrange
// id = 0
func addActeurZero() {
	ctx := ctxt.NewContext()
	db := ctx.DB
	_,_ = db.Exec("insert into acteur values(0,0,'','','','','','','','','','','','',false,false,false,'')")
	fmt.Println("Créé l'acteur 0")
}

// *********************************************************
// Ajoute un acteur SCTL
// id = 1
// Nécessaire au fonctionnement du logiciel BDL (car c'est un propriétaire)
// Un exploitant "SCTL" est présent dans la base SCTL, mais pas importé car non agricole
func addActeurSCTL() {
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
		"SCTL",
		true,
		true).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Crée acteur SCTL, id = %d\n", id)
}

// *********************************************************
// Ajoute un acteur BDL
// id = 2
// Nécessaire au fonctionnement du logiciel BDL car c'est un fournisseur
func addActeurBDL() {
    
    var adresse1, cp, ville, tel, email string
    if PRIVACY {
		adresse1 = ""
		cp = ""
		ville = ""
		tel = ""
		email = ""
    } else {
		adresse1 = "Montredon"
		cp = "12230"
		ville = "La Roque-Sainte-Marguerite"
		tel = "05 65 62 13 39"
		email = "lesboisdularzac@larzac.org"
    }
	
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
		adresse1,
		cp,
		ville,
		tel,
		email,
		true,
		true).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Crée acteur BDL, id = %d\n", id)
}

// *********************************************************
// Ajoute un acteur GFA
// id = 3
// Nécessaire au fonctionnement du logiciel BDL car c'est un propriétaire
// Un exploitant "GFA" est présent dans la base SCTL, mais pas importé car non agricole
func addActeurGFA() {
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
		"GFA Larzac",
		true,
		true).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Crée acteur GFA, id = %d\n", id)
}

// *********************************************************
// Ajoute les acteurs saisis dans un fichier csv de Bastien
// pour importer les acteurs de BDL au moment du démarrage de la base
func AddActeursFromCSV() {
	table := "acteur"
	csvfile := "acteurs-bdl-bastien.csv"
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
    var actif bool
    var adresse1, adresse2, cp, ville string
    for _, line := range(csv){
        if PRIVACY {
            adresse1 = ""
            adresse2 = ""
            cp       = ""
            ville    = ""
        } else {
            adresse1 = line["adresse1"]
            adresse2 = line["adresse2"]
            cp       = line["cp"]
            ville    = line["ville"]
        }
        if line["actif"] == "oui" {
            actif = true
        } else {
            actif = false
        }
        err := db.QueryRow(
            query,
            line["nom"],
            line["prenom"],
            adresse1,
            adresse2,
            cp,
            ville,
            actif,
            line["notes"]).Scan(&id)
        if err != nil {
            panic(err)
        }
        n++
    }
	fmt.Println("Insère", n, "lignes dans", table, "à partir de", csvfile)
}


