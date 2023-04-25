/*
Calcul d'activités regroupées par saison

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-04-25 18:44:57+02:00, Thierry Graff : Creation
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"time"
"fmt"
)

/* Liste d'activités ayant lieu dans une période donnée */
type ActiviteParSaison struct {
	Datedeb   time.Time
	Datefin   time.Time
	Activites []*Activite
}

type BilanActivitesParSaison struct {
	Datedeb                       time.Time
	Datefin                       time.Time
	TotalActivitesParValo []*TotalActivitesParValo
}

type TotalActivitesParValo struct {
	TypeValo string
	Volume   float64
	PrixHT   float64
}

func ComputeBilanActivitesParSaison(db *sqlx.DB, debutSaison string, activites []*Activite) (result []*BilanActivitesParSaison, err error) {
    activitesParSaison, err := ComputeActivitesParSaison(db, debutSaison, activites)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel ComputeActivitesParSaison()")
	}
	result = []*BilanActivitesParSaison{}
	for _, activiteParSaison := range(activitesParSaison){
	    if len(activiteParSaison.Activites) == 0 {
	        continue
	    }
	    currentRes := BilanActivitesParSaison{Datedeb: activiteParSaison.Datedeb, Datefin: activiteParSaison.Datefin, TotalActivitesParValo:[]*TotalActivitesParValo{}}
	    // map intermédiaire
        mapValos := map[string]TotalActivitesParValo{}
        for _, activite := range(activiteParSaison.Activites){
            valo := activite.TypeValo
            if _, ok := mapValos[valo]; !ok {
                mapValos[valo] = TotalActivitesParValo{TypeValo: valo}
            }
            mapValos[valo].Volume += activite.Volume       /////////////////// ici faire conversion d'unité ///////////////////
            mapValos[valo].PrixHT += activite.PrixHT
        }
        // utilise map pour remplir le tableau
        for valo, total := range(mapValos){
            newRes := TotalActivitesParValo{TypeValo: valo, Volume: total.Volume, PrixHT: total.PrixHT}
            currentRes = append(currentRes, newRes)
        }
        result = append(result, currentRes)
	}
fmt.Printf("ComputeBilansActivitesParSaison()\n")

fmt.Printf("result = %+v\n",result)
	return result, nil
}

func ComputeActivitesParSaison(db *sqlx.DB, debutSaison string, activites []*Activite) (result []*ActiviteParSaison, err error) {
fmt.Printf("ComputeBilansActivitesParSaison()\n")
    limites, _, err := ComputeLimitesSaisons(db, debutSaison)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel ComputeLimitesSaisons()")
	}
	result = []*ActiviteParSaison{}
    for _, limite := range(limites){
        newRes := ActiviteParSaison{Datedeb: limite[0], Datefin: limite[0], Activites: []*Activite{}}
        result = append(result, &newRes)
    }
    for _, activite := range(activites){
        for _, res := range(result){
            if (activite.DateActivite.After(res.Datedeb) || activite.DateActivite.Equal(res.Datedeb)) && activite.DateActivite.Before(res.Datefin){
                res.Activites = append(res.Activites, activite)
                break
            }
        }
    }
	return result, nil
}

