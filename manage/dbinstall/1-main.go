/******************************************************************************

    Installation et initialisation de la base BDL
    
    Lancer l'exécution en utilisant des variables d'environnement :
    ENV_CONFIG_FILE='../../config.env' APPLI_CONFIG_FILE='../../config.yml' go run 1-main.go

    Utilisation :
    -i : install
    -f : fixture
    -s : version de la base SCTL à utiliser
        Si l'option -s contient 2020-12-23
        Alors les exports de la base Access situés dans manage/sctl-data/csv-2020-12-23 seront utilisés.
        Ces exports sont des fichiers csv obtenus avec mdb-export

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-09-26 17:41:35+02:00, Thierry Graff : Creation
    @history    2023-01-12 10:32:08+01:00, Thierry Graff : Grosse refactorisation
********************************************************************************/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"bdl.dbinstall/bdl/install"
	"bdl.dbinstall/bdl/fixture"
	"flag"
	"fmt"
	"os"
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

	strSctlDataSource := "Répertoire contenant les dumps SCTL"

	flagInstall = flag.String("i", "", strInstall)
	flagFixture = flag.String("f", "", strFixture)
	flagSctlDataSource = flag.String("s", "", strSctlDataSource)
}

// *********************************************************
func main() {

	model.MustLoadEnv()
	ctxt.MustLoadConfig()
	ctxt.MustInitDB()
	ctxt.MustInitTemplates()

	flag.Parse()

	// check que un seul flag est utilisé
	i := (*flagInstall != "")
	f := (*flagFixture != "")
	// si aucun flag ou 2 flags sont utilisés
	if (!i && !f) || (i && f) {
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
		dirCsv := install.GetSCTLDataDir(ctx, *flagSctlDataSource)
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
	}

}

// *********************************************************
func handleInstall(ctx *ctxt.Context) {

	if *flagInstall == "all" {

		/*
		//fmt.Printf("db = %+v\n",ctx)
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
		*/

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
	install.CreateTable(ctx, "typessence")
	install.CreateTable(ctx, "typexploitation")
	install.CreateTable(ctx, "typeop")
	install.CreateTable(ctx, "typeunite")
	install.CreateTable(ctx, "typevalorisation")
	install.CreateTable(ctx, "typevente")
	install.CreateTable(ctx, "typegranulo")
	install.CreateTable(ctx, "typestockfrais")
}
func installCommune(ctx *ctxt.Context) {
	install.CreateTable(ctx, "commune")
	install.CreateTable(ctx, "lieudit")
	install.CreateTable(ctx, "commune_lieudit")
	install.FillCommune(ctx)
	install.FillLieudit(ctx, *flagSctlDataSource)
	install.FillLiensCommuneLieudit(ctx, *flagSctlDataSource)
	install.CreateTable(ctx, "lieudit_mot")
	install.FillLieuditMot(ctx)
}
func installActeur(ctx *ctxt.Context) {
	install.CreateTable(ctx, "acteur")
	install.AddActeursInitiaux(ctx)
	install.AddActeursFromCSV(ctx)
}
func installFermier(ctx *ctxt.Context) {
	install.CreateTable(ctx, "fermier")
	install.FillFermier(ctx, *flagSctlDataSource)
}
func installParcelle(ctx *ctxt.Context) {
	install.CreateTable(ctx, "parcelle")
	install.CreateTable(ctx, "parcelle_lieudit")
	install.CreateTable(ctx, "parcelle_fermier")
	install.FillParcelle(ctx, *flagSctlDataSource)
	install.FillLiensParcelleFermier(ctx, *flagSctlDataSource)
	install.FillLiensParcelleLieudit(ctx, *flagSctlDataSource)
}
func installUG(ctx *ctxt.Context) {
	install.CreateTable(ctx, "ug")
	install.CreateTable(ctx, "parcelle_ug")
	install.FillUG(ctx)
	install.FillLiensParcelleUG(ctx)
}
func installStockage(ctx *ctxt.Context) {
	install.CreateTable(ctx, "stockage")
	install.CreateTable(ctx, "stockfrais")
	install.CreateTable(ctx, "plaq")
	install.CreateTable(ctx, "tas")
	install.CreateTable(ctx, "humid")
	install.CreateTable(ctx, "humid_acteur")
	install.FillHangarsInitiaux(ctx)
}
func installChantier(ctx *ctxt.Context) {
	// liens pour plaq, chautre
	install.CreateTable(ctx, "chantier_ug")
	install.CreateTable(ctx, "chantier_fermier")
	install.CreateTable(ctx, "chantier_lieudit")
	// plaquettes
	install.CreateTable(ctx, "plaqop")
	install.CreateTable(ctx, "plaqtrans")
	install.CreateTable(ctx, "plaqrange")
	// autres valorisations
	install.CreateTable(ctx, "chautre")
	// chauffage fermier
	install.CreateTable(ctx, "chaufer")
	install.CreateTable(ctx, "chaufer_parcelle")
}
func installVente(ctx *ctxt.Context) {
	install.CreateTable(ctx, "venteplaq")
	install.CreateTable(ctx, "ventelivre")
	install.CreateTable(ctx, "ventecharge")
}
func installRecent(ctx *ctxt.Context) {
	install.CreateTable(ctx, "recent")
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
