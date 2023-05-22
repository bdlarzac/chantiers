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

func ComputeUGsFromFiltres(db *sqlx.DB, filtres map[string][]string) (result []*UG, err error) {
fmt.Printf("=== model.ComputeUGsFromFiltres() - filtres = %+v\n", filtres)
	result = []*UG{}
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
		result, err = filtreUG_fermier(db, result, filtres["fermier"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel filtreUG_fermier()")
		}
	}
	// Compute liés à l'affichage si pas déjà calculé pour les filtres
	if len(filtres["essence"]) == 0 {
        for _, ug := range(result){
            err = ug.ComputeEssences(db)
            if err != nil {
                return result, werr.Wrapf(err, "Erreur appel UG.ComputeEssences()")
            }
        }
    }
    //
	if len(filtres["fermier"]) == 0 {
        for _, ug := range(result){
            err = ug.ComputeFermiers(db)
            if err != nil {
                return result, werr.Wrapf(err, "Erreur appel UG.ComputeFermiers()")
            }
        }
	}
	//
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

