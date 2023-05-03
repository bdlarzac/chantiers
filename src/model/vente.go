/*
*****************************************************************************

		Vente générique
	    - vente de plaquettes forestières
	    - vente issue d'un chantier autres valorisations

		@copyright  BDL, Bois du Larzac.
		@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
		@history    2023-05-03 08:14:58+02:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type Vente struct {
	Id           int
	TypeVente    string // "plaq" ou "autre"
	Titre        string
	URL          string // Chaîne vide ou URL du détail de l'entité, ex "/plaq/32"
	DateVente time.Time
	TypeValo     string
	Volume       float64
	Unite        string // pour le volume
	CodeEssence  string
	PrixHT       float64
	PUHT         float64
	TVA          float64
	NumFacture   string
	DateFacture  time.Time
	Notes        string
	//
	Details interface{}
}

// ************************** Nom *******************************

func (v *Vente) String() string {
	return v.Titre
}

// ************************** Get many *******************************

func ComputeVentesFromFiltres(db *sqlx.DB, filtres map[string][]string) (result []*Vente, err error) {
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
	//
	// Filtres suivants
	//
	if len(filtres["essence"]) != 0 {
		result = filtreVente_essence(db, result, filtres["essence"])
	}
	//
	if len(filtres["valo"]) != 0 {
		result = filtreVente_valo(db, result, filtres["valo"])
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
	return result, nil
}

// ****************************************************************************************************
// ************************** Auxiliaires de ComputeVentesFromFiltres() *******************************
// ****************************************************************************************************

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
		vente, err := ventePlaq2Vente(db, vp)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel ventePlaq2Vente()")
		}
		result = append(result, vente)
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
		vente, err := chautre2Vente(db, ch)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel chautre2Vente()")
		}
		result = append(result, vente)
	}
	return result, nil
}

// ************************** Application des filtres *******************************
// En entrée : liste de ventes
// En sortie : liste de ventes qui satisfont au filtre

func filtreVente_essence(db *sqlx.DB, input []*Vente, filtre []string) (result []*Vente) {
	result = []*Vente{}
	for _, v := range input {
		for _, f := range filtre {
			if v.CodeEssence == f {
				result = append(result, v)
				break
			}
		}
	}
	return result
}

func filtreVente_valo(db *sqlx.DB, input []*Vente, filtre []string) (result []*Vente) {
	result = []*Vente{}
	for _, v := range input {
		for _, f := range filtre {
			if v.TypeValo == f {
				result = append(result, v)
				break
			}
		}
	}
	return result
}

// TODO implement
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

// ************************** Conversion de struct vers une Vente *******************************

func ventePlaq2Vente(db *sqlx.DB, vp *VentePlaq) (v *Vente, err error) {
    v = &Vente{}
	err = vp.ComputeClient(db) // Obligatoire pour pouvoir utiliser String()
	if err != nil {
		return v, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeClient()")
	}
	v.Id = vp.Id
	v.TypeVente = "plaq"
	v.Titre = vp.String()
	v.URL = "/vente/" + strconv.Itoa(vp.Id)
	v.DateVente = vp.DateVente
	v.TypeValo = "PQ"
	err = vp.ComputeQte(db)
	if err != nil {
		return v, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeQte()")
	}
	v.Volume = vp.Qte
	v.Unite = "MA"
	v.CodeEssence = "PS"
	v.PrixHT =  vp.PUHT * vp.Qte
	v.PUHT = vp.PUHT
	v.TVA = vp.TVA
	v.NumFacture = vp.NumFacture
	v.DateFacture = vp.DateFacture
	v.Notes = vp.Notes
	return v, nil
}

func chautre2Vente(db *sqlx.DB, ch *Chautre) (v *Vente, err error) {
    v = &Vente{}
	v.Id = ch.Id
	v.TypeVente = "autre"
	v.Titre = ch.Titre
//	v.URL = "/chantier/autre/" + strconv.Itoa(idChautre)
	v.DateVente = ch.DateContrat
	v.TypeValo = ch.TypeValo
	v.Volume = ch.VolumeRealise
	v.Unite = ch.Unite
	v.CodeEssence = ch.Essence
	v.PrixHT =  ch.PUHT * ch.VolumeRealise
	v.PUHT = ch.PUHT
	v.TVA = ch.TVA
	v.NumFacture = ch.NumFacture
	v.DateFacture = ch.DateFacture
	v.Notes = ch.Notes
	return v, nil
}
