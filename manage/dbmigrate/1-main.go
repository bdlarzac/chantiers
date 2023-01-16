    /******************************************************************************

    Modifications (migrations) de la base BDL
    Code pas utilisé en fonctionnement normal.
    
    Lancer l'exécution en utilisant des variables d'environnement et en utilisant *.go :
    ENV_CONFIG_FILE='../../../config.env' APPLI_CONFIG_FILE='../../../config.yml' go run *.go


    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-09-26 17:41:35+02:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"

	// "io/ioutil"
	// "gopkg.in/yaml.v3"
	"bdl.local/bdl/model"
)

var possibleMigrations []string

// *********************************************************
func main() {

	if len(os.Args) != 2 {
	    fmt.Println("Cette commande a besoin d'un seul argument - voir fichier README")
	    return
	}
	
	migration := os.Args[1]
	possibleMigrations = computeMigrations()
	if !tiglib.InArrayString(migration, possibleMigrations) {
		fmt.Println("MIGRATION INEXISTANTE : " + migration)
		fmt.Println("Migrations possibles : " + strings.Join(possibleMigrations, " "))
		return
	}
	
	model.MustLoadEnv()
	ctxt.MustLoadConfig()
	ctxt.MustInitDB()
	ctx := ctxt.NewContext()

	switch migration {
	case "Migrate_2021_03_01_exemple":
		Migrate_2021_03_01_exemple(ctx)
	case "Migrate_2021_11_10_note_plaq":
		Migrate_2021_11_10_note_plaq(ctx)
	case "Migrate_2022_01_10_facture_vente_km_map":
		Migrate_2022_01_10_facture_vente_km_map(ctx)
	case "Migrate_2022_02_07_unite_piquets":
		Migrate_2022_02_07_unite_piquets(ctx)
	case "Migrate_2022_09_24_km_livraison":
		Migrate_2022_09_24_km_livraison(ctx)
	case "Migrate_2023_01_16_fix_parcelle":
	    Migrate_2023_01_16_fix_parcelle(ctx)
	    
	case "Migrate_2023_01_chantier_parcelle":
		Migrate_2023_01_chantier_parcelle(ctx)
	}
    
}

// *********************************************************
// Renvoie la liste des migrations possibles
// = liste des fonctions du répertoire courant commençant par Migrate_
// Ça n'a pas l'air possible avec reflect => bidouille avec regex
func computeMigrations() (res []string) {
	out, err := exec.Command("grep", "-rn", "func Migrate_", ".").Output()
	if err != nil {
		panic(err)
	}
	r := regexp.MustCompile(`func (Migrate_.*?)\s*\(`)
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		m := r.FindStringSubmatch(line)
		res = append(res, m[1])
	}
	return res
}
