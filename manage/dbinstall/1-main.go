/******************************************************************************

    Installation de la base BDL
    Code pas utilisé en fonctionnement normal.
    
    Lancer l'exécution en utilisant des variables d'environnement et en utilisant *.go :
    ENV_CONFIG_FILE='../../config.env' APPLI_CONFIG_FILE='../../config.yml' go run *.go

    Utilisation :
    -i : install
    -f : fixture

    -s est utilisée en lien avec la config (dev / sctl-data).
    Correspond à un sous-répertoire de la valeur de la config.
    Ex: Si la config contient
    dev:
      sctl-data: /path/to/db-sctl
    Et que l'option -s contient 2020-12-23
    Alors les exports de la base Access doivent se trouver dans /path/to/db-sctl/csv-2020-12-23
    Ces exports sont des fichiers csv obtenus avec mdb-export

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-09-26 17:41:35+02:00, Thierry Graff : Creation
    @history    2023-01-12 10:32:08+01:00, Thierry Graff : Grosse refactorisation
********************************************************************************/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
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
		dirCsv := GetSCTLDataDir(ctx, *flagSctlDataSource)
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
	CreateTable(ctx, "typessence")
	CreateTable(ctx, "typexploitation")
	CreateTable(ctx, "typeop")
	CreateTable(ctx, "typeunite")
	CreateTable(ctx, "typevalorisation")
	CreateTable(ctx, "typevente")
	CreateTable(ctx, "typegranulo")
	CreateTable(ctx, "typestockfrais")
}
func installCommune(ctx *ctxt.Context) {
	CreateTable(ctx, "commune")
	CreateTable(ctx, "lieudit")
	CreateTable(ctx, "commune_lieudit")
	FillCommune(ctx)
	FillLieudit(ctx, *flagSctlDataSource)
	FillLiensCommuneLieudit(ctx, *flagSctlDataSource)
	CreateTable(ctx, "lieudit_mot")
	FillLieuditMot(ctx)
}
func installActeur(ctx *ctxt.Context) {
	CreateTable(ctx, "acteur")
	AddActeursInitiaux(ctx)
	AddActeursFromCSV(ctx)
}
func installFermier(ctx *ctxt.Context) {
	CreateTable(ctx, "fermier")
	FillFermier(ctx, *flagSctlDataSource)
}
func installParcelle(ctx *ctxt.Context) {
	CreateTable(ctx, "parcelle")
	CreateTable(ctx, "parcelle_lieudit")
	CreateTable(ctx, "parcelle_fermier")
	AddParcelleCode11(ctx, *flagSctlDataSource)
	FillParcelle(ctx, *flagSctlDataSource)
	FillLiensParcelleFermier(ctx, *flagSctlDataSource)
	FillLiensParcelleLieudit(ctx, *flagSctlDataSource)
}
func installUG(ctx *ctxt.Context) {
	// CreateTable(ctx, "ug")
	// CreateTable(ctx, "parcelle_ug")
	FillUG(ctx)
	FillLiensParcelleUG(ctx)
}
func installStockage(ctx *ctxt.Context) {
	CreateTable(ctx, "stockage")
	CreateTable(ctx, "stockfrais")
	CreateTable(ctx, "plaq")
	CreateTable(ctx, "tas")
	CreateTable(ctx, "humid")
	CreateTable(ctx, "humid_acteur")
	FillHangarsInitiaux(ctx)
}
func installChantier(ctx *ctxt.Context) {
	// liens pour plaq, chautre
	CreateTable(ctx, "chantier_ug")
	CreateTable(ctx, "chantier_fermier")
	CreateTable(ctx, "chantier_lieudit")
	// plaquettes
	CreateTable(ctx, "plaqop")
	CreateTable(ctx, "plaqtrans")
	CreateTable(ctx, "plaqrange")
	// autres valorisations
	CreateTable(ctx, "chautre")
	// chauffage fermier
	CreateTable(ctx, "chaufer")
	CreateTable(ctx, "chaufer_parcelle")
}
func installVente(ctx *ctxt.Context) {
	CreateTable(ctx, "venteplaq")
	CreateTable(ctx, "ventelivre")
	CreateTable(ctx, "ventecharge")
}
func installRecent(ctx *ctxt.Context) {
	CreateTable(ctx, "recent")
}

// *********************************************************
func handleFixture(ctx *ctxt.Context) {
	if *flagFixture == "stockage" {
		FillStockage(ctx)
	} else if *flagFixture == "acteur" {
		AnonymizeActeurs(ctx)
	} else {
		fmt.Println(errorMsg)
	}
}
