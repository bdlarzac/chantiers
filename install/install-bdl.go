/******************************************************************************
    Initialisation de l'environnement nécessaire au fonctionnement de l'application.
    Code pas utilisé en fonctionnement normal.
    Option -i pour installer la base de données et mettre les valeurs de départ
    Option -f pour fabriquer des données de test (fixtures)
    Exemples d'utilisation :
    go run install-bdl.go -i commune
    go run install-bdl.go -f stockage

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

var msgInstall, msgFixture string
var flagInstall, flagFixture *string

// *********************************************************
func init() {
	possibleInstall := []string{
		"all",
		"acteur",
		"bspied",
		"chaufer",
		"chautre",
		"commune",
		"plaquette",
		"parcelle",
		"stockage",
		"type",
		"vente",
		"ug",
	}
	msgInstall = strings.Join(possibleInstall, ", ")
	flagInstall = flag.String("i", "", "Install - valeurs possibles : "+msgInstall)

	possibleFixture := []string{
		"stockage",
		"acteur",
	}
	msgFixture = strings.Join(possibleFixture, ", ")
	flagFixture = flag.String("f", "", "Fixture - valeurs possibles : "+msgFixture)
}

// *********************************************************
func main() {

	flag.Parse()

	if (*flagInstall == "" && *flagFixture == "") || (*flagInstall != "" && *flagFixture != "") {
		fmt.Println("Utiliser avec -i (install) ou -f (fixture)")
		fmt.Println("Exemples :\n  go run install-bdl.go -i commune\n  go run install-bdl.go -f stockage")
		return
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
		installParcelle()
		installUG()
		installStockage()
		installPlaquette()
		installVente()
		installChaufer()
		installBSPied()
		installChautre()
	} else if *flagInstall == "type" {
		installTypes()
	} else if *flagInstall == "commune" {
		installCommune()
	} else if *flagInstall == "acteur" {
		installActeur()
	} else if *flagInstall == "parcelle" {
		installParcelle()
	} else if *flagInstall == "ug" {
		installUG()
	} else if *flagInstall == "stockage" {
		installStockage()
	} else if *flagInstall == "plaquette" {
		installPlaquette()
	} else if *flagInstall == "vente" {
		installVente()
	} else if *flagInstall == "chaufer" {
		installChaufer()
	} else if *flagInstall == "bspied" {
		installBSPied()
	} else if *flagInstall == "chautre" {
		installChautre()
	} else {
		fmt.Println("COMMANDE INVALIDE - Valeurs possibles pour -i : " + msgInstall)
	}
}
func installTypes() {
	initialize.CreateTable("typessence")
	initialize.CreateTable("typexploitation")
	initialize.CreateTable("typeop")
	initialize.CreateTable("typeunite")
	initialize.CreateTable("typevalorisation")
}
func installCommune() {
	initialize.CreateTable("commune")
	initialize.CreateTable("lieudit")
	initialize.CreateTable("commune_lieudit")
	initialize.FillCommuneOuLieudit("commune")
	initialize.FillCommuneOuLieudit("lieudit")
	initialize.FillLiensCommuneLieudit()
	initialize.CreateTable("lieudit_mot")
	initialize.FillLieuditMot()
}
func installActeur() {
	initialize.CreateTable("acteur")
	initialize.FillActeur()
	initialize.AddActeurBDL()
	initialize.AddActeurGFA()
	initialize.FillProprietaire()
}
func installParcelle() {
	initialize.CreateTable("parcelle")
	initialize.CreateTable("parcelle_lieudit")
	initialize.CreateTable("parcelle_exploitant")
	initialize.FillParcelle()
	initialize.FillLiensParcelleExploitant()
	initialize.FillLiensParcelleLieudit()
}
func installUG() {
	initialize.CreateTable("ug")
	initialize.CreateTable("parcelle_ug")
	initialize.FillUG()
	initialize.FillLiensParcelleUG()
}
func installStockage() {
	initialize.CreateTable("stockage")
	initialize.CreateTable("stockloyer")
	initialize.CreateTable("plaq")
	initialize.CreateTable("tas")
	initialize.CreateTable("humid")
	initialize.CreateTable("humid_acteur")
}
func installPlaquette() {
	initialize.CreateTable("plaqop")
	initialize.CreateTable("plaqtrans")
	initialize.CreateTable("plaqrange")
}
func installVente() {
	initialize.CreateTable("venteplaq")
	initialize.CreateTable("ventelivre")
	initialize.CreateTable("ventecharge")
}
func installChaufer() {
	initialize.CreateTable("chaufer")
	initialize.CreateTable("chaufer_parcelle")
}
func installBSPied() {
	initialize.CreateTable("bspied")
	initialize.CreateTable("bspied_parcelle")
}
func installChautre() {
	initialize.CreateTable("chautre")
}

// *********************************************************
func handleFixture() {
	if *flagFixture == "stockage" {
		fixture.FillStockage()
	} else if *flagFixture == "acteur" {
		fixture.AnonymizeActeurs()
	} else {
		fmt.Println("COMMANDE INVALIDE - Valeurs possibles pour -f : " + msgFixture)
	}
}
