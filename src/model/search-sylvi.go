/*
*****************************************************************************

	Code relatif à la recherche sylviculture.
	Manipule des UGs.

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-05-12 17:56:06+02:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	//"strings"
	"fmt"
)

// ************************** Get many *******************************

// fonction à optimiser pour n'appeler filtreUG_sansfiltre() que s'il n'y a aucun filtre
// mais attendre que la demande de BDL soit stabilisée
func ComputeUGsFromFiltres(db *sqlx.DB, filtres map[string][]string) (result []*UG, err error) {
fmt.Printf("=== model.ComputeUGsFromFiltres() - filtres = %+v\n", filtres)
	result = []*UG{}
	// Booléens indiquant si les Compute*() ont été appelés (pas si les filtres ont été appliqués)
	essenceDone := false
	fermierDone := false
	//communeDone := false // commenté car bug bizarre à la compil
	//
	// Filtres sur les champs de la table ug
	//
    result, err = filtreUG_sansfiltre(db)
    if err != nil {
        return result, werr.Wrapf(err, "Erreur appel filtreUG_sansfiltre()")
    }
	//
	// Filtres suivants
	//
	if len(filtres["essence"]) != 0 {
        for _, ug := range(result){
            err = ug.ComputeEssences(db)
            if err != nil {
                return result, werr.Wrapf(err, "Erreur appel UG.ComputeEssences()")
            }
        }
        essenceDone = true
        result, err = filtreUG_essence(db, result, filtres["essence"])
        if err != nil {
            return result, werr.Wrapf(err, "Erreur appel filtreUG_essence()")
        }
    }
    //
	if len(filtres["fermier"]) != 0 {
        for _, ug := range(result){
            err = ug.ComputeFermiers(db)
            if err != nil {
                return result, werr.Wrapf(err, "Erreur appel UG.ComputeFermiers()")
            }
        }
        fermierDone = true
		result, err = filtreUG_fermier(db, result, filtres["fermier"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel filtreUG_fermier()")
		}
	}
    //
	if len(filtres["commune"]) != 0 {
        for _, ug := range(result){
            err = ug.ComputeCommunes(db)
            if err != nil {
                return result, werr.Wrapf(err, "Erreur appel UG.ComputeCommunes()")
            }
        }
        //communeDone = true // commenté car bug bizarre à la compil
		result, err = filtreUG_commune(db, result, filtres["commune"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel filtreUG_commune()")
		}
	}
	//
	// Compute nécessaires pour l'affichage si pas déjà calculé pour les filtres
	//
	if !essenceDone {
        for _, ug := range(result){
            err = ug.ComputeEssences(db)
            if err != nil {
                return result, werr.Wrapf(err, "Erreur appel UG.ComputeEssences()")
            }
        }
        essenceDone = true
    }
    //
	if !fermierDone {
        for _, ug := range(result){
            err = ug.ComputeFermiers(db)
            if err != nil {
                return result, werr.Wrapf(err, "Erreur appel UG.ComputeFermiers()")
            }
        }
        fermierDone = true
	}
	// Calcul d'activités, à faire à la fin
    for _, ug := range(result){
        err = ug.ComputeActivites(db)
        if err != nil {
            return result, werr.Wrapf(err, "Erreur appel UG.ComputeActivites()")
        }
    }
//fmt.Printf("result = %+v\n",result)
	return result, nil
}

// ****************************************************************************************************
// ************************** Auxiliaires de ComputeActivitesFromFiltres() ****************************
// ****************************************************************************************************

// ************************** Selection initiale, par champs de la table ug *******************************

func filtreUG_sansfiltre(db *sqlx.DB) (result []*UG, err error) {
	result = []*UG{}
	query := "select * from ug"
    err = db.Select(&result, query)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	return result, nil
}

// ************************** Filtres *******************************
// En entrée : liste d'UGs
// En sortie : liste d'UGs qui satisfont au filtre

func filtreUG_essence(db *sqlx.DB, input []*UG, filtre []string) (result []*UG, err error) {
	result = []*UG{}
	for _, ug := range input {
        for _, codeEssence := range(ug.CodesEssence){
            if tiglib.InArray(codeEssence, filtre){
                result = append(result, ug)
                break
            }
        }
	}
	return result, err
}

func filtreUG_fermier(db *sqlx.DB, input []*UG, filtre []string) (result []*UG, err error) {
	result = []*UG{}
	idFermier, _ := strconv.Atoi(filtre[0])
	for _, ug := range input {
        for _, fermier := range(ug.Fermiers){
            if fermier.Id == idFermier {
                result = append(result, ug)
                break
            }
        }
	}
	return result, err
}

func filtreUG_commune(db *sqlx.DB, input []*UG, filtre []string) (result []*UG, err error) {
	result = []*UG{}
	idCommune, _ := strconv.Atoi(filtre[0])
	for _, ug := range input {
        for _, commune := range(ug.Communes){
            if commune.Id == idCommune {
                result = append(result, ug)
                break
            }
        }
	}
	return result, err
}

