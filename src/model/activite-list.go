/*
Calcul d'activités, en appliquant des filtres

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-04-21 09:00:39+02:00, Thierry Graff : Creation
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

// ************************** Fonction principale *******************************

func ComputeActivitesFromFiltres(db *sqlx.DB, filtres map[string][]string) (result []*Activite, err error) {
	fmt.Printf("ComputeActivitesFromFiltres() - filtres = %+v\n", filtres)
	result = []*Activite{}
	//
	// Première sélection, par filtre période
	//
	var tmp []*Activite
	// if plaq dans le filtre activite
	tmp, err = computePlaqFromFiltrePeriode(db, filtres["periode"])
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel computePlaqFromFiltrePeriode(\"plaq\")")
	}
	result = append(result, tmp...)
	//
	tmp, err = computeChautreFromFiltrePeriode(db, filtres["periode"])
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel computeChautreFromFiltrePeriode(\"plaq\")")
	}
	result = append(result, tmp...)
	//
	tmp, err = computeChauferFromFiltrePeriode(db, filtres["periode"])
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel computeChauferFromFiltrePeriode(\"plaq\")")
	}
	result = append(result, tmp...)
	//
	// Filtres suivants
	//
	if len(filtres["essence"]) != 0 {
		result = filtreEssence(db, result, filtres["essence"])
	}
	//
	if len(filtres["fermier"]) != 0 {
		for _, activite := range result {
			activite.ComputeFermiers(db)
		}
		result = filtreFermier(db, result, filtres["fermier"])
	}
	if len(filtres["ug"]) != 0 {
		for _, activite := range result {
			activite.ComputeUGs(db)
		}
		result = filtreFermier(db, result, filtres["ug"])
	}
	//
	// préparation (faire le plus tard possible pour optimiser)
	//
	if len(filtres["proprio"]) != 0 || len(filtres["parcelle"]) != 0 {
		for _, activite := range result {
			activite.ComputeLiensParcelles(db)
		}
	}
	//
	if len(filtres["parcelle"]) != 0 {
		result = filtreParcelle(db, result, filtres["parcelle"])
	}
	//
	if len(filtres["proprio"]) != 0 {
		result, err = filtreProprio(db, result, filtres["proprio"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel filtreProprio()")
		}
	}
	//fmt.Printf("result = %+v\n",result)
	// TODO éventuellement, trier par date
	return result, nil
}

// ************************** Selection initiale, par période *******************************

func computePlaqFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (result []*Activite, err error) {
	var query string
	chantiers := []*Plaq{}
	query = "select * from plaq"
	if len(filtrePeriode) == 2 {
		query += " where datedeb >= $1 and datedeb <= $2"
		err = db.Select(&chantiers, query, filtrePeriode[0], filtrePeriode[1])
	} else {
		err = db.Select(&chantiers, query)
	}
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chantier := range chantiers {
		tmp, err := GetActivite(db, "plaq", chantier.Id)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel GetActivite(\"plaq\", "+strconv.Itoa(chantier.Id)+")")
		}
		result = append(result, tmp)
	}
	return result, nil
}

func computeChautreFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (result []*Activite, err error) {
	var query string
	chantiers := []*Chautre{}
	query = "select * from chautre"
	if len(filtrePeriode) == 2 {
		query += " where datecontrat >= $1 and datecontrat <= $2"
		err = db.Select(&chantiers, query, filtrePeriode[0], filtrePeriode[1])
	} else {
		err = db.Select(&chantiers, query)
	}
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chantier := range chantiers {
		tmp, err := GetActivite(db, "chautre", chantier.Id)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel GetActivite(\"chautre\", "+strconv.Itoa(chantier.Id)+")")
		}
		result = append(result, tmp)
	}
	return result, nil
}

func computeChauferFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (result []*Activite, err error) {
	var query string
	chantiers := []*Chaufer{}
	query = "select * from chaufer"
	if len(filtrePeriode) == 2 {
		query += " where datechantier >= $1 and datechantier <= $2"
		err = db.Select(&chantiers, query, filtrePeriode[0], filtrePeriode[1])
	} else {
		err = db.Select(&chantiers, query)
	}
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chantier := range chantiers {
		tmp, err := GetActivite(db, "chaufer", chantier.Id)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel GetActivite(\"chaufer\", "+strconv.Itoa(chantier.Id)+")")
		}
		result = append(result, tmp)
	}
	return result, nil
}

// ************************** Filtres *******************************

func filtreEssence(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			if a.CodeEssence == f {
				result = append(result, a)
				break
			}
		}
	}
	return result
}

func filtreFermier(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, fermier := range a.Fermiers {
				if fermier.Id == id {
					result = append(result, a)
					break
				}
			}
		}
	}
	return result
}

func filtreUGs(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, ug := range a.UGs {
				if ug.Id == id {
					result = append(result, a)
					break
				}
			}
		}
	}
	return result
}

func filtreParcelle(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, lienParcelle := range a.LiensParcelles {
				if lienParcelle.IdParcelle == id {
					result = append(result, a)
					break
				}
			}
		}
	}
	return result
}

func filtreProprio(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite, err error) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, lienParcelle := range a.LiensParcelles {
				parcelle, err := GetParcelle(db, lienParcelle.IdParcelle)
				if err != nil {
					return result, werr.Wrapf(err, "Erreur appel GetParcelle()")
				}

				if parcelle.IdProprietaire == id {
					result = append(result, a)
					break
				}
			}
		}
	}
	return result, nil
}
