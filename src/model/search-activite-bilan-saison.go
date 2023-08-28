/*
Calcul d'activités regroupées par saison

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-04-25 18:44:57+02:00, Thierry Graff : Creation
*/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"time"
//	"fmt"
)

// Liste d'activités ayant lieu dans une période donnée
type ActiviteParSaison struct {
	Datedeb   time.Time
	Datefin   time.Time
	Activites []*Activite
}

// Contient les bilans pour une saison donnée
// Pour simplifier l'affichage, les activités plaquettes sont stockées séparément
type BilanActivitesParSaison struct {
	Datedeb                            time.Time
	Datefin                            time.Time
	// toutes activités sauf plaquettes
	TotalActivitesParValoEtProprio     map[string]map[int]VolumePrixHT // map[code valo][id proprio]
	// uniquement plaquettes
	TotalActivitesPlaquettesParProprio map[int]VolumePrixHT // key = id proprio
	TotalVentePlaquettesParProprio     map[int]float64     // key = id proprio - value = total vendu
}

type VolumePrixHT struct {
    Volume float64
    PrixHT float64
}

func ComputeBilansActivitesParSaison(db *sqlx.DB, debutSaison string, activites []*Activite) (result []*BilanActivitesParSaison, err error) {
	//
	activitesParSaison, err := computeActivitesParSaison(db, debutSaison, activites)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel computeActivitesParSaison()")
	}
	//
	result = []*BilanActivitesParSaison{}
	for _, activiteParSaison := range activitesParSaison {
		if len(activiteParSaison.Activites) == 0 {
			continue // exclut les saisons sans activités des bilans
		}
        newRes := BilanActivitesParSaison{
            Datedeb:                            activiteParSaison.Datedeb,
            Datefin:                            activiteParSaison.Datefin,
        }
        newRes.TotalActivitesParValoEtProprio = make(map[string]map[int]VolumePrixHT)
        newRes.TotalActivitesPlaquettesParProprio = make(map[int]VolumePrixHT)
        newRes.TotalVentePlaquettesParProprio = make(map[int]float64)
        //
		for _, activite := range activiteParSaison.Activites {
			valo := activite.TypeValo
			// initialisation
			if _, ok := newRes.TotalActivitesParValoEtProprio[valo]; !ok {
				newRes.TotalActivitesParValoEtProprio[valo] = map[int]VolumePrixHT{}
			}
			// on répartit systématiquement le prix et le volume par proprio (même si un seul proprio, tant pis)
			err = activite.ComputeSurfaceParProprio(db)
			if err != nil {
				return result, werr.Wrapf(err, "Erreur appel activite.ComputeSurfaceParProprio()")
			}
			for idProprio, surface := range activite.SurfaceParProprio {
			    // initialisation
                if _, ok := newRes.TotalActivitesParValoEtProprio[valo][idProprio]; !ok {
                    newRes.TotalActivitesParValoEtProprio[valo][idProprio] = VolumePrixHT{}
                }
				// ici, prix par proprio et volume par proprio proportionnels à la surface
				entry := newRes.TotalActivitesParValoEtProprio[valo][idProprio]
				entry.PrixHT += activite.PrixHT * surface / activite.SurfaceTotale
				entry.Volume += activite.Volume * surface / activite.SurfaceTotale
				newRes.TotalActivitesParValoEtProprio[valo][idProprio] = entry
			}
		} // end loop sur activiteParSaison.Activites
		// si besoin, trasfère les activités plaquettes dans le bon champ
        if _, ok := newRes.TotalActivitesParValoEtProprio["PQ"]; ok {
            newRes.TotalActivitesPlaquettesParProprio = newRes.TotalActivitesParValoEtProprio["PQ"]
            delete(newRes.TotalActivitesParValoEtProprio, "PQ")
        }
        // ventes
        newRes.TotalVentePlaquettesParProprio, err = ComputeQuantiteVenteParProprio(db, activiteParSaison.Datedeb, activiteParSaison.Datefin)
        if err != nil {
            return result, werr.Wrapf(err, "Erreur appel activite.ComputeSurfaceParProprio()")
        }
        //
		result = append(result, &newRes)
	} // end loop sur activitesParSaison
	return result, nil
}

// Auxiliaire de ComputeBilansActivitesParSaison()
// Parmi les activités passées en paramètre, ne retient que les activités ayant lieu dans une saison donnée.
func computeActivitesParSaison(db *sqlx.DB, debutSaison string, activites []*Activite) (result []*ActiviteParSaison, err error) {
	limites, _, err := ComputeLimitesSaisons(db, debutSaison)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel ComputeLimitesSaisons()")
	}
	tiglib.ArrayReverse(limites)
	result = []*ActiviteParSaison{}
	for _, limite := range limites {
		newRes := ActiviteParSaison{Datedeb: limite[0], Datefin: limite[1], Activites: []*Activite{}}
		result = append(result, &newRes)
	}
	for _, activite := range activites {
		for i, _ := range result {
			if (activite.DateActivite.After(result[i].Datedeb) || activite.DateActivite.Equal(result[i].Datedeb)) && activite.DateActivite.Before(result[i].Datefin) {
				result[i].Activites = append(result[i].Activites, activite)
				break
			}
		}
	}
	return result, nil
}
