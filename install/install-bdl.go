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
    mdb-export -d ';' -Q Sctl-Gfa-2020-02-27.mdb LieuDit > csv-2020-02-27/LieuDit.csv
    
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
    errorMsg += "  go run install-bdl.go -i all -s 2021-07-27\n"
    errorMsg += "  go run install-bdl.go -i chantier\n"
    errorMsg += "  go run install-bdl.go -f stockage\n"
    errorMsg += "  go run install-bdl.go -i commune -s 2021-07-27\n"
    
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
    strMigrate := strings.Join(possibleMigrations, " ")
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
	
    ctx := ctxt.NewContext()
        
	// options ayant besoin de la version de la base SCTL utilisée
	needFlagS := *flagInstall == "fermier" || *flagInstall == "commune" || *flagInstall == "parcelle" || *flagInstall == "all"
	if needFlagS {
	    if *flagSctlDataSource == "" {
            fmt.Println(errorMsg)
            fmt.Println("PARAMETRE MANQUANT : -s")
            return
	    }
	    // check que le répertoire existe
        dirCsv := dbcreate.GetSCTLDataDir(ctx, *flagSctlDataSource)
        _, err := os.Stat(dirCsv)
        if os.IsNotExist(err) {
            fmt.Println("REPERTOIRE SCTL INEXISTANT :", dirCsv)
            return
        }
	}
	
	if *flagInstall != "" {
		handleInstall(ctx)
	} else if *flagFixture != "" {
		handleFixture(ctx)
	} else if *flagMigrate != "" {
		handleMigration(ctx, *flagMigrate)
	}

}

// *********************************************************
func handleInstall(ctx *ctxt.Context) {
	if *flagInstall == "all" {
	    
        db := ctx.DB
        var err error
        _, err = db.Exec(fmt.Sprintf("drop schema if exists %s cascade", ctx.Config.Database.Schema))
        if err != nil {
            panic(err)
        }
        _, err = db.Exec(fmt.Sprintf("create schema %s", ctx.Config.Database.Schema))
        if err != nil {
            panic(err)
        }
        _, err = db.Exec(fmt.Sprintf(`set search_path='%s'`, ctx.Config.Database.Schema))
	    
		installTypes(ctx)
		installCommune(ctx)
		installActeur(ctx)
		installFermier(ctx)
		installParcelle(ctx)
		installUG(ctx)
		installStockage(ctx)
		installChantier(ctx)
		installVente(ctx)
		installRecent(ctx)
		for _, migration := range(possibleMigrations){
		    handleMigration(ctx, migration)
		}
	} else if *flagInstall == "type" {
		installTypes(ctx)
	} else if *flagInstall == "commune" {
		installCommune(ctx)
	} else if *flagInstall == "acteur" {
		installActeur(ctx)
	} else if *flagInstall == "fermier" {
		installFermier(ctx)
	} else if *flagInstall == "parcelle" {
		installParcelle(ctx)
	} else if *flagInstall == "ug" {
		installUG(ctx)
	} else if *flagInstall == "stockage" {
		installStockage(ctx)
	} else if *flagInstall == "chantier" {
		installChantier(ctx)
	} else if *flagInstall == "vente" {
		installVente(ctx)
	} else if *flagInstall == "recent" {
		installRecent(ctx)
	} else {
		fmt.Println(errorMsg)
	}
}
func installTypes(ctx *ctxt.Context) {
	dbcreate.CreateTable(ctx, "typessence")
	dbcreate.CreateTable(ctx, "typexploitation")
	dbcreate.CreateTable(ctx, "typeop")
	dbcreate.CreateTable(ctx, "typeunite")
	dbcreate.CreateTable(ctx, "typevalorisation")
	dbcreate.CreateTable(ctx, "typevente")
	dbcreate.CreateTable(ctx, "typegranulo")
	dbcreate.CreateTable(ctx, "typestockfrais")
}
func installCommune(ctx *ctxt.Context) {
	dbcreate.CreateTable(ctx, "commune")
	dbcreate.CreateTable(ctx, "lieudit")
	dbcreate.CreateTable(ctx, "commune_lieudit")
	dbcreate.FillCommune(ctx)
	dbcreate.FillLieudit(ctx, *flagSctlDataSource)
	dbcreate.FillLiensCommuneLieudit(ctx, *flagSctlDataSource)
	dbcreate.CreateTable(ctx, "lieudit_mot")
	dbcreate.FillLieuditMot(ctx)
}
func installActeur(ctx *ctxt.Context) {
	dbcreate.CreateTable(ctx, "acteur")
	dbcreate.AddActeursInitiaux(ctx)
	dbcreate.AddActeursFromCSV(ctx)
}
func installFermier(ctx *ctxt.Context){
	dbcreate.CreateTable(ctx, "fermier")
	dbcreate.FillFermier(ctx, *flagSctlDataSource)
}
func installParcelle(ctx *ctxt.Context) {
	dbcreate.CreateTable(ctx, "parcelle")
	dbcreate.CreateTable(ctx, "parcelle_lieudit")
	dbcreate.CreateTable(ctx, "parcelle_fermier")
	dbcreate.FillParcelle(ctx, *flagSctlDataSource)
	dbcreate.FillLiensParcelleFermier(ctx, *flagSctlDataSource)
	dbcreate.FillLiensParcelleLieudit(ctx, *flagSctlDataSource)
}
func installUG(ctx *ctxt.Context) {
	dbcreate.CreateTable(ctx, "ug")
	dbcreate.CreateTable(ctx, "parcelle_ug")
	dbcreate.FillUG(ctx)
	dbcreate.FillLiensParcelleUG(ctx)
}
func installStockage(ctx *ctxt.Context) {
	dbcreate.CreateTable(ctx, "stockage")
	dbcreate.CreateTable(ctx, "stockfrais")
	dbcreate.CreateTable(ctx, "plaq")
	dbcreate.CreateTable(ctx, "tas")
	dbcreate.CreateTable(ctx, "humid")
	dbcreate.CreateTable(ctx, "humid_acteur")
	dbcreate.FillHangarsInitiaux(ctx)
}
func installChantier(ctx *ctxt.Context) {
    // liens pour plaq, chautre
	dbcreate.CreateTable(ctx, "chantier_ug")
	dbcreate.CreateTable(ctx, "chantier_fermier")
	dbcreate.CreateTable(ctx, "chantier_lieudit")
    // plaquettes
	dbcreate.CreateTable(ctx, "plaqop")
	dbcreate.CreateTable(ctx, "plaqtrans")
	dbcreate.CreateTable(ctx, "plaqrange")
	// autres valorisations
	dbcreate.CreateTable(ctx, "chautre")
	// chauffage fermier
	dbcreate.CreateTable(ctx, "chaufer")
	dbcreate.CreateTable(ctx, "chaufer_parcelle")
}
func installVente(ctx *ctxt.Context) {
	dbcreate.CreateTable(ctx, "venteplaq")
	dbcreate.CreateTable(ctx, "ventelivre")
	dbcreate.CreateTable(ctx, "ventecharge")
}
func installRecent(ctx *ctxt.Context) {
	dbcreate.CreateTable(ctx, "recent")
}

// *********************************************************
func handleFixture(ctx *ctxt.Context) {
	if *flagFixture == "stockage" {
		fixture.FillStockage(ctx)
	} else if *flagFixture == "acteur" {
		fixture.AnonymizeActeurs(ctx)
	} else {
		fmt.Println(errorMsg)
	}
}

// *********************************************************
//  @param migration Nom de la migration = nom de la fonction de dbmigrate à exécuter 
//                   pour effectuer la migration
func handleMigration(ctx *ctxt.Context, migration string) {
    // check que la migration existe
    if !tiglib.InArrayString(migration, possibleMigrations){
        fmt.Println("MIGRATION INEXISTANTE : " + migration)
        fmt.Println("Migrations possibles : " + strings.Join(possibleMigrations, " "))
        return
    }
    switch(migration){
    case "Migrate_2021_03_01_exemple":
        dbmigrate.Migrate_2021_03_01_exemple(ctx)
    case "Migrate_2021_11_10_note_plaq":
        dbmigrate.Migrate_2021_11_10_note_plaq(ctx)
    }
}


// *********************************************************
// Renvoie la liste des migrations possibles 
// = liste des fonctions du package dbmigrate commençant par Migrate_
// Ça n'a pas l'air possible avec reflect => bidouille avec regex
func computeMigrations() (res []string) {
    out, err := exec.Command("grep", "-rn", "func Migrate_", "dbmigrate").Output()
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
