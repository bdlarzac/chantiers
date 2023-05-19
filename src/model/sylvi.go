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
	//	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	//	"strconv"
	//	"strings"
	"fmt"
)

// ************************** Get many *******************************

func ComputeUGsFromFiltres(db *sqlx.DB, filtres map[string][]string) (result []*UG, err error) {
	fmt.Printf("=== model.ComputeUGsFromFiltres() - filtres = %+v\n", filtres)
	result = []*UG{}
	//
	// Filtres sur les champs de la table ug
	//
	//
	// Filtres suivants
	//
	if len(filtres["essence"]) != 0 {
		result, err = filtreUG_essence(db, result, filtres["essence"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel filtreUG_essence()")
		}
	}
	//fmt.Printf("result = %+v\n",result)
	// TODO éventuellement, trier par date
	return result, nil
}

// ************************** Application des filtres *******************************
// En entrée : liste de ventes
// En sortie : liste de ventes qui satisfont au filtre

// TODO implement
func filtreUG_essence(db *sqlx.DB, input []*UG, filtre []string) (result []*UG, err error) {
	result = []*UG{}
	/*
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
	*/
	return result, nil
}
