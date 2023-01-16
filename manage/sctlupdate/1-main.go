/******************************************************************************

    Mise à jour des données SCTL (uniquement les fermiers)
    Exemple d'utilisation : voir README

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2023-01-11 05:04:02+01:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"os"
	"fmt"
	"strconv"
	"bdl.local/bdl/model"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
)

func main() {
	if len(os.Args) != 2 {
	    fmt.Println("Cette commande a besoin d'un seul argument - voir fichier README")
	    return
	}
	var err error
	dirname := ".." + string(os.PathSeparator) + "sctl-data" + string(os.PathSeparator) + "csv-"+os.Args[1]
    _, err = os.Stat(dirname)
    if os.IsNotExist(err) {
        fmt.Printf("Directory %s/ does not exist - voir fichier README.\n", dirname)
        return
    }
    filename := dirname + string(os.PathSeparator) + "Exploita.csv"
	records, err := tiglib.CsvMap(filename, ';')
	if err != nil {
		fmt.Println("Erreur de lecture de Exploita.csv avec tiglib.CsvMap()")
		return
	}
	
	model.MustLoadEnv()
	ctxt.MustInitDB()
	ctx := ctxt.NewContext()
	
	report := ""
	nUpdate, nInsert := 0, 0
	for _, record := range records {
		idExploitant, err := strconv.Atoi(record["IdExploitant"])
		if idExploitant == 1 {
		    continue // "PERSONNE"
		}
		if err != nil {
			fmt.Printf("Erreur appel strconv.Atoi(%s) pour IdExploitant", record["IdExploitant"])
			return
		}
		cp := record["CPExp"]
		if len(cp) > 5 {
			cp = cp[:5] // fix une typo dans la base SCTL
		}
		csvFermier := model.Fermier{
			Id:       idExploitant,
			Nom:      record["NOMEXP"],
			Prenom:   record["Prenom"],
			Adresse:  record["AdresseExp"],
			Cp:       cp,
			Ville:    record["VilleExp"],
			Tel:      record["Telephone"],
			Email:    record["Mail"],
		}
		dbFermier, err := model.GetFermier(ctx.DB, idExploitant)
		if dbFermier.Id == 0 {
		    err = model.InsertFermier(ctx.DB, &csvFermier)
            if err != nil {
                panic(err)
            }
		    nInsert++
		    report += "INSERTED " + record["NOMEXP"] + " " + record["Prenom"] + "\n"
		} else {
		    err = model.UpdateFermier(ctx.DB, &csvFermier)
            if err != nil {
                panic(err)
            }
		    nUpdate++
		    report += "Updated " + record["IdExploitant"] + " - "+ record["NOMEXP"] + " " + record["Prenom"] + "\n"
		}
	}
	report += "Inserted " + strconv.Itoa(nInsert) + "\n"
	report += "Updated  " + strconv.Itoa(nUpdate) + "\n"
	report += "Total    " + strconv.Itoa(nInsert + nUpdate) + "\n"
	fmt.Println(report)
}