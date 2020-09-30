/******************************************************************************
    Initialisation de l'environnement nécessaire au fonctionnement de l'application.
    Code pas utilisé en fonctionnement normal.

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
var flagInstall, flagFixture *string

// *********************************************************
// toujours appelé par go au chargement du package
func init() {
    errorMsg = "COMMANDE INVALIDE\n"
    errorMsg += "Utiliser avec -i (install) ou -f (fixture)\n"
    errorMsg += "Exemples :\n"
    errorMsg += "  go run install-bdl.go -i commune\n"
    errorMsg += "  go run install-bdl.go -f stockage\n"
    
	possibleInstall := []string{
		"all",
		"acteur",
		"bspied",
		"chaufer",
		"chautre",
		"commune",
		"plaquette",
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
	
	flagInstall = flag.String("i", "", strInstall)
	flagFixture = flag.String("f", "", strFixture)
}

// *********************************************************
func main() {

	flag.Parse()

	if (*flagInstall == "" && *flagFixture == "") || (*flagInstall != "" && *flagFixture != "") {
		fmt.Println(errorMsg)
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
		installRecent()
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
	initialize.AddActeursInitiaux()
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
	initialize.CreateTable("stockfrais")
	initialize.CreateTable("plaq")
	initialize.CreateTable("tas")
	initialize.CreateTable("humid")
	initialize.CreateTable("humid_acteur")
	initialize.FillHangarsInitiaux()
}
func installPlaquette() {
	initialize.CreateTable("plaqop")
	initialize.CreateTable("plaqtrans")
	initialize.CreateTable("plaqrange")
	//initialize.FillChantiersPlaquettesFromXls()
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
