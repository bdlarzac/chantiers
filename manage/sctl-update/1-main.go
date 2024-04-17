/*
Mise à jour des données SCTL (uniquement les fermiers)
MAJ les tables
- fermier
- parcelle_fermier

Exemple d'utilisation : voir README

@copyright  BDL, Bois du Larzac
@license    GPL
@history    2023-01-11 05:04:02+01:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/model"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Cette commande a besoin d'un seul argument - voir fichier README")
		return
	}
	var err error
	versionSCTL := os.Args[1]
	dirname := ".." + string(os.PathSeparator) + "sctl-data" + string(os.PathSeparator) + "csv-" + versionSCTL
	_, err = os.Stat(dirname)
	if os.IsNotExist(err) {
		fmt.Printf("Directory %s/ does not exist - voir fichier README.\n", dirname)
		return
	}

	model.MustLoadEnv()
	ctxt.MustInitDB()
	ctx := ctxt.NewContext()

	report := ""
	report += updateFermiers(ctx, dirname)

	install.FillLiensParcelleFermier(ctx, versionSCTL)

	fmt.Println(report)
}

func updateFermiers(ctx *ctxt.Context, dirname string) (report string) {
	filename := dirname + string(os.PathSeparator) + "Exploita.csv"
	records, err := tiglib.CsvMap(filename, ';')
	if err != nil {
		return "Erreur de lecture de Exploita.csv avec tiglib.CsvMap()"
	}

	nUpdate, nInsert := 0, 0
	for _, record := range records {
		idExploitant, err := strconv.Atoi(record["IdExploitant"])
		if idExploitant == 1 {
			continue // "PERSONNE"
		}
		if err != nil {
			return "Erreur appel strconv.Atoi(%s) pour IdExploitant " + record["IdExploitant"]

		}
		if record["Agricole"] == "0" {
			//fmt.Printf("%d,",idExploitant)
			// pour vérifier qu'il n'y a pas de chantier associé à des fermiers non agricoles
			// select * from chantier_fermier where id_fermier in(28,37,42,43,46,48,52,53,54,58,59,63,64,66,69,70,128,129,130,73,74,77,79,85,86,90,91,92,96,97,100,102,103,105,110,113,115,116,117,119,120,121,122,134);
			//select * from parcelle_fermier where id_fermier in(28,37,42,43,46,48,52,53,54,58,59,63,64,66,69,70,128,129,130,73,74,77,79,85,86,90,91,92,96,97,100,102,103,105,110,113,115,116,117,119,120,121,122,134);
			continue // Ne garder que les agricoles
		}
		cp := record["CPExp"]
		if len(cp) > 5 {
			cp = cp[:5] // fix une typo dans la base SCTL
		}
		csvFermier := model.Fermier{
			Id:      idExploitant,
			Nom:     record["NOMEXP"],
			Prenom:  record["Prenom"],
			Adresse: record["AdresseExp"],
			Cp:      cp,
			Ville:   record["VilleExp"],
			Tel:     record["Telephone"],
			Email:   record["Mail"],
		}
		dbFermier, err := model.GetFermier(ctx.DB, idExploitant)
		if dbFermier.Id == 0 {
			err = model.InsertFermier(ctx.DB, &csvFermier)
			if err != nil {
				panic(err)
			}
			nInsert++
			report += "INSERTED " + record["IdExploitant"] + " " + record["NOMEXP"] + " " + record["Prenom"] + "\n"
		} else {
			err = model.UpdateFermier(ctx.DB, &csvFermier)
			if err != nil {
				panic(err)
			}
			nUpdate++
			report += "Updated " + record["IdExploitant"] + " - " + record["NOMEXP"] + " " + record["Prenom"] + "\n"
		}
	}
	report += "Inserted " + strconv.Itoa(nInsert) + "\n"
	report += "Updated  " + strconv.Itoa(nUpdate) + "\n"
	report += "Total    " + strconv.Itoa(nInsert+nUpdate) + "\n"
	return report
}
