/*
Génère un fichier CSV avec la liste des ugs et leurs chantiers.
A la demande de Maël, stagiaire Forêt Evolution.

Voir https://github.com/bdlarzac/chantiers/issues/26
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
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/model"
	"cmp"
	"fmt"
	"io/ioutil"
	"slices"
	"strconv"
	"strings"
	"time"
)

type ligneUG struct {
	CodeUG        string
	CodeTypo      string
	Coupe         string
	Valo          string
	TitreChantier string
	DateChantier  time.Time
	Essence       string
	Volume        float64
	Unite         string
	Prix          float64
	Exploitation  string
}

func liste_ugs(ctx *ctxt.Context) {
	res := []ligneUG{}
	db := ctx.DB
	stmt, _ := db.Prepare(`select * from chantier_ug`)
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
			ch.ComputePrixTotalHT()
			newLine.Valo = ch.TypeValo
			newLine.TitreChantier = ch.String()
			newLine.DateChantier = ch.DateContrat
			newLine.Essence = ch.Essence
			newLine.Volume = ch.VolumeRealise
			newLine.Unite = ch.Unite
			newLine.Prix = ch.PrixTotalHT
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
			ch, _ = model.GetPlaqFull(db, id_chantier)
			_ = ch.ComputeVolume(db)
			_ = ch.ComputeCouts(db, ctx.Config)
			newLine.Valo = "PQ"
			newLine.TitreChantier = ch.String()
			newLine.DateChantier = ch.DateDebut
			newLine.Essence = "PS"
			newLine.Volume = ch.Volume
			newLine.Unite = "MA"
			newLine.Prix = ch.CoutTotal.Total
			newLine.Exploitation = ch.Exploitation
		}
		res = append(res, newLine)
	}

	// Tri par code UG
	slices.SortFunc[[]ligneUG, ligneUG](res, func(a, b ligneUG) int {
		tmpA := strings.Split(a.CodeUG, "-")
		tmpB := strings.Split(b.CodeUG, "-")
		// teste chiffres romains
		idxA := tiglib.ArraySearch(model.RomanNumbers, tmpA[0])
		idxB := tiglib.ArraySearch(model.RomanNumbers, tmpB[0])
		if idxA < idxB {
			return -1
		}
		if idxA > idxB {
			return 1
		}
		// idxA = idxB => chiffres romains identiques
		nA, _ := strconv.Atoi(tmpA[1])
		nB, _ := strconv.Atoi(tmpB[1])
		return cmp.Compare(nA, nB)
	})

	// Génère fichier CSV
	SEP := ";"
	OUTFILE := "../../tmp/ug-chantier.csv"

	csv := strings.Join([]string{"UG", "TYPO", "COUPE", "VALO", "TITRE CHANTIER", "DATE", "ESS", "VOLUME", "UNITE", "PRIX", "EX"}, SEP) + "\n"
	for _, line := range res {
		csv += strings.Join([]string{
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
	_ = ioutil.WriteFile(OUTFILE, []byte(csv), 0755)
	fmt.Println("Fichier csv généré : " + OUTFILE)
}
