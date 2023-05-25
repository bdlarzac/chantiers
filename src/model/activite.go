/*
*****************************************************************************

		Activité générique - Représente toute activité = entité avec une date et souvent un prix.
		Stocké dans les tables = types d'activité concernée
	        chaufer
	        chautre
	        plaq

		@copyright  BDL, Bois du Larzac.
		@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
		@history    2023-03-09 14:54:36+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/generic/tiglib"
	"github.com/jmoiron/sqlx"
	"sort"
	"strconv"
	"time"
"fmt"
)

type Activite struct {
	Id           int
    TypeActivite string // "plaq", "chautre" ou "chaufer"
    Titre        string
    URL          string // Chaîne vide ou URL du détail de l'entité, ex "/plaq/32"
    DateActivite time.Time
    TypeValo     string
    CodeEssence  string
    Volume       float64
    Unite        string // pour le volume
	PUHT         float64
	PrixHT       float64
	// TVA          float64            // pas utilisé
	// NumFacture   string             // pas utilisé
	// DateFacture  time.Time          // pas utilisé
	// relations n-n - utiles pour l'application de certains filtres
	LiensParcelles []*ChantierParcelle
	UGs            []*UG
	Fermiers       []*Fermier
}

/*
Renvoie une map code activité => nom
*/
func GetActivitesMap() map[string]string {
	return map[string]string{
		"chaufer": "Chauffage fermier",
		"chautre": "Autre valorisation",
		"plaq":    "Ch. plaquettes",
	}
}

// ************************** Nom *******************************

func (a *Activite) String() string {
	return a.Titre
}

// ************************** Instance methods *******************************

func (a *Activite) ComputeLiensParcelles(db *sqlx.DB) (err error) {
    a.LiensParcelles, err = computeLiensParcellesOfChantier(db, a.TypeActivite, a.Id)
    if err != nil {
        return werr.Wrapf(err, "Erreur appel computeLiensParcellesOfChantier()")
    }
    return nil
}

func (a *Activite) ComputeFermiers(db *sqlx.DB) (err error) {
    a.Fermiers, err = computeFermiersOfChantier(db, a.TypeActivite, a.Id)
    if err != nil {
        return werr.Wrapf(err, "Erreur appel computeFermiersOfChantier()")
    }
    return nil
}

func (a *Activite) ComputeUGs(db *sqlx.DB) (err error) {
    a.UGs, err = computeUGsOfChantier(db, a.TypeActivite, a.Id)
    if err != nil {
        return werr.Wrapf(err, "Erreur appel computeUGsOfChantier()")
    }
    return nil
}

// ************************** Get many *******************************
func ComputeActivitesFromFiltres(db *sqlx.DB, filtres map[string][]string) (res []*Activite, err error) {
fmt.Printf("filtres = %+v\n",filtres)
	res = []*Activite{}
	//
	// Première sélection, par filtre période
	//
	var tmp []*Activite
	//
	if len(filtres["valo"]) == 0 || tiglib.InArray("PQ", filtres["valo"]){
        tmp, err = computePlaqActivitesFromFiltrePeriode(db, filtres["periode"])
        if err != nil {
            return res, werr.Wrapf(err, "Erreur appel computePlaqActivitesFromFiltrePeriode()")
        }
        res = append(res, tmp...)
    }
	//
	// todo: calculer chautre que si len(filtre valo) == 0 ou contient une ou + valo parmi BO PL CH PP PI
	tmp, err = computeChautreActivitesFromFiltrePeriode(db, filtres["periode"])
	if err != nil {
		return res, werr.Wrapf(err, "Erreur appel computeChautreActivitesFromFiltrePeriode()")
	}
	res = append(res, tmp...)
	//
	if len(filtres["valo"]) == 0 || tiglib.InArray("CF", filtres["valo"]){
        tmp, err = computeChauferActivitesFromFiltrePeriode(db, filtres["periode"])
        if err != nil {
            return res, werr.Wrapf(err, "Erreur appel computeChauferActivitesFromFiltrePeriode()")
        }
        res = append(res, tmp...)
    }
	//
	// Filtres suivants
	//
	if len(filtres["essence"]) != 0 {
		res = filtreActivite_essence(db, res, filtres["essence"])
	}
	//
	if len(filtres["valo"]) != 0 {
		res = filtreActivite_valo(db, res, filtres["valo"])
	}
	//
	if len(filtres["fermier"]) != 0 {
		for _, activite := range res {
			activite.ComputeFermiers(db)
		}
		res = filtreActivite_fermier(db, res, filtres["fermier"])
	}
	//
	//
	// préparation (faire le plus tard possible pour optimiser)
	//
	// on calcule toujours les UGs puisqu'on affiche une liste par UG
	for _, activite := range res {
		err = activite.ComputeUGs(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel ComputeUGs()")
		}
	}
	//
	if len(filtres["ug"]) != 0 {
		res = filtreActivite_ug(db, res, filtres["ug"])
	}
	if len(filtres["proprio"]) != 0 || len(filtres["parcelle"]) != 0 {
		for _, activite := range res {
			activite.ComputeLiensParcelles(db)
		}
	}
	//
	if len(filtres["parcelle"]) != 0 {
		res = filtreActivite_parcelle(db, res, filtres["parcelle"])
	}
	//
	if len(filtres["proprio"]) != 0 {
		res, err = filtreActivite_proprio(db, res, filtres["proprio"])
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel filtreActivite_proprio()")
		}
	}
	// Tri par date
	sortedRes := make(activiteSlice, 0, len(res))
	for _, elt := range res {
		sortedRes = append(sortedRes, elt)
	}
	sort.Sort(sortedRes)
	return sortedRes, nil
}

// ****************************************************************************************************
// ************************** Auxiliaires de ComputeActivitesFromFiltres() ****************************
// ****************************************************************************************************

// ************************** Auxiliaire pour trier par date *******************************

type activiteSlice []*Activite
func (p activiteSlice) Len() int           { return len(p) }
func (p activiteSlice) Less(i, j int) bool { return p[i].DateActivite.Before(p[j].DateActivite) }
func (p activiteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }


// ************************** Selection initiale, par période *******************************
// Fabriquent des activités

func computePlaqActivitesFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (res []*Activite, err error) {
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
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chantier := range chantiers {
		tmp, err := plaq2Activite(db, chantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel plaq2Activite()")
		}
		res = append(res, tmp)
	}
	return res, nil
}

func computeChautreActivitesFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (res []*Activite, err error) {
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
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chantier := range chantiers {
		tmp, err := chautre2Activite(db, chantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel chautre2Activite()")
		}
		res = append(res, tmp)
	}
	return res, nil
}

func computeChauferActivitesFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (res []*Activite, err error) {
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
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chantier := range chantiers {
		tmp, err := chaufer2Activite(db, chantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel chaufer2Activite()")
		}
		res = append(res, tmp)
	}
	return res, nil
}

// ************************** Conversion de struct vers une Activite *******************************
// Auxiliaires des fonctions compute*ActivitesFromFiltrePeriode()

func plaq2Activite(db *sqlx.DB, ch *Plaq) (a *Activite, err error) {
	a = &Activite{}
	a.Id = ch.Id
	a.TypeActivite = "plaq"
	a.Titre = ch.Titre
	a.URL = "/chantier/plaquette/" + strconv.Itoa(a.Id)
	a.DateActivite = ch.DateDebut
	a.TypeValo = "PQ"
	err = ch.ComputeVolume(db)
	a.CodeEssence = ch.Essence
	if err != nil {
		return a, werr.Wrapf(err, "Erreur appel ComputeVolume()")
	}
	a.Volume = ch.Volume
	a.Unite = "MA"
	return a, nil
}

func chautre2Activite(db *sqlx.DB, ch *Chautre) (a *Activite, err error) {
	a = &Activite{}
	a.Id = ch.Id
	a.TypeActivite = "chautre"
	a.Titre = ch.Titre
	a.URL = "/chantier/autre/" + strconv.Itoa(a.Id)
	a.DateActivite = ch.DateContrat
	a.TypeValo = ch.TypeValo
	a.CodeEssence = ch.Essence
	a.Volume = ch.VolumeRealise
	a.Unite = ch.Unite
	a.PUHT = ch.PUHT
	a.PrixHT = ch.PUHT * ch.VolumeRealise
	// inutile, à supprimer ?
	// a.TVA = ch.TVA
	// a.NumFacture = ch.NumFacture
	// a.DateFacture = ch.DateFacture
	return a, nil
}

func chaufer2Activite(db *sqlx.DB, ch *Chaufer) (a *Activite, err error) {
	a = &Activite{}
	a.Id = ch.Id
	a.TypeActivite = "chaufer"
	a.Titre = ch.Titre
	a.URL = "/chantier/chauffage-fermier/" + strconv.Itoa(a.Id)
	a.DateActivite = ch.DateChantier
	a.TypeValo = "CF"
	a.CodeEssence = ch.Essence
	a.Volume = ch.Volume
	a.Unite = ch.Unite
	return a, nil
}

// ************************** Filtres *******************************
// En entrée : liste d'activités
// En sortie : liste d'activités qui satisfont au filtre

func filtreActivite_essence(db *sqlx.DB, input []*Activite, filtre []string) (res []*Activite) {
	res = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			if a.CodeEssence == f {
				res = append(res, a)
				break
			}
		}
	}
	return res
}

func filtreActivite_valo(db *sqlx.DB, input []*Activite, filtre []string) (res []*Activite) {
	res = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			if a.TypeValo == f {
				res = append(res, a)
				break
			}
		}
	}
	return res
}

func filtreActivite_fermier(db *sqlx.DB, input []*Activite, filtre []string) (res []*Activite) {
	res = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, fermier := range a.Fermiers {
				if fermier.Id == id {
					res = append(res, a)
					break
				}
			}
		}
	}
	return res
}

func filtreActivite_ug(db *sqlx.DB, input []*Activite, filtre []string) (res []*Activite) {
	res = []*Activite{}
	// map pour ne pas inclure des activités en double.
	// Se produit si on demande des ugs voisines,
	// avec des activités communes à plusieurs ugs demandées
	m := map[int]bool{}
    ActiviteLoop:
	for _, a := range input {
	    idActivite := a.Id
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, ug := range a.UGs {
				if ug.Id == id {
				    if _, ok := m[id]; !ok {
                        res = append(res, a)
                        m[idActivite] = true
                        continue ActiviteLoop
                    }
				}
			}
		}
	}
	return res
}

func filtreActivite_parcelle(db *sqlx.DB, input []*Activite, filtre []string) (res []*Activite) {
	res = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, lienParcelle := range a.LiensParcelles {
				if lienParcelle.IdParcelle == id {
					res = append(res, a)
					break
				}
			}
		}
	}
	return res
}

func filtreActivite_proprio(db *sqlx.DB, input []*Activite, filtre []string) (res []*Activite, err error) {
	res = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, lienParcelle := range a.LiensParcelles {
				parcelle, err := GetParcelle(db, lienParcelle.IdParcelle)
				if err != nil {
					return res, werr.Wrapf(err, "Erreur appel GetParcelle()")
				}
				if parcelle.IdProprietaire == id {
					res = append(res, a)
					break
				}
			}
		}
	}
	return res, nil
}
