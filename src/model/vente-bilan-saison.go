/*
Calcul de ventes regroupées par saison

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
)

/* Liste de ventes ayant lieu dans une période donnée */
type VenteParSaison struct {
	Datedeb   time.Time
	Datefin   time.Time
	Ventes    []*Vente
}

type BilanVentesParSaison struct {
	Datedeb               time.Time
	Datefin               time.Time
	TotalVentesParValo []*TotalVentesParValo
}

type TotalVentesParValo struct {
	TypeValo string
	Volume   float64
	Unite    string
	PrixHT   float64
}

func ComputeBilansVentesParSaison(db *sqlx.DB, debutSaison string, ventes []*Vente) (result []*BilanVentesParSaison, err error) {
    ventesParSaison, err := ComputeVentesParSaison(db, debutSaison, ventes)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel ComputeVentesParSaison()")
	}
	result = []*BilanVentesParSaison{}
	for _, venteParSaison := range(ventesParSaison){
	    if len(venteParSaison.Ventes) == 0 {
	        continue // exclut les saisons sans ventes des bilans
	    }
	    currentRes := BilanVentesParSaison{
	        Datedeb: venteParSaison.Datedeb,
	        Datefin: venteParSaison.Datefin,
	        TotalVentesParValo:[]*TotalVentesParValo{},
        }
	    // map intermédiaire
        mapValos := map[string]TotalVentesParValo{}
        for _, vente := range(venteParSaison.Ventes){
            valo := vente.TypeValo
            if _, ok := mapValos[valo]; !ok {
                mapValos[valo] = TotalVentesParValo{TypeValo: valo}
            }
            entry := mapValos[valo]
            entry.Volume += vente.Volume
            entry.Unite = vente.Unite /////////////////// ici faire conversion d'unité pour certaines valos ? ///////////////////
            entry.PrixHT += vente.PrixHT
            mapValos[valo] = entry
        }
        // utilise map pour remplir currentRes
        for valo, total := range(mapValos){
            newRes := TotalVentesParValo{
                TypeValo: valo,
                Volume: total.Volume,
                Unite: total.Unite,
                PrixHT: total.PrixHT,
            }
            currentRes.TotalVentesParValo = append(currentRes.TotalVentesParValo, &newRes)
        }
        result = append(result, &currentRes)
	}
	return result, nil
}

func ComputeVentesParSaison(db *sqlx.DB, debutSaison string, ventes []*Vente) (result []*VenteParSaison, err error) {
    limites, _, err := ComputeLimitesSaisons(db, debutSaison)
    tiglib.ArrayReverse(limites)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel ComputeLimitesSaisons()")
	}
	result = []*VenteParSaison{}
    for _, limite := range(limites){
        newRes := VenteParSaison{Datedeb: limite[0], Datefin: limite[1], Ventes: []*Vente{}}
        result = append(result, &newRes)
    }
    for _, vente := range(ventes){
        for i, _ := range(result){
            if (vente.DateVente.After(result[i].Datedeb) || vente.DateVente.Equal(result[i].Datedeb)) && vente.DateVente.Before(result[i].Datefin){
                result[i].Ventes = append(result[i].Ventes, vente)
                break
            }
        }
    }
	return result, nil
}

