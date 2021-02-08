/******************************************************************************
    Initialisation de l'environnement nécessaire au fonctionnement de l'application.
    Code pas utilisé en fonctionnement normal.
    
    L'option -s est utilisée en lien avec la config (dev / sctl-data).
    Correspond à un sous-répertoire de la valeur de la config.
    Ex: Si la config contient
    dev:
      sctl-data: /path/to/db-sctl
    Et que l'option -s contient 2020-12-23
    Alors les exports de la base Access doivent se trouver dans /path/to/db-sctl/csv-2020-12-23
    Ces exports sont des fichiers csv obtenus avec mdb-export
    
    Ex de commande à exécuter depuis /path/to/db-sctl :
    mdb-export -d ';' -Q Sctl-Gfa-2020-02-27.mdb LieuDit > csv-2020-03-06/LieuDit.csv
    
    Pour installer mdb-export :
    sudo apt install mdbtools
    
    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-09-26 17:41:35+02:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"bdl.local/install/fixture"
	"bdl.local/install/initialize"
	"flag"
	"fmt"
	"strings"
)

var errorMsg string
var flagInstall, flagFixture, flagSctlDataSource *string

// *********************************************************
// toujours appelé par go au chargement du package
func init() {
    errorMsg = "COMMANDE INVALIDE\n"
    errorMsg += "Utiliser avec -i (install) ou -f (fixture)\n"
    errorMsg += "Certaines commandes ont aussi besoin de -s (source des données SCTL) :\n"
    errorMsg += "  -i fermier\n"
    errorMsg += "  -i commune\n"
    errorMsg += "  -i parcelle\n"
    errorMsg += "  -i all\n"
    errorMsg += "Exemples :\n"
    errorMsg += "  go run install-bdl.go -i chantier\n"
    errorMsg += "  go run install-bdl.go -f stockage\n"
    errorMsg += "  go run install-bdl.go -i commune -s 2020-12-23\n"
    
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
	
	strSctlDataSource := "Répertoire contenant les dumps SCTL"
	
	flagInstall = flag.String("i", "", strInstall)
	flagFixture = flag.String("f", "", strFixture)
	flagSctlDataSource = flag.String("s", "", strSctlDataSource)
}

// *********************************************************
func main() {

	flag.Parse()

	if (*flagInstall == "" && *flagFixture == "") || (*flagInstall != "" && *flagFixture != "") {
		fmt.Println(errorMsg)
		return
	}
	
	// options ayant besoin de la version de la base SCTL utilisée
	if *flagInstall == "fermier" || *flagInstall == "commune" || *flagInstall == "parcelle" || *flagInstall == "all" {
	    if *flagSctlDataSource == "" {
            fmt.Println(errorMsg)
            fmt.Println("PARAMETRE MANQUANT : -s")
            return
	    }
	}
	
	if *flagInstall != "" {
		handleInstall()
	} else if *flagFixture != "" {
		handleFixture()
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
	initialize.CreateTable("typessence")
	initialize.CreateTable("typexploitation")
	initialize.CreateTable("typeop")
	initialize.CreateTable("typeunite")
	initialize.CreateTable("typevalorisation")
	initialize.CreateTable("typegranulo")
	initialize.CreateTable("typestockfrais")
}
func installCommune() {
	initialize.CreateTable("commune")
	initialize.CreateTable("lieudit")
	initialize.CreateTable("commune_lieudit")
	initialize.FillCommune()
	initialize.FillLieudit(*flagSctlDataSource)
	initialize.FillLiensCommuneLieudit(*flagSctlDataSource)
	initialize.CreateTable("lieudit_mot")
	initialize.FillLieuditMot()
}
func installActeur() {
	initialize.CreateTable("acteur")
	initialize.AddActeursInitiaux()
	initialize.AddActeursFromCSV()
}
func installFermier(){
	initialize.CreateTable("fermier")
	initialize.FillFermier(*flagSctlDataSource)
}
func installParcelle() {
	initialize.CreateTable("parcelle")
	initialize.CreateTable("parcelle_lieudit")
	initialize.CreateTable("parcelle_fermier")
	initialize.FillParcelle(*flagSctlDataSource)
	initialize.FillLiensParcelleFermier(*flagSctlDataSource)
	initialize.FillLiensParcelleLieudit(*flagSctlDataSource)
}
func installUG() {
	initialize.CreateTable("ug")
	initialize.CreateTable("parcelle_ug")
	initialize.FillUG()
	initialize.FillLiensParcelleUG()
}
func installStockage() {
	initialize.CreateTable("stockage")
	initialize.CreateTable("stockfrais")
	initialize.CreateTable("plaq")
	initialize.CreateTable("tas")
	initialize.CreateTable("humid")
	initialize.CreateTable("humid_acteur")
	initialize.FillHangarsInitiaux()
}
func installChantier() {
    // liens pour plaq, bspied, chautre
	initialize.CreateTable("chantier_ug")
	initialize.CreateTable("chantier_fermier")
	initialize.CreateTable("chantier_lieudit")
    // plaquettes
	initialize.CreateTable("plaqop")
	initialize.CreateTable("plaqtrans")
	initialize.CreateTable("plaqrange")
	// bois sur pied
	initialize.CreateTable("bspied")
	// autres valorisations
	initialize.CreateTable("chautre")
	// chauffage fermier
	initialize.CreateTable("chaufer")
	initialize.CreateTable("chaufer_parcelle")
}
func installVente() {
	initialize.CreateTable("venteplaq")
	initialize.CreateTable("ventelivre")
	initialize.CreateTable("ventecharge")
}
func installRecent() {
	initialize.CreateTable("recent")
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
