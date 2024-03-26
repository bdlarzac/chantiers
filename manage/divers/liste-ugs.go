/*
Génère un fichier CSV avec la liste des ugs et leurs chantiers.
- code UG
- typo
- type de coupe
- valorisation
- titre chantier
- date chantier (format ISO)
- essence
- volume coupé
- unité volume
- prix
- exploitation (1, 2...)
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2024-03-26, Thierry Graff : Creation
*/
package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"time"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"bdl.local/bdl/generic/tiglib"
)

type ligneUG struct {
    CodeUG string
    CodeTypo string
    Coupe string
    Valo string
    TitreChantier string
    DateChantier time.Time
    Essence string
    Volume float64
    Unite string
    Prix float64
    Exploitation string
}

func liste_ugs(ctx *ctxt.Context){
    
    lines := []ligneUG{}
	db := ctx.DB
	stmt, _ := db.Prepare(`select * from chantier_ug order by id_ug`)
	defer stmt.Close()
    
	rows, _ := stmt.Query()
	defer rows.Close()
	var type_chantier string
	var id_chantier, id_ug int
	var ug *model.UG
	for rows.Next() {
	    var newLine ligneUG
		_ = rows.Scan(&type_chantier, &id_chantier, &id_ug)
		//
		ug, _ = model.GetUG(db, id_ug)
		newLine.CodeUG = ug.Code
		newLine.CodeTypo = ug.CodeTypo
		newLine.Coupe = ug.Coupe
		//
		switch type_chantier {
		case "chautre":
		    var ch *model.Chautre
		    ch, _ = model.GetChautre(db, id_chantier)
		    newLine.Valo = ch.TypeValo
		    newLine.TitreChantier = ch.String()
		    newLine.DateChantier = ch.DateContrat
		    newLine.Essence = ch.Essence
		    newLine.Volume = ch.VolumeRealise
		    newLine.Unite = ch.Unite
		    newLine.Prix = ch.PUHT
		    newLine.Exploitation = ch.Exploitation
		case "chaufer":
		    var ch *model.Chaufer
		    ch, _ = model.GetChaufer(db, id_chantier)
		    newLine.Valo = "CF"
		    newLine.TitreChantier = ch.String()
		    newLine.DateChantier = ch.DateChantier
		    newLine.Essence = ch.Essence
		    newLine.Volume = ch.Volume
		    newLine.Unite = ch.Unite
		    newLine.Prix = 0
		    newLine.Exploitation = ch.Exploitation
		case "plaq":
		    var ch *model.Plaq
		    ch, _ = model.GetPlaq(db, id_chantier)
		    _ = ch.ComputeVolume(db)
		    _ = ch.ComputeCouts(db, ctx.Config)
		    newLine.Valo = "PQ"
		    newLine.TitreChantier = ch.String()
		    newLine.DateChantier = ch.DateDebut
		    newLine.Essence = "PS"
		    newLine.Volume = ch.Volume
		    newLine.Unite = "MA"
		    //newLine.Prix = ch.CoutTotal.Total
		    newLine.Exploitation = ch.Exploitation
		}
		lines = append(lines, newLine)
//break
/* 
calcul prix total pour chautre
trier ugs
*/
	}
	
    SEP := ";"
    OUTFILE := "../../tmp/ug-chantier.csv"
//    OUTFILE := "ug-chantier.csv"

	res := strings.Join([]string{"UG", "TYPO", "COUPE", "VALO", "CHANTIER", "DATE", "ESS", "VOLUME", "UNITE", "PRIX", "EX"}, SEP) + "\n"
	for _, line := range(lines){
	    res += strings.Join([]string{
            line.CodeUG,
            line.CodeTypo,
            line.Coupe,
            line.Valo,
            line.TitreChantier,
            tiglib.DateIso(line.DateChantier),
            line.Essence,
            strconv.FormatFloat(line.Volume, 'f', 2, 64),
            line.Unite,
            strconv.FormatFloat(line.Prix, 'f', 2, 64),
            line.Exploitation,
	    }, SEP) + "\n"
	}
//fmt.Println(res)
//	fmt.Printf("%+v\n",lines)
    _ = ioutil.WriteFile(OUTFILE, []byte(res), 0755)
	fmt.Println("Fichier csv généré : " + OUTFILE)
}
