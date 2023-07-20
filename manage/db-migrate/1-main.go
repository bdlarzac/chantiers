/*
*****************************************************************************

	Modifications (migrations) de la base BDL
	Code pas utilisé en fonctionnement normal.

	Voir fichier README pour ajouter ou exécuter une migration

	@copyright  BDL, Bois du Larzac
	@license    GPL
	@history    2019-09-26 17:41:35+02:00, Thierry Graff : Creation

*******************************************************************************
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/model"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

var possibleMigrations []string

// *********************************************************
func main() {
	possibleMigrations = computeMigrations()
	msgPossibles := "Migrations possibles : \n    " + strings.Join(possibleMigrations, "\n    ")
	if len(os.Args) != 2 {
		fmt.Println("Cette commande a besoin d'un argument, la migration à exécuter")
		fmt.Println(msgPossibles)
		return
	}
	migration := os.Args[1]
	if !tiglib.InArray(migration, possibleMigrations) {
		fmt.Println("MIGRATION INEXISTANTE : " + migration)
		fmt.Println("Modifier 1.main.go pour la rajouter dans le switch : ")
		fmt.Println(msgPossibles)
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
	case "Migrate_2023_01_20_fix_non_agricoles":
		Migrate_2023_01_20_fix_non_agricoles(ctx)
	case "Migrate_2023_01_23_chantier_parcelle":
		Migrate_2023_01_23_chantier_parcelle(ctx)
	case "Migrate_2023_02_20_titre_chantier":
		Migrate_2023_02_20_titre_chantier(ctx)
	case "Migrate_2023_02_21_num_facture":
		Migrate_2023_02_21_num_facture(ctx)
		// A partir d'ici, nouvelle convention
		// le nom de la fonction est suivi du numéro d'issue sur github
	case "Migrate_2023_02_24_details_ug__15":
		Migrate_2023_02_24_details_ug__15(ctx)
	case "Migrate_2023_04_03_role_acteur__16":
		Migrate_2023_04_03_role_acteur__16(ctx)
	case "Migrate_2023_05_18_clean_types__19":
		Migrate_2023_05_18_clean_types__19(ctx)
	case "Migrate_2023_05_19_ug_again__18":
		Migrate_2023_05_19_ug_again__18(ctx)
	case "Migrate_2023_05_22_non_agricoles__20":
		Migrate_2023_05_22_non_agricoles__20(ctx)
	case "Migrate_2023_05_25_fix_code_ug__21":
		Migrate_2023_05_25_fix_code_ug__21(ctx)
	case "Migrate_2023_05_26_date_paiement__22":
		Migrate_2023_05_26_date_paiement__22(ctx)
	case "Migrate_2023_06_21_ajout_roles":
		Migrate_2023_06_21_ajout_roles(ctx)
	default:
		fmt.Println("Migration inconnue : " + migration)
		fmt.Println("Modifier 1.main.go pour la rajouter dans le switch")
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
		if !strings.Contains(m[1], "Output") { // virer la ligne avec exec.Command du grep
			res = append(res, m[1])
		}
	}
	sort.Strings(res)
	return res
}
