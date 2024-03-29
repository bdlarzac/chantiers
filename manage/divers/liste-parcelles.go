/*
Génère un fichier CSV avec la liste des parcelles et leurs chantiers.

Voir https://github.com/bdlarzac/chantiers/issues/26
- code parcelle 11
- date format ISO
- type chantier (valorisation)
- titre chantier
- essence
- surface (ha)
- volume coupé
- unité volume
- prix
- exploitation (1, 2...)
- proprietaire
- id SCTL fermier
- fermiers

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2024-03-28, Thierry Graff : Creation
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

type ligneParcelle struct {
	CodeParcelle11 string
	DateChantier   time.Time
	Valo           string
	TitreChantier  string
	Essence        string
	Surface        float64
	Volume         float64
	Unite          string
	Prix           float64
	Exploitation   string
	Proprietaire   string
	Fermiers       string
}

func liste_parcelles(ctx *ctxt.Context) {
	res := []ligneParcelle{}
	db := ctx.DB
	stmt, _ := db.Prepare(`select * from chantier_parcelle`)
	defer stmt.Close()
	rows, _ := stmt.Query()
	defer rows.Close()

	var type_chantier string
	var id_chantier, id_parcelle int
	var parcelle *model.Parcelle
	var entiere bool
	var surface float64

	for rows.Next() {
		var newLine ligneParcelle
		_ = rows.Scan(&id_chantier, &id_parcelle, &entiere, &surface, &type_chantier)
		//
		parcelle, _ = model.GetParcelle(db, id_parcelle)
		_ = parcelle.ComputeCommune(db)
		_ = parcelle.ComputeProprietaire(db)
		_ = parcelle.ComputeFermiers(db)
		_ = parcelle.ComputeCode11(db)
		newLine.CodeParcelle11 = parcelle.Code11
		newLine.Proprietaire = parcelle.Proprietaire.Nom
		if entiere {
			newLine.Surface = parcelle.Surface
		} else {
			newLine.Surface = surface
		}
		fermiers := []string{}
		for _, f := range parcelle.Fermiers {
			fermiers = append(fermiers, "["+strconv.Itoa(f.Id)+"] "+model.NomEtAdresseFermier(f))
		}
		newLine.Fermiers = strings.Join(fermiers, " + ")
		//
		switch type_chantier {
		case "chautre":
			var ch *model.Chautre
			ch, _ = model.GetChautre(db, id_chantier)
			ch.ComputePrixTotalHT()
			newLine.DateChantier = ch.DateContrat
			newLine.Valo = ch.TypeValo
			newLine.TitreChantier = ch.String()
			newLine.Essence = ch.Essence
			newLine.Volume = ch.VolumeRealise
			newLine.Unite = ch.Unite
			newLine.Prix = ch.PrixTotalHT
			newLine.Exploitation = ch.Exploitation
		case "chaufer":
			var ch *model.Chaufer
			ch, _ = model.GetChaufer(db, id_chantier)
			newLine.DateChantier = ch.DateChantier
			newLine.Valo = "CF"
			newLine.TitreChantier = ch.String()
			newLine.Essence = "PS"
			newLine.Volume = ch.Volume
			newLine.Unite = ch.Unite
			newLine.Prix = 0
			newLine.Exploitation = ch.Exploitation
			ch, _ = model.GetChaufer(db, id_chantier)
		case "plaq":
			var ch *model.Plaq
			ch, _ = model.GetPlaqFull(db, id_chantier)
			_ = ch.ComputeVolume(db)
			_ = ch.ComputeCouts(db, ctx.Config)
			newLine.DateChantier = ch.DateDebut
			newLine.Valo = "PQ"
			newLine.TitreChantier = ch.String()
			newLine.Essence = ch.Essence
			newLine.Volume = ch.Volume
			newLine.Unite = "MA"
			newLine.Prix = ch.CoutTotal.Total
			newLine.Exploitation = ch.Exploitation
		}
		res = append(res, newLine)
	}
	// Tri par code Code parcelle 11
	slices.SortFunc[[]ligneParcelle, ligneParcelle](res, func(a, b ligneParcelle) int {
		return cmp.Compare(a.CodeParcelle11, b.CodeParcelle11)
	})

	// Génère fichier CSV
	SEP := ";"
	OUTFILE := "../../tmp/parcelle-chantier.csv"

	csv := strings.Join([]string{"PARCELLE", "DATE", "VALO", "TITRE CHANTIER", "ESS", "SURF", "VOLUME", "UNITE", "PRIX", "EX", "PROPRIO", "FERMIERS"}, SEP) + "\n"
	for _, line := range res {
		csv += strings.Join([]string{
			line.CodeParcelle11,
			tiglib.DateIso(line.DateChantier),
			line.Valo,
			line.TitreChantier,
			line.Essence,
			strconv.FormatFloat(line.Surface, 'f', 2, 64),
			strconv.FormatFloat(line.Volume, 'f', 2, 64),
			line.Unite,
			strconv.FormatFloat(line.Prix, 'f', 2, 64),
			line.Exploitation,
			line.Proprietaire,
			line.Fermiers,
		}, SEP) + "\n"
	}
	_ = ioutil.WriteFile(OUTFILE, []byte(csv), 0755)
	fmt.Println("Fichier csv généré : " + OUTFILE)
}
