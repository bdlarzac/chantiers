/******************************************************************************
    Initialisation de l'environnement nécessaire au fonctionnement de l'application.
    - Installation de la base (package dbcreate)
    - Modifications de la base (package dbmigrate)
    
    Code pas utilisé en fonctionnement normal.
    
    Utilisation :
    -i : install
    -f : fixture
    -m : migrate
    
    -s est utilisée en lien avec la config (dev / sctl-data).
    Correspond à un sous-répertoire de la valeur de la config.
    Ex: Si la config contient
    dev:
      sctl-data: /path/to/db-sctl
    Et que l'option -s contient 2020-12-23
    Alors les exports de la base Access doivent se trouver dans /path/to/db-sctl/csv-2020-12-23
    Ces exports sont des fichiers csv obtenus avec mdb-export
    
    Ex de commande mdb-export à exécuter depuis /path/to/db-sctl :
    mdb-export -d ';' -Q Sctl-Gfa-2020-02-27.mdb LieuDit > csv-2020-03-06/LieuDit.csv
    
    Pour installer mdb-export :
    sudo apt install mdbtools
    
    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-09-26 17:41:35+02:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/install/fixture"
	"bdl.local/install/dbcreate"
	"bdl.local/install/dbmigrate"
	"bdl.local/bdl/generic/tiglib"
	"flag"
	"fmt"
	"strings"
	"os"
	"os/exec"
	"regexp"
)

var errorMsg string
var flagInstall, flagFixture, flagMigrate, flagSctlDataSource *string

var possibleMigrations []string

// *********************************************************
// toujours appelé par go au chargement du package
func init() {
    errorMsg = "COMMANDE INVALIDE\n"
    errorMsg += "Utiliser avec -i (install) ou -m (migration) ou -f (fixture)\n"
    errorMsg += "Certaines commandes ont aussi besoin de -s (source des données SCTL) :\n"
    errorMsg += "  -i fermier\n"
    errorMsg += "  -i commune\n"
    errorMsg += "  -i parcelle\n"
    errorMsg += "  -i all\n"
    errorMsg += "Exemples :\n"
    errorMsg += "  go run install-bdl.go -i all -s 2020-12-16\n"
    errorMsg += "  go run install-bdl.go -i chantier\n"
    errorMsg += "  go run install-bdl.go -f stockage\n"
    errorMsg += "  go run install-bdl.go -i commune -s 2020-12-16\n"
    
	possibleInstall := []string{
		"all",
		"acteur",
		"chantier",
		"commune",
		"fermier",
		"parcelle",
		"recent",
		"stockage",
		"type",
		"vente",
		"ug",
	}
    errorMsg += "Valeurs possibles pour -i :\n  "
    strInstall := strings.Join(possibleInstall, ", ")
	errorMsg += strInstall + "\n"

	possibleFixture := []string{
		"stockage",
		"acteur",
	}
    errorMsg += "Valeurs possibles pour -f :\n  "
    strFixture := strings.Join(possibleFixture, ", ")
	errorMsg += strFixture + "\n"
	
	possibleMigrations = computeMigrations()
    errorMsg += "Valeurs possibles pour -m :\n  "
    strMigrate := strings.Join(possibleMigrations, ", ")
	errorMsg += strMigrate + "\n"
	
	strSctlDataSource := "Répertoire contenant les dumps SCTL"
	
	flagInstall = flag.String("i", "", strInstall)
	flagFixture = flag.String("f", "", strFixture)
	flagMigrate = flag.String("m", "", strMigrate)
	flagSctlDataSource = flag.String("s", "", strSctlDataSource)
}

// *********************************************************
func main() {

	flag.Parse()
	
	// check que un seul flag est utilisé
	i := (*flagInstall != "")
	f := (*flagFixture != "")
	m := (*flagMigrate != "")
    // si aucun flag ou 2 flags sont utilisés
	if ( (!i && !f && !m) || (i && f) || (i && m) || (f && m) ) {
		fmt.Println(errorMsg)
		return
	}
	
	// options ayant besoin de la version de la base SCTL utilisée
	needFlagS := *flagInstall == "fermier" || *flagInstall == "commune" || *flagInstall == "parcelle" || *flagInstall == "all"
	if needFlagS {
	    if *flagSctlDataSource == "" {
            fmt.Println(errorMsg)
            fmt.Println("PARAMETRE MANQUANT : -s")
            return
	    }
	    // check que le répertoire existe
        ctx := ctxt.NewContext()
        dirCsv := dbcreate.GetSCTLDataDir(ctx, *flagSctlDataSource)
        _, err := os.Stat(dirCsv)
        if os.IsNotExist(err) {
            fmt.Println("REPERTOIRE SCTL INEXISTANT :", dirCsv)
            return
        }
	}
	
	if *flagInstall != "" {
		handleInstall()
	} else if *flagFixture != "" {
		handleFixture()
	} else if *flagMigrate != "" {
		handleMigration(*flagMigrate)
	}

}

// *********************************************************
func handleInstall() {
	if *flagInstall == "all" {
		installTypes()
		installCommune()
		installActeur()
		installFermier()
		installParcelle()
		installUG()
		installStockage()
		installChantier()
		installVente()
		installRecent()
		for _, migration := range(possibleMigrations){
fmt.Printf("possibleMigrations %+v\n", migration)
		    handleMigration(migration)
		}
	} else if *flagInstall == "type" {
		installTypes()
	} else if *flagInstall == "commune" {
		installCommune()
	} else if *flagInstall == "acteur" {
		installActeur()
	} else if *flagInstall == "fermier" {
		installFermier()
	} else if *flagInstall == "parcelle" {
		installParcelle()
	} else if *flagInstall == "ug" {
		installUG()
	} else if *flagInstall == "stockage" {
		installStockage()
	} else if *flagInstall == "chantier" {
		installChantier()
	} else if *flagInstall == "vente" {
		installVente()
	} else if *flagInstall == "recent" {
		installRecent()
	} else {
		fmt.Println(errorMsg)
	}
}
func installTypes() {
	dbcreate.CreateTable("typessence")
	dbcreate.CreateTable("typexploitation")
	dbcreate.CreateTable("typeop")
	dbcreate.CreateTable("typeunite")
	dbcreate.CreateTable("typevalorisation")
	dbcreate.CreateTable("typevente")
	dbcreate.CreateTable("typegranulo")
	dbcreate.CreateTable("typestockfrais")
}
func installCommune() {
	dbcreate.CreateTable("commune")
	dbcreate.CreateTable("lieudit")
	dbcreate.CreateTable("commune_lieudit")
	dbcreate.FillCommune()
	dbcreate.FillLieudit(*flagSctlDataSource)
	dbcreate.FillLiensCommuneLieudit(*flagSctlDataSource)
	dbcreate.CreateTable("lieudit_mot")
	dbcreate.FillLieuditMot()
}
func installActeur() {
	dbcreate.CreateTable("acteur")
	dbcreate.AddActeursInitiaux()
	dbcreate.AddActeursFromCSV()
}
func installFermier(){
	dbcreate.CreateTable("fermier")
	dbcreate.FillFermier(*flagSctlDataSource)
}
func installParcelle() {
	dbcreate.CreateTable("parcelle")
	dbcreate.CreateTable("parcelle_lieudit")
	dbcreate.CreateTable("parcelle_fermier")
	dbcreate.FillParcelle(*flagSctlDataSource)
	dbcreate.FillLiensParcelleFermier(*flagSctlDataSource)
	dbcreate.FillLiensParcelleLieudit(*flagSctlDataSource)
}
func installUG() {
	dbcreate.CreateTable("ug")
	dbcreate.CreateTable("parcelle_ug")
	dbcreate.FillUG()
	dbcreate.FillLiensParcelleUG()
}
func installStockage() {
	dbcreate.CreateTable("stockage")
	dbcreate.CreateTable("stockfrais")
	dbcreate.CreateTable("plaq")
	dbcreate.CreateTable("tas")
	dbcreate.CreateTable("humid")
	dbcreate.CreateTable("humid_acteur")
	dbcreate.FillHangarsInitiaux()
}
func installChantier() {
    // liens pour plaq, chautre
	dbcreate.CreateTable("chantier_ug")
	dbcreate.CreateTable("chantier_fermier")
	dbcreate.CreateTable("chantier_lieudit")
    // plaquettes
	dbcreate.CreateTable("plaqop")
	dbcreate.CreateTable("plaqtrans")
	dbcreate.CreateTable("plaqrange")
	// autres valorisations
	dbcreate.CreateTable("chautre")
	// chauffage fermier
	dbcreate.CreateTable("chaufer")
	dbcreate.CreateTable("chaufer_parcelle")
}
func installVente() {
	dbcreate.CreateTable("venteplaq")
	dbcreate.CreateTable("ventelivre")
	dbcreate.CreateTable("ventecharge")
}
func installRecent() {
	dbcreate.CreateTable("recent")
}

// *********************************************************
func handleFixture() {
	if *flagFixture == "stockage" {
		fixture.FillStockage()
	} else if *flagFixture == "acteur" {
		fixture.AnonymizeActeurs()
	} else {
		fmt.Println(errorMsg)
	}
}

// *********************************************************
//  @param migration Nom de la migration = nom de la fonction de dbmigrate à exécuter 
//                   pour effectuer la migration
func handleMigration(migration string) {
    // check que la migration existe
    if !tiglib.InArrayString(migration, possibleMigrations){
        fmt.Println("MIGRATION INEXISTANTE : " + migration)
        fmt.Println("Migrations possibles : " + strings.Join(possibleMigrations, ", "))
    }
    switch(migration){
    case "Migrate_2021_03_01":
        dbmigrate.Migrate_2021_03_01()
    }
}


// *********************************************************
// Renvoie la liste des migrations possibles 
// = liste des fonctions du package dbmigrate commençant par Migrate_
// Ça n'a pas l'air possible avec reflect => bidouille avec regex
func computeMigrations() (res []string) {
    out, err := exec.Command("grep", "-rn", "func Migrate_", "/home/thierry/dev/jobs/bdl/appli/install/dbmigrate").Output()
    if err != nil {
        panic(err)
    }
    r := regexp.MustCompile(`func (Migrate_.*?)\s*\(`)
    for _, line := range(strings.Split(strings.TrimSpace(string(out)), "\n")) {
        m := r.FindStringSubmatch(line)
        res = append(res, m[1])
    }
    return res
}
