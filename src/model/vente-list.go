/*
Calcul d'activités, en appliquant des filtres

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-04-21 09:00:39+02:00, Thierry Graff : Creation
*/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

// ************************** Fonction principale *******************************

func ComputeVentesFromFiltres(db *sqlx.DB, filtres map[string][]string) (result []*Vente, err error) {
fmt.Printf("ComputeVentesFromFiltres() - filtres = %+v\n", filtres)
	result = []*Vente{}
	//
	// Première sélection, par filtre période
	//
	var tmp []*Vente
	// Si les ventes plaquettes sont demandées
	if len(filtres["valo"]) == 0 || tiglib.InArrayString("PQ", filtres["valo"]) {
        tmp, err = computeVentePlaqFromFiltrePeriode(db, filtres["periode"])
        if err != nil {
            return result, werr.Wrapf(err, "Erreur appel computePlaqFromFiltrePeriode()")
        }
        result = append(result, tmp...)
    }
	// Si d'autres valorisations que les plaquettes sont demandées
	if len(filtres["valo"]) == 0 || filtreValoContientAutreValo(filtres["valo"]) {
        tmp, err = computeVenteChautreFromFiltrePeriode(db, filtres["periode"])
        if err != nil {
            return result, werr.Wrapf(err, "Erreur appel computeVenteChautreFromFiltrePeriode()")
        }
        result = append(result, tmp...)
    }
/* 
	//
	// Filtres suivants
	//
	if len(filtres["essence"]) != 0 {
		result = filtreVente_essence(db, result, filtres["essence"])
	}
	//
	if len(filtres["proprio"]) != 0 {
		result, err = filtreVente_proprio(db, result, filtres["proprio"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel filtreVente_proprio()")
		}
	}
	//fmt.Printf("result = %+v\n",result)
	// TODO éventuellement, trier par date
*/
	return result, nil
}

// ************************** Selection initiale, par période *******************************

func computeVentePlaqFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (result []*Vente, err error) {
	var query string
	ventes := []*VentePlaq{}
	query = "select * from venteplaq"
	if len(filtrePeriode) == 2 {
		query += " where datevente >= $1 and datevente <= $2"
		err = db.Select(&ventes, query, filtrePeriode[0], filtrePeriode[1])
	} else {
		err = db.Select(&ventes, query)
	}
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, vp := range ventes {
		tmp, err := GetVente(db, "plaq", vp.Id)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel GetVente(\"plaq\", "+strconv.Itoa(vp.Id)+")")
		}
		result = append(result, tmp)
	}
	return result, nil
}

func computeVenteChautreFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (result []*Vente, err error) {
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
	for _, ch := range chantiers {
		tmp, err := GetVente(db, "autre", ch.Id)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel GetVente(\"autre\", "+strconv.Itoa(ch.Id)+")")
		}
		result = append(result, tmp)
	}
	return result, nil
}

// ************************** Filtres *******************************

func filtreVente_essence(db *sqlx.DB, input []*Vente, filtre []string) (result []*Vente) {
	result = []*Vente{}
	/* 
	for _, a := range input {
		for _, f := range filtre {
			if a.CodeEssence == f {
				result = append(result, a)
				break
			}
		}
	}
	*/
	return result
}

func filtreVente_proprio(db *sqlx.DB, input []*Vente, filtre []string) (result []*Vente, err error) {
	result = []*Vente{}
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
